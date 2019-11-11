package model

//虚拟设备点位
type State interface {
	DBEntry
	EnableEntry
	OptionEntry
	Resource
	Profile

	Measure() Measure
	Equipment() Equipment

	SetMeasure(measure interface{})

	Title() string
	SetTitle(string)

	Desc() string
	SetDesc(string)

	IsAlarmEnabled() bool
	EnableAlarm()
	DisableAlarm()

	AlarmDeadBand() float32
	SetAlarmDeadBand(v float32)

	AlarmDelaySecond() int
	SetAlarmDelay(seconds int)

	GetAlarmEntries() map[string]float32
	GetAlarmEntry(name string) (float32, bool)
	SetAlarmEntry(name string, value float32)
	EnableAlarmEntry(name string)
	DisableAlarmEntry(name string)
}
