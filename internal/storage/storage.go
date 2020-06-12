package storage

import (
	"github.com/sirupsen/logrus"
	"longhorn/proxy/internal/constants/enum"
	"longhorn/proxy/internal/global"
)

type Storage interface {
	Close() error
	Create(prefix string, e Element) (uint64, error)
	Update(prefix string, condition *Condition, e Element) error
	Delete(prefix string, condition *Condition) error
	Get(prefix string, idField string, idVal uint64, target Element) error
	Walk(prefix string, condition *Condition, startField string, start uint64, limit int64, elementFactory func() Element, walking func(e Element) error) (uint64, error)
}

type Element interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	GetIdentity() uint64
	SetIdentity(uint64)
}

type UnionElement interface {
	Element
	GetUnionIdentity() string
}

var Database = &Delegate{}

type Delegate struct {
	driver Storage
}

func (d *Delegate) Create(prefix string, e Element) (uint64, error) {
	return d.driver.Create(prefix, e)
}

func (d *Delegate) Update(prefix string, condition *Condition, e Element) error {
	return d.driver.Update(prefix, condition, e)
}

func (d *Delegate) Delete(prefix string, condition *Condition) error {
	return d.driver.Delete(prefix, condition)
}

func (d *Delegate) Get(prefix string, idField string, idVal uint64, target Element) error {
	return d.driver.Get(prefix, idField, idVal, target)
}

func (d *Delegate) Walk(prefix string, condition *Condition, startField string, start uint64, limit int64, elementFactory func() Element, walking func(e Element) error) (uint64, error) {
	return d.driver.Walk(prefix, condition, startField, start, limit, elementFactory, walking)
}

func (d *Delegate) Close() error {
	return d.driver.Close()
}

func (d *Delegate) Init(dbConfig global.DBConfig) {
	var err error
	if dbConfig.DBType == enum.DB_TYPE__ETCD {
		d.driver, err = NewDBEtcd(dbConfig)
	} else if dbConfig.DBType == enum.DB_TYPE__MYSQL {
		d.driver, err = NewDBMysql(dbConfig)
	} else if dbConfig.DBType == enum.DB_TYPE__MONGODB {
		d.driver, err = NewDBMongo(dbConfig)
	}

	if err != nil {
		logrus.Panic(err)
	}
}
