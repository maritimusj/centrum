package mysqlStore

import (
	"time"

	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/dirty"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Config struct {
	id        int64
	name      string
	extra     []byte
	createdAt time.Time
	updateAt  time.Time

	dirty *dirty.Dirty
	store *mysqlStore
}

func NewConfig(store *mysqlStore, id int64) *Config {
	return &Config{
		id:    id,
		dirty: dirty.New(),
		store: store,
	}
}

func (config *Config) GetID() int64 {
	return config.id
}

func (config *Config) Name() string {
	return config.name
}

func (config *Config) UpdateAt() time.Time {
	return config.updateAt
}

func (config *Config) Option() map[string]interface{} {
	return gjson.ParseBytes(config.extra).Value().(map[string]interface{})
}

func (config *Config) GetOption(key string) gjson.Result {
	if config != nil {
		return gjson.GetBytes(config.extra, key)
	}
	return gjson.Result{}
}

func (config *Config) SetOption(key string, value interface{}) error {
	if config != nil {
		data, err := sjson.SetBytes(config.extra, key, value)
		if err != nil {
			return err
		}

		config.extra = data
		config.dirty.Set("extra", func() interface{} {
			return config.extra
		})

		return nil
	}
	return lang.ErrConfigNotFound.Error()
}

func (config *Config) CreatedAt() time.Time {
	if config != nil {
		return config.createdAt
	}
	return time.Time{}
}

func (config *Config) Destroy() error {
	if config == nil {
		return lang.ErrConfigNotFound.Error()
	}

	return config.store.RemoveConfig(config.id)
}

func (config *Config) Save() error {
	if config != nil {
		if config.dirty.Any() {
			err := SaveData(config.store.db, TbConfig, config.dirty.Data(true), "id=?", config.id)
			if err != nil {
				return lang.InternalError(err)
			}
		}
		return nil
	}
	return lang.ErrConfigNotFound.Error()
}
