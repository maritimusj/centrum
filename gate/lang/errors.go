package lang

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/maritimusj/centrum/synchronized"
)

var (
	errMap    = map[string]error{}
	errStrMap = map[int]map[ErrIndex]string{}
)

type ErrIndex int

func (index ErrIndex) Str(params ...interface{}) string {
	return ErrorStr(index, params...)
}

func (index ErrIndex) Error(params ...interface{}) error {
	return Error(index, params...)
}

const (
	Ok ErrIndex = iota
	ErrUnknown
	ErrUnknownLang
	ErrInternal
	ErrNetworkFail
	ErrEdgeInvokeFail
	ErrServerIsBusy
	ErrInvalidDBConnStr
	ErrInvalidRequestData
	ErrTokenExpired
	ErrNoPermission
	ErrCacheNotFound
	ErrInvalidUser
	ErrInvalidUserName
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

	ErrCommentNotFound

	ErrNoStatisticsData
	ErrExportNotExists

	ErrRegFirst
	ErrInvalidRegCode

	ErrDeviceDisconnected
	ErrDeviceNotExistsOrActive

	ErrNoEdgeAvailable
)

func ErrorStr(index ErrIndex, params ...interface{}) string {
	str := <-synchronized.Do("error.str", func() interface{} {
		var str string
		if region, ok := errStrMap[regionIndex]; ok {
			if str, ok = region[index]; !ok {
				str = region[ErrUnknown]
			}
		} else {
			str = region[ErrUnknownLang]
		}
		return fmt.Sprintf(str, params...)
	})
	return str.(string)
}

func Error(index ErrIndex, params ...interface{}) error {
	err := <-synchronized.Do("error.str", func() interface{} {
		var str string
		if region, ok := errStrMap[regionIndex]; ok {
			if str, ok = region[index]; !ok {
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
	})

	return err.(error)
}

func InternalError(err error) error {
	_, file, line, _ := runtime.Caller(1)
	return Error(ErrInternal, err, file, line)
}
