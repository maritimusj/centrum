package resource

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

type Class int

const (
	Default Class = iota
	Api
	Group
	Device
	Equipment
	Measure
	State
)

type Resource interface {
	ResourceUID() string
	ResourceClass() Class
	ResourceTitle() string
	ResourceDesc() string
}
