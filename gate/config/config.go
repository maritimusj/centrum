package config

import (
	"sync"
	"time"

	"github.com/maritimusj/centrum/gate/web/model"
	"github.com/maritimusj/centrum/gate/web/resource"
	"github.com/maritimusj/centrum/gate/web/store"
	"github.com/maritimusj/centrum/util"
)

const (
	baseConfigPath = "__base"

	ApiAddrPath                = "api.addr"
	ApiPortPath                = "api.port"
	InversePortPath            = "inverse.port"
	DefaultUserNamePath        = "default.username"
	DefaultOrganizationPath    = "default.organization"
	DefaultEffectPath          = "default.effect"
	DefaultPageSizePath        = "default.pagesize"
	DefaultTokenExpirationPath = "default.token.expiration"
	LogLevelPath               = "log.level"
	LogFileNamePath            = "log.filename"

	InfluxDBUrl      = "influxdb.url"
	InfluxDBUserName = "influxdb.username"
	InfluxDBPassword = "influxdb.password"
)

const (
	DefaultApiPort         = 9090
	DefaultInversePort     = 10502
	DefaultOrganization    = "default"
	DefaultUserName        = "admin"
	DefaultPageSize        = 20
	DefaultTokenExpiration = 3600 // second
	DefaultLogLevel        = "info"
	DefaultLogFileName     = "./log.data"
)

type Config struct {
	store store.Store

	BaseConfig model.Config
	Status     sync.Map
}

func New(store store.Store) *Config {
	return &Config{
		store: store,
	}
}

func (c *Config) Load() error {
	var err error
	base, err := c.store.GetConfig(baseConfigPath)
	if err != nil {
		base, err = c.store.CreateConfig(baseConfigPath, nil)
		if err != nil {
			return err
		}
	}

	c.BaseConfig = base
	return nil
}

func (c *Config) Save() error {
	if c.BaseConfig != nil {
		return c.BaseConfig.Save()
	}
	return nil
}

func (c *Config) APIAddr() string {
	addr := c.BaseConfig.GetOption(ApiAddrPath)
	if addr.Exists() {
		return addr.Str
	}
	return ""
}

func (c *Config) APIPort() int {
	port := c.BaseConfig.GetOption(ApiPortPath)
	if port.Exists() {
		return int(port.Int())
	}
	return DefaultApiPort
}

func (c *Config) InversePort() int {
	port := c.BaseConfig.GetOption(InversePortPath)
	if port.Exists() {
		return int(port.Int())
	}
	return DefaultInversePort
}

func (c *Config) DefaultEffect() resource.Effect {
	effect := c.BaseConfig.GetOption(DefaultEffectPath)
	if effect.Exists() {
		return resource.Effect(effect.Int())
	}
	return resource.Deny
}

func (c *Config) DefaultOrganization() string {
	org := c.BaseConfig.GetOption(DefaultOrganizationPath)
	if org.Exists() {
		return org.Str
	}
	return DefaultOrganization
}

func (c *Config) DefaultUserName() string {
	username := c.BaseConfig.GetOption(DefaultUserNamePath)
	if username.Exists() {
		return username.Str
	}
	return DefaultUserName
}

func (c *Config) DefaultPageSize() int64 {
	pagesize := c.BaseConfig.GetOption(DefaultPageSizePath)
	if pagesize.Exists() {
		return pagesize.Int()
	}
	return DefaultPageSize
}

func (c *Config) JwtTokenKey() []byte {
	if v, ok := c.Status.Load("jwt"); ok {
		return v.([]byte)
	}
	jwt := []byte(util.RandStr(32, util.RandAll))
	c.Status.Store("jwt", jwt)
	return jwt
}

func (c *Config) DefaultTokenExpiration() time.Duration {
	exp := c.BaseConfig.GetOption(DefaultTokenExpirationPath)
	if exp.Exists() {
		return time.Duration(exp.Int()) * time.Second
	}
	return time.Duration(DefaultTokenExpiration) * time.Second
}

func (c *Config) LogLevel() string {
	level := c.BaseConfig.GetOption(LogLevelPath)
	if level.Exists() {
		return level.Str
	}
	return DefaultLogLevel
}

func (c *Config) SetLogLevel(level string) {
	_ = c.BaseConfig.SetOption(LogLevelPath, level)
}

func (c *Config) LogFileName() string {
	filename := c.BaseConfig.GetOption(LogFileNamePath)
	if filename.Exists() {
		return filename.Str
	}
	return DefaultLogFileName
}

func (c *Config) InfluxDBConfig() map[string]string {
	return map[string]string{
		"url":      c.BaseConfig.GetOption(InfluxDBUrl).String(),
		"username": c.BaseConfig.GetOption(InfluxDBUserName).String(),
		"password": c.BaseConfig.GetOption(InfluxDBPassword).String(),
	}
}
