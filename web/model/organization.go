package model

type Organization interface {
	DBEntry
	EnableEntry
	OptionEntry
	Profile

	Name() string
	Title() string
	SetTitle(string)
}
