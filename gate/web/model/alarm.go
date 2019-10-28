package model

import "time"

type Alarm interface {
	DBEntry
	OptionEntry
	Profile

	DeviceID() int64
	MeasureID() int64

	Device() (Device, error)
	Measure() (Measure, error)

	Status() (int, string)
	Confirm(map[string]interface{}) error

	UpdatedAt() time.Time
	Updated()
}
