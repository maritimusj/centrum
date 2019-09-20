package store

import (
	"github.com/maritimusj/centrum/helper"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
)

type Store interface {
	IsOrganizationExists(org interface{}) (bool, error)
	GetOrganization(org interface{}) (model.Organization, error)
	CreateOrganization(name string, title string) (model.Organization, error)
	RemoveOrganization(org interface{}) error
	GetOrganizationList(options ...helper.OptionFN) ([]model.Organization, int64, error)

	IsUserExists(user interface{}) (bool, error)
	GetUser(user interface{}) (model.User, error)
	CreateUser(org interface{}, name string, password []byte, roles... interface{}) (model.User, error)
	RemoveUser(user interface{}) error
	GetUserList(options ...helper.OptionFN) ([]model.User, int64, error)

	IsRoleExists(role interface{}) (bool, error)
	GetRole(role interface{}) (model.Role, error)
	CreateRole(org interface{}, name, title string) (model.Role, error)
	RemoveRole(role interface{}) error
	GetRoleList(options ...helper.OptionFN) ([]model.Role, int64, error)

	GetPolicy(roleID int64) (model.Policy, error)
	GetPolicyFrom(roleID int64, res model.Resource, action resource.Action) (model.Policy, error)
	RemovePolicy(policyID int64) error
	GetPolicyList(res model.Resource, options ...helper.OptionFN) ([]model.Policy, int64, error)

	GetGroup(groupID int64) (model.Group, error)
	CreateGroup(org interface{}, title, desc string, parentID int64) (model.Group, error)
	RemoveGroup(groupID int64) error
	GetGroupList(options ...helper.OptionFN) ([]model.Group, int64, error)

	GetDevice(deviceID int64) (model.Device, error)
	CreateDevice(org interface{}, title string, data map[string]interface{}) (model.Device, error)
	RemoveDevice(deviceID int64) error
	GetDeviceList(options ...helper.OptionFN) ([]model.Device, int64, error)

	CreateMeasure(deviceID int64, title string, tag string, kind resource.MeasureKind) (model.Measure, error)
	GetMeasure(measureID int64) (model.Measure, error)
	RemoveMeasure(measureID int64) error
	GetMeasureList(options ...helper.OptionFN) ([]model.Measure, int64, error)

	GetEquipment(equipmentID int64) (model.Equipment, error)
	CreateEquipment(org interface{}, title, desc string) (model.Equipment, error)
	RemoveEquipment(equipmentID int64) error
	GetEquipmentList(options ...helper.OptionFN) ([]model.Equipment, int64, error)

	GetState(stateID int64) (model.State, error)
	CreateState(equipmentID, measureID int64, title, desc, script string) (model.State, error)
	RemoveState(stateID int64) error
	GetStateList(options ...helper.OptionFN) ([]model.State, int64, error)

	GetResourceGroupList() []interface{}
	GetResourceList(class resource.Class, options ...helper.OptionFN) ([]model.Resource, int64, error)
	GetResource(class resource.Class, resourceID int64) (model.Resource, error)

	GetApiResourceList(options ...helper.OptionFN) ([]model.ApiResource, int64, error)
	GetApiResource(res interface{}) (model.ApiResource, error)
	InitApiResource() error

	InitDefaultRoles(org interface{}) error
}
