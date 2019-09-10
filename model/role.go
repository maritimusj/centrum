package model

import "github.com/maritimusj/centrum/resource"

//角色
type Role interface {
	DBEntry
	EnableEntry
	Profile

	Title() string
	SetTitle(title string) error

	SetPolicy(res resource.Resource, action resource.Action, effect resource.Effect) (Policy, error)

	//对于每个资源，都应该返回一组Policy，表示对该资源的访问权限
	GetPolicy(res resource.Resource) (map[resource.Action]Policy, error)

	IsAllowed(res resource.Resource, action resource.Action) error
}
