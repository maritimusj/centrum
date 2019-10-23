package lang

import (
	"errors"
	"fmt"
	"runtime"
)

var (
	errMap    = map[string]error{}
	errStrMap = map[int]map[ErrorCode]string{}
)

type ErrorCode int

const (
	Ok ErrorCode = iota
	ErrUnknown
	ErrUnknownLang
	ErrInternal
	ErrServerIsBusy
	ErrInvalidDBConnStr
	ErrInvalidRequestData
	ErrTokenExpired
	ErrNoPermission
	ErrCacheNotFound
	ErrInvalidUser
	ErrUserDisabled
	ErrPasswordWrong
	ErrConfirmCodeWrong
	ErrInvalidResourceClassID
	ErrApiResourceNotFound
	ErrUnknownRole
	ErrOrganizationNotFound
	ErrOrganizationDifferent
	ErrOrganizationExists
	ErrFailedRemoveDefaultOrganization
	ErrUserNotFound
	ErrUserExists
	ErrFailedDisableDefaultUser
	ErrFailedRemoveDefaultUser
	ErrFailedEditDefaultUser
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
	ErrInvalidDeviceConnStr
	ErrDeviceExists
	ErrConfigNotFound

	ErrAlarmNotFound
	ErrNotStatisticsData
)

func ErrorStr(code ErrorCode, params ...interface{}) string {
	var str string
	if region, ok := errStrMap[regionIndex]; ok {
		if str, ok = region[code]; !ok {
			str = region[ErrUnknown]
		}
	} else {
		str = region[ErrUnknownLang]
	}
	return str
}

func Error(code ErrorCode, params ...interface{}) error {
	var str string
	if region, ok := errStrMap[regionIndex]; ok {
		if str, ok = region[code]; !ok {
			str = region[ErrUnknown]
		}
	} else {
		str = region[ErrUnknownLang]
	}
	errStr := fmt.Sprintf(str, params...)
	if err, ok := errMap[errStr]; ok {
		return err
	}

	err := errors.New(errStr)
	errMap[errStr] = err
	return err
}

func InternalError(err error) error {
	_, file, line, _ := runtime.Caller(1)
	return Error(ErrInternal, err, file, line)
}
