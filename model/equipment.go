package model

import (
	"github.com/maritimusj/centrum/helper"
	"github.com/maritimusj/centrum/resource"
)

//虚拟设备
type Equipment interface {
	DBEntry
	EnableEntry
	LogEntry
	Profile

	resource.Resource

	Title() string
	SetTitle(title string)

	Desc() string
	SetDesc(string)

	SetGroups(groups ...interface{}) error
	Groups() ([]Group, error)

	GetStateList(options ...helper.OptionFN) ([]State, int64, error)
	CreateState(title, desc string, measure interface{}, script string) (State, error)
}
