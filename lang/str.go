package lang

import (
	"errors"
	"github.com/maritimusj/centrum/resource"
)

const (
	_ = iota
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

	RoleSystemAdminTitle
	RoleOrganizationAdminTitle
	RoleGuestTitle

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

	ResourceLogListTitle
	ResourceLogListDesc
	ResourceLogDeleteTitle
	ResourceLogDeleteDesc
	ResourceLogLevelListTitle
	ResourceLogLevelListDesc

	DefaultUserPasswordResetOk
	LogDeletedByUser

	CreateOrgOk
	CreateOrgFail
	DeleteOrgOk
	DeleteOrgFail

	CreateDeviceOk
	CreateDeviceFail
	UpdateDeviceOk
	UpdateDeviceFail
	DeleteDeviceOk
	DeleteDeviceFail
)

type ErrorCode int

const (
	Ok ErrorCode = iota
	ErrUnknown
	ErrUnknownLang
	ErrInternal
	ErrInvalidConnStr
	ErrInvalidRequestData
	ErrTokenExpired
	ErrNoPermission
	ErrCacheNotFound
	ErrInvalidUser
	ErrUserDisabled
	ErrPasswordWrong
	ErrInvalidResourceClassID
	ErrApiResourceNotFound
	ErrUnknownRole
	ErrOrganizationNotFound
	ErrOrganizationExists
	ErrFailedRemoveDefaultOrganization
	ErrUserNotFound
	ErrUserExists
	ErrFailedDisableDefaultUser
	ErrFailedRemoveDefaultUser
	ErrFailedEditDefaultUserPerm
	ErrFailedRemoveUserSelf
	ErrFailedDisableUserSelf
	ErrRoleNotFound
	ErrRoleExists
	ErrPolicyNotFound
	ErrGroupNotFound
	ErrDeviceNotFound
	ErrMeasureNotFound
	ErrEquipmentNotFound
	ErrStateNotFound
	ErrDeviceOrganizationDifferent
	ErrEquipmentOrganizationDifferent
	ErrRecursiveDetected
)

var (
	resourceGroupsMap map[resource.Class]string
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
}

var (
	RoleSystemAdminName       = "__sys__"
	RoleOrganizationAdminName = "__admin__"
	RoleGuestName             = "__guest__"
)

func DefaultRoles() map[[2]string][]string {
	return map[[2]string][]string{
		{RoleSystemAdminName, Str(RoleSystemAdminTitle)}:             resource.SystemAdmin,
		{RoleOrganizationAdminName, Str(RoleOrganizationAdminTitle)}: resource.OrganizationAdmin,
		{RoleGuestName, Str(RoleGuestTitle)}:                         resource.Guest,
	}
}

func ApiResourcesMap() [][3]string {
	return [][3]string{
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

		{resource.LogList, Str(ResourceLogListTitle), Str(ResourceLogListDesc)},
		{resource.LogDelete, Str(ResourceLogDeleteTitle), Str(ResourceLogDeleteDesc)},
		{resource.LogLevelList, Str(ResourceLogLevelListTitle), Str(ResourceLogLevelListDesc)},
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
