package modules

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"longhorn/proxy/internal/constants/enum"
	"longhorn/proxy/internal/global"
	"longhorn/proxy/internal/storage"
)

type Server struct {
	// 唯一标识
	ID uint64 `json:"id,string" default:""`
	// 服务器名称
	Name string `json:"name"`
	// 地址
	Host string `json:"host"`
	// 端口
	Port uint16 `json:"port"`
	// 服务器类型
	ServerType enum.ServerType `json:"serverType"`
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

func (v *Server) GetHost() string {
	return fmt.Sprintf("%s:%d", v.Host, v.Port)
}

func GetServer(id uint64, db storage.Storage) (c *Server, err error) {
	err = db.Get(global.Config.ServerPrefix, id, c)
	return
}

type WebServiceServer struct {
	Server
}

func (v WebServiceServer) Marshal() (result []byte, err error) {
	buf := bytes.NewBuffer(result)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(v)
	return buf.Bytes(), err
}

func (v *WebServiceServer) Unmarshal(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(v)
	return
}

func CreateWebServiceServer(c *WebServiceServer, db storage.Storage) (id uint64, err error) {
	id, err = db.Create(global.Config.ServerPrefix, c)
	return
}

func GetWebServiceServer(id uint64, db storage.Storage) (c *WebServiceServer, err error) {
	err = db.Get(global.Config.ServerPrefix, id, c)
	return
}

func WalkWebServiceServers(start uint64, limit int64, walking func(e storage.Element) error, db storage.Storage) (nextID uint64, err error) {
	nextID, err = db.Walk(global.Config.ServerPrefix, start, limit, func() storage.Element {
		return &WebServiceServer{}
	}, walking)
	return
}

func UpdateWebServiceServer(c *WebServiceServer, db storage.Storage) (err error) {
	err = db.Update(global.Config.ServerPrefix, c)
	return
}

func DeleteWebServiceServer(id uint64, db storage.Storage) (err error) {
	err = db.Delete(global.Config.ServerPrefix, fmt.Sprintf("%d", id))
	return
}

type DatabaseServer struct {
	Server

	// 数据库用户名
	UserName string `json:"userName"`
	// 数据库密码
	Password string `json:"password"`
	// 数据库配置扩展
	Extends map[string]string `json:"extends"`
}

func (v DatabaseServer) Marshal() (result []byte, err error) {
	buf := bytes.NewBuffer(result)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(v)
	return buf.Bytes(), err
}

func (v *DatabaseServer) Unmarshal(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(v)
	return
}

func CreateDatabaseServer(c *DatabaseServer, db storage.Storage) (id uint64, err error) {
	id, err = db.Create(global.Config.ServerPrefix, c)
	return
}

func GetDatabaseServer(id uint64, db storage.Storage) (c *DatabaseServer, err error) {
	err = db.Get(global.Config.ServerPrefix, id, c)
	return
}

func WalkDatabaseServers(start uint64, limit int64, walking func(e storage.Element) error, db storage.Storage) (nextID uint64, err error) {
	nextID, err = db.Walk(global.Config.ServerPrefix, start, limit, func() storage.Element {
		return &DatabaseServer{}
	}, walking)
	return
}

func UpdateDatabaseServer(c *DatabaseServer, db storage.Storage) (err error) {
	err = db.Update(global.Config.ServerPrefix, c)
	return
}

func DeleteDatabaseServer(id uint64, db storage.Storage) (err error) {
	err = db.Delete(global.Config.ServerPrefix, fmt.Sprintf("%d", id))
	return
}
