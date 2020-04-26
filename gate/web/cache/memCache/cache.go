package memCache

import (
	"fmt"
	"strconv"
	"time"

	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/model"
	goCache "github.com/patrickmn/go-cache"
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
	prefixAlarm       = "A."
	prefixComment     = "B."
)

type cache struct {
	client *goCache.Cache
}

func New() *cache {
	return &cache{client: goCache.New(6*time.Minute, 10*time.Minute)}
}

func getKeys(obj interface{}) []string {
	var pref string
	switch v := obj.(type) {
	case string:
		return []string{v}
	case model.Organization:
		pref = PrefixOrg
	case model.User:
		pref = PrefixUser
	case model.Role:
		pref = PrefixRole
	case model.Policy:
		pref = prefixPolicy
	case model.Group:
		pref = prefixGroup
	case model.Device:
		pref = prefixDevice
	case model.Measure:
		pref = prefixMeasure
	case model.Equipment:
		pref = prefixEquipment
	case model.State:
		pref = prefixState
	case model.ApiResource:
		pref = prefixApiResource
	case model.Alarm:
		pref = prefixAlarm
	case model.Comment:
		pref = prefixComment
	}

	keys := make([]string, 0)
	type getID interface {
		GetID() int64
	}

	type getName interface {
		Name() string
	}

	if v, ok := obj.(getID); ok {
		keys = append(keys, pref+strconv.FormatInt(v.GetID(), 10))
	}

	if v, ok := obj.(getName); ok {
		name := v.Name()
		if name != "" {
			keys = append(keys, pref+name)
		}
	}

	return keys
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
	for _, key := range getKeys(obj) {
		c.client.SetDefault(key, obj)
	}

	return nil
}

func (c *cache) Remove(obj interface{}) {
	for _, key := range getKeys(obj) {
		if o, ok := c.client.Get(key); ok {
			for _, key := range getKeys(o) {
				c.client.Delete(key)
			}
			c.client.Delete(key)
		}
	}
}

func (c *cache) getUID(v interface{}) string {
	switch vv := v.(type) {
	case int64:
		return strconv.FormatInt(vv, 10)
	case float64:
		return strconv.FormatInt(int64(vv), 10)
	case string:
		return vv
	case model.DBEntry:
		return strconv.FormatInt(vv.GetID(), 10)
	case model.LogEntry:
		return vv.UID()
	}
	panic(fmt.Errorf("cache: unknown uid: %#v", v))
}

func (c *cache) LoadConfig(config interface{}) (model.Config, error) {
	if v, ok := c.client.Get(PrefixConfig + c.getUID(config)); ok {
		if u, ok := v.(model.Config); ok {
			return u, nil
		}
	}
	return nil, lang.ErrCacheNotFound.Error()
}

func (c *cache) LoadOrganization(o interface{}) (model.Organization, error) {
	if v, ok := c.client.Get(PrefixOrg + c.getUID(o)); ok {
		if u, ok := v.(model.Organization); ok {
			return u, nil
		}
	}
	return nil, lang.ErrCacheNotFound.Error()
}

func (c *cache) LoadUser(u interface{}) (model.User, error) {
	if v, ok := c.client.Get(PrefixUser + c.getUID(u)); ok {
		if u, ok := v.(model.User); ok {
			return u, nil
		}
	}
	return nil, lang.ErrCacheNotFound.Error()
}

func (c *cache) LoadRole(r interface{}) (model.Role, error) {
	if v, ok := c.client.Get(PrefixRole + c.getUID(r)); ok {
		if u, ok := v.(model.Role); ok {
			return u, nil
		}
	}
	return nil, lang.ErrCacheNotFound.Error()
}

func (c *cache) LoadPolicy(p interface{}) (model.Policy, error) {
	if v, ok := c.client.Get(prefixPolicy + c.getUID(p)); ok {
		if u, ok := v.(model.Policy); ok {
			return u, nil
		}
	}
	return nil, lang.ErrCacheNotFound.Error()
}

func (c *cache) LoadGroup(g interface{}) (model.Group, error) {
	if v, ok := c.client.Get(prefixGroup + c.getUID(g)); ok {
		if u, ok := v.(model.Group); ok {
			return u, nil
		}
	}
	return nil, lang.ErrCacheNotFound.Error()
}

func (c *cache) LoadDevice(device interface{}) (model.Device, error) {
	if v, ok := c.client.Get(prefixDevice + c.getUID(device)); ok {
		if u, ok := v.(model.Device); ok {
			return u, nil
		}
	}
	return nil, lang.ErrCacheNotFound.Error()
}

func (c *cache) LoadMeasure(m interface{}) (model.Measure, error) {
	if v, ok := c.client.Get(prefixMeasure + c.getUID(m)); ok {
		if u, ok := v.(model.Measure); ok {
			return u, nil
		}
	}
	return nil, lang.ErrCacheNotFound.Error()
}

func (c *cache) LoadEquipment(e interface{}) (model.Equipment, error) {
	if v, ok := c.client.Get(prefixEquipment + c.getUID(e)); ok {
		if u, ok := v.(model.Equipment); ok {
			return u, nil
		}
	}
	return nil, lang.ErrCacheNotFound.Error()
}

func (c *cache) LoadState(state interface{}) (model.State, error) {
	if v, ok := c.client.Get(prefixState + c.getUID(state)); ok {
		if u, ok := v.(model.State); ok {
			return u, nil
		}
	}
	return nil, lang.ErrCacheNotFound.Error()
}

func (c *cache) LoadApiResource(api interface{}) (model.ApiResource, error) {
	if v, ok := c.client.Get(prefixApiResource + c.getUID(api)); ok {
		if u, ok := v.(model.ApiResource); ok {
			return u, nil
		}
	}
	return nil, lang.ErrCacheNotFound.Error()
}

func (c *cache) LoadAlarm(alarm interface{}) (model.Alarm, error) {
	if v, ok := c.client.Get(prefixAlarm + c.getUID(alarm)); ok {
		if u, ok := v.(model.Alarm); ok {
			return u, nil
		}
	}
	return nil, lang.ErrCacheNotFound.Error()
}
func (c *cache) LoadComment(comment interface{}) (model.Comment, error) {
	if v, ok := c.client.Get(prefixComment + c.getUID(comment)); ok {
		if u, ok := v.(model.Comment); ok {
			return u, nil
		}
	}
	return nil, lang.ErrCommentNotFound.Error()
}
