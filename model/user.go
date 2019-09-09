package model

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
	CheckPassword(password string) error

	Update(profile Map) error

	SetRoles(roles ...interface{}) error
	GetRoles() ([]Role, error)

	IsAllowed(request Request) error
}
