package model

import (
	"github.com/maritimusj/centrum/helper"
	"github.com/maritimusj/centrum/resource"
)

//物理设备，网关等
type Device interface {
	DBEntry
	EnableEntry
	OptionEntry
	LogEntry
	Profile

	resource.Resource

	Organization() (Organization, error)

	Title() string
	SetTitle(title string)

	SetGroups(groups ...interface{}) error
	Groups() ([]Group, error)

	GetMeasureList(options ...helper.OptionFN) ([]Measure, int64, error)
	CreateMeasure(title string, tag string, kind resource.MeasureKind) (Measure, error)
}
