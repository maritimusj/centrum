package config

import (
	"github.com/maritimusj/centrum/web/resource"
	"github.com/maritimusj/centrum/web/store"
	"time"
)

type Config struct {
	store store.Store

	APIAddr string
	APIPort int

	DefaultEffect resource.Effect

	DefaultOrganization string
	DefaultUserName     string

	DefaultPageSize int64

	JwtTokenKey            []byte
	DefaultTokenExpiration time.Duration

	LogLevel    string
	LogFileName string
}

func New(store store.Store) *Config {
	return &Config{
		APIAddr:         "",
		APIPort:         9090,
		DefaultUserName: "admin",

		DefaultOrganization: "default",
		DefaultEffect:       resource.Deny,

		DefaultPageSize:        20,
		DefaultTokenExpiration: time.Hour * 10,

		JwtTokenKey: []byte("util.RandStr(32, util.RandAll)"),

		LogLevel:    "error",
		LogFileName: "./log.data",

		store: store,
	}
}

func (c *Config) Load() error {
	base, err := c.store.GetConfig("base")
	if err != nil {
		return err
	}
	c.APIAddr = base.GetOption("api.addr").Str
	c.APIPort = int(base.GetOption("api.port").Int())
	c.DefaultUserName = base.GetOption("default.username").Str
	c.DefaultOrganization = base.GetOption("default.organization").Str
	c.DefaultEffect = resource.Effect(base.GetOption("default.effect").Int())
	c.DefaultPageSize = base.GetOption("default.pagesize").Int()
	c.DefaultTokenExpiration = time.Duration(base.GetOption("default.token.expiration").Int())

	c.LogLevel = base.GetOption("log.level").Str
	c.LogFileName = base.GetOption("log.filename").Str

	return nil
}

func (c *Config) Save() error {
	panic("implement me")
}
