package modules

import (
	"bytes"
	"encoding/gob"
	"longhorn/proxy/internal/constants/enum"
)

type Cluster struct {
	// 唯一标识
	ID uint64 `json:"id" default:""`
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

func CreateCluster(c *Cluster, db DB) (id uint64, err error) {
	id, err = db.CreateCluster(c)
	return
}

func GetCluster(id uint64, db DB) (c *Cluster, err error) {
	c, err = db.GetCluster(id)
	return
}

func GetClusters(db DB) (c []*Cluster, err error) {
	c, err = db.GetClusters()
	return
}

func UpdateCluster(c *Cluster, db DB) (err error) {
	err = db.UpdateCluster(c)
	return
}

func DeleteCluster(id uint64, db DB) (err error) {
	err = db.DeleteCluster(id)
	return
}
