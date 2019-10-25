package lang

import (
	"errors"
	"fmt"
	"github.com/maritimusj/centrum/synchronized"
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

	ErrDeviceNotExists
	ErrDeviceNotConnected
)

func ErrorStr(code ErrorCode, params ...interface{}) string {
	str := <-synchronized.Do("error.str", func() interface{} {
		var str string
		if region, ok := errStrMap[regionIndex]; ok {
			if str, ok = region[code]; !ok {
				str = region[ErrUnknown]
			}
		} else {
			str = region[ErrUnknownLang]
		}
		return str
	})
	return str.(string)
}

func Error(code ErrorCode, params ...interface{}) error {
	err := <-synchronized.Do("error.str", func() interface{} {
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
	})
	return err.(error)
}

func InternalError(err error) error {
	_, file, line, _ := runtime.Caller(1)
	return Error(ErrInternal, err, file, line)
}
