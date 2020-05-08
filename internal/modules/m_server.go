package modules

import (
	"bytes"
	"encoding/gob"
	"longhorn/proxy/internal/global"
	"longhorn/proxy/internal/storage"
)

type Server struct {
	// 唯一标识
	ID uint64 `json:"id" default:""`
	// 集群名称
	Name string `json:"name"`
	// 地址
	Host string `json:"host"`
	// 端口
	Port uint16 `json:"port"`
}

func (v *Server) SetIdentity(id uint64) {
	v.ID = id
}

func (v Server) GetIdentity() uint64 {
	return v.ID
}

func (v Server) Marshal() (result []byte, err error) {
	buf := bytes.NewBuffer(result)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(v)
	return buf.Bytes(), err
}

func (v *Server) Unmarshal(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(v)
	return
}

func CreateServer(c *Server, db storage.Storage) (id uint64, err error) {
	id, err = db.Create(global.Config.ServerPrefix, c)
	return
}

func GetServer(id uint64, db storage.Storage) (c *Server, err error) {
	err = db.Get(global.Config.ServerPrefix, id, c)
	return
}

func WalkServers(start uint64, limit int64, walking func(e storage.Element) error, db storage.Storage) (nextID uint64, err error) {
	nextID, err = db.Walk(global.Config.ServerPrefix, start, limit, func() storage.Element {
		return &Server{}
	}, walking)
	return
}

func UpdateServer(c *Server, db storage.Storage) (err error) {
	err = db.Update(global.Config.ServerPrefix, c)
	return
}

func DeleteServer(id uint64, db storage.Storage) (err error) {
	err = db.Delete(global.Config.ServerPrefix, id)
	return
}
