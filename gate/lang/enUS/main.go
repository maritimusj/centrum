package enUS

import (
	"github.com/maritimusj/centrum/gate/lang"
)

func init() {
	lang.Register(lang.EnUS, strMap, errStrMap)
}

var (
	strMap = map[int]string{
		lang.DefaultGroupTitle: "Default",
		lang.DefaultGroupDesc:  "Default group",

		lang.ResourceDefault:   "default",
		lang.ResourceApi:       "api",
		lang.ResourceGroup:     "group",
		lang.ResourceDevice:    "device",
		lang.ResourceMeasure:   "measure",
		lang.ResourceEquipment: "equipment",
		lang.ResourceState:     "state",

		lang.LogTrace: "trace",
		lang.LogDebug: "debug",
		lang.LogInfo:  "info",
		lang.LogWarn:  "warning",
		lang.LogError: "error",
		lang.LogFatal: "fatal",
		lang.LogPanic: "panic",

		lang.Years:        "y",
		lang.Weeks:        "w",
		lang.Days:         "d",
		lang.Hours:        "h",
		lang.Minutes:      "m",
		lang.Seconds:      "s",
		lang.Milliseconds: "ms",

		lang.AlarmUnconfirmed: "unconfirmed",
		lang.AlarmConfirmed:   "confirmed",

		lang.RoleSystemAdminTitle:       "sysadmin",
		lang.RoleOrganizationAdminTitle: "admin",
		lang.RoleGuestTitle:             "guest",

		lang.MenuRoleGalleryTitle:        "Overview",
		lang.MenuRoleDevicesTitle:        "Devices",
		lang.MenuRoleAlertTitle:          "Alarm",
		lang.MenuRoleStatsTitle:          "Graph",
		lang.MenuRoleUsersTitle:          "Users",
		lang.MenuRoleExportTitle:         "Report",
		lang.MenuRoleSystemSettingsTitle: "Setting",
		lang.MenuRoleSysLogsTitle:        "Logs",

		lang.UserDefaultRoleDesc: "",

		lang.RoleSystemAdminDesc:       "",
		lang.RoleOrganizationAdminDesc: "",
		lang.RoleGuestDesc:             "",

		lang.MenuRoleGalleryDesc:        "Viewing real-time status of devices",
		lang.MenuRoleDevicesDesc:        "Adding, deleting & editing devices",
		lang.MenuRoleAlertDesc:          "Viewing & managing alarm information",
		lang.MenuRoleStatsDesc:          "Viewing & managing statistical data",
		lang.MenuRoleExportDesc:         "Exporting data of system devices and customized devices",
		lang.MenuRoleUsersDesc:          "Managing, adding, deleting & editing users",
		lang.MenuRoleSystemSettingsDesc: "Viewing & modifying system configuration",
		lang.MenuRoleSysLogsDesc:        "Viewing & deleting system logs",

		lang.ResourceConfigBaseDetailTitle: "",
		lang.ResourceConfigBaseDetailDesc:  "",

		lang.ResourceConfigBaseUpdateTitle: "",
		lang.ResourceConfigBaseUpdateDesc:  "",

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

		lang.ConfirmAdminPassword:       "Input [%s] to flush database:",
		lang.FlushDBOk:                  "Flush ok.",
		lang.DefaultUserPasswordResetOk: "Default user password reset ok.",
		lang.LogDeletedByUser:           "Logs deleted by user: %s.",
		lang.DeviceLogDeletedByUser:     "User %s erased device %s 's logs.",

		lang.CreateOrgOk:   "Creating organization %s(%s) ok.",
		lang.CreateOrgFail: "Creating organization %s(%s) failed：%s",
		lang.DeleteOrgOk:   "Deleting organization %s(%s) ok.",
		lang.DeleteOrgFail: "Deleting organization %s(%s) failed：%s",

		lang.AdminCreateUserOk: "%s creating user %s ok.",
		lang.AdminUpdateUserOk: "%s updating user %s ok.",
		lang.AdminDeleteUserOk: "%s deleting user %s ok.",

		lang.UserCreateDeviceOk: "%s adding device %s ok.",
		lang.UserUpdateDeviceOk: "%s updating device %s ok.",
		lang.UserDeleteDeviceOk: "%s deleting device %s ok.",

		lang.UserCreateEquipmentOk: "%s creating customized device %s ok.",
		lang.UserUpdateEquipmentOk: "%s updating customized device %s ok.",
		lang.UserDeleteEquipmentOk: "%s deleting customized device %s ok.",

		lang.UserLoginOk:                       "%s logging in ok from %s.",
		lang.UserLoginFailedCauseDisabled:      "%s logging in failed from %s: disabled.",
		lang.UserLoginFailedCausePasswordWrong: "%s logging in failed from %s: invalid password.",
		lang.UserProfileUpdateOk:               "%s profile updating ok.",
		lang.ExportInitialized:                 "Initializing...",
		lang.ExportingData:                     "Exporting %s => %s...",
		lang.ArrangingData:                     "Arrange data...",
		lang.WritingData:                       "Writing data %d%%...",
		lang.ExportReady:                       "Export ok.",
	}

	errStrMap = map[lang.ErrorCode]string{
		lang.Ok:                                 "Success",
		lang.ErrUnknown:                         "Unknown error.",
		lang.ErrUnknownLang:                     "Unknown language region.",
		lang.ErrInternal:                        "Internal error: %s, file: %s, line: %d.",
		lang.ErrNetworkFail:                     "Network error: %s.",
		lang.ErrServerIsBusy:                    "Server is busy.",
		lang.ErrInvalidDBConnStr:                "Invalid parameters.",
		lang.ErrInvalidRequestData:              "Bad request data.",
		lang.ErrConfirmCodeWrong:                "Confirm code is wrong.",
		lang.ErrTokenExpired:                    "Please login first.",
		lang.ErrNoPermission:                    "Err：access denied.",
		lang.ErrCacheNotFound:                   "Missed in cache.",
		lang.ErrOrganizationNotFound:            "Organization does not exists.",
		lang.ErrOrganizationDifferent:           "Organization unmatched.",
		lang.ErrOrganizationExists:              "Organization existed.",
		lang.ErrFailedRemoveDefaultOrganization: "Delete default organization does not allowed.",
		lang.ErrInvalidUser:                     "Invalid user.",
		lang.ErrUserDisabled:                    "User disabled.",
		lang.ErrPasswordWrong:                   "Password is wrong.",
		lang.ErrInvalidResourceClassID:          "Invalid resource class id.",
		lang.ErrApiResourceNotFound:             "Invalid api resource.",
		lang.ErrUnknownRole:                     "Unknown role.",
		lang.ErrUserNotFound:                    "User does not exists.",
		lang.ErrUserExists:                      "User existed.",
		lang.ErrFailedDisableDefaultUser:        "Failed to disable default user.",
		lang.ErrFailedRemoveDefaultUser:         "Failed to remove default user.",
		lang.ErrFailedEditDefaultUser:           "Failed to edit default user.",
		lang.ErrFailedRemoveUserSelf:            "Failed to remove current user.",
		lang.ErrFailedDisableUserSelf:           "Failed to disable current user.",
		lang.ErrRoleNotFound:                    "Role does not exists.",
		lang.ErrRoleExists:                      "Role existed.",
		lang.ErrPolicyNotFound:                  "Policy does not exists.",
		lang.ErrGroupNotFound:                   "Group does not exists.",
		lang.ErrDeviceNotFound:                  "Device does not exists.",
		lang.ErrDeviceExists:                    "Device existed.",
		lang.ErrMeasureNotFound:                 "Tag does not exist.",
		lang.ErrEquipmentNotFound:               "Customized device does not exists.",
		lang.ErrStateNotFound:                   "Customized state does not exists.",
		lang.ErrDeviceOrganizationDifferent:     "Err: device and organization unmatched.",
		lang.ErrEquipmentOrganizationDifferent:  "Err: equipment and organization unmatched.",
		lang.ErrRecursiveDetected:               "Err: recursive detected.",
		lang.ErrInvalidDeviceConnStr:            "Err: invalid device connStr.",
		lang.ErrConfigNotFound:                  "Err: config entry does not exists.",
		lang.ErrAlarmNotFound:                   "Alarm does not exists.",
		lang.ErrCommentNotFound:                 "Comment does not exists.",
		lang.ErrNoStatisticsData:                "No statistics data.",
		lang.ErrExportNotExists:                 "Export task does not exists.",
		lang.ErrRegFirst:                        "Please register first.",
		lang.ErrInvalidRegCode:                  "Invalid registry code.",
		lang.ErrDeviceDisconnected:              "Device disconnected.",
	}
)
