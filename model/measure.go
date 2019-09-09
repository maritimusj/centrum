package model

type MeasureKind int

const (
	AllKind MeasureKind = iota
	AI
	AO
	DI
	DO
)

//点位
type Measure interface {
	DBEntry
	EnableEntry
	Profile

	Device() Device

	Title() string
	SetTitle(title string) error

	Tag() string
	Kind() MeasureKind
}
