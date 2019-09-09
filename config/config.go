package config

type Config interface {
	DefaultPageSize() int64
}

type config struct {
}

func (c *config) DefaultPageSize() int64 {
	return 10
}

func New() Config {
	return &config{}
}
