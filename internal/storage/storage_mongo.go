package storage

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"longhorn/proxy/internal/global"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type StorageMongo struct {
	sync.RWMutex

	timeout time.Duration
	client  *mongo.Client
	db      *mongo.Database
}

func NewDBMongo(config global.DBConfig) (*StorageMongo, error) {
	db := &StorageMongo{}
	err := db.init(config)
	return db, err
}

func (s *StorageMongo) init(config global.DBConfig) error {
	var uri = "mongodb://" + strings.Join(config.Endpoints, ",")
	var opts = options.Client().ApplyURI(uri)

	ctx, _ := context.WithTimeout(context.Background(), config.ConnectionTimeout)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	s.client = client
	s.db = client.Database(config.DatabaseName)
	s.timeout = config.ConnectionTimeout
	return nil
}

func (s *StorageMongo) Close() error {
	ctx, _ := context.WithTimeout(context.Background(), s.timeout)
	return s.client.Disconnect(ctx)
}

func (s *StorageMongo) Create(prefix string, e Element) (uint64, error) {
	s.Lock()
	defer s.Unlock()

	return s.putElement(prefix, e)
}

func (s *StorageMongo) Update(prefix string, condition *Condition, e Element) error {
	s.Lock()
	defer s.Unlock()

	collection := s.db.Collection(prefix)
	ctx, _ := context.WithTimeout(context.Background(), s.timeout)

	opts := options.Update()
	filter := bson.M{
		condition.Key: condition.Val,
	}
	_, err := collection.UpdateOne(ctx, filter, e, opts)
	return err
}

func (s *StorageMongo) Delete(prefix string, condition *Condition) error {
	s.Lock()
	defer s.Unlock()

	collection := s.db.Collection(prefix)
	ctx, _ := context.WithTimeout(context.Background(), s.timeout)

	opts := options.Delete()
	filter := bson.M{
		condition.Key: condition.Val,
	}
	_, err := collection.DeleteOne(ctx, filter, opts)
	return err
}

func (s *StorageMongo) Get(prefix string, idField string, idVal uint64, target Element) error {
	s.RLock()
	defer s.RUnlock()

	collection := s.db.Collection(prefix)
	ctx, _ := context.WithTimeout(context.Background(), s.timeout)

	opts := options.FindOne()
	filter := bson.M{
		idField: idVal,
	}
	result := collection.FindOne(ctx, filter, opts)
	if result.Err() != nil {
		return result.Err()
	}

	return result.Decode(target)
}

func (s *StorageMongo) Walk(prefix string, condition *Condition, startField string, start uint64, limit int64, elementFactory func() Element, walking func(e Element) error) (nextID uint64, err error) {
	s.RLock()
	defer s.RUnlock()

	collection := s.db.Collection(prefix)
	ctx, _ := context.WithTimeout(context.Background(), s.timeout)

	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(limit)
	}
	opts.SetSort(bson.D{
		{
			Key:   startField,
			Value: -1,
		},
	})
	var filter = bson.D{}
	if start != 0 {
		filter = bson.D{
			{startField, bson.M{"$lt": start}},
		}
	}
	if condition != nil {
		filter = append(filter, bson.E{condition.Key, condition.Val})
	}

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return 0, err
	}

	for cursor.Next(context.TODO()) {
		element := elementFactory()
		err = cursor.Decode(element)
		if err != nil {
			return 0, err
		}
		err = walking(element)
		nextID = element.GetIdentity() - 1
	}

	return
}

func (s *StorageMongo) putElement(prefix string, value Element) (uint64, error) {
	if value.GetIdentity() == 0 {
		return 0, fmt.Errorf("must set identity")
	}

	collection := s.db.Collection(prefix)
	ctx, _ := context.WithTimeout(context.Background(), s.timeout)
	_, err := collection.InsertOne(ctx, value)
	if err != nil {
		return 0, err
	}

	return value.GetIdentity(), nil
}
