package store

import (
	"context"
	"github.com/maritimusj/centrum/helper"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
)

type Store interface {
	Open(context context.Context, option map[string]interface{}) error
	Close()

	GetUser(user interface{}) (model.User, error)
	CreateUser(name string, password []byte, role model.Role) (model.User, error)
	RemoveUser(userID int64) error
	GetUserList(options ...helper.OptionFN) ([]model.User, int64, error)

	GetRole(roleID int64) (model.Role, error)
	CreateRole(title string) (model.Role, error)
	RemoveRole(roleID int64) error
	GetRoleList(options ...helper.OptionFN) ([]model.Role, int64, error)

	//CreatePolicy(roleID int64, res resource.Resource, action resource.Action, effect resource.Effect) (model.Policy, error)
	GetPolicy(roleID int64) (model.Policy, error)
	CreatePolicyIsNotExists(roleID int64, res resource.Resource, action resource.Action, defaultEffect resource.Effect) (model.Policy, error)
	RemovePolicy(policyID int64) error
	GetPolicyList(res resource.Resource, options ...helper.OptionFN) ([]model.Policy, int64, error)

	GetGroup(groupID int64) (model.Group, error)
	CreateGroup(title string, parentID int64) (model.Group, error)
	RemoveGroup(groupID int64) error
	GetGroupList(options ...helper.OptionFN) ([]model.Group, int64, error)

	GetDevice(deviceID int64) (model.Device, error)
	CreateDevice(title string, data map[string]interface{}) (model.Device, error)
	RemoveDevice(deviceID int64) error
	GetDeviceList(options ...helper.OptionFN) ([]model.Device, int64, error)

	CreateMeasure(deviceID int64, title string, tag string, kind model.MeasureKind) (model.Measure, error)
	GetMeasure(measureID int64) (model.Measure, error)
	RemoveMeasure(measureID int64) error
	GetMeasureList(options ...helper.OptionFN) ([]model.Measure, int64, error)

	GetEquipment(equipmentID int64) (model.Equipment, error)
	CreateEquipment(title, desc string) (model.Equipment, error)
	RemoveEquipment(equipmentID int64) error
	GetEquipmentList(options ...helper.OptionFN) ([]model.Equipment, int64, error)

	GetState(stateID int64) (model.State, error)
	CreateState(equipmentID int64, measureID int64, title string, script string) (model.State, error)
	RemoveState(stateID int64) error
	GetStateList(options ...helper.OptionFN) ([]model.State, int64, error)

	GetResourceGroupList() []interface{}
	GetResourceList(class resource.Class, options ...helper.OptionFN) ([]resource.Resource, int64, error)
	GetResource(class resource.Class, resourceID int64) (resource.Resource, error)

	GetApiResourceList(options ...helper.OptionFN) ([]model.ApiResource, int64, error)
	GetApiResource(res interface{}) (model.ApiResource, error)
	InitApiResource() error
}
