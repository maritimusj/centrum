package model

import "time"

type DBEntry interface {
	GetID() int64
	CreatedAt() time.Time

	Save() error
	Destroy() error
}

type EnableEntry interface {
	Enable() error
	Disable() error
	IsEnabled() bool
}

type Profile interface {
	Simple() Map
	Brief() Map
	Detail() Map
}
