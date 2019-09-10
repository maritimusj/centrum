package model

import "github.com/maritimusj/centrum/resource"

//虚拟设备
type Equipment interface {
	DBEntry
	EnableEntry
	Profile

	resource.Resource

	Title() string
	SetTitle(title string) error

	Desc() string
	SetDesc(string) error

	SetGroups(groups ...Group) error
	Groups() ([]Group, error)

	CreateState(title string, measure interface{}, script string) (State, error)
}
