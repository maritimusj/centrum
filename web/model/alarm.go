package model

import "time"

type Alarm interface {
	DBEntry
	OptionEntry

	DeviceID() int64
	MeasureID() int64

	Device() (Device, error)
	Measure() (Measure, error)

	UpdatedAt() time.Time
	Updated()
}
