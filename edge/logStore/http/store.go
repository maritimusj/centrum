package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/maritimusj/centrum/edge/logStore"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	uid       string
	cache     chan []byte
	done      chan struct{}
	wg        sync.WaitGroup
	logLevels []logrus.Level
}

var (
	defaultHttpClient = &http.Client{}
)

func DefaultHttpClient() *http.Client {
	return defaultHttpClient
}

func New() logStore.Store {
	return &Logger{
		done:  make(chan struct{}),
		cache: make(chan []byte, 1000),
	}
}

func (logger *Logger) SetUID(uid string) {
	logger.uid = uid
}

func (logger *Logger) write(url string, data []byte) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return
	}
	_, err = defaultHttpClient.Do(req)
	if err != nil {
		return
	}
}

func (logger *Logger) Open(ctx context.Context, url string, level logrus.Level) error {
	for _, lv := range []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	} {
		if lv <= level {
			logger.logLevels = append(logger.logLevels, lv)
		}
	}
	logger.wg.Add(1)
	go func() {
		defer func() {
			logger.wg.Done()
		}()
		for {
			select {
			case <-logger.done:
				return
			case <-ctx.Done():
				return
			case data := <-logger.cache:
				if data != nil {
					logger.write(url, data)
				}
			}
		}
	}()

	return nil
}

func (logger *Logger) Close() {
	close(logger.done)
	logger.wg.Wait()
}

func (logger *Logger) Levels() []logrus.Level {
	return logger.logLevels
}

func (logger *Logger) Fire(entry *logrus.Entry) error {
	data, err := json.Marshal(map[string]interface{}{
		"log": map[string]interface{}{
			"uid":   logger.uid,
			"level": entry.Level.String(),
			"msg":   entry.Message,
			"time":  entry.Time,
		},
	})
	if err != nil {
		return err
	}

	select {
	case <-logger.done:
		return nil
	default:
		logger.cache <- data
	}
	return nil
}
