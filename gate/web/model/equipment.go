package model

import (
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
)

//虚拟设备
type Equipment interface {
	DBEntry
	EnableEntry
	LogEntry
	Profile

	Resource

	Organization() (Organization, error)

	Title() string
	SetTitle(title string)

	Desc() string
	SetDesc(string)

	SetGroups(groups ...interface{}) error
	Groups() ([]Group, error)

	GetStateList(options ...helper2.OptionFN) ([]State, int64, error)
	CreateState(title, desc string, measure interface{}, script string) (State, error)
}
