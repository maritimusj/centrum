package model

import "github.com/maritimusj/centrum/resource"

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

	AddDevice(devices ...interface{}) error
	RemoveDevice(devices ...interface{}) error

	AddEquipment(equipments ...interface{}) error
	RemoveEquipment(equipments ...interface{}) error
}
