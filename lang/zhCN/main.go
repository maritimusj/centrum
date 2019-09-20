package zhCN

import "github.com/maritimusj/centrum/lang"

func init() {
	lang.Register(lang.ZhCN, strMap, errStrMap)
}

var (
	strMap = map[int]string{
		lang.ResourceDefault:   "默认分组",
		lang.ResourceApi:       "后台权限",
		lang.ResourceGroup:     "设备分组",
		lang.ResourceDevice:    "设备",
		lang.ResourceMeasure:   "点位",
		lang.ResourceEquipment: "自定义设备",
		lang.ResourceState:     "自定义点位",

		lang.LogTrace: "跟踪",
		lang.LogDebug: "调试",
		lang.LogInfo:  "信息",
		lang.LogWarn:  "警告",
		lang.LogError: "错误",
		lang.LogFatal: "严重",
		lang.LogPanic: "异常",

		lang.RoleSystemAdminTitle:       "系统管理员",
		lang.RoleOrganizationAdminTitle: "管理员",
		lang.RoleGuestTitle:             "普通用户",

		lang.ResourceOrganizationCreateTitle: "",
		lang.ResourceOrganizationListTitle:   "",
		lang.ResourceOrganizationDetailTitle: "",
		lang.ResourceOrganizationUpdateTitle: "",
		lang.ResourceOrganizationDeleteTitle: "",
		lang.ResourceResourceListTitle:       "",
		lang.ResourceResourceDetailTitle:     "",

		lang.OrganizationCreateDesc: "",
		lang.OrganizationListDesc:   "",
		lang.OrganizationDetailDesc: "",
		lang.OrganizationUpdateDesc: "",
		lang.OrganizationDeleteDesc: "",
		lang.ResourceListDesc:       "",
		lang.ResourceDetailDesc:     "",

		lang.ResourceMyProfileDetailTitle:    "",
		lang.ResourceMyProfileDetailDesc:     "",
		lang.ResourceMyProfileUpdateTitle:    "",
		lang.ResourceMyProfileUpdateDesc:     "",
		lang.ResourceMyPermTitle:             "",
		lang.ResourceMyPermDesc:              "",
		lang.ResourceMyPermMultiTitle:        "",
		lang.ResourceMyPermMultiDesc:         "",
		lang.ResourceUserLogListTitle:        "",
		lang.ResourceUserLogListDesc:         "",
		lang.ResourceUserLogDeleteTitle:      "",
		lang.ResourceUserLogDeleteDesc:       "",
		lang.ResourceDeviceLogListTitle:      "",
		lang.ResourceDeviceLogListDesc:       "",
		lang.ResourceDeviceLogDeleteTitle:    "",
		lang.ResourceDeviceLogDeleteDesc:     "",
		lang.ResourceEquipmentLogListTitle:   "",
		lang.ResourceEquipmentLogListDesc:    "",
		lang.ResourceEquipmentLogDeleteTitle: "",
		lang.ResourceEquipmentLogDeleteDesc:  "",

		lang.ResourceUserListTitle:   "用户列表",
		lang.ResourceUserListDesc:    "查看用户列表",
		lang.ResourceUserCreateTitle: "创建用户",
		lang.ResourceUserCreateDesc:  "创建新用户",
		lang.ResourceUserUpdateTitle: "编辑用户",
		lang.ResourceUserUpdateDesc:  "修改用户资料",
		lang.ResourceUserDetailTitle: "查看资料",
		lang.ResourceUserDetailDesc:  "查看用户资料",
		lang.ResourceUserDeleteTitle: "删除用户",
		lang.ResourceUserDeleteDesc:  "删除用户",

		lang.ResourceRoleListTitle:   "角色列表",
		lang.ResourceRoleListDesc:    "查看角色列表",
		lang.ResourceRoleCreateTitle: "创建角色",
		lang.ResourceRoleCreateDesc:  "创建新的角色",
		lang.ResourceRoleUpdateTitle: "编辑角色",
		lang.ResourceRoleUpdateDesc:  "修改角色权限",
		lang.ResourceRoleDetailTitle: "查看角色",
		lang.ResourceRoleDetailDesc:  "查看角色详情",
		lang.ResourceRoleDeleteTitle: "删除角色",
		lang.ResourceRoleDeleteDesc:  "删除角色",

		lang.ResourceGroupListTitle:   "分组列表",
		lang.ResourceGroupListDesc:    "查看全部设备分组",
		lang.ResourceGroupCreateTitle: "创建分组",
		lang.ResourceGroupCreateDesc:  "创建新的设备分组",
		lang.ResourceGroupDetailTitle: "查看分组",
		lang.ResourceGroupDetailDesc:  "查看分组详情",
		lang.ResourceGroupUpdateTitle: "编辑分组",
		lang.ResourceGroupUpdateDesc:  "修改分组信息",
		lang.ResourceGroupDeleteTitle: "删除分级",
		lang.ResourceGroupDeleteDesc:  "删除设备分组",

		lang.ResourceDeviceListTitle:   "设备列表",
		lang.ResourceDeviceListDesc:    "查看设备列表",
		lang.ResourceDeviceCreateTitle: "添加设备",
		lang.ResourceDeviceCreateDesc:  "添加新设备",
		lang.ResourceDeviceDetailTitle: "查看设备",
		lang.ResourceDeviceDetailDesc:  "查看设备详情",
		lang.ResourceDeviceUpdateTitle: "编辑设备",
		lang.ResourceDeviceUpdateDesc:  "修改设备详情",
		lang.ResourceDeviceDeleteTitle: "删除设备",
		lang.ResourceDeviceDeleteDesc:  "移除一个设备",

		lang.ResourceMeasureListTitle:   "点位列表",
		lang.ResourceMeasureListDesc:    "查看点位列表",
		lang.ResourceMeasureCreateTitle: "添加点位",
		lang.ResourceMeasureCreateDesc:  "添加一个新的点位",
		lang.ResourceMeasureDetailTitle: "查看点位",
		lang.ResourceMeasureDetailDesc:  "查看点位详情",
		lang.ResourceMeasureUpdateTitle: "编辑点位",
		lang.ResourceMeasureUpdateDesc:  "修改点位",
		lang.ResourceMeasureDeleteTitle: "删除点位",
		lang.ResourceMeasureDeleteDesc:  "删除点位",

		lang.ResourceEquipmentListTitle:   "自定义设备列表",
		lang.ResourceEquipmentListDesc:    "查看自定义设备列表",
		lang.ResourceEquipmentCreateTitle: "创建自定义设备",
		lang.ResourceEquipmentCreateDesc:  "创建一个新的自定义设备",
		lang.ResourceEquipmentDetailTitle: "查看自定义设备",
		lang.ResourceEquipmentDetailDesc:  "查看自定义设备详情",
		lang.ResourceEquipmentUpdateTitle: "编辑自定义设备",
		lang.ResourceEquipmentUpdateDesc:  "修改自定义设备",
		lang.ResourceEquipmentDeleteTitle: "删除自定义设备",
		lang.ResourceEquipmentDeleteDesc:  "删除一个自定设备",

		lang.ResourceStateListTitle:   "自定义点位列表",
		lang.ResourceStateListDesc:    "查看自定义点位列表",
		lang.ResourceStateCreateTitle: "添加自定义点位",
		lang.ResourceStateCreateDesc:  "添加一个新的自定义点位",
		lang.ResourceStateDetailTitle: "查看自定义点位",
		lang.ResourceStateDetailDesc:  "查看片定义点位详情",
		lang.ResourceStateUpdateTitle: "编辑自定义点位",
		lang.ResourceStateUpdateDesc:  "修改自定义点位",
		lang.ResourceStateDeleteTitle: "删除自定义点位",
		lang.ResourceStateDeleteDesc:  "删除一个自定义点位",

		lang.ResourceLogListTitle:      "系统日志",
		lang.ResourceLogListDesc:       "查看系统日志",
		lang.ResourceLogDeleteTitle:    "删除系统日志",
		lang.ResourceLogDeleteDesc:     "清空系统日志记录",
		lang.ResourceLogLevelListTitle: "",
		lang.ResourceLogLevelListDesc:  "",

		lang.DefaultUserPasswordResetOk: "默认用户密码已重置！",
		lang.LogDeletedByUser:           "管理 %s 清空日志！",

		lang.CreateOrgOk:   "创建组织 %s(%s) 成功！",
		lang.CreateOrgFail: "创建组织 %s(%s) 失败：%s",
		lang.DeleteOrgOk:   "删除组织 %s(%s) 成功！",
		lang.DeleteOrgFail: "删除组织 %s(%s) 失败：%s",

		lang.CreateDeviceOk:   "创建设备 %s 成功！",
		lang.CreateDeviceFail: "创建设备失败：%s",
		lang.UpdateDeviceOk:   "更新设备 %s 成功！",
		lang.UpdateDeviceFail: "更新设备 %s 失败：%s",
		lang.DeleteDeviceOk:   "删除设备 %s 成功！",
		lang.DeleteDeviceFail: "删除设备 %s 失败：%s",
	}

	errStrMap = map[lang.ErrorCode]string{
		lang.Ok:                                 "成功！",
		lang.ErrUnknown:                         "未知错误！",
		lang.ErrUnknownLang:                     "未知语言区域！",
		lang.ErrInternal:                        "系统错误: %s",
		lang.ErrInvalidConnStr:                  "数据库连接参数不正确！",
		lang.ErrInvalidRequestData:              "不正确的请求数据！",
		lang.ErrTokenExpired:                    "请先登录！",
		lang.ErrNoPermission:                    "没有权限",
		lang.ErrCacheNotFound:                   "缓存中没有数据！",
		lang.ErrOrganizationNotFound:            "组织机构不存在！",
		lang.ErrOrganizationDifferent:           "不同的组织机构！",
		lang.ErrOrganizationExists:              "组织机构已存在！",
		lang.ErrFailedRemoveDefaultOrganization: "不能删除默认的组织机构！",
		lang.ErrInvalidUser:                     "当前用户不可用或者登录超时！",
		lang.ErrUserDisabled:                    "用户已经被禁用！",
		lang.ErrPasswordWrong:                   "密码不正确！",
		lang.ErrInvalidResourceClassID:          "资源类型不正确！",
		lang.ErrApiResourceNotFound:             "错误的api资源！",
		lang.ErrUnknownRole:                     "用户角色不正确！",
		lang.ErrUserNotFound:                    "用户不存在！",
		lang.ErrUserExists:                      "用户已经存在！",
		lang.ErrFailedDisableDefaultUser:        "不能禁用系统默认用户！",
		lang.ErrFailedRemoveDefaultUser:         "不能删除系统默认用户！",
		lang.ErrFailedEditDefaultUserPerm:       "不能编辑默认用户角色！",
		lang.ErrFailedRemoveUserSelf:            "不能删除当前登录的用户账号！",
		lang.ErrFailedDisableUserSelf:           "不能禁用当前登录的用户账号！",
		lang.ErrRoleNotFound:                    "用户角色不存在！",
		lang.ErrRoleExists:                      "用户角色已存在！",
		lang.ErrPolicyNotFound:                  "策略不存在！",
		lang.ErrGroupNotFound:                   "设备分组不存在！",
		lang.ErrDeviceNotFound:                  "设备不存在！",
		lang.ErrMeasureNotFound:                 "点位不存在！",
		lang.ErrEquipmentNotFound:               "自定设备不存在！",
		lang.ErrStateNotFound:                   "自定义点位不存在！",
		lang.ErrDeviceOrganizationDifferent:     "设备不属于同一组织机构！",
		lang.ErrEquipmentOrganizationDifferent:  "自定义设备不属于同一组织机构！",
		lang.ErrRecursiveDetected:               "检测到循环关系链，请联系管理员！",
	}
)
