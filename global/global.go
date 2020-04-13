package global

import (
	"sync"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var (
	Params   = New()
	Stats    = New()
	Messages = &messages{
		msgBoxes: map[string]*msgBox{},
	}
)

type msgBox struct {
	lastUse time.Time
	list    []*msg
}

type msg struct {
	CreatedAt time.Time
	Data      interface{}
}

type messages struct {
	sync.Mutex
	msgBoxes map[string]*msgBox
}

func (m *messages) Add(data interface{}, fn func(string) bool) {
	m.Lock()
	defer m.Unlock()

	for uid, box := range m.msgBoxes {
		//删除长时间不用的信箱
		if time.Now().Sub(box.lastUse) > 30*time.Minute {
			delete(m.msgBoxes, uid)
			continue
		}

		if fn != nil && !fn(uid) {
			continue
		}

		//删除超过3分钟的消息
		for i := 0; i < len(box.list); {
			if time.Now().Sub(box.list[i].CreatedAt) > 3*time.Minute {
				box.list = append(box.list[:i], box.list[i+1:]...)
			} else {
				i++
			}
		}

		if len(box.list) > 6 {
			box.list = box.list[len(box.list)-6:]
		}

		box.list = append(box.list, &msg{
			CreatedAt: time.Now(),
			Data:      data,
		})

		m.msgBoxes[uid] = box
	}

}

func (m *messages) Create(uid string) {
	m.Lock()
	defer m.Unlock()

	m.msgBoxes[uid] = &msgBox{
		lastUse: time.Now(),
	}
}

func (m *messages) Close(uid string) {
	m.Lock()
	defer m.Unlock()

	delete(m.msgBoxes, uid)
}

func (m *messages) GetAll(uid string) []*msg {
	m.Lock()
	defer m.Unlock()

	box, ok := m.msgBoxes[uid]
	if ok {
		list := box.list
		box.list = []*msg{}
		box.lastUse = time.Now()
		return list
	}

	return []*msg{}
}

type stats struct {
	data []byte
	sync.RWMutex
}

func New() *stats {
	return &stats{}
}

func (stats *stats) Set(path string, val interface{}) error {
	stats.Lock()
	defer stats.Unlock()

	data, err := sjson.SetBytes(stats.data, path, val)
	if err != nil {
		return err
	}
	stats.data = data
	return nil
}

func (stats *stats) Exists(path string) bool {
	stats.RLock()
	defer stats.RUnlock()

	return gjson.GetBytes(stats.data, path).Exists()
}

func (stats *stats) Get(path string) (interface{}, bool) {
	stats.RLock()
	defer stats.RUnlock()

	v := gjson.GetBytes(stats.data, path)
	return v.Value(), v.Exists()
}

func (stats *stats) MustGet(path string) interface{} {
	stats.RLock()
	defer stats.RUnlock()

	v := gjson.GetBytes(stats.data, path)
	if !v.Exists() {
		panic("stats not exists")
	}

	return v.Value()
}

func (stats *stats) Remove(path string) {
	stats.Lock()
	defer stats.Unlock()

	v, err := sjson.DeleteBytes(stats.data, path)
	if err != nil {
		return
	}
	stats.data = v
}

func (stats *stats) Reset() {
	stats.data = []byte{}
}
