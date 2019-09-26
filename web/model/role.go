package model

import (
	"github.com/maritimusj/centrum/web/helper"
	"github.com/maritimusj/centrum/web/resource"
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
	SetPolicy(res Resource, action resource.Action, effect resource.Effect, recursiveMap map[Resource]struct{}) (Policy, error)

	//对于每个资源，都应该返回一组Policy，表示对该资源的访问权限
	GetPolicy(res Resource) (map[resource.Action]Policy, error)

	RemovePolicy(res Resource) error

	IsAllow(res Resource, action resource.Action) (bool, error)

	GetUserList(options ...helper.OptionFN) ([]User, int64, error)
}
