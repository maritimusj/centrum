package model

import "time"

type Config interface {
	DBEntry
	OptionEntry

	Name() string
	UpdateAt() time.Time
}
