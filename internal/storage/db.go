package storage

import (
	"github.com/sirupsen/logrus"
	"longhorn/proxy/internal/constants/enum"
	"longhorn/proxy/internal/global"
)

type DB interface {
	Close() error
	Create(e Element) (uint64, error)
	Update(e Element) error
	Delete(id uint64) error
	Get(id uint64, target Element) error
	Walk(start, limit int64, elementFactory func() Element, walking func(e Element) error) (int64, error)
}

type Element interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	GetIdentity() uint64
	SetIdentity(uint64)
}

var Database = &Delegate{}

type Delegate struct {
	driver DB
}

func (d *Delegate) Create(e Element) (uint64, error) {
	return d.driver.Create(e)
}

func (d *Delegate) Update(e Element) error {
	return d.driver.Update(e)
}

func (d *Delegate) Delete(id uint64) error {
	return d.driver.Delete(id)
}

func (d *Delegate) Get(id uint64, target Element) error {
	return d.Get(id, target)
}

func (d *Delegate) Walk(start, limit int64, elementFactory func() Element, walking func(e Element) error) (int64, error) {
	return d.Walk(start, limit, elementFactory, walking)
}

func (d *Delegate) Close() error {
	return d.driver.Close()
}

func (d *Delegate) Init(dbConfig global.DBConfig, idConfig global.SnowflakeConfig) {
	var err error
	if dbConfig.DBType == enum.DB_TYPE__ETCD {
		d.driver, err = NewDBEtcd(dbConfig.Endpoints, dbConfig.Prefix, idConfig)
	}

	if err != nil {
		logrus.Panic(err)
	}
}
