package SysInfo

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Info struct {
	l         sync.Mutex
	lastFetch time.Time
	data      interface{}
}

func (info *Info) Have() bool {
	return info.data != nil && info.lastFetch.Second() > 0
}
func (info *Info) Expired(duration time.Duration) bool {
	return info.data == nil || time.Now().Sub(info.lastFetch) > duration
}

func (info *Info) Fetch(fn func() (interface{}, error)) {
	defer func() {
		recover()
		info.data = nil
	}()

	info.l.Lock()
	defer info.l.Unlock()
	if fn != nil {
		data, err := fn()
		if err != nil {
			log.Warn(err)
		} else {
			info.data = data
			info.lastFetch = time.Now()
		}
	}
}

func (info *Info) Data() interface{} {
	if info.data == nil {
		return map[string]interface{}{}
	}
	return info.data
}
