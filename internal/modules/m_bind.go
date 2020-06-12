package modules

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"longhorn/proxy/internal/global"
	"longhorn/proxy/internal/storage"
)

type BindRequestBody struct {
	// 集群ID
	ClusterID uint64 `json:"clusterID,string,omitempty"`
	// 服务器ID
	ServerID uint64 `json:"serverID,string,omitempty"`
}

type Bind struct {
	// 集群ID
	ClusterID uint64 `json:"clusterID,string"`
	// 服务器ID
	ServerID uint64 `json:"serverID,string"`
}

func (b *Bind) GetUnionIdentity() string {
	return fmt.Sprintf("%d/%d", b.ClusterID, b.ServerID)
}

func (b *Bind) SetIdentity(id uint64) {}

func (b Bind) GetIdentity() uint64 {
	return b.ClusterID
}

func (b Bind) Marshal() (result []byte, err error) {
	buf := bytes.NewBuffer(result)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(b)
	return buf.Bytes(), err
}

func (b *Bind) Unmarshal(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(b)
	return
}

func CreateBind(c *Bind, db storage.Storage) (id uint64, err error) {
	id, err = db.Create(global.Config.BindPrefix, c)
	return
}

func WalkBinds(clusterID uint64, start uint64, limit int64, walking func(e storage.Element) error, db storage.Storage) (nextID uint64, err error) {
	condition := storage.WithConditionKey("clusterid").Eq(clusterID)
	nextID, err = db.Walk(global.Config.BindPrefix, condition, "clusterid", start, limit, func() storage.Element {
		return &Bind{}
	}, walking)
	return
}

func DeleteBind(clusterID uint64, serverID uint64, db storage.Storage) (err error) {
	condition := storage.WithConditionKey("clusterid").Eq(clusterID)
	err = db.Delete(global.Config.BindPrefix, condition)
	return
}
