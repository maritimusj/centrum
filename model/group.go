package model

import (
	"github.com/maritimusj/centrum/helper"
	"github.com/maritimusj/centrum/resource"
)

//设备分组
type Group interface {
	DBEntry
	Profile

	resource.Resource

	Parent() Group
	SetParent(group Group) error

	Title() string
	SetTitle(title string) error

	Desc() string
	SetDesc(desc string) error

	GetDeviceList(options ...helper.OptionFN) ([]Device, int64, error)
	AddDevice(devices ...interface{}) error
	RemoveDevice(devices ...interface{}) error

	GetEquipmentList(options ...helper.OptionFN) ([]Equipment, int64, error)
	AddEquipment(equipments ...interface{}) error
	RemoveEquipment(equipments ...interface{}) error
}
