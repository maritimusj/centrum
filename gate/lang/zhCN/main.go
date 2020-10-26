package zhCN

import (
	"github.com/maritimusj/centrum/gate/lang"
)

func init() {
	lang.Register(lang.ZhCN, strMap, errStrMap)
}

var (
	strMap = map[lang.StrIndex]string{
		lang.DatetimeFormatterStr: "2006-01-02 15:04:05",
		lang.DefaultGroupTitle:    "默认分组",
		lang.DefaultGroupDesc:     "系统默认分组",

		lang.ResourceDefault:   "默认分组",
		lang.ResourceApi:       "后台权限",
		lang.ResourceGroup:     "设备分组",
		lang.ResourceDevice:    "设备",
		lang.ResourceMeasure:   "点位",
		lang.ResourceEquipment: "自定义设备",
		lang.ResourceState:     "自定义点位",

		lang.LogTrace:   "跟踪",
		lang.LogDebug:   "调试",
		lang.LogInfo:    "信息",
		lang.LogWarning: "警告",
		lang.LogError:   "错误",
		lang.LogFatal:   "严重",
		lang.LogPanic:   "异常",

		lang.Years:        "年",
		lang.Weeks:        "星期",
		lang.Days:         "天",
		lang.Hours:        "小时",
		lang.Minutes:      "分钟",
		lang.Seconds:      "秒",
		lang.Milliseconds: "毫秒",

		lang.AlarmUnconfirmed: "未确认",
		lang.AlarmConfirmed:   "已确认",

		lang.RoleSystemAdminTitle:       "系统管理员",
		lang.RoleOrganizationAdminTitle: "管理员",
		lang.RoleGuestTitle:             "普通用户",

		lang.AlarmNotifyTitle:   "警报通知",
		lang.AlarmNotifyContent: "点位：%s, 名称：%s，数值：%s 异常，请关注！",

		lang.MenuRoleGalleryTitle:        "设备总览",
		lang.MenuRoleDevicesTitle:        "设备管理",
		lang.MenuRoleAlertTitle:          "报警查询",
		lang.MenuRoleStatsTitle:          "趋势图",
		lang.MenuRoleUsersTitle:          "权限设定",
		lang.MenuRoleExportTitle:         "导出报表",
		lang.MenuRoleSystemSettingsTitle: "系统设定",
		lang.MenuRoleSysLogsTitle:        "系统日志",
		lang.MenuWebViewTitle:            "云视图",
		lang.MenuStreamViewTitle:         "云监控",
		lang.MenuWebViewDesc:             "查看、设定加载特定的网页",
		lang.MenuStreamViewDesc:          "查看、编辑远程监控视频",

		lang.UserDefaultRoleDesc: "",

		lang.RoleSystemAdminDesc:       "",
		lang.RoleOrganizationAdminDesc: "",
		lang.RoleGuestDesc:             "",

		lang.MenuRoleGalleryDesc:        "查看设备实时状态",
		lang.MenuRoleDevicesDesc:        "可以对设备进行管理，添加删除和编辑设备",
		lang.MenuRoleAlertDesc:          "可以查看管理警报信息",
		lang.MenuRoleStatsDesc:          "可以查看，管理统计数据",
		lang.MenuRoleExportDesc:         "可以导出系统设备和自定义设备的数据报表",
		lang.MenuRoleUsersDesc:          "可以管理用户，添加删除和编辑用户",
		lang.MenuRoleSystemSettingsDesc: "可以查看，修改系统设置",
		lang.MenuRoleSysLogsDesc:        "可以查看，删除系统日志",

		lang.ResourceConfigBaseDetailTitle: "",
		lang.ResourceConfigBaseDetailDesc:  "",

		lang.ResourceConfigBaseUpdateTitle: "",
		lang.ResourceConfigBaseUpdateDesc:  "",

		lang.ResourceConfigWebDetailTitle: "",
		lang.ResourceConfigWebDetailDesc:  "",

		lang.ResourceConfigWebUpdateTitle: "",
		lang.ResourceConfigWebUpdateDesc:  "",

		lang.ResourceConfigStreamDetailTitle: "",
		lang.ResourceConfigStreamDetailDesc:  "",

		lang.ResourceConfigStreamUpdateTitle: "",
		lang.ResourceConfigStreamUpdateDesc:  "",

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

		lang.ResourceMyProfileDetailTitle: "",
		lang.ResourceMyProfileDetailDesc:  "",

		lang.ResourceMyProfileUpdateTitle: "",
		lang.ResourceMyProfileUpdateDesc:  "",

		lang.ResourceMyPermTitle: "",
		lang.ResourceMyPermDesc:  "",

		lang.ResourceMyPermMultiTitle: "",
		lang.ResourceMyPermMultiDesc:  "",

		lang.ResourceUserLogListTitle: "",
		lang.ResourceUserLogListDesc:  "",

		lang.ResourceUserLogDeleteTitle: "",
		lang.ResourceUserLogDeleteDesc:  "",

		lang.ResourceDeviceLogListTitle: "",
		lang.ResourceDeviceLogListDesc:  "",

		lang.ResourceDeviceLogDeleteTitle: "",
		lang.ResourceDeviceLogDeleteDesc:  "",

		lang.ResourceEquipmentLogListTitle: "",
		lang.ResourceEquipmentLogListDesc:  "",

		lang.ResourceEquipmentLogDeleteTitle: "",
		lang.ResourceEquipmentLogDeleteDesc:  "",

		lang.ResourceUserListTitle: "",
		lang.ResourceUserListDesc:  "",

		lang.ResourceUserCreateTitle: "",
		lang.ResourceUserCreateDesc:  "",

		lang.ResourceUserUpdateTitle: "",
		lang.ResourceUserUpdateDesc:  "",

		lang.ResourceUserDetailTitle: "",
		lang.ResourceUserDetailDesc:  "",

		lang.ResourceUserDeleteTitle: "",
		lang.ResourceUserDeleteDesc:  "",

		lang.ResourceRoleListTitle: "",
		lang.ResourceRoleListDesc:  "",

		lang.ResourceRoleCreateTitle: "",
		lang.ResourceRoleCreateDesc:  "",

		lang.ResourceRoleUpdateTitle: "",
		lang.ResourceRoleUpdateDesc:  "",

		lang.ResourceRoleDetailTitle: "",
		lang.ResourceRoleDetailDesc:  "",

		lang.ResourceRoleDeleteTitle: "",
		lang.ResourceRoleDeleteDesc:  "",

		lang.ResourceGroupListTitle: "",
		lang.ResourceGroupListDesc:  "",

		lang.ResourceGroupCreateTitle: "",
		lang.ResourceGroupCreateDesc:  "",

		lang.ResourceGroupDetailTitle: "",
		lang.ResourceGroupDetailDesc:  "",

		lang.ResourceGroupUpdateTitle: "",
		lang.ResourceGroupUpdateDesc:  "",

		lang.ResourceGroupDeleteTitle: "",
		lang.ResourceGroupDeleteDesc:  "",

		lang.ResourceDeviceListTitle: "",
		lang.ResourceDeviceListDesc:  "",

		lang.ResourceDeviceCreateTitle: "",
		lang.ResourceDeviceCreateDesc:  "",

		lang.ResourceDeviceDetailTitle: "",
		lang.ResourceDeviceDetailDesc:  "",

		lang.ResourceDeviceUpdateTitle: "",
		lang.ResourceDeviceUpdateDesc:  "",

		lang.ResourceDeviceDeleteTitle: "",
		lang.ResourceDeviceDeleteDesc:  "",

		lang.ResourceMeasureListTitle: "",
		lang.ResourceMeasureListDesc:  "",

		lang.ResourceMeasureCreateTitle: "",
		lang.ResourceMeasureCreateDesc:  "",

		lang.ResourceMeasureDetailTitle: "",
		lang.ResourceMeasureDetailDesc:  "",

		lang.ResourceMeasureUpdateTitle: "",
		lang.ResourceMeasureUpdateDesc:  "",

		lang.ResourceMeasureDeleteTitle: "",
		lang.ResourceMeasureDeleteDesc:  "",

		lang.ResourceEquipmentListTitle: "",
		lang.ResourceEquipmentListDesc:  "",

		lang.ResourceEquipmentCreateTitle: "",
		lang.ResourceEquipmentCreateDesc:  "",

		lang.ResourceEquipmentDetailTitle: "",
		lang.ResourceEquipmentDetailDesc:  "",

		lang.ResourceEquipmentUpdateTitle: "",
		lang.ResourceEquipmentUpdateDesc:  "",

		lang.ResourceEquipmentDeleteTitle: "",
		lang.ResourceEquipmentDeleteDesc:  "",

		lang.ResourceDeviceStatusTitle:     "",
		lang.ResourceDeviceDataTitle:       "",
		lang.ResourceDeviceCtrlTitle:       "",
		lang.ResourceDeviceCHValueTitle:    "",
		lang.ResourceDeviceStatisticsTitle: "",

		lang.ResourceDeviceStatusDesc:     "",
		lang.ResourceDeviceDataDesc:       "",
		lang.ResourceDeviceCtrlDesc:       "",
		lang.ResourceDeviceCHValueDesc:    "",
		lang.ResourceDeviceStatisticsDesc: "",

		lang.ResourceEquipmentStatusTitle:     "",
		lang.ResourceEquipmentDataTitle:       "",
		lang.ResourceEquipmentCtrlTitle:       "",
		lang.ResourceEquipmentCHValueTitle:    "",
		lang.ResourceEquipmentStatisticsTitle: "",

		lang.ResourceEquipmentStatusDesc:     "",
		lang.ResourceEquipmentDataDesc:       "",
		lang.ResourceEquipmentCtrlDesc:       "",
		lang.ResourceEquipmentCHValueDesc:    "",
		lang.ResourceEquipmentStatisticsDesc: "",

		lang.ResourceStateListTitle: "",
		lang.ResourceStateListDesc:  "",

		lang.ResourceStateCreateTitle: "",
		lang.ResourceStateCreateDesc:  "",

		lang.ResourceStateDetailTitle: "",
		lang.ResourceStateDetailDesc:  "",

		lang.ResourceStateUpdateTitle: "",
		lang.ResourceStateUpdateDesc:  "",

		lang.ResourceStateDeleteTitle: "",
		lang.ResourceStateDeleteDesc:  "",

		lang.AlarmListTitle:    "",
		lang.AlarmListDesc:     "",
		lang.AlarmConfirmTitle: "",
		lang.AlarmConfirmDesc:  "",
		lang.AlarmDeleteTitle:  "",
		lang.AlarmDeleteDesc:   "",
		lang.AlarmDetailTitle:  "",
		lang.AlarmDetailDesc:   "",

		lang.CommentListTitle:   "",
		lang.CommentListDesc:    "",
		lang.CommentDetailTitle: "",
		lang.CommentDetailDesc:  "",
		lang.CommentCreateTitle: "",
		lang.CommentCreateDesc:  "",
		lang.CommentDeleteTitle: "",
		lang.CommentDeleteDesc:  "",

		lang.ResourceLogListTitle: "",
		lang.ResourceLogListDesc:  "",

		lang.ResourceLogDeleteTitle:    "",
		lang.ResourceLogDeleteDesc:     "",
		lang.ResourceLogLevelListTitle: "",
		lang.ResourceLogLevelListDesc:  "",

		lang.SysBriefTitle: "",
		lang.SysBriefDesc:  "",

		lang.DataExportTitle: "",
		lang.DataExportDesc:  "",

		lang.ConfirmAdminPassword:       "请输入[ %s ]确认重置数据库：",
		lang.FlushDBOk:                  "数据库已重置！",
		lang.DefaultUserPasswordResetOk: "默认用户密码已重置！",
		lang.LogDeletedByUser:           "管理 %s 清空日志！",
		lang.DeviceLogDeletedByUser:     "管理员 %s 清空设备 %s 日志！",

		lang.CreateOrgOk:   "创建组织 %s(%s) 成功！",
		lang.CreateOrgFail: "创建组织 %s(%s) 失败：%s",
		lang.DeleteOrgOk:   "删除组织 %s(%s) 成功！",
		lang.DeleteOrgFail: "删除组织 %s(%s) 失败：%s",

		lang.AdminCreateUserOk: "管理员 %s 创建用户 %s 成功！",
		lang.AdminUpdateUserOk: "管理员 %s 更新用户 %s 成功！",
		lang.AdminDeleteUserOk: "管理员 %s 删除用户 %s 成功！",

		lang.UserCreateDeviceOk: "管理员 %s 创建设备 %s 成功！",
		lang.UserUpdateDeviceOk: "管理员 %s 更新设备 %s 成功！",
		lang.UserDeleteDeviceOk: "管理员 %s 删除设备 %s 成功！",

		lang.UserCreateEquipmentOk: "管理员 %s 创建自定义设备 %s 成功！",
		lang.UserUpdateEquipmentOk: "管理员 %s 更新自定义设备 %s 成功！",
		lang.UserDeleteEquipmentOk: "管理员 %s 删除自定义设备 %s 成功！",

		lang.UserLoginOk:                       "用户 %s 从 %s 登录成功！",
		lang.UserLoginFailedCauseDisabled:      "用户 %s 从 %s 登录失败：已禁用！",
		lang.UserLoginFailedCausePasswordWrong: "用户 %s 从 %s 登录失败：密码错误！",
		lang.UserProfileUpdateOk:               "用户 %s 更新资料成功！",
		lang.GeTuiRegisterUserFailed:           "个推注册用户%s失败！",
		lang.GeTuiSendMessageFailed:            "无法推送警报消息：%s",
		lang.GeTuiNotInitialized:               "个推没有正确配置！",
		lang.ExportInitialized:                 "正在准备导出...",
		lang.ExportingData:                     "正在导出数据 %s => %s...",
		lang.ArrangingData:                     "正在整理数据...",
		lang.WritingData:                       "正在写入数据 %d%%...",
		lang.ExportReady:                       "导出完成！",

		lang.CVSHeaderDevice:      "设备",
		lang.CVSHeaderPoint:       "点位",
		lang.CVSHeaderVal:         "值",
		lang.CVSHeaderThreshold:   "阈值",
		lang.CVSHeaderAlarm:       "警报",
		lang.CVSHeaderCreatedAt:   "创建时间",
		lang.CVSHeaderUpdatedAt:   "更新时间",
		lang.CVSHeaderUser:        "用户",
		lang.CVSHeaderConfirmedBy: "确认人",

		lang.DeviceConnected: "设备已连接！",
	}

	errStrMap = map[lang.ErrIndex]string{
		lang.Ok:                                 "成功！",
		lang.ErrUnknown:                         "未知错误！",
		lang.ErrUnknownLang:                     "未知语言区域！",
		lang.ErrInternal:                        "系统错误: %s，文件：%s，行：%d",
		lang.ErrNetworkFail:                     "网络错误：%s",
		lang.ErrEdgeInvokeFail:                  "请求数据失败 [error: %d]，请稍后再试！",
		lang.ErrServerIsBusy:                    "服务器忙，请稍后再试！",
		lang.ErrInvalidDBConnStr:                "数据库连接参数不正确！",
		lang.ErrInvalidRequestData:              "提交的数据不正确，请检查后再试！",
		lang.ErrConfirmCodeWrong:                "确认操作输入错误！",
		lang.ErrTokenExpired:                    "请先登录！",
		lang.ErrNoPermission:                    "请求失败：没有权限！",
		lang.ErrCacheNotFound:                   "缓存中没有数据！",
		lang.ErrOrganizationNotFound:            "组织机构不存在！",
		lang.ErrOrganizationDifferent:           "不同的组织机构！",
		lang.ErrOrganizationExists:              "组织机构已存在！",
		lang.ErrFailedRemoveDefaultOrganization: "不能删除默认的组织机构！",
		lang.ErrInvalidUser:                     "当前用户不可用或者登录超时！",
		lang.ErrInvalidUserName:                 "无效的用户名！",
		lang.ErrUserDisabled:                    "用户已经被禁用！",
		lang.ErrPasswordWrong:                   "密码不正确！",
		lang.ErrInvalidResourceClassID:          "资源类型不正确！",
		lang.ErrApiResourceNotFound:             "错误的api资源！",
		lang.ErrUnknownRole:                     "用户角色不正确！",
		lang.ErrUserNotFound:                    "用户不存在！",
		lang.ErrUserExists:                      "用户已经存在！",
		lang.ErrFailedDisableDefaultUser:        "不能禁用系统默认用户！",
		lang.ErrFailedRemoveDefaultUser:         "不能删除系统默认用户！",
		lang.ErrFailedEditDefaultUser:           "不能编辑系统默认管理员！",
		lang.ErrFailedRemoveUserSelf:            "不能删除当前登录的用户账号！",
		lang.ErrFailedDisableUserSelf:           "不能禁用当前登录的用户账号！",
		lang.ErrRoleNotFound:                    "用户角色不存在！",
		lang.ErrRoleExists:                      "用户角色已存在！",
		lang.ErrPolicyNotFound:                  "策略不存在！",
		lang.ErrGroupNotFound:                   "设备分组不存在！",
		lang.ErrDeviceNotFound:                  "设备不存在！",
		lang.ErrDeviceNotExistsOrActive:         "设备不存在或者还没有加载，请稍后再试！",
		lang.ErrDeviceExists:                    "设备已经存在！",
		lang.ErrMeasureNotFound:                 "点位不存在！",
		lang.ErrEquipmentNotFound:               "自定设备不存在！",
		lang.ErrStateNotFound:                   "自定义点位不存在！",
		lang.ErrDeviceOrganizationDifferent:     "设备不属于同一组织机构！",
		lang.ErrEquipmentOrganizationDifferent:  "自定义设备不属于同一组织机构！",
		lang.ErrRecursiveDetected:               "检测到循环关系链，请联系管理员！",
		lang.ErrInvalidDeviceConnStr:            "设备连接参数不正确！",
		lang.ErrConfigNotFound:                  "没找到这个配置项！",
		lang.ErrAlarmNotFound:                   "没有找到这个警报数据！",
		lang.ErrCommentNotFound:                 "没有找到这个备注！",
		lang.ErrNoStatisticsData:                "没有任何数据！",
		lang.ErrExportNotExists:                 "导出任务不存在！",
		lang.ErrRegFirst:                        "请先注册软件！",
		lang.ErrInvalidRegCode:                  "无效的注册码！",
		lang.ErrDeviceDisconnected:              "断开连接！",
		lang.ErrNoEdgeAvailable:                 "没有可用的edge程序，请重启系统！",
	}
)
