package model

import "github.com/maritimusj/centrum/resource"

//点位
type Measure interface {
	DBEntry
	EnableEntry
	Profile

	Resource

	Device() Device

	Title() string
	SetTitle(title string)

	Tag() string
	Kind() resource.MeasureKind
}
