package model

import (
	"github.com/maritimusj/centrum/resource"
	"github.com/tidwall/gjson"
)

//物理设备，网关等
type Device interface {
	DBEntry
	EnableEntry
	Profile

	resource.Resource

	Title() string
	SetTitle(title string) error

	GetOption(path string) gjson.Result
	SetOption(path string, value interface{}) error

	SetGroups(groups ...interface{}) error
	Groups() ([]Group, error)

	GetMeasureList(keyword string, kind MeasureKind, page, pageSize int64) ([]Measure, int64, error)
	CreateMeasure(title string, tag string, kind MeasureKind) (Measure, error)
}
