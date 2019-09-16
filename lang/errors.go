package lang

import (
	"errors"
	"fmt"
)

type ErrorCode int

const (
	Ok ErrorCode = iota
	ErrUnknown
	ErrInternal

	ErrInvalidConnStr

	ErrInvalidRequestData
	ErrTokenExpired
	ErrNoPermission

	ErrCacheNotFound

	ErrInvalidUser
	ErrUserNotFound
	ErrUserExists

	ErrFailedDisableDefaultUser
	ErrFailedRemoveDefaultUser
	ErrFailedEditDefaultUserPerm

	ErrInvalidResourceClassID

	ErrPasswordWrong

	ErrApiResourceNotFound

	ErrUnknownRole
	ErrRoleNotFound

	ErrPolicyNotFound
	ErrGroupNotFound

	ErrDeviceNotFound
	ErrMeasureNotFound
	ErrEquipmentNotFound
	ErrStateNotFound
)

var (
	errMap    = make(map[string]error)
	errStrMap = map[ErrorCode]string{
		Ok:          "Ok",
		ErrUnknown:  "未知错误！",
		ErrInternal: "系统错误: %s",

		ErrInvalidConnStr: "数据库连接参数不正确！",

		ErrInvalidRequestData: "不正确的请求数据！",
		ErrTokenExpired:       "请先登录！",
		ErrNoPermission:       "没有权限！",

		ErrCacheNotFound: "缓存中没有数据！",

		ErrInvalidUser:   "当前用户不可用或者登录超时！",
		ErrPasswordWrong: "密码不正确！",

		ErrInvalidResourceClassID: "资源类型不正确！",

		ErrApiResourceNotFound: "错误的api资源！",

		ErrUnknownRole: "用户角色不正确！",

		ErrUserNotFound:              "用户不存在！",
		ErrUserExists:                "用户已经存在！",
		ErrFailedDisableDefaultUser:  "不能禁用系统默认用户！",
		ErrFailedRemoveDefaultUser:   "不能删除系统默认用户！",
		ErrFailedEditDefaultUserPerm: "不能编辑默认用户角色！",

		ErrRoleNotFound:      "用户角色不存在！",
		ErrPolicyNotFound:    "策略不存在！",
		ErrGroupNotFound:     "设备分组不存在！",
		ErrDeviceNotFound:    "设备不存在！",
		ErrMeasureNotFound:   "点位不存在！",
		ErrEquipmentNotFound: "自定设备不存在！",
		ErrStateNotFound:     "自定义点位不存在！",
	}
)

func ErrorStr(code ErrorCode) string {
	if str, ok := errStrMap[code]; ok {
		return str
	}
	return errStrMap[ErrUnknown]
}

func InternalError(err error) error {
	return Error(ErrInternal, err)
}

func Error(code ErrorCode, params ...interface{}) error {
	str, ok := errStrMap[code]
	if !ok {
		str = errStrMap[ErrUnknown]
	}

	errStr := fmt.Sprintf(str, params...)
	if err, ok := errMap[errStr]; ok {
		return err
	}
	err := errors.New(errStr)
	errMap[errStr] = err
	return err
}
