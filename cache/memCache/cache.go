package memCache

import (
	"errors"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	goCache "github.com/patrickmn/go-cache"
	"strconv"
	"time"
)

const (
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
	var id int64
	switch v := obj.(type) {
	case model.Organization:
		pref = PrefixOrg
		id = v.GetID()
	case model.User:
		pref = PrefixUser
		id = v.GetID()
	case model.Role:
		pref = PrefixRole
		id = v.GetID()
	case model.Policy:
		pref = prefixPolicy
		id = v.GetID()
	case model.Group:
		pref = prefixGroup
		id = v.GetID()
	case model.Device:
		pref = prefixDevice
		id = v.GetID()
	case model.Measure:
		pref = prefixMeasure
		id = v.GetID()
	case model.Equipment:
		pref = prefixEquipment
		id = v.GetID()
	case model.State:
		pref = prefixState
		id = v.GetID()
	case model.ApiResource:
		pref = prefixApiResource
		id = v.GetID()
	default:
		panic(errors.New("cache save: unknown obj"))
	}
	return pref + strconv.FormatInt(id, 10)
}

func (c *cache) Flush() {
	c.client.Flush()
}

func (c *cache) Save(obj interface{}) error {
	c.client.SetDefault(getKey(obj), obj)
	return nil
}

func (c *cache) Remove(obj interface{}) {
	c.client.Delete(getKey(obj))
}

func (c *cache) LoadOrganization(id int64) (model.Organization, error) {
	if v, ok := c.client.Get(PrefixOrg + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model.Organization); ok {
			return u, nil
		}
	}
	return nil, lang.Error(lang.ErrCacheNotFound)
}

func (c *cache) LoadUser(id int64) (model.User, error) {
	if v, ok := c.client.Get(PrefixUser + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model.User); ok {
			return u, nil
		}
	}
	return nil, lang.Error(lang.ErrCacheNotFound)
}

func (c *cache) LoadRole(id int64) (model.Role, error) {
	if v, ok := c.client.Get(PrefixRole + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model.Role); ok {
			return u, nil
		}
	}
	return nil, lang.Error(lang.ErrCacheNotFound)
}

func (c *cache) LoadPolicy(id int64) (model.Policy, error) {
	if v, ok := c.client.Get(prefixPolicy + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model.Policy); ok {
			return u, nil
		}
	}
	return nil, lang.Error(lang.ErrCacheNotFound)
}

func (c *cache) LoadGroup(id int64) (model.Group, error) {
	if v, ok := c.client.Get(prefixGroup + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model.Group); ok {
			return u, nil
		}
	}
	return nil, lang.Error(lang.ErrCacheNotFound)
}

func (c *cache) LoadDevice(id int64) (model.Device, error) {
	if v, ok := c.client.Get(prefixDevice + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model.Device); ok {
			return u, nil
		}
	}
	return nil, lang.Error(lang.ErrCacheNotFound)
}

func (c *cache) LoadMeasure(id int64) (model.Measure, error) {
	if v, ok := c.client.Get(prefixMeasure + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model.Measure); ok {
			return u, nil
		}
	}
	return nil, lang.Error(lang.ErrCacheNotFound)
}

func (c *cache) LoadEquipment(id int64) (model.Equipment, error) {
	if v, ok := c.client.Get(prefixEquipment + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model.Equipment); ok {
			return u, nil
		}
	}
	return nil, lang.Error(lang.ErrCacheNotFound)
}

func (c *cache) LoadState(id int64) (model.State, error) {
	if v, ok := c.client.Get(prefixState + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model.State); ok {
			return u, nil
		}
	}
	return nil, lang.Error(lang.ErrCacheNotFound)
}

func (c *cache) LoadApiResource(id int64) (model.ApiResource, error) {
	if v, ok := c.client.Get(prefixApiResource + strconv.FormatInt(id, 10)); ok {
		if u, ok := v.(model.ApiResource); ok {
			return u, nil
		}
	}
	return nil, lang.Error(lang.ErrCacheNotFound)
}
