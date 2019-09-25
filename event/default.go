package event

import (
	"context"
)

const (
	_ = iota
	User
	Device
	Equipment
)

const (
	_ = iota
	Created
	Updated
	Deleted
)

var (
	defaultEvent = New()
)

func NewData(code int) *Data {
	return &Data{
		Code:   code,
		Values: map[string]interface{}{},
	}
}

func Wait() {
	defaultEvent.Wait()
}

func Sub(ctx context.Context, codes ...int) <-chan *Data {
	return defaultEvent.Sub(ctx, codes...)
}

func Register(ctx context.Context, fn CallbackFN, codes ...int) {
	defaultEvent.Register(ctx, fn, codes...)
}

func Fire(data *Data) {
	defaultEvent.Fire(data)
}
