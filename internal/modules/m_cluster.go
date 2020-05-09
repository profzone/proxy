package modules

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"longhorn/proxy/internal/constants/enum"
	"longhorn/proxy/internal/global"
	"longhorn/proxy/internal/storage"
)

type Cluster struct {
	// 唯一标识
	ID uint64 `json:"id,string" default:""`
	// 集群名称
	Name string `json:"name"`
	// 负载均衡类型
	LoadBalanceType enum.LoadBalanceType `json:"loadBalanceType"`
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

func CreateCluster(c *Cluster, db storage.Storage) (id uint64, err error) {
	id, err = db.Create(global.Config.ClusterPrefix, c)
	return
}

func GetCluster(id uint64, db storage.Storage) (c *Cluster, err error) {
	c = &Cluster{}
	err = db.Get(global.Config.ClusterPrefix, id, c)
	return
}

func WalkClusters(start uint64, limit int64, walking func(e storage.Element) error, db storage.Storage) (nextID uint64, err error) {
	nextID, err = db.Walk(global.Config.ClusterPrefix, start, limit, func() storage.Element {
		return &Cluster{}
	}, walking)
	return
}

func UpdateCluster(c *Cluster, db storage.Storage) (err error) {
	err = db.Update(global.Config.ClusterPrefix, c)
	return
}

func DeleteCluster(id uint64, db storage.Storage) (err error) {
	err = db.Delete(global.Config.ClusterPrefix, fmt.Sprintf("%d", id))
	return
}
