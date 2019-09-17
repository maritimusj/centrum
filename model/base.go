package model

import (
	"github.com/sirupsen/logrus"
	"time"
)

type DBEntry interface {
	GetID() int64
	CreatedAt() time.Time

	Save() error
	Destroy() error
}

type EnableEntry interface {
	Enable()
	Disable()
	IsEnabled() bool
}

type Profile interface {
	Simple() Map
	Brief() Map
	Detail() Map
}

type LogEntry interface {
	LogUID() string
	Logger() *logrus.Entry
}
