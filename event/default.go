package event

const (
	UserLog      = "user::log"
	DeviceLog    = "device::log"
	EquipmentLog = "Equipment::log"
)

const (
	_ = iota
	Created
	Updated
	Deleted
)
