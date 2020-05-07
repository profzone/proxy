package modules

import (
	"bytes"
	"encoding/gob"
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

func CreateAPI(c *API, db storage.DB) (id uint64, err error) {
	id, err = db.Create(c)
	return
}

func GetAPI(id uint64, db storage.DB) (c *API, err error) {
	err = db.Get(id, c)
	return
}

func WalkAPIs(start, limit int64, walking func(e storage.Element) error, db storage.DB) (count int64, err error) {
	count, err = db.Walk(start, limit, func() storage.Element {
		return &API{}
	}, walking)
	return
}

func UpdateAPI(c *API, db storage.DB) (err error) {
	err = db.Update(c)
	return
}

func DeleteAPI(id uint64, db storage.DB) (err error) {
	err = db.Delete(id)
	return
}
