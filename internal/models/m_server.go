package models

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"longhorn/proxy/internal/constants/enum"
	"longhorn/proxy/internal/global"
	"longhorn/proxy/internal/storage"
)

type ServerContract interface {
	storage.Element
	GetHost() string
	GetType() enum.ServerType
}

type GeneralServer struct {
	Server
}

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

func (v *Server) GetType() enum.ServerType {
	return v.ServerType
}

func GetServer(id uint64, db storage.Storage) (c ServerContract, err error) {
	c = &GeneralServer{}
	err = db.Get(global.Config.ServerPrefix, "server.id", id, c)
	return
}

func WalkServers(start uint64, limit int64, walking func(e storage.Element) error, db storage.Storage) (nextID uint64, err error) {
	nextID, err = db.Walk(global.Config.ServerPrefix, nil, "server.id", start, limit, func() storage.Element {
		return &GeneralServer{}
	}, walking)
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

func GetWebServiceServer(id uint64, db storage.Storage) (c *WebServiceServer, err error) {
	c = &WebServiceServer{}
	err = db.Get(global.Config.ServerPrefix, "server.id", id, c)
	return
}

func WalkWebServiceServers(start uint64, limit int64, walking func(e storage.Element) error, db storage.Storage) (nextID uint64, err error) {
	nextID, err = db.Walk(global.Config.ServerPrefix, nil, "server.id", start, limit, func() storage.Element {
		return &WebServiceServer{}
	}, walking)
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

func GetDatabaseServer(id uint64, db storage.Storage) (c *DatabaseServer, err error) {
	c = &DatabaseServer{}
	err = db.Get(global.Config.ServerPrefix, "server.id", id, c)
	return
}

func WalkDatabaseServers(start uint64, limit int64, walking func(e storage.Element) error, db storage.Storage) (nextID uint64, err error) {
	nextID, err = db.Walk(global.Config.ServerPrefix, nil, "server.id", start, limit, func() storage.Element {
		return &DatabaseServer{}
	}, walking)
	return
}
