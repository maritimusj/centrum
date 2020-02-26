package lang

import (
	"errors"

	"github.com/maritimusj/centrum/gate/web/resource"
	"github.com/maritimusj/centrum/gate/web/status"
)

const (
	_ = iota

	DefaultGroupTitle
	DefaultGroupDesc

	ResourceDefault
	ResourceApi
	ResourceGroup
	ResourceDevice
	ResourceMeasure
	ResourceEquipment
	ResourceState

	LogTrace
	LogDebug
	LogInfo
	LogWarn
	LogError
	LogFatal
	LogPanic

	Years
	Weeks
	Days
	Hours
	Minutes
	Seconds
	Milliseconds

	UserDefaultRoleDesc

	RoleSystemAdminTitle
	RoleOrganizationAdminTitle
	RoleGuestTitle

	MenuRoleGalleryTitle
	MenuRoleDevicesTitle
	MenuRoleAlertTitle
	MenuRoleStatsTitle
	MenuRoleExportTitle
	MenuRoleUsersTitle
	MenuRoleSystemSettingsTitle
	MenuRoleSysLogsTitle

	RoleSystemAdminDesc
	RoleOrganizationAdminDesc
	RoleGuestDesc
	MenuRoleGalleryDesc
	MenuRoleDevicesDesc
	MenuRoleAlertDesc
	MenuRoleStatsDesc
	MenuRoleExportDesc
	MenuRoleUsersDesc
	MenuRoleSystemSettingsDesc
	MenuRoleSysLogsDesc

	ResourceConfigBaseDetailTitle
	ResourceConfigBaseDetailDesc

	ResourceConfigBaseUpdateTitle
	ResourceConfigBaseUpdateDesc

	ResourceOrganizationCreateTitle
	ResourceOrganizationListTitle
	ResourceOrganizationDetailTitle
	ResourceOrganizationUpdateTitle
	ResourceOrganizationDeleteTitle
	ResourceResourceListTitle
	ResourceResourceDetailTitle

	OrganizationCreateDesc
	OrganizationListDesc
	OrganizationDetailDesc
	OrganizationUpdateDesc
	OrganizationDeleteDesc
	ResourceListDesc
	ResourceDetailDesc

	ResourceMyProfileDetailTitle
	ResourceMyProfileDetailDesc
	ResourceMyProfileUpdateTitle
	ResourceMyProfileUpdateDesc
	ResourceMyPermTitle
	ResourceMyPermDesc
	ResourceMyPermMultiTitle
	ResourceMyPermMultiDesc
	ResourceUserLogListTitle
	ResourceUserLogListDesc
	ResourceUserLogDeleteTitle
	ResourceUserLogDeleteDesc
	ResourceDeviceLogListTitle
	ResourceDeviceLogListDesc
	ResourceDeviceLogDeleteTitle
	ResourceDeviceLogDeleteDesc
	ResourceEquipmentLogListTitle
	ResourceEquipmentLogListDesc
	ResourceEquipmentLogDeleteTitle
	ResourceEquipmentLogDeleteDesc

	ResourceDeviceStatusTitle
	ResourceDeviceDataTitle
	ResourceDeviceCtrlTitle
	ResourceDeviceCHValueTitle
	ResourceDeviceStatisticsTitle

	ResourceDeviceStatusDesc
	ResourceDeviceDataDesc
	ResourceDeviceCtrlDesc
	ResourceDeviceCHValueDesc
	ResourceDeviceStatisticsDesc

	ResourceEquipmentStatusTitle
	ResourceEquipmentDataTitle
	ResourceEquipmentCtrlTitle
	ResourceEquipmentCHValueTitle
	ResourceEquipmentStatisticsTitle

	ResourceEquipmentStatusDesc
	ResourceEquipmentDataDesc
	ResourceEquipmentCtrlDesc
	ResourceEquipmentCHValueDesc
	ResourceEquipmentStatisticsDesc

	ResourceUserListTitle
	ResourceUserListDesc
	ResourceUserCreateTitle
	ResourceUserCreateDesc
	ResourceUserUpdateTitle
	ResourceUserUpdateDesc
	ResourceUserDetailTitle
	ResourceUserDetailDesc
	ResourceUserDeleteTitle
	ResourceUserDeleteDesc
	ResourceRoleListTitle
	ResourceRoleListDesc
	ResourceRoleCreateTitle
	ResourceRoleCreateDesc
	ResourceRoleUpdateTitle
	ResourceRoleUpdateDesc
	ResourceRoleDetailTitle
	ResourceRoleDetailDesc
	ResourceRoleDeleteTitle
	ResourceRoleDeleteDesc
	ResourceGroupListTitle
	ResourceGroupListDesc
	ResourceGroupCreateTitle
	ResourceGroupCreateDesc
	ResourceGroupDetailTitle
	ResourceGroupDetailDesc
	ResourceGroupUpdateTitle
	ResourceGroupUpdateDesc
	ResourceGroupDeleteTitle
	ResourceGroupDeleteDesc
	ResourceDeviceListTitle
	ResourceDeviceListDesc
	ResourceDeviceCreateTitle
	ResourceDeviceCreateDesc
	ResourceDeviceDetailTitle
	ResourceDeviceDetailDesc
	ResourceDeviceUpdateTitle
	ResourceDeviceUpdateDesc
	ResourceDeviceDeleteTitle
	ResourceDeviceDeleteDesc
	ResourceMeasureListTitle
	ResourceMeasureListDesc
	ResourceMeasureCreateTitle
	ResourceMeasureCreateDesc
	ResourceMeasureDetailTitle
	ResourceMeasureDetailDesc
	ResourceMeasureUpdateTitle
	ResourceMeasureUpdateDesc
	ResourceMeasureDeleteTitle
	ResourceMeasureDeleteDesc
	ResourceEquipmentListTitle
	ResourceEquipmentListDesc
	ResourceEquipmentCreateTitle
	ResourceEquipmentCreateDesc
	ResourceEquipmentDetailTitle
	ResourceEquipmentDetailDesc
	ResourceEquipmentUpdateTitle
	ResourceEquipmentUpdateDesc
	ResourceEquipmentDeleteTitle
	ResourceEquipmentDeleteDesc
	ResourceStateListTitle
	ResourceStateListDesc
	ResourceStateCreateTitle
	ResourceStateCreateDesc
	ResourceStateDetailTitle
	ResourceStateDetailDesc
	ResourceStateUpdateTitle
	ResourceStateUpdateDesc
	ResourceStateDeleteTitle
	ResourceStateDeleteDesc

	AlarmListTitle
	AlarmListDesc
	AlarmConfirmTitle
	AlarmConfirmDesc
	AlarmDeleteTitle
	AlarmDeleteDesc
	AlarmDetailTitle
	AlarmDetailDesc

	CommentListTitle
	CommentDetailTitle
	CommentCreateTitle
	CommentDeleteTitle

	CommentListDesc
	CommentDetailDesc
	CommentCreateDesc
	CommentDeleteDesc

	ResourceLogListTitle
	ResourceLogListDesc
	ResourceLogDeleteTitle
	ResourceLogDeleteDesc
	ResourceLogLevelListTitle
	ResourceLogLevelListDesc

	ConfirmAdminPassword

	DefaultUserPasswordResetOk
	FlushDBOk
	LogDeletedByUser
	DeviceLogDeletedByUser

	CreateOrgOk
	CreateOrgFail
	DeleteOrgOk
	DeleteOrgFail

	AdminCreateUserOk
	AdminUpdateUserOk
	AdminDeleteUserOk

	UserCreateDeviceOk
	UserUpdateDeviceOk
	UserDeleteDeviceOk

	UserCreateEquipmentOk
	UserUpdateEquipmentOk
	UserDeleteEquipmentOk

	AlarmUnconfirmed
	AlarmConfirmed

	SysBriefTitle
	SysBriefDesc

	DataExportTitle
	DataExportDesc

	UserLoginOk
	UserLoginFailedCauseDisabled
	UserLoginFailedCausePasswordWrong
	UserProfileUpdateOk

	ExportInitialized
	ExportingData
	ArrangingData
	WritingData
	ExportReady
)

