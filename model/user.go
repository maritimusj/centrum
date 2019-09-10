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

	ResetPassword(password string) error
	CheckPassword(password string) bool

	Update(profile Map) error

	SetRoles(roles ...interface{}) error
	GetRoles() ([]Role, error)

	IsAllowed(resource resource.Resource, action resource.Action) error
}
