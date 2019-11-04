package resource

const (
	Unknown = ""

	ConfigBaseDetail = "config.detail"
	ConfigBaseUpdate = "config.update"

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

	DeviceStatus  = "device.status"
	DeviceData    = "device.data"
	DeviceCtrl    = "device.ctrl"
	DeviceCHValue = "device.val"

	EquipmentStatus  = "equipment.status"
	EquipmentData    = "equipment.data"
	EquipmentCtrl    = "equipment.ctrl"
	EquipmentCHValue = "equipment.val"

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

	AlarmList    = "alarm.list"
	AlarmConfirm = "alarm.confirm"
	AlarmDelete  = "alarm.delete"
	AlarmDetail  = "alarm.detail"

	CommentList   = "comment.list"
	CommentDetail = "comment.detail"
	CommentCreate = "comment.create"
	CommentDelete = "comment.delete"

	LogLevelList = "log.level.list"
	LogList      = "log.list"
	LogDelete    = "log.delete"

	SysBrief = "sys.brief"
)

var (
	Guest = []string{
		MyProfileDetail,
		MyProfileUpdate,
		MyPerm,
		MyPermMulti,
		SysBrief,
	}

	OrganizationAdmin = append(Guest,
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

		DeviceStatus,
		DeviceData,
		DeviceCtrl,
		DeviceCHValue,

		EquipmentStatus,
		EquipmentData,
		EquipmentCtrl,
		EquipmentCHValue,

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

		AlarmList,
		AlarmConfirm,
		AlarmDelete,
		AlarmDetail,

		CommentList,
		CommentDetail,
		CommentCreate,
		CommentDelete,

		LogLevelList,
		LogList,
		LogDelete,

		SysBrief,
	)

	SystemAdmin = append(OrganizationAdmin,
		OrganizationCreate,
		OrganizationList,
		OrganizationDetail,
		OrganizationUpdate,
		OrganizationDelete,
	)
)

//后台菜单角色权限列表
var (
	MenuRoleGallery = []string{
		DeviceStatus,
		DeviceData,
		DeviceCtrl,
		DeviceCHValue,

		EquipmentStatus,
		EquipmentData,
		EquipmentCtrl,
		EquipmentCHValue,
	}

	MenuRoleDevices = []string{
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

		DeviceStatus,
		DeviceData,
		DeviceCtrl,
		DeviceCHValue,

		EquipmentStatus,
		EquipmentData,
		EquipmentCtrl,
		EquipmentCHValue,

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
	}

	MenuRoleAlert = []string{
		AlarmList,
		AlarmConfirm,
		AlarmDelete,
		AlarmDetail,

		CommentList,
		CommentDetail,
		CommentCreate,
	}

	MenuRoleStats = []string{
		Unknown,
	}

	MenuRoleSMSSettings = []string{
		Unknown,
	}

	MenuRoleUsers = []string{
		UserList,
		UserCreate,
		UserDetail,
		UserUpdate,
		UserDelete,
	}

	MenuRoleSystemSettings = []string{
		LogLevelList,
		LogList,
		LogDelete,
	}

	MenuRoleSysLogs = []string{
		LogLevelList,
		LogList,
		LogDelete,
	}
)
