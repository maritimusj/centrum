package model

//动作
type Action int

const (
	View Action = iota
	Ctrl
)

//结果
type Effect int

const (
	Allow Effect = iota
	Deny
)

//策略
type Policy interface {
	DBEntry
	EnableEntry
	Profile

	Role() Role

	Resource() Resource
	Action() Action
	Effect() Effect
}
