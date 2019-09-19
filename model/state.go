package model

//虚拟设备点位
type State interface {
	DBEntry
	EnableEntry
	Profile

	Resource

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
