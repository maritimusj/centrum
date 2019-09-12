package resource

import "reflect"

//动作
type Action int8

const (
	Invoke Action = 0
)
const (
	View Action = iota
	Ctrl
)

//结果
type Effect int8

const (
	Allow Effect = iota
	Deny
)

type Class int8

const (
	Default Class = iota
	Api
	Group
	Device
	Measure
	Equipment
	State
)

func IsValidClass(class interface{}) bool {
	v := reflect.ValueOf(class)
	if v.IsValid() {
		switch v.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val := Class(v.Int())
			return val > Default && val <= State
		}
	}
	return false
}

type Resource interface {
	ResourceClass() Class
	ResourceID() int64
	ResourceTitle() string
	ResourceDesc() string
}
