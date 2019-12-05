package logStore

import (
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
	return json.Marshal(map[string]interface{}{
		"level":      entry.Level,
		"message":    entry.Message,
		"fields":     entry.Fields,
		"created_at": entry.CreatedAt.Format("2006-01-02 15:04:05"),
	})
}

type Store interface {
	Open(option map[string]interface{}) error
	Close()

	GetList(orgID int64, src, level string, start *uint64, offset, limit uint64) (result []*Data, total uint64, err error)
	Delete(orgID int64, src string) error
	Stats(orgID int64) map[string]interface{}

	//interface for logrus hook
	Levels() []logrus.Level
	Fire(entry *logrus.Entry) error
}
