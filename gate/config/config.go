package config

import (
	"sync"
	"time"

	"github.com/kataras/iris"

	"github.com/maritimusj/centrum/gate/web/model"
	"github.com/maritimusj/centrum/gate/web/resource"
	"github.com/maritimusj/centrum/gate/web/store"
	"github.com/maritimusj/centrum/util"
)

const (
	baseConfigPath   = "__base"
	streamConfigPath = "__stream"

	SysRegCodePath             = "sys.reg.code"
	SysRegOwnerPath            = "sys.reg.owner"
	SysTitlePath               = "sys.title"
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
	StreamURLsPath             = "stream.urls"
	WebViewURLsPath            = "web.urls"
	GeTuiAppIDPath             = "getui.app.id"
	GeTuiAppKeyPath            = "getui.app.key"
	GeTuiAppSecretPath         = "getui.app.secret"
	GeTuiMasterSecretPath      = "getui.master.secret"

	InfluxDBUrl      = "influxdb.url"
	InfluxDBUserName = "influxdb.username"
	InfluxDBPassword = "influxdb.password"
)

const (
	GeTuiAppIDEnvStr        = "GETUI_APP_ID"
	GeTuiAppKeyEnvStr       = "GETUI_APP_KEY"
	GeTuiAppSecretEnvStr    = "GETUI_APP_SECRET"
	GeTuiMasterSecretEnvStr = "GETUI_MASTER_SECRET"
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

	BaseConfig  model.Config
	ExtraConfig model.Config

	Status sync.Map
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

	stream, err := c.store.GetConfig(streamConfigPath)
	if err != nil {
		stream, err = c.store.CreateConfig(streamConfigPath, nil)
		if err != nil {
			return err
		}
	}

	c.ExtraConfig = stream
	return nil
}

func (c *Config) Save() error {
	if c.BaseConfig != nil {
		return c.BaseConfig.Save()
	}
	return nil
}

func (c *Config) StreamURLs() interface{} {
	stream := c.ExtraConfig.GetOption(StreamURLsPath)
	if stream.Exists() {
		return stream.Value()
	}
	return iris.Map{}
}

func (c *Config) SaveStreamURLs(streams interface{}) error {
	err := c.ExtraConfig.SetOption(StreamURLsPath, streams)
	if err != nil {
		return err
	}
	return c.ExtraConfig.Save()
}

func (c *Config) WebViewUrls() interface{} {
	stream := c.ExtraConfig.GetOption(WebViewURLsPath)
	if stream.Exists() {
		return stream.Value()
	}
	return []string{}
}

func (c *Config) SaveWebViewUrls(urls interface{}) error {
	err := c.ExtraConfig.SetOption(WebViewURLsPath, urls)
	if err != nil {
		return err
	}
	return c.ExtraConfig.Save()
}

func (c *Config) RegCode() string {
	code := c.BaseConfig.GetOption(SysRegCodePath)
	if code.Exists() {
		return code.Str
	}
	return ""
}

func (c *Config) RegOwner() string {
	owner := c.BaseConfig.GetOption(SysRegOwnerPath)
	if owner.Exists() {
		return owner.Str
	}
	return ""
}

func (c *Config) SysTitle() string {
	title := c.BaseConfig.GetOption(SysTitlePath)
	if title.Exists() {
		return title.Str
	}
	return ""
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

func (c *Config) GeTuiConfig() map[string]string {
	return map[string]string{
		"app_id":        c.BaseConfig.GetOption(GeTuiAppIDPath).String(),
		"app_key":       c.BaseConfig.GetOption(GeTuiAppKeyPath).String(),
		"app_secret":    c.BaseConfig.GetOption(GeTuiAppSecretPath).String(),
		"master_secret": c.BaseConfig.GetOption(GeTuiMasterSecretPath).String(),
	}
}
