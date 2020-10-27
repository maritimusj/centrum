package zhCN

import (
	"github.com/maritimusj/centrum/gate/lang"
)

func init() {
	lang.Register(lang.ZhTW, strMap, errStrMap)
}

var (
	strMap = map[lang.StrIndex]string{
		lang.DatetimeFormatterStr: "2006-01-02 15:04:05",
		lang.DefaultGroupTitle:    "默認分組",
		lang.DefaultGroupDesc:     "系統默認分組",

		lang.ResourceDefault:   "默認分組",
		lang.ResourceApi:       "後臺權限",
		lang.ResourceGroup:     "設備分組",
		lang.ResourceDevice:    "設備",
		lang.ResourceMeasure:   "點位",
		lang.ResourceEquipment: "自定義設備",
		lang.ResourceState:     "自定義點位",

		lang.LogTrace:   "跟蹤",
		lang.LogDebug:   "調試",
		lang.LogInfo:    "信息",
		lang.LogWarning: "警告",
		lang.LogError:   "錯誤",
		lang.LogFatal:   "嚴重",
		lang.LogPanic:   "異常",

		lang.Years:        "年",
		lang.Weeks:        "星期",
		lang.Days:         "天",
		lang.Hours:        "小時",
		lang.Minutes:      "分鐘",
		lang.Seconds:      "秒",
		lang.Milliseconds: "毫秒",

		lang.AlarmUnconfirmed: "未確認",
		lang.AlarmConfirmed:   "已確認",

		lang.RoleSystemAdminTitle:       "系統管理員",
		lang.RoleOrganizationAdminTitle: "管理員",
		lang.RoleGuestTitle:             "普通用戶",

		lang.AlarmNotifyTitle:   "警報通知",
		lang.AlarmNotifyContent: "發生新的警報，請及時處理！",

		lang.MenuRoleGalleryTitle:        "設備總覽",
		lang.MenuRoleDevicesTitle:        "設備管理",
		lang.MenuRoleAlertTitle:          "報警查詢",
		lang.MenuRoleStatsTitle:          "趨勢圖",
		lang.MenuRoleUsersTitle:          "權限設定",
		lang.MenuRoleExportTitle:         "導出報表",
		lang.MenuRoleSystemSettingsTitle: "系統設定",
		lang.MenuRoleSysLogsTitle:        "系統日誌",
		lang.MenuWebViewTitle:            "雲視圖",
		lang.MenuStreamViewTitle:         "雲監控",
		lang.MenuWebViewDesc:             "查看、設定加載特定的網頁",
		lang.MenuStreamViewDesc:          "查看、編輯遠程監控視頻",

		lang.UserDefaultRoleDesc: "",

		lang.RoleSystemAdminDesc:       "",
		lang.RoleOrganizationAdminDesc: "",
		lang.RoleGuestDesc:             "",

		lang.MenuRoleGalleryDesc:        "查看設備實時狀態",
		lang.MenuRoleDevicesDesc:        "可以對設備進行管理，添加刪除和編輯設備",
		lang.MenuRoleAlertDesc:          "可以查看管理警報信息",
		lang.MenuRoleStatsDesc:          "可以查看，管理統計數據",
		lang.MenuRoleExportDesc:         "可以導出系統設備和自定義設備的數據報表",
		lang.MenuRoleUsersDesc:          "可以管理用戶，添加刪除和編輯用戶",
		lang.MenuRoleSystemSettingsDesc: "可以查看，修改系統設置",
		lang.MenuRoleSysLogsDesc:        "可以查看，刪除系統日誌",

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

		lang.ConfirmAdminPassword:       "請輸入[ %s ]確認重置數據庫：",
		lang.FlushDBOk:                  "數據庫已重置！",
		lang.DefaultUserPasswordResetOk: "默認用戶密碼已重置！",
		lang.LogDeletedByUser:           "管理 %s 清空日誌！",
		lang.DeviceLogDeletedByUser:     "管理員 %s 清空設備 %s 日誌！",

		lang.CreateOrgOk:   "創建組織 %s(%s) 成功！",
		lang.CreateOrgFail: "創建組織 %s(%s) 失敗：%s",
		lang.DeleteOrgOk:   "刪除組織 %s(%s) 成功！",
		lang.DeleteOrgFail: "刪除組織 %s(%s) 失敗：%s",

		lang.AdminCreateUserOk: "管理員 %s 創建用戶 %s 成功！",
		lang.AdminUpdateUserOk: "管理員 %s 更新用戶 %s 成功！",
		lang.AdminDeleteUserOk: "管理員 %s 刪除用戶 %s 成功！",

		lang.UserCreateDeviceOk: "管理員 %s 創建設備 %s 成功！",
		lang.UserUpdateDeviceOk: "管理員 %s 更新設備 %s 成功！",
		lang.UserDeleteDeviceOk: "管理員 %s 刪除設備 %s 成功！",

		lang.UserCreateEquipmentOk: "管理員 %s 創建自定義設備 %s 成功！",
		lang.UserUpdateEquipmentOk: "管理員 %s 更新自定義設備 %s 成功！",
		lang.UserDeleteEquipmentOk: "管理員 %s 刪除自定義設備 %s 成功！",

		lang.UserLoginOk:                       "用戶 %s 從 %s 登錄成功！",
		lang.UserLoginFailedCauseDisabled:      "用戶 %s 從 %s 登錄失敗：已禁用！",
		lang.UserLoginFailedCausePasswordWrong: "用戶 %s 從 %s 登錄失敗：密碼錯誤！",
		lang.UserProfileUpdateOk:               "用戶 %s 更新資料成功！",

		lang.ExportInitialized: "正在準備導出...",
		lang.ExportingData:     "正在導出數據 %s => %s...",
		lang.ArrangingData:     "正在整理數據...",
		lang.WritingData:       "正在寫入數據 %d%%...",
		lang.ExportReady:       "導出完成！",

		lang.CVSHeaderDevice:      "設備",
		lang.CVSHeaderPoint:       "點位",
		lang.CVSHeaderVal:         "值",
		lang.CVSHeaderThreshold:   "閾值",
		lang.CVSHeaderAlarm:       "警報",
		lang.CVSHeaderCreatedAt:   "創建時間",
		lang.CVSHeaderUpdatedAt:   "更新時間",
		lang.CVSHeaderUser:        "用戶",
		lang.CVSHeaderConfirmedBy: "確認人",

		lang.DeviceConnected: "設備已連接！",
	}

	errStrMap = map[lang.ErrIndex]string{
		lang.Ok:                                 "成功！",
		lang.ErrUnknown:                         "未知錯誤！",
		lang.ErrUnknownLang:                     "未知語言區域！",
		lang.ErrInternal:                        "系統錯誤: %s，文件：%s，行：%d",
		lang.ErrNetworkFail:                     "網絡錯誤：%s",
		lang.ErrEdgeInvokeFail:                  "請求數據失敗 [error: %d]，請稍後再試！",
		lang.ErrServerIsBusy:                    "服務器忙，請稍後再試！",
		lang.ErrInvalidDBConnStr:                "數據庫連接參數不正確！",
		lang.ErrInvalidRequestData:              "提交的數據不正確，請檢查後再試！",
		lang.ErrConfirmCodeWrong:                "確認操作輸入錯誤！",
		lang.ErrTokenExpired:                    "請先登錄！",
		lang.ErrNoPermission:                    "請求失敗：沒有權限！",
		lang.ErrCacheNotFound:                   "緩存中沒有數據！",
		lang.ErrOrganizationNotFound:            "組織機構不存在！",
		lang.ErrOrganizationDifferent:           "不同的組織機構！",
		lang.ErrOrganizationExists:              "組織機構已存在！",
		lang.ErrFailedRemoveDefaultOrganization: "不能刪除默認的組織機構！",
		lang.ErrInvalidUser:                     "當前用戶不可用或者登錄超時！",
		lang.ErrInvalidUserName:                 "無效的用戶名！",
		lang.ErrUserDisabled:                    "用戶已經被禁用！",
		lang.ErrPasswordWrong:                   "密碼不正確！",
		lang.ErrInvalidResourceClassID:          "資源類型不正確！",
		lang.ErrApiResourceNotFound:             "錯誤的api資源！",
		lang.ErrUnknownRole:                     "用戶角色不正確！",
		lang.ErrUserNotFound:                    "用戶不存在！",
		lang.ErrUserExists:                      "用戶已經存在！",
		lang.ErrFailedDisableDefaultUser:        "不能禁用系統默認用戶！",
		lang.ErrFailedRemoveDefaultUser:         "不能刪除系統默認用戶！",
		lang.ErrFailedEditDefaultUser:           "不能編輯系統默認管理員！",
		lang.ErrFailedRemoveUserSelf:            "不能刪除當前登錄的用戶賬號！",
		lang.ErrFailedDisableUserSelf:           "不能禁用當前登錄的用戶賬號！",
		lang.ErrRoleNotFound:                    "用戶角色不存在！",
		lang.ErrRoleExists:                      "用戶角色已存在！",
		lang.ErrPolicyNotFound:                  "策略不存在！",
		lang.ErrGroupNotFound:                   "設備分組不存在！",
		lang.ErrDeviceNotFound:                  "設備不存在！",
		lang.ErrDeviceNotExistsOrActive:         "設備不存在或者還沒有加載，請稍後再試！",
		lang.ErrDeviceExists:                    "設備已經存在！",
		lang.ErrMeasureNotFound:                 "點位不存在！",
		lang.ErrEquipmentNotFound:               "自定設備不存在！",
		lang.ErrStateNotFound:                   "自定義點位不存在！",
		lang.ErrDeviceOrganizationDifferent:     "設備不屬於同壹組織機構！",
		lang.ErrEquipmentOrganizationDifferent:  "自定義設備不屬於同壹組織機構！",
		lang.ErrRecursiveDetected:               "檢測到循環關系鏈，請聯系管理員！",
		lang.ErrInvalidDeviceConnStr:            "設備連接參數不正確！",
		lang.ErrConfigNotFound:                  "沒找到這個配置項！",
		lang.ErrAlarmNotFound:                   "沒有找到這個警報數據！",
		lang.ErrCommentNotFound:                 "沒有找到這個備註！",
		lang.ErrNoStatisticsData:                "沒有任何數據！",
		lang.ErrExportNotExists:                 "導出任務不存在！",
		lang.ErrRegFirst:                        "請先註冊軟件！",
		lang.ErrInvalidRegCode:                  "無效的註冊碼！",
		lang.ErrDeviceDisconnected:              "斷開連接！",
		lang.ErrNoEdgeAvailable:                 "沒有可用的edge程序，請重啟系統！",
		lang.ErrGeTuiRegisterUserFailed:         "個推註冊用戶%s失敗！",
		lang.ErrGeTuiSendMessageFailed:          "無法推送警報消息：%s",
		lang.ErrGeTuiNotInitialized:             "個推沒有正確配置！",
	}
)
