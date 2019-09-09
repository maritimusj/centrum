package store

import (
	"context"
	"github.com/maritimusj/centrum/model"
)

type Store interface {
	Open(context context.Context, option map[string]interface{}) error
	Close()

	GetUser(userID int64) (model.User, error)
	CreateUser(name string, password []byte) (model.User, error)
	RemoveUser(userID int64) error
	GetUserList(options ...OptionFN) ([]model.User, int64, error)

	GetRole(roleID int64) (model.Role, error)
	CreateRole(title string) (model.Role, error)
	RemoveRole(roleID int64) error
	GetRoleList(options ...OptionFN) ([]model.Role, int64, error)

	CreatePolicy(roleID int64, resourceID int64, resourceGroupID int64, action model.Action, effect model.Effect) (model.Policy, error)
	RemovePolicy(policyID int64) error
	GetPolicyList(roleID int64, group model.ResourceClass, resourceID int64, options ...OptionFN) ([]model.Policy, int64, error)

	GetGroup(groupID int64) (model.Group, error)
	CreateGroup(title string, parentID int64) (model.Group, error)
	RemoveGroup(groupID int64) error
	GetGroupList(parentID int64, options ...OptionFN) ([]model.Group, int64, error)

	GetDevice(deviceID int64) (model.Device, error)
	CreateDevice(title string, data map[string]interface{}) (model.Device, error)
	RemoveDevice(deviceID int64) error
	GetDeviceList(options ...OptionFN) ([]model.Device, int64, error)

	CreateMeasure(deviceID int64, title string, tag string, kind model.MeasureKind) (model.Measure, error)
	GetMeasure(measureID int64) (model.Measure, error)
	RemoveMeasure(measureID int64) error
	GetMeasureList(deviceID int64, options ...OptionFN) ([]model.Measure, int64, error)

	GetEquipment(equipmentID int64) (model.Equipment, error)
	CreateEquipment(title, desc string) (model.Equipment, error)
	RemoveEquipment(equipmentID int64) error
	GetEquipmentList(options ...OptionFN) ([]model.Equipment, int64, error)

	GetState(stateID int64) (model.State, error)
	CreateState(equipmentID int64, measureID int64, title string, script string) (model.State, error)
	GetStateList(equipmentID int64, options ...OptionFN) ([]model.State, int64, error)
	RemoveState(stateID int64) error

	GetResourceList(group model.ResourceClass, options ...OptionFN) ([]model.Resource, int64, error)
	GetResource(groupID int, resourceID int64) (model.Resource, error)
	GetApiResource(routerName string, httpMethod string) (model.Resource, error)
}
