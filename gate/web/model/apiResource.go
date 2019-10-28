package model

type ApiResource interface {
	Resource

	GetID() int64
	Title() string
	Desc() string
}
