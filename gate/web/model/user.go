package model

import (
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
)

type Map map[string]interface{}

//用户
type User interface {
	DBEntry
	EnableEntry
	LogEntry
	Profile

	OrganizationID() int64
	Organization() (Organization, error)

	Name() string
	Title() string
	Mobile() string
	Email() string

	ResetPassword(password string)
	CheckPassword(password string) bool

	Update(profile Map)

	SetRoles(roles ...interface{}) error
	GetRoles() ([]Role, error)
	Is(role interface{}) (bool, error)

	SetAllow(res Resource, actions ...resource2.Action) error
	SetDeny(res Resource, actions ...resource2.Action) error

	RemovePolicies(res Resource) error

	IsAllow(res Resource, action resource2.Action) (bool, error)
}
