package store

import (
	"github.com/kataras/iris"
	cache2 "github.com/maritimusj/centrum/gate/web/cache"
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	"time"
)

type Store interface {
	Close()
	Cache() cache2.Cache
	EraseAllData() error

	CreateConfig(name string, data interface{}) (model2.Config, error)
	RemoveConfig(cfg interface{}) error
	GetConfig(cfg interface{}) (model2.Config, error)
	GetConfigList(options ...helper2.OptionFN) ([]model2.Config, int64, error)

	MustGetUserFromContext(ctx iris.Context) model2.User
	InitDefaultRoles(org interface{}) error

	IsOrganizationExists(org interface{}) (bool, error)
	GetOrganization(org interface{}) (model2.Organization, error)
	CreateOrganization(name string, title string) (model2.Organization, error)
	RemoveOrganization(org interface{}) error
	GetOrganizationList(options ...helper2.OptionFN) ([]model2.Organization, int64, error)

	IsUserExists(user interface{}) (bool, error)
	GetUser(user interface{}) (model2.User, error)
	CreateUser(org interface{}, name string, password []byte, roles ...interface{}) (model2.User, error)
	RemoveUser(user interface{}) error
	GetUserList(options ...helper2.OptionFN) ([]model2.User, int64, error)

	IsRoleExists(role interface{}) (bool, error)
	GetRole(role interface{}) (model2.Role, error)
	CreateRole(org interface{}, name, title, desc string) (model2.Role, error)
	RemoveRole(role interface{}) error
	GetRoleList(options ...helper2.OptionFN) ([]model2.Role, int64, error)

	GetPolicy(roleID int64) (model2.Policy, error)
	GetPolicyFrom(roleID int64, res model2.Resource, action resource2.Action) (model2.Policy, error)
	RemovePolicy(policyID int64) error
	GetPolicyList(res model2.Resource, options ...helper2.OptionFN) ([]model2.Policy, int64, error)

	GetDeviceGroups(deviceID int64) ([]model2.Group, error)
	GetEquipmentGroups(deviceID int64) ([]model2.Group, error)

	GetGroup(groupID int64) (model2.Group, error)
	CreateGroup(org interface{}, title, desc string, parentID int64) (model2.Group, error)
	RemoveGroup(groupID int64) error
	GetGroupList(options ...helper2.OptionFN) ([]model2.Group, int64, error)

	GetDevice(deviceID int64) (model2.Device, error)
	CreateDevice(org interface{}, title string, data map[string]interface{}) (model2.Device, error)
	RemoveDevice(deviceID int64) error
	GetDeviceList(options ...helper2.OptionFN) ([]model2.Device, int64, error)

	CreateMeasure(deviceID int64, title string, tag string, kind resource2.MeasureKind) (model2.Measure, error)
	GetMeasure(measureID int64) (model2.Measure, error)
	GetMeasureFromTagName(deviceID int64, tagName string) (model2.Measure, error)
	RemoveMeasure(measureID int64) error
	GetMeasureList(options ...helper2.OptionFN) ([]model2.Measure, int64, error)

	GetEquipment(equipmentID int64) (model2.Equipment, error)
	CreateEquipment(org interface{}, title, desc string) (model2.Equipment, error)
	RemoveEquipment(equipmentID int64) error
	GetEquipmentList(options ...helper2.OptionFN) ([]model2.Equipment, int64, error)

	GetState(stateID int64) (model2.State, error)
	CreateState(equipmentID, measureID int64, title, desc, script string) (model2.State, error)
	RemoveState(stateID int64) error
	GetStateList(options ...helper2.OptionFN) ([]model2.State, int64, error)

	GetAlarm(alarmID int64) (model2.Alarm, error)
	CreateAlarm(device model2.Device, measureID int64, data map[string]interface{}) (model2.Alarm, error)
	RemoveAlarm(alarmID int64) error
	GetAlarmList(start, end *time.Time, options ...helper2.OptionFN) ([]model2.Alarm, int64, error)
	GetLastUnconfirmedAlarm(device model2.Device, measureID int64) (model2.Alarm, error)

	GetResourceGroupList() []interface{}
	GetResourceList(class resource2.Class, options ...helper2.OptionFN) ([]model2.Resource, int64, error)
	GetResource(class resource2.Class, resourceID int64) (model2.Resource, error)

	GetApiResourceList(options ...helper2.OptionFN) ([]model2.ApiResource, int64, error)
	GetApiResource(res interface{}) (model2.ApiResource, error)
	InitApiResource() error
}
