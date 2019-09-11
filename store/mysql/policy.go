package mysqlStore

import (
	"github.com/maritimusj/centrum/dirty"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"time"
)

type Policy struct {
	id            int64
	roleID        int64
	resourceClass resource.Class
	resourceID    int64
	action        resource.Action
	effect        resource.Effect
	createdAt     time.Time

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
	return p.id
}

func (p *Policy) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Policy) Destroy() error {
	return p.store.RemovePolicy(p.id)
}

func (p *Policy) Save() error {
	if p.dirty.Any() {
		err := SaveData(p.store.db, TbPolicies, p.dirty.Data(true), "id=?", p.id)
		if err != nil {
			return lang.InternalError(err)
		}
	}

	return nil
}

func (p *Policy) Role() model.Role {
	role, _ := p.store.GetRole(p.roleID)
	return role
}

func (p *Policy) Resource() resource.Resource {
	res, _ := p.store.GetResource(p.resourceClass, p.resourceID)
	return res
}

func (p *Policy) Action() resource.Action {
	return p.action
}

func (p *Policy) Effect() resource.Effect {
	return p.effect
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
