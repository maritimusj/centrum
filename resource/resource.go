package resource

//动作
type Action int

const (
	Invoke Action = 0
)
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
	Measure
	Equipment
	State
)

type Resource interface {
	ResourceUID() string
	ResourceClass() Class
	ResourceTitle() string
	ResourceDesc() string
}
