package modules

import (
	"longhorn/proxy/internal/models"
)

type Organization struct {
	// 唯一标识
	ID uint64
	// 集群名称
	Name string
}

func NewOrganization(model *models.Organization) *Organization {
	return &Organization{
		ID:   model.ID,
		Name: model.Name,
	}
}

func (v *Organization) SetIdentity(id uint64) {
	v.ID = id
}

func (v Organization) GetIdentity() uint64 {
	return v.ID
}
