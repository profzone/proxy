package modules

import (
	"bytes"
	"encoding/gob"
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

func CreateRouter(c *Router, db storage.DB) (id uint64, err error) {
	id, err = db.Create(c)
	return
}

func GetRouter(id uint64, db storage.DB) (c *API, err error) {
	err = db.Get(id, c)
	return
}

func WalkRouters(start, limit int64, walking func(e storage.Element) error, db storage.DB) (count int64, err error) {
	count, err = db.Walk(start, limit, func() storage.Element {
		return &Router{}
	}, walking)
	return
}

func UpdateRouter(c *API, db storage.DB) (err error) {
	err = db.Update(c)
	return
}

func DeleteRouter(id uint64, db storage.DB) (err error) {
	err = db.Delete(id)
	return
}
