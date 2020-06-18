package models

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"longhorn/proxy/internal/global"
	"longhorn/proxy/internal/storage"
)

type Organization struct {
	// 唯一标识
	ID uint64 `json:"id" default:""`
	// 集群名称
	Name string `json:"name"`
}

func (v *Organization) SetIdentity(id uint64) {
	v.ID = id
}

func (v Organization) GetIdentity() uint64 {
	return v.ID
}

func (v Organization) Marshal() (result []byte, err error) {
	buf := bytes.NewBuffer(result)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(v)
	return buf.Bytes(), err
}

func (v *Organization) Unmarshal(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(v)
	return
}

func CreateOrganization(c *Organization, db storage.Storage) (id uint64, err error) {
	id, err = db.Create(global.Config.OrganizationPrefix, c)
	return
}

func GetOrganization(id uint64, db storage.Storage) (c *Organization, err error) {
	err = db.Get(global.Config.OrganizationPrefix, "id", id, c)
	return
}

func WalkOrganizations(start uint64, limit int64, walking func(e storage.Element) error, db storage.Storage) (nextID uint64, err error) {
	nextID, err = db.Walk(global.Config.OrganizationPrefix, nil, "id", start, limit, func() storage.Element {
		return &Organization{}
	}, walking)
	return
}

func UpdateOrganization(c *Organization, condition *storage.Condition, db storage.Storage) (err error) {
	err = db.Update(global.Config.OrganizationPrefix, condition, c)
	return
}

func DeleteOrganization(id uint64, db storage.Storage) (err error) {
	condition := storage.WithConditionKey("id").Eq(fmt.Sprintf("%d", id))
	err = db.Delete(global.Config.OrganizationPrefix, condition)
	return
}
