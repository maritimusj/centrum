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

	SetMeasure(measure interface{})

	Title() string
	SetTitle(string)

	Desc() string
	SetDesc(string)

	Script() string
	SetScript(string)
}
