package mysqlStore

import (
	"github.com/maritimusj/centrum/dirty"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"time"
)

type Policy struct {
	id          int64
	enable      int8
	roleID      int64
	resourceUID string
	action      resource.Action
	effect      resource.Effect
	createdAt   time.Time

	dirty *dirty.Dirty
	store *mysqlStore
}

func (p *Policy) SetEffect(effect resource.Effect) error {
	panic("implement me")
}

func (p *Policy) IsAllow() bool {
	panic("implement me")
}

func (p *Policy) IsDeny() bool {
	panic("implement me")
}

func NewPolicy(s *mysqlStore, id int64) *Policy {
	return &Policy{
		id:    id,
		dirty: dirty.New(),
		store: s,
	}
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

func (p *Policy) Resource() resource.Resource {
	panic("implement me")
}

func (p *Policy) ResourceUID() string {
	panic("implement me")
}

func (p *Policy) Action() resource.Action {
	panic("implement me")
}

func (p *Policy) Effect() resource.Effect {
	panic("implement me")
}
