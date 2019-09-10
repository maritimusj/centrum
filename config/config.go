package config

import (
	"time"
)

type Config interface {
	DefaultPageSize() int64
	DBConnStr() string
	JwtTokenKey() []byte
	MaxTokenExpiration() time.Duration
	DefaultUserName() string
}

type config struct {
	jwtTokenKey []byte
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
	return "root:123456@/chuanyan?charset=utf8mb4&parseTime=true&loc=Local"
}

func (c *config) DefaultPageSize() int64 {
	return 10
}

func New() Config {
	return &config{
		jwtTokenKey: []byte("util.RandStr(32, util.RandAll)"),
	}
}
