package model

import (
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
)

//策略
type Policy interface {
	DBEntry
	Profile

	Role() Role

	SetEffect(effect resource2.Effect)

	Resource() Resource
	Action() resource2.Action
	Effect() resource2.Effect
}
