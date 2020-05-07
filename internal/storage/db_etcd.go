package storage

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"longhorn/proxy/internal/global"
	"longhorn/proxy/internal/pkg"
	"math"
	"sync"
	"time"
)

type DBEtcd struct {
	sync.RWMutex

	client   *clientv3.Client
	kvClient clientv3.KV

	globalPrefix  string
	clusterPrefix string
	idPrefix      string

	idLock      sync.Mutex
	idGenerator *pkg.GeneratorSnowFlake
}

func (d *DBEtcd) Create(e Element) (uint64, error) {
	d.Lock()
	defer d.Unlock()

	return d.putElement(d.clusterPrefix, e)
}

func (d *DBEtcd) Update(e Element) error {
	panic("implement me")
}

func (d *DBEtcd) Delete(id uint64) error {
	panic("implement me")
}

func (d *DBEtcd) Get(id uint64, target Element) error {
	d.RLock()
	defer d.RUnlock()

	err := d.getElement(d.clusterPrefix, id, target)
	return err
}

func (d *DBEtcd) Walk(start, limit int64, elementFactory func() Element, walking func(e Element) error) (int64, error) {
	d.RLock()
	defer d.RUnlock()

	nextStart, err := d.getElements(d.clusterPrefix, start, limit, elementFactory, walking)

	return nextStart, err
}

func (d *DBEtcd) Close() error {
	return d.client.Close()
}

func NewDBEtcd(endpoints []string, prefix string, idConfig global.SnowflakeConfig) (*DBEtcd, error) {
	db := &DBEtcd{
		globalPrefix:  prefix,
		clusterPrefix: fmt.Sprintf("%s/clusters", prefix),

		idGenerator: pkg.NewSnowflake(idConfig),
	}
	err := db.init(endpoints)

	return db, err
}

func (d *DBEtcd) init(endpoints []string) (err error) {
	d.client, err = clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return
	}

	d.kvClient = clientv3.NewKV(d.client)
	return
}

func (d *DBEtcd) getKey(prefix string, id uint64) string {
	return fmt.Sprintf("%s/%d", prefix, id)
}

func (d *DBEtcd) withTxn() (clientv3.Txn, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(d.client.Ctx(), 10*time.Second)
	return d.kvClient.Txn(ctx), cancel
}

func (d *DBEtcd) getResponse(key string, options ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(d.client.Ctx(), 10*time.Second)
	defer cancel()

	return d.kvClient.Get(ctx, key, options...)
}

func (d *DBEtcd) get(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(d.client.Ctx(), 10*time.Second)
	defer cancel()

	resp, err := d.kvClient.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, nil
	}

	return resp.Kvs[0].Value, nil
}

func (d *DBEtcd) getElement(prefix string, id uint64, value Element) error {
	data, err := d.get(d.getKey(prefix, id))
	if err != nil {
		return err
	}

	err = value.Unmarshal(data)
	return err
}

func (d *DBEtcd) getElements(prefix string, start int64, limit int64, elementFactory func() Element, walking func(element Element) error) (int64, error) {
	// if start is negative then get all elements
	isGetAll := false
	if start < 0 {
		start = 0
		isGetAll = true
	}
	startKey := d.getKey(prefix, uint64(start))
	endKey := d.getKey(prefix, math.MaxUint64)

	withRange := clientv3.WithRange(endKey)
	withLimit := clientv3.WithLimit(limit)

	for {
		resp, err := d.getResponse(startKey, withRange, withLimit)
		if err != nil {
			return 0, err
		}

		for _, v := range resp.Kvs {
			ele := elementFactory()
			err = ele.Unmarshal(v.Value)
			if err != nil {
				return 0, err
			}

			walking(ele)

			start = int64(ele.GetIdentity()) + 1
		}

		// if start is negative then get all elements or all element have got
		if !isGetAll || int64(len(resp.Kvs)) < limit {
			break
		}
	}

	return start, nil
}

func (d *DBEtcd) put(key, value string, options ...clientv3.OpOption) error {
	txn, cancel := d.withTxn()
	defer cancel()
	_, err := txn.Then(clientv3.OpPut(key, value, options...)).Commit()
	return err
}

func (d *DBEtcd) putElement(prefix string, value Element) (uint64, error) {
	if value.GetIdentity() == 0 {
		id, err := d.idGenerator.GenerateUniqueID()
		if err != nil {
			return 0, err
		}
		value.SetIdentity(id)
	}

	data, err := value.Marshal()
	if err != nil {
		return 0, err
	}

	return value.GetIdentity(), d.put(d.getKey(prefix, value.GetIdentity()), string(data))
}
