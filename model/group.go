package model

//设备分组
type Group interface {
	DBEntry
	Profile

	Parent() Group
	Title() string

	SetTitle(title string) error
	SetParent(group Group) error

	AddDevice(devices ...interface{}) error
	RemoveDevice(devices ...interface{}) error

	AddEquipment(equipments ...interface{}) error
	RemoveEquipment(equipments ...interface{}) error
}
