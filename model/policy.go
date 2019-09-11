package model

import "github.com/maritimusj/centrum/resource"

//策略
type Policy interface {
	DBEntry
	Profile

	Role() Role

	SetEffect(effect resource.Effect) error

	IsAllow() bool
	IsDeny() bool

	Resource() resource.Resource
	Action() resource.Action
	Effect() resource.Effect
}
