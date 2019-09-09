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

	ErrInvalidUser

	ErrorUnknownRole

	ErrInvalidResourceGroup
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

		ErrInvalidUser: "当前用户不可用或者登录超时！",

		ErrorUnknownRole:        "用户角色不正确！",
		ErrInvalidResourceGroup: "无效的分组！",
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
