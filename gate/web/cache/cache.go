package cache

import (
	"github.com/maritimusj/centrum/gate/web/model"
)

type Cache interface {
	Flush()
	Foreach(func(key string, obj interface{}))
	Save(obj interface{}) error
	Remove(obj interface{})

	LoadConfig(interface{}) (model.Config, error)
	LoadOrganization(interface{}) (model.Organization, error)
	LoadUser(interface{}) (model.User, error)
	LoadRole(interface{}) (model.Role, error)
	LoadPolicy(interface{}) (model.Policy, error)
	LoadGroup(interface{}) (model.Group, error)
	LoadDevice(interface{}) (model.Device, error)
	LoadMeasure(interface{}) (model.Measure, error)
	LoadEquipment(interface{}) (model.Equipment, error)
	LoadState(interface{}) (model.State, error)
	LoadApiResource(interface{}) (model.ApiResource, error)
	LoadAlarm(interface{}) (model.Alarm, error)
	LoadComment(interface{}) (model.Comment, error)
}
