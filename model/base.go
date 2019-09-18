package model

import (
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"time"
)

type DBEntry interface {
	GetID() int64
	CreatedAt() time.Time

	Save() error
	Destroy() error
}

type OptionEntry interface {
	GetOption(path string) gjson.Result
	SetOption(path string, value interface{}) error
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
