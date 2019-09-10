package dirty

import "sync"

type Dirty struct {
	data map[string]func() interface{}
	sync.RWMutex
}

func New() *Dirty {
	return &Dirty{
		data: make(map[string]func() interface{}),
	}
}

func (dirty *Dirty) Data(reset bool) map[string]interface{} {
	if reset {
		dirty.Lock()
		defer dirty.Unlock()
	} else {
		dirty.RLock()
		defer dirty.RUnlock()
	}

	var result = make(map[string]interface{})
	for name, getter := range dirty.data {
		if getter != nil {
			result[name] = getter()
		}
	}

	if reset {
		dirty.reset()
	}
	return result
}

func (dirty *Dirty) Set(name string, fn func() interface{}) {
	dirty.Lock()
	defer dirty.Unlock()

	dirty.data[name] = fn
}

func (dirty *Dirty) Is(name string) bool {
	dirty.RLock()
	defer dirty.RUnlock()

	if _, ok := dirty.data[name]; ok {
		return true
	}
	return false
}

func (dirty *Dirty) Any() bool {
	dirty.RLock()
	defer dirty.RUnlock()

	return len(dirty.data) > 0
}

func (dirty *Dirty) Unset(names ...string) {
	dirty.Lock()
	defer dirty.Unlock()

	for _, name := range names {
		delete(dirty.data, name)
	}
}

func (dirty *Dirty) reset() {
	dirty.data = make(map[string]func() interface{})
}

func (dirty *Dirty) Reset() {
	dirty.Lock()
	defer dirty.Unlock()
	dirty.reset()
}
