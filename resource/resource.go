package resource

import (
	"reflect"
	"strings"
)

//动作
type Action int8

const (
	View Action = iota
	Ctrl
)

const (
	Invoke = Ctrl
)

//结果
type Effect int8

const (
	Deny Effect = iota
	Allow
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

type MeasureKind int

const (
	AllKind MeasureKind = iota
	AI
	AO
	DI
	DO
)

var (
	classTitle = map[string]Class{
		"api":       Api,
		"group":     Group,
		"device":    Device,
		"measure":   Measure,
		"equipment": Equipment,
		"state":     State,
	}
)

func ParseClass(class string) Class {
	if v, ok := classTitle[strings.ToLower(class)]; ok {
		return v
	}
	return Default
}

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
