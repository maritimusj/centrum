package model

import "github.com/maritimusj/centrum/resource"

//请求
type Request interface {
	Resource() resource.Resource
	Action() resource.Action
}