var (
	resourceGroupsMap map[resource.Class]string
	alarmStatusDesc   map[int]string
)

func load() {
	resourceGroupsMap = map[resource.Class]string{
		resource.Default:   Str(ResourceDefault),
		resource.Api:       Str(ResourceApi),
		resource.Group:     Str(ResourceGroup),
		resource.Device:    Str(ResourceDevice),
		resource.Measure:   Str(ResourceMeasure),
		resource.Equipment: Str(ResourceEquipment),
		resource.State:     Str(ResourceState),
	}

	alarmStatusDesc = map[int]string{
		status.Unconfirmed: Str(AlarmUnconfirmed),
		status.Confirmed:   Str(AlarmConfirmed),
	}
}

var (
	RoleSystemAdminName       = "__sys__"
	RoleOrganizationAdminName = "__admin__"
	RoleGuestName             = "__guest__"
)

//后台菜单对应的角色
var (
	MenuRoleGallery = "__menu_gallery__"  //设备总览
	MenuRoleDevices = "__menu_devices__"  //设备管理
	MenuRoleAlert   = "__menu_alert__"    //报警查询
	MenuRoleStats   = "__menu_stats__"    //趋势图
	MenuRoleExport  = "__menu_export__"   //导出报表
	MenuRoleUsers   = "__menu_users__"    //用户、权限管理
	MenuRoleSysLogs = "__menu_sys_logs__" //系统日志
)

