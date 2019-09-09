package model

//物理设备，网关等
type Device interface {
	DBEntry
	EnableEntry
	Profile

	Title() string
	SetTitle(title string) error

	Option() Map
	SetOption(option Map) error

	SetGroups(groups ...Group) error
	Groups() ([]Group, error)

	CreateMeasure(title string, tag string, kind MeasureKind) (Measure, error)
}
