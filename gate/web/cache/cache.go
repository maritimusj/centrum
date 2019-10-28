package cache

import (
	model2 "github.com/maritimusj/centrum/gate/web/model"
)

type Cache interface {
	Flush()
	Foreach(func(key string, obj interface{}))
	Save(obj interface{}) error
	Remove(obj interface{})

	LoadConfig(id int64) (model2.Config, error)
	LoadOrganization(id int64) (model2.Organization, error)
	LoadUser(id int64) (model2.User, error)
	LoadRole(id int64) (model2.Role, error)
	LoadPolicy(id int64) (model2.Policy, error)
	LoadGroup(id int64) (model2.Group, error)
	LoadDevice(id int64) (model2.Device, error)
	LoadMeasure(id int64) (model2.Measure, error)
	LoadEquipment(id int64) (model2.Equipment, error)
	LoadState(id int64) (model2.State, error)
	LoadApiResource(id int64) (model2.ApiResource, error)
	LoadAlarm(id int64) (model2.Alarm, error)
}