//默认角色，名称，说明及权限集合
func DefaultRoles() map[[3]string][]string {
	return map[[3]string][]string{
		{RoleSystemAdminName, Str(RoleSystemAdminTitle), Str(RoleSystemAdminDesc)}:                   resource.SystemAdmin,
		{RoleOrganizationAdminName, Str(RoleOrganizationAdminTitle), Str(RoleOrganizationAdminDesc)}: resource.OrganizationAdmin,
		{RoleGuestName, Str(RoleGuestTitle), Str(RoleGuestDesc)}:                                     resource.Guest,

		//后台菜单角色
		{MenuRoleGallery, Str(MenuRoleGalleryTitle), Str(MenuRoleGalleryDesc)}: resource.MenuRoleGallery,
		{MenuRoleDevices, Str(MenuRoleDevicesTitle), Str(MenuRoleDevicesDesc)}: resource.MenuRoleDevices,
		{MenuRoleAlert, Str(MenuRoleAlertTitle), Str(MenuRoleAlertDesc)}:       resource.MenuRoleAlert,
		{MenuRoleStats, Str(MenuRoleStatsTitle), Str(MenuRoleStatsDesc)}:       resource.MenuRoleStats,
		{MenuRoleUsers, Str(MenuRoleUsersTitle), Str(MenuRoleUsersDesc)}:       resource.MenuRoleUsers,
		{MenuRoleExport, Str(MenuRoleExportTitle), Str(MenuRoleExportDesc)}:    resource.MenuRoleExport,
		{MenuRoleSysLogs, Str(MenuRoleSysLogsTitle), Str(MenuRoleSysLogsDesc)}: resource.MenuRoleSysLogs,
	}
}

