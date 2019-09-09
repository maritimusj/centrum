package model

//角色
type Role interface {
	DBEntry
	EnableEntry
	Profile

	Title() string
	SetTitle(title string) error

	SetPolicy(resource interface{}, action Action, effect Effect) (Policy, error)
	//对于每个资源，都应该返回一组Policy，表示对该资源的访问权限
	GetPolicy(resource interface{}) (map[Action]Policy, error)

	IsAllowed(request Request) error
}
