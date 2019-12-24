package logStore

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Store interface {
	Open(ctx context.Context, url string, level logrus.Level) error
	Close()

	SetUID(uid string)

	//interface for logrus hook
	Levels() []logrus.Level
	Fire(entry *logrus.Entry) error
}
