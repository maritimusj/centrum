package cache

import (
	"github.com/maritimusj/centrum/model"
)

type Cache interface {
	Save(obj interface{}) error
	Remove(obj interface{})

	LoadOrganization(id int64) (model.Organization, error)
	LoadUser(id int64) (model.User, error)
	LoadRole(id int64) (model.Role, error)
	LoadPolicy(id int64) (model.Policy, error)
	LoadGroup(id int64) (model.Group, error)
	LoadDevice(id int64) (model.Device, error)
	LoadMeasure(id int64) (model.Measure, error)
	LoadEquipment(id int64) (model.Equipment, error)
	LoadState(id int64) (model.State, error)
	LoadApiResource(id int64) (model.ApiResource, error)
}
