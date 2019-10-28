package memCache

import (
	"errors"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	goCache "github.com/patrickmn/go-cache"
	"strconv"
	"time"
)

const (
	PrefixConfig      = "x."
	PrefixOrg         = "0."
	PrefixUser        = "1."
	PrefixRole        = "2."
	prefixPolicy      = "3."
	prefixGroup       = "4."
	prefixDevice      = "5."
	prefixMeasure     = "6."
	prefixEquipment   = "7."
	prefixState       = "8."
	prefixApiResource = "9."
	prefixAlarm       = "10."
)

type cache struct {
	client *goCache.Cache
}

type ID interface {
	GetID() int64
}

func New() *cache {
	return &cache{client: goCache.New(6*time.Minute, 10*time.Minute)}
}

func getKey(obj interface{}) string {
	var pref string
	switch v := obj.(type) {
	case string:
		return v
	case model2.Organization:
		pref = PrefixOrg
	case model2.User:
		pref = PrefixUser
	case model2.Role:
		pref = PrefixRole
	case model2.Policy:
		pref = prefixPolicy
	case model2.Group:
		pref = prefixGroup
	case model2.Device:
		pref = prefixDevice
	case model2.Measure:
		pref = prefixMeasure
	case model2.Equipment:
		pref = prefixEquipment
	case model2.State:
		pref = prefixState
	case model2.ApiResource:
		pref = prefixApiResource
	case model2.Alarm:
		pref = prefixAlarm
	}
	type getID interface {
		GetID() int64
	}
	if v, ok := obj.(getID); ok {
		return pref + strconv.FormatInt(v.GetID(), 10)
	}
	panic(errors.New("cache save: unknown obj"))
}

func (c *cache) Flush() {
	c.client.Flush()
}

func (c *cache) Foreach(fn func(key string, obj interface{})) {
	for k, v := range c.client.Items() {
		fn(k, v)
	}
}

func (c *cache) Save(obj interface{}) error {
	c.client.SetDefault(getKey(obj), obj)
	return nil
}

func (c *cache) Remove(obj interface{}) {
	c.client.Delete(getKey(obj))
}

func (c *cache) LoadConfig(id int64) (model2.Config, error) {
	if v, ok := c.client.Get(PrefixConfig + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model2.Config); ok {
			return u, nil
		}
	}
	return nil, lang2.Error(lang2.ErrCacheNotFound)
}

func (c *cache) LoadOrganization(id int64) (model2.Organization, error) {
	if v, ok := c.client.Get(PrefixOrg + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model2.Organization); ok {
			return u, nil
		}
	}
	return nil, lang2.Error(lang2.ErrCacheNotFound)
}

func (c *cache) LoadUser(id int64) (model2.User, error) {
	if v, ok := c.client.Get(PrefixUser + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model2.User); ok {
			return u, nil
		}
	}
	return nil, lang2.Error(lang2.ErrCacheNotFound)
}

func (c *cache) LoadRole(id int64) (model2.Role, error) {
	if v, ok := c.client.Get(PrefixRole + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model2.Role); ok {
			return u, nil
		}
	}
	return nil, lang2.Error(lang2.ErrCacheNotFound)
}

func (c *cache) LoadPolicy(id int64) (model2.Policy, error) {
	if v, ok := c.client.Get(prefixPolicy + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model2.Policy); ok {
			return u, nil
		}
	}
	return nil, lang2.Error(lang2.ErrCacheNotFound)
}

func (c *cache) LoadGroup(id int64) (model2.Group, error) {
	if v, ok := c.client.Get(prefixGroup + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model2.Group); ok {
			return u, nil
		}
	}
	return nil, lang2.Error(lang2.ErrCacheNotFound)
}

func (c *cache) LoadDevice(id int64) (model2.Device, error) {
	if v, ok := c.client.Get(prefixDevice + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model2.Device); ok {
			return u, nil
		}
	}
	return nil, lang2.Error(lang2.ErrCacheNotFound)
}

func (c *cache) LoadMeasure(id int64) (model2.Measure, error) {
	if v, ok := c.client.Get(prefixMeasure + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model2.Measure); ok {
			return u, nil
		}
	}
	return nil, lang2.Error(lang2.ErrCacheNotFound)
}

func (c *cache) LoadEquipment(id int64) (model2.Equipment, error) {
	if v, ok := c.client.Get(prefixEquipment + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model2.Equipment); ok {
			return u, nil
		}
	}
	return nil, lang2.Error(lang2.ErrCacheNotFound)
}

func (c *cache) LoadState(id int64) (model2.State, error) {
	if v, ok := c.client.Get(prefixState + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model2.State); ok {
			return u, nil
		}
	}
	return nil, lang2.Error(lang2.ErrCacheNotFound)
}

func (c *cache) LoadApiResource(id int64) (model2.ApiResource, error) {
	if v, ok := c.client.Get(prefixApiResource + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model2.ApiResource); ok {
			return u, nil
		}
	}
	return nil, lang2.Error(lang2.ErrCacheNotFound)
}

func (c *cache) LoadAlarm(id int64) (model2.Alarm, error) {
	if v, ok := c.client.Get(prefixAlarm + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model2.Alarm); ok {
			return u, nil
		}
	}
	return nil, lang2.Error(lang2.ErrCacheNotFound)
}
