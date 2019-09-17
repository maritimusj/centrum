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

	DefaultUserPasswordResetOk
	LogDeletedByUser
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
	ErrUserNotFound
	ErrUserExists
	ErrFailedDisableDefaultUser
	ErrFailedRemoveDefaultUser
	ErrFailedEditDefaultUserPerm
	ErrFailedRemoveUserSelf
	ErrFailedDisableUserSelf
	ErrRoleNotFound
	ErrPolicyNotFound
	ErrGroupNotFound
	ErrDeviceNotFound
	ErrMeasureNotFound
	ErrEquipmentNotFound
	ErrStateNotFound
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

func ApiResourcesMap() [][3]string {
	return [][3]string{
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
