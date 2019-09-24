package logStore

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	SystemLog = "system"
)

type Data struct {
	ID      uint64 `json:"id"`
	Content []byte `json:"content"`
}

type Entry struct {
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields"`
	CreatedAt time.Time              `json:"created_at"`
}

func (entry *Entry) Marshal() ([]byte, error) {
	return json.Marshal(entry)
}

type Store interface {
	Open(ctx context.Context, option map[string]interface{}) error
	Wait()

	Get(orgID int64, src, level string, start *uint64, offset, limit uint64) (result []*Data, total uint64, err error)
	Delete(orgID int64, src string) error
	Stats(orgID int64) map[string]interface{}

	//interface for logrus hook
	Levels() []logrus.Level
	Fire(entry *logrus.Entry) error
}

