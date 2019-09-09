package model

//虚拟设备点位
type State interface {
	DBEntry
	EnableEntry
	Profile

	Measure() Measure
	Equipment() Equipment

	SetMeasure(measure interface{}) error

	Title() string
	SetTitle(string) error

	Script() string
	SetScript(string) error
}
