package modules

import (
	"github.com/sirupsen/logrus"
	"longhorn/proxy/internal/constants/enum"
	"longhorn/proxy/internal/global"
)

type DB interface {
	Close() error
	CreateCluster(c *Cluster) (uint64, error)
	GetCluster(id uint64) (*Cluster, error)
	GetClusters() ([]*Cluster, error)
	WalkClusters(start int64, limit int64, walking func(element Element) error) (int64, error)
	UpdateCluster(c *Cluster) error
	DeleteCluster(id uint64) error
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

func (d *Delegate) CreateCluster(c *Cluster) (uint64, error) {
	return d.driver.CreateCluster(c)
}

func (d *Delegate) GetCluster(id uint64) (*Cluster, error) {
	return d.driver.GetCluster(id)
}

func (d *Delegate) GetClusters() ([]*Cluster, error) {
	return d.driver.GetClusters()
}

func (d *Delegate) WalkClusters(start int64, limit int64, walking func(element Element) error) (int64, error) {
	return d.driver.WalkClusters(start, limit, walking)
}

func (d *Delegate) UpdateCluster(c *Cluster) error {
	return d.driver.UpdateCluster(c)
}

func (d *Delegate) DeleteCluster(id uint64) error {
	return d.driver.DeleteCluster(id)
}

