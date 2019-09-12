package model

import "github.com/maritimusj/centrum/resource"

type Map map[string]interface{}

//用户
type User interface {
	DBEntry
	EnableEntry
	Profile

	Name() string
	Title() string
	Mobile() string
	Email() string

	ResetPassword(password string)
	CheckPassword(password string) bool

	Update(profile Map)

	SetRoles(roles ...interface{}) error
	GetRoles() ([]Role, error)

	IsAllow(resource resource.Resource, action resource.Action) (bool, error)
}
