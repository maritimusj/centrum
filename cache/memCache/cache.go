package memCache

import (
	"github.com/maritimusj/centrum/model"
	goCache "github.com/patrickmn/go-cache"
	"time"
)

type cache struct {
	client *goCache.Cache
}

func New() *cache {
	return &cache{client: goCache.New(6*time.Minute, 10*time.Minute)}
}

func (c *cache) Save(obj interface{}) {

}

func (c *cache) Remove(obj interface{}) {

}

func (c *cache) LoadUser(id int64) (model.User, error) {
	return nil, nil
}
