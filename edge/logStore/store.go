package logStore

import (
	"context"
	"github.com/sirupsen/logrus"
)

type Store interface {
	Open(ctx context.Context, url string) error
	Close()

	//interface for logrus hook
	Levels() []logrus.Level
	Fire(entry *logrus.Entry) error
}
