package model

import "github.com/maritimusj/centrum/web/resource"

//点位
type Measure interface {
	DBEntry
	EnableEntry
	Profile

	Resource

	Device() Device

	Title() string
	SetTitle(title string)

	TagName() string
	Kind() resource.MeasureKind
}
