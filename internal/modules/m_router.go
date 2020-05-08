package modules

import (
	"bytes"
	"encoding/gob"
	"longhorn/proxy/internal/global"
	"longhorn/proxy/internal/storage"
)

type Router struct {
	// 唯一标识
	ID uint64 `json:"id" default:""`
	// 集群名称
	Name string `json:"name"`
	// 地址
	Host string `json:"host"`
	// 端口
	Port uint16 `json:"port"`
}

func (v *Router) SetIdentity(id uint64) {
	v.ID = id
}

func (v Router) GetIdentity() uint64 {
	return v.ID
}

func (v Router) Marshal() (result []byte, err error) {
	buf := bytes.NewBuffer(result)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(v)
	return buf.Bytes(), err
}

func (v *Router) Unmarshal(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(v)
	return
}

func CreateRouter(c *Router, db storage.Storage) (id uint64, err error) {
	id, err = db.Create(global.Config.RouterPrefix, c)
	return
}

func GetRouter(id uint64, db storage.Storage) (c *API, err error) {
	err = db.Get(global.Config.RouterPrefix, id, c)
	return
}

func WalkRouters(start uint64, limit int64, walking func(e storage.Element) error, db storage.Storage) (nextID uint64, err error) {
	nextID, err = db.Walk(global.Config.RouterPrefix, start, limit, func() storage.Element {
		return &Router{}
	}, walking)
	return
}

func UpdateRouter(c *API, db storage.Storage) (err error) {
	err = db.Update(global.Config.RouterPrefix, c)
	return
}

func DeleteRouter(id uint64, db storage.Storage) (err error) {
	err = db.Delete(global.Config.RouterPrefix, id)
	return
}