//api资源的名称，标题和说明
func ApiResourcesMap() [][3]string {
	return [][3]string{

		{resource.ConfigBaseDetail, Str(ResourceConfigBaseDetailTitle), Str(ResourceConfigBaseDetailDesc)},
		{resource.ConfigBaseUpdate, Str(ResourceConfigBaseUpdateTitle), Str(ResourceConfigBaseUpdateDesc)},

		{resource.OrganizationCreate, Str(ResourceOrganizationCreateTitle), Str(OrganizationCreateDesc)},
		{resource.OrganizationList, Str(ResourceOrganizationListTitle), Str(OrganizationListDesc)},
		{resource.OrganizationDetail, Str(ResourceOrganizationDetailTitle), Str(OrganizationDetailDesc)},
		{resource.OrganizationUpdate, Str(ResourceOrganizationUpdateTitle), Str(OrganizationUpdateDesc)},
		{resource.OrganizationDelete, Str(ResourceOrganizationDeleteTitle), Str(OrganizationDeleteDesc)},

		{resource.ResourceList, Str(ResourceResourceListTitle), Str(ResourceListDesc)},
		{resource.ResourceDetail, Str(ResourceResourceDetailTitle), Str(ResourceDetailDesc)},

		{resource.MyProfileDetail, Str(ResourceMyProfileDetailTitle), Str(ResourceMyProfileDetailDesc)},
		{resource.MyProfileUpdate, Str(ResourceMyProfileUpdateTitle), Str(ResourceMyProfileUpdateDesc)},
		{resource.MyPerm, Str(ResourceMyPermTitle), Str(ResourceMyPermDesc)},
		{resource.MyPermMulti, Str(ResourceMyPermMultiTitle), Str(ResourceMyPermMultiDesc)},
		{resource.UserLogList, Str(ResourceUserLogListTitle), Str(ResourceUserLogListDesc)},
		{resource.UserLogDelete, Str(ResourceUserLogDeleteTitle), Str(ResourceUserLogDeleteDesc)},
		{resource.DeviceLogList, Str(ResourceDeviceLogListTitle), Str(ResourceDeviceLogListDesc)},
		{resource.DeviceLogDelete, Str(ResourceDeviceLogDeleteTitle), Str(ResourceDeviceLogDeleteDesc)},
		{resource.EquipmentLogList, Str(ResourceEquipmentLogListTitle), Str(ResourceEquipmentLogListDesc)},
		{resource.EquipmentLogDelete, Str(ResourceEquipmentLogDeleteTitle), Str(ResourceEquipmentLogDeleteDesc)},

		{resource.DeviceStatus, Str(ResourceDeviceStatusTitle), Str(ResourceDeviceStatusDesc)},
		{resource.DeviceData, Str(ResourceDeviceDataTitle), Str(ResourceDeviceDataDesc)},
		{resource.DeviceCtrl, Str(ResourceDeviceCtrlTitle), Str(ResourceDeviceCtrlDesc)},
		{resource.DeviceCHValue, Str(ResourceDeviceCHValueTitle), Str(ResourceDeviceCHValueDesc)},
		{resource.DeviceStatistics, Str(ResourceDeviceStatisticsTitle), Str(ResourceDeviceStatisticsDesc)},

		{resource.EquipmentStatus, Str(ResourceEquipmentStatusTitle), Str(ResourceEquipmentStatusDesc)},
		{resource.EquipmentData, Str(ResourceEquipmentDataTitle), Str(ResourceEquipmentDataDesc)},
		{resource.EquipmentCtrl, Str(ResourceEquipmentCtrlTitle), Str(ResourceEquipmentCtrlDesc)},
		{resource.EquipmentCHValue, Str(ResourceEquipmentCHValueTitle), Str(ResourceEquipmentCHValueDesc)},
		{resource.EquipmentStatistics, Str(ResourceEquipmentStatisticsTitle), Str(ResourceEquipmentStatisticsDesc)},

		{resource.UserList, Str(ResourceUserListTitle), Str(ResourceUserListDesc)},
		{resource.UserCreate, Str(ResourceUserCreateTitle), Str(ResourceUserCreateDesc)},
		{resource.UserUpdate, Str(ResourceUserUpdateTitle), Str(ResourceUserUpdateDesc)},
		{resource.UserDetail, Str(ResourceUserDetailTitle), Str(ResourceUserDetailDesc)},
		{resource.UserDelete, Str(ResourceUserDeleteTitle), Str(ResourceUserDeleteDesc)},

		{resource.RoleList, Str(ResourceRoleListTitle), Str(ResourceRoleListDesc)},
		{resource.RoleCreate, Str(ResourceRoleCreateTitle), Str(ResourceRoleCreateDesc)},
		{resource.RoleUpdate, Str(ResourceRoleUpdateTitle), Str(ResourceRoleUpdateDesc)},
		{resource.RoleDetail, Str(ResourceRoleDetailTitle), Str(ResourceRoleDetailDesc)},
		{resource.RoleDelete, Str(ResourceRoleDeleteTitle), Str(ResourceRoleDeleteDesc)},

		{resource.GroupList, Str(ResourceGroupListTitle), Str(ResourceGroupListDesc)},
		{resource.GroupCreate, Str(ResourceGroupCreateTitle), Str(ResourceGroupCreateDesc)},
		{resource.GroupDetail, Str(ResourceGroupDetailTitle), Str(ResourceGroupDetailDesc)},
		{resource.GroupUpdate, Str(ResourceGroupUpdateTitle), Str(ResourceGroupUpdateDesc)},
		{resource.GroupDelete, Str(ResourceGroupDeleteTitle), Str(ResourceGroupDeleteDesc)},

		{resource.DeviceList, Str(ResourceDeviceListTitle), Str(ResourceDeviceListDesc)},
		{resource.DeviceCreate, Str(ResourceDeviceCreateTitle), Str(ResourceDeviceCreateDesc)},
		{resource.DeviceDetail, Str(ResourceDeviceDetailTitle), Str(ResourceDeviceDetailDesc)},
		{resource.DeviceUpdate, Str(ResourceDeviceUpdateTitle), Str(ResourceDeviceUpdateDesc)},
		{resource.DeviceDelete, Str(ResourceDeviceDeleteTitle), Str(ResourceDeviceDeleteDesc)},

		{resource.MeasureList, Str(ResourceMeasureListTitle), Str(ResourceMeasureListDesc)},
		{resource.MeasureCreate, Str(ResourceMeasureCreateTitle), Str(ResourceMeasureCreateDesc)},
		{resource.MeasureDetail, Str(ResourceMeasureDetailTitle), Str(ResourceMeasureDetailDesc)},
		{resource.MeasureUpdate, Str(ResourceMeasureUpdateTitle), Str(ResourceMeasureUpdateDesc)},
		{resource.MeasureDelete, Str(ResourceMeasureDeleteTitle), Str(ResourceMeasureDeleteDesc)},

		{resource.EquipmentList, Str(ResourceEquipmentListTitle), Str(ResourceEquipmentListDesc)},
		{resource.EquipmentCreate, Str(ResourceEquipmentCreateTitle), Str(ResourceEquipmentCreateDesc)},
		{resource.EquipmentDetail, Str(ResourceEquipmentDetailTitle), Str(ResourceEquipmentDetailDesc)},
		{resource.EquipmentUpdate, Str(ResourceEquipmentUpdateTitle), Str(ResourceEquipmentUpdateDesc)},
		{resource.EquipmentDelete, Str(ResourceEquipmentDeleteTitle), Str(ResourceEquipmentDeleteDesc)},

		{resource.StateList, Str(ResourceStateListTitle), Str(ResourceStateListTitle)},
		{resource.StateCreate, Str(ResourceStateCreateTitle), Str(ResourceStateCreateDesc)},
		{resource.StateDetail, Str(ResourceStateDetailTitle), Str(ResourceStateDetailDesc)},
		{resource.StateUpdate, Str(ResourceStateUpdateTitle), Str(ResourceStateUpdateDesc)},
		{resource.StateDelete, Str(ResourceStateDeleteTitle), Str(ResourceStateDeleteDesc)},

		{resource.AlarmList, Str(AlarmListTitle), Str(AlarmListDesc)},
		{resource.AlarmConfirm, Str(AlarmConfirmTitle), Str(AlarmConfirmDesc)},
		{resource.AlarmDelete, Str(AlarmDeleteTitle), Str(AlarmDeleteDesc)},
		{resource.AlarmDetail, Str(AlarmDetailTitle), Str(AlarmDetailDesc)},

		{resource.CommentList, Str(CommentListTitle), Str(CommentListDesc)},
		{resource.CommentDetail, Str(CommentDetailTitle), Str(CommentDetailDesc)},
		{resource.CommentCreate, Str(CommentCreateTitle), Str(CommentCreateDesc)},
		{resource.CommentDelete, Str(CommentDeleteTitle), Str(CommentDeleteDesc)},

		{resource.LogList, Str(ResourceLogListTitle), Str(ResourceLogListDesc)},
		{resource.LogDelete, Str(ResourceLogDeleteTitle), Str(ResourceLogDeleteDesc)},

		{resource.SysBrief, Str(SysBriefTitle), Str(SysBriefDesc)},
		{resource.DataExport, Str(DataExportTitle), Str(DataExportDesc)},
	}
}

func ResourceClassTitle(class resource.Class) string {
	if len(resourceGroupsMap) == 0 {
		load()
	}

	if v, ok := resourceGroupsMap[class]; ok {
		return v
	}

	panic(errors.New("unknown resource class"))
}

func AlarmStatusDesc(stats int) string {
	if len(alarmStatusDesc) == 0 {
		load()
	}
	if v, ok := alarmStatusDesc[stats]; ok {
		return v
	}
	panic(errors.New("unknown alarm status"))
}
