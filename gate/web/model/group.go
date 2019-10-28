package model

import (
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
)

//设备分组
type Group interface {
	DBEntry
	Profile

	Resource

	Organization() (Organization, error)

	Parent() Group
	SetParent(group interface{})

	Title() string
	SetTitle(title string)

	Desc() string
	SetDesc(desc string)

	GetDeviceList(options ...helper2.OptionFN) ([]Device, int64, error)
	AddDevice(devices ...interface{}) error
	RemoveDevice(devices ...interface{}) error

	GetEquipmentList(options ...helper2.OptionFN) ([]Equipment, int64, error)
	AddEquipment(equipments ...interface{}) error
	RemoveEquipment(equipments ...interface{}) error
}
