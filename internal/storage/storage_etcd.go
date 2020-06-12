package storage

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"longhorn/proxy/internal/global"
	"math"
	"sync"
	"time"
)

type StorageEtcd struct {
	sync.RWMutex

	client   *clientv3.Client
	kvClient clientv3.KV

	idLock sync.Mutex
}

func NewDBEtcd(config global.DBConfig) (*StorageEtcd, error) {
	db := &StorageEtcd{}
	err := db.init(config.Endpoints)

	return db, err
}

func (s *StorageEtcd) init(endpoints []string) (err error) {
	s.client, err = clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return
	}

	s.kvClient = clientv3.NewKV(s.client)
	return
}

func (s *StorageEtcd) Create(prefix string, e Element) (uint64, error) {
	s.Lock()
	defer s.Unlock()

	return s.putElement(prefix, e)
}

func (s *StorageEtcd) Update(prefix string, condition *Condition, e Element) error {
	s.Lock()
	defer s.Unlock()

	_, err := s.putElement(prefix, e)
	return err
}

func (s *StorageEtcd) Delete(prefix string, condition *Condition) error {
	s.Lock()
	defer s.Unlock()

	var id string
	if str, ok := condition.Val.(string); ok {
		id = str
	} else {
		id = fmt.Sprintf("%v", condition.Val)
	}
	return s.deleteElement(prefix, id)
}

func (s *StorageEtcd) Get(prefix string, idField string, idVal uint64, target Element) error {
	s.RLock()
	defer s.RUnlock()

	err := s.getElement(prefix, idVal, target)
	return err
}

func (s *StorageEtcd) Walk(prefix string, condition *Condition, startField string, start uint64, limit int64, elementFactory func() Element, walking func(e Element) error) (uint64, error) {
	s.RLock()
	defer s.RUnlock()

	if condition == nil {
		condition = WithConditionKey("").Eq(0)
	}
	// TODO optimize condition
	prefix = fmt.Sprintf("%s/%d", prefix, condition.Val)
	nextStartID, err := s.getElements(prefix, start, limit, elementFactory, walking)

	return nextStartID, err
}

func (s *StorageEtcd) Close() error {
	return s.client.Close()
}

func (s *StorageEtcd) getKey(prefix string, id uint64) string {
	return fmt.Sprintf("%s/%d", prefix, id)
}

func (s *StorageEtcd) getUnionKey(prefix string, unionID string) string {
	return fmt.Sprintf("%s/%s", prefix, unionID)
}

func (s *StorageEtcd) withTxn() (clientv3.Txn, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(s.client.Ctx(), 10*time.Second)
	return s.kvClient.Txn(ctx), cancel
}

func (s *StorageEtcd) getResponse(key string, options ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(s.client.Ctx(), 10*time.Second)
	defer cancel()

	return s.kvClient.Get(ctx, key, options...)
}

func (s *StorageEtcd) get(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(s.client.Ctx(), 10*time.Second)
	defer cancel()

	resp, err := s.kvClient.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("not found")
	}

	return resp.Kvs[0].Value, nil
}

func (s *StorageEtcd) getElement(prefix string, id uint64, value Element) error {
	data, err := s.get(s.getKey(prefix, id))
	if err != nil {
		return err
	}

	if data == nil {

	}

	err = value.Unmarshal(data)
	return err
}

func (s *StorageEtcd) getElements(prefix string, start uint64, limit int64, elementFactory func() Element, walking func(element Element) error) (nextID uint64, err error) {
	// if start is negative then get all elements
	isGetAll := false
	if limit < 0 {
		limit = 10
		isGetAll = true
	}
	startKey := s.getKey(prefix, start)
	endKey := s.getKey(prefix, math.MaxUint64)

	withRange := clientv3.WithRange(endKey)
	withLimit := clientv3.WithLimit(limit)

	var resp *clientv3.GetResponse
	for {
		resp, err = s.getResponse(startKey, withRange, withLimit)
		if err != nil {
			return
		}

		for _, v := range resp.Kvs {
			ele := elementFactory()
			err = ele.Unmarshal(v.Value)
			if err != nil {
				return
			}

			err = walking(ele)

			nextID = ele.GetIdentity() + 1
		}

		// if start is negative then get all elements or all element have got
		if !isGetAll || int64(len(resp.Kvs)) < limit {
			break
		}
	}

	return
}

func (s *StorageEtcd) put(key, value string, options ...clientv3.OpOption) error {
	txn, cancel := s.withTxn()
	defer cancel()
	_, err := txn.Then(clientv3.OpPut(key, value, options...)).Commit()
	return err
}

func (s *StorageEtcd) putElement(prefix string, value Element) (uint64, error) {
	union, ok := value.(UnionElement)
	var key string
	if !ok {
		if value.GetIdentity() == 0 {
			return 0, fmt.Errorf("must set identity")
		}
		key = s.getKey(prefix, value.GetIdentity())
	} else {
		key = s.getUnionKey(prefix, union.GetUnionIdentity())
	}

	data, err := value.Marshal()
	if err != nil {
		return 0, err
	}

	return value.GetIdentity(), s.put(key, string(data))
}

func (s *StorageEtcd) delete(key string, options ...clientv3.OpOption) error {
	txn, cancel := s.withTxn()
	defer cancel()
	_, err := txn.Then(clientv3.OpDelete(key, options...)).Commit()
	return err
}

func (s *StorageEtcd) deleteElement(prefix string, id string) error {
	err := s.delete(s.getUnionKey(prefix, id))
	return err
}
