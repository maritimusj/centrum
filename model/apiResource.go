package model

import (
	"github.com/maritimusj/centrum/resource"
)

type ApiResource interface {
	resource.Resource

	GetID() int64
	Title() string
	Desc() string
}
