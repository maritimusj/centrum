package model

import (
	"github.com/maritimusj/centrum/gate/web/resource"
)

//策略
type Policy interface {
	DBEntry
	Profile

	Role() Role

	SetEffect(effect resource.Effect)

	Resource() Resource
	Action() resource.Action
	Effect() resource.Effect
}
