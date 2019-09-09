package mysqlStore

import (
	"github.com/maritimusj/centrum/model"
	"time"
)

type Policy struct {
	id            int64
	enable        int8
	roleID        int64
	resourceGroup int
	resourceID    int64
	action        int8
	effect        int8
	createdAt     time.Time

	store *mysqlStore
}

func (p *Policy) GetID() int64 {
	panic("implement me")
}

func (p *Policy) CreatedAt() time.Time {
	panic("implement me")
}

func (p *Policy) Destroy() error {
	panic("implement me")
}

func (p *Policy) Save() error {
	panic("implement me")
}

func (p *Policy) Enable() error {
	panic("implement me")
}

func (p *Policy) Disable() error {
	panic("implement me")
}

func (p *Policy) IsEnabled() bool {
	panic("implement me")
}

func (p *Policy) Simple() model.Map {
	panic("implement me")
}

func (p *Policy) Brief() model.Map {
	panic("implement me")
}

func (p *Policy) Detail() model.Map {
	panic("implement me")
}

func (p *Policy) Role() model.Role {
	panic("implement me")
}

func (p *Policy) Resource() model.Resource {
	panic("implement me")
}

func (p *Policy) Action() model.Action {
	panic("implement me")
}

func (p *Policy) Effect() model.Effect {
	panic("implement me")
}
