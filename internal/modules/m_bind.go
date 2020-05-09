package modules

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"longhorn/proxy/internal/global"
	"longhorn/proxy/internal/storage"
)

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
	id, err = db.Create(global.Config.ServerPrefix, c)
	return
}

func GetBind(id uint64, db storage.Storage) (c *Bind, err error) {
	err = db.Get(global.Config.ServerPrefix, id, c)
	return
}

func WalkBinds(start uint64, limit int64, walking func(e storage.Element) error, db storage.Storage) (nextID uint64, err error) {
	nextID, err = db.Walk(global.Config.ServerPrefix, start, limit, func() storage.Element {
		return &Bind{}
	}, walking)
	return
}

func UpdateBind(c *Bind, db storage.Storage) (err error) {
	err = db.Update(global.Config.ServerPrefix, c)
	return
}

func DeleteBind(id uint64, db storage.Storage) (err error) {
	err = db.Delete(global.Config.ServerPrefix, fmt.Sprintf("%d", id))
	return
}
