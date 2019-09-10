package model

import "github.com/maritimusj/centrum/resource"

//虚拟设备点位
type State interface {
	DBEntry
	EnableEntry
	Profile

	resource.Resource

	Measure() Measure
	Equipment() Equipment

	SetMeasure(measure interface{}) error

	Title() string
	SetTitle(string) error

	Desc() string
	SetDesc(string) error

	Script() string
	SetScript(string) error
}
