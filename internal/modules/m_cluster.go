package modules

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"longhorn/proxy/internal/constants/enum"
	"longhorn/proxy/internal/global"
	"longhorn/proxy/internal/storage"
	"longhorn/proxy/pkg"
)

type Cluster struct {
	// 唯一标识
	ID uint64 `json:"id,string" default:""`
	// 集群名称
	Name string `json:"name"`
	// 负载均衡类型
	LoadBalanceType enum.LoadBalanceType `json:"loadBalanceType"`
	// 负载均衡器
	loadBalancer LoadBalancer
}

func (v *Cluster) SetIdentity(id uint64) {
	v.ID = id
}

func (v Cluster) GetIdentity() uint64 {
	return v.ID
}

func (v Cluster) Marshal() (result []byte, err error) {
	buf := bytes.NewBuffer(result)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(v)
	return buf.Bytes(), err
}

func (v *Cluster) Unmarshal(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(v)
	return
}

func (v *Cluster) InitLoadBalancer() {
	switch v.LoadBalanceType {
	case enum.LOAD_BALANCE_TYPE__ROUND_ROBIN:
		v.loadBalancer = NewRoundRobin()
	default:
	}
}

func (v *Cluster) GetLoadBalancer() LoadBalancer {
	return v.loadBalancer
}

func CreateCluster(c *Cluster, db storage.Storage) (id uint64, err error) {
	if c.ID == 0 {
		id, err = pkg.Generator.GenerateUniqueID()
		if err != nil {
			return
		}
		c.SetIdentity(id)
	}
	id, err = db.Create(global.Config.ClusterPrefix, c)
	return
}

func GetCluster(id uint64, db storage.Storage) (c *Cluster, err error) {
	c = &Cluster{}
	err = db.Get(global.Config.ClusterPrefix, "id", id, c)
	return
}

func WalkClusters(start uint64, limit int64, walking func(e storage.Element) error, db storage.Storage) (nextID uint64, err error) {
	nextID, err = db.Walk(global.Config.ClusterPrefix, nil, "id", start, limit, func() storage.Element {
		return &Cluster{}
	}, walking)
	return
}

func UpdateCluster(c *Cluster, db storage.Storage) (err error) {
	condition := storage.WithConditionKey("id").Eq(c.ID)
	err = db.Update(global.Config.ClusterPrefix, condition, c)
	return
}

func DeleteCluster(id uint64, db storage.Storage) (err error) {
	condition := storage.WithConditionKey("id").Eq(fmt.Sprintf("%d", id))
	err = db.Delete(global.Config.ClusterPrefix, condition)
	return
}
