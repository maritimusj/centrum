package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/maritimusj/centrum/edge/logStore"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

type Logger struct {
	cache chan []byte
	url   string
	done  chan struct{}
	wg    sync.WaitGroup
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
		cache: make(chan []byte, 100),
	}
}

func (logger *Logger) write(url string, data []byte) error {
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	_, err = defaultHttpClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func (logger *Logger) Open(ctx context.Context, url string) error {
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
				err := logger.write(logger.url, data)
				if err != nil {
					fmt.Println(err)
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
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	}
}

func (logger *Logger) Fire(entry *logrus.Entry) error {
	data, err := json.Marshal(entry.Data)
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
