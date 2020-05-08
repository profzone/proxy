package modules

import (
	"bytes"
	"encoding/gob"
	"longhorn/proxy/internal/global"
	"longhorn/proxy/internal/storage"
)

type API struct {
	// 唯一标识
	ID uint64 `json:"id" default:""`
	// 集群名称
	Name string `json:"name"`
	// 地址
	Host string `json:"host"`
	// 端口
	Port uint16 `json:"port"`
}

func (v *API) SetIdentity(id uint64) {
	v.ID = id
}

func (v API) GetIdentity() uint64 {
	return v.ID
}

func (v API) Marshal() (result []byte, err error) {
	buf := bytes.NewBuffer(result)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(v)
	return buf.Bytes(), err
}

func (v *API) Unmarshal(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(v)
	return
}

func CreateAPI(c *API, db storage.Storage) (id uint64, err error) {
	id, err = db.Create(global.Config.ApiPrefix, c)
	return
}

func GetAPI(id uint64, db storage.Storage) (c *API, err error) {
	err = db.Get(global.Config.ApiPrefix, id, c)
	return
}

func WalkAPIs(start uint64, limit int64, walking func(e storage.Element) error, db storage.Storage) (nextID uint64, err error) {
	nextID, err = db.Walk(global.Config.ApiPrefix, start, limit, func() storage.Element {
		return &API{}
	}, walking)
	return
}

func UpdateAPI(c *API, db storage.Storage) (err error) {
	err = db.Update(global.Config.ApiPrefix, c)
	return
}

func DeleteAPI(id uint64, db storage.Storage) (err error) {
	err = db.Delete(global.Config.ApiPrefix, id)
	return
}
