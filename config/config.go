package config

import (
	"github.com/maritimusj/centrum/resource"
	"time"
)

type Config interface {
	LogLevel() string
	SetLogLevel(string)

	LogFileName() string

	DefaultPageSize() int64
	DBConnStr() string
	JwtTokenKey() []byte
	MaxTokenExpiration() time.Duration
	DefaultUserName() string

	DefaultEffect() resource.Effect
	IsRoleEnabled() bool
}

type config struct {
	jwtTokenKey []byte
	logLevel    string
	logFileName string
}

func (c *config) LogFileName() string {
	return c.logFileName
}

func (c *config) IsRoleEnabled() bool {
	return false
}

func (c *config) LogLevel() string {
	return c.logLevel
}

func (c *config) SetLogLevel(level string) {
	c.logLevel = level
}

func (c *config) DefaultEffect() resource.Effect {
	return resource.Allow
}

func (c *config) DefaultUserName() string {
	return "admin"
}

func (c *config) MaxTokenExpiration() time.Duration {
	return time.Hour * 10
}

func (c *config) JwtTokenKey() []byte {
	return c.jwtTokenKey
}

func (c *config) DBConnStr() string {
	return "root:12345678@/chuanyan?charset=utf8mb4&parseTime=true&loc=Local"
}

func (c *config) DefaultPageSize() int64 {
	return 10
}

func New() Config {
	return &config{
		jwtTokenKey: []byte("util.RandStr(32, util.RandAll)"),
		logLevel:    "error",
		logFileName: "./log.data",
	}
}
