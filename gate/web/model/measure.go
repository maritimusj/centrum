package model

import (
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
)

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
	Kind() resource2.MeasureKind
}
