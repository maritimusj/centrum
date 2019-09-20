package resource

const (
	OrganizationCreate = "org.create"
	OrganizationList   = "org.list"
	OrganizationDetail = "org.detail"
	OrganizationUpdate = "org.update"
	OrganizationDelete = "org.delete"

	ResourceList   = "resource.list"
	ResourceDetail = "resource.detail"

	MyProfileDetail = "my.profile.detail"
	MyProfileUpdate = "my.profile.update"
	MyPerm          = "my.perm"
	MyPermMulti     = "my.perm.multi"

	UserLogList        = "user.log.list"
	UserLogDelete      = "user.log.delete"
	DeviceLogList      = "device.log.list"
	DeviceLogDelete    = "device.log.delete"
	EquipmentLogList   = "equip.log.list"
	EquipmentLogDelete = "equip.log.delete"

	UserList   = "user.list"
	UserCreate = "user.create"
	UserDetail = "user.detail"
	UserUpdate = "user.update"
	UserDelete = "user.delete"

	RoleList   = "role.list"
	RoleCreate = "role.create"
	RoleDetail = "role.detail"
	RoleUpdate = "role.update"
	RoleDelete = "role.delete"

	GroupList   = "group.list"
	GroupCreate = "group.create"
	GroupDetail = "group.detail"
	GroupUpdate = "group.update"
	GroupDelete = "group.delete"

	DeviceList   = "device.list"
	DeviceCreate = "device.create"
	DeviceDetail = "device.detail"
	DeviceUpdate = "device.update"
	DeviceDelete = "device.delete"

	MeasureList   = "measure.list"
	MeasureCreate = "measure.create"
	MeasureDetail = "measure.detail"
	MeasureUpdate = "measure.update"
	MeasureDelete = "measure.delete"

	EquipmentList   = "equipment.list"
	EquipmentCreate = "equipment.create"
	EquipmentDetail = "equipment.detail"
	EquipmentUpdate = "equipment.update"
	EquipmentDelete = "equipment.delete"

	StateList   = "state.list"
	StateCreate = "state.create"
	StateDetail = "state.detail"
	StateUpdate = "state.update"
	StateDelete = "state.delete"

	LogLevelList = "log.level.list"
	LogList      = "log.list"
	LogDelete    = "log.delete"
)

var (
	Guest = []string{
		MyProfileDetail,
		MyProfileUpdate,
		MyPerm,
		MyPermMulti,
		ResourceList,
		ResourceDetail,
	}

	OrganizationAdmin = append(Guest, []string{
		UserLogList,
		UserLogDelete,
		DeviceLogList,
		DeviceLogDelete,
		EquipmentLogList,
		EquipmentLogDelete,

		ResourceList,
		ResourceDetail,

		UserList,
		UserCreate,
		UserDetail,
		UserUpdate,
		UserDelete,

		RoleList,
		RoleCreate,
		RoleDetail,
		RoleUpdate,
		RoleDelete,

		GroupList,
		GroupCreate,
		GroupDetail,
		GroupUpdate,
		GroupDelete,

		DeviceList,
		DeviceCreate,
		DeviceDetail,
		DeviceUpdate,
		DeviceDelete,

		MeasureList,
		MeasureCreate,
		MeasureDetail,
		MeasureUpdate,
		MeasureDelete,

		EquipmentList,
		EquipmentCreate,
		EquipmentDetail,
		EquipmentUpdate,
		EquipmentDelete,

		StateList,
		StateCreate,
		StateDetail,
		StateUpdate,
		StateDelete,

		LogLevelList,
		LogList,
		LogDelete,
	}...)

	SystemAdmin = append(OrganizationAdmin, []string{
		OrganizationCreate,
		OrganizationList,
		OrganizationDetail,
		OrganizationUpdate,
		OrganizationDelete,
	}...)
)
