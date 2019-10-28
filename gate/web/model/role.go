package model

import (
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
)

//角色
type Role interface {
	DBEntry
	EnableEntry
	Profile

	OrganizationID() int64

	Name() string

	Title() string
	SetTitle(title string)

	Desc() string
	SetDesc(desc string)

	//设置指定资源的对于指定action的effect，传入recursiveMap则对子资源进行递归设置，recursiveMap为nil则只设置当前资源
	SetPolicy(res Resource, action resource2.Action, effect resource2.Effect, recursiveMap map[Resource]struct{}) (Policy, error)

	//对于每个资源，都应该返回一组Policy，表示对该资源的访问权限
	GetPolicy(res Resource) (map[resource2.Action]Policy, error)

	RemovePolicy(res Resource) error

	IsAllow(res Resource, action resource2.Action) (bool, error)

	GetUserList(options ...helper2.OptionFN) ([]User, int64, error)
}
