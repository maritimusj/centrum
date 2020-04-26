package mysqlStore

import (
	"time"

	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/dirty"
	"github.com/maritimusj/centrum/gate/web/model"
	"github.com/maritimusj/centrum/gate/web/resource"
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

func (p *Policy) Resource() model.Resource {
	res, _ := p.store.GetResource(p.resourceClass, p.resourceID)
	return res
}

func (p *Policy) Action() resource.Action {
	return p.action
}

func (p *Policy) SetEffect(effect resource.Effect) {
	if p.effect != effect {
		p.effect = effect
		p.dirty.Set("effect", func() interface{} {
			return p.effect
		})
	}
}

func (p *Policy) Effect() resource.Effect {
	return p.effect
}

func (p *Policy) Simple() model.Map {
	if p == nil {
		return model.Map{}
	}
	return model.Map{
		"id": p.id,
		"resource": model.Map{
			"class": p.resourceClass,
			"id":    p.resourceID,
		},
		"action": p.action,
		"effect": p.effect,
	}
}

func (p *Policy) Brief() model.Map {
	if p == nil {
		return model.Map{}
	}
	return model.Map{
		"id": p.id,
		"resource": model.Map{
			"class": p.resourceClass,
			"id":    p.resourceID,
		},
		"action":     p.action,
		"effect":     p.effect,
		"created_at": p.createdAt.Format(lang.DatetimeFormatterStr.Str()),
	}
}

func (p *Policy) Detail() model.Map {
	if p == nil {
		return model.Map{}
	}
	return model.Map{
		"id": p.id,
		"resource": model.Map{
			"class": p.resourceClass,
			"id":    p.resourceID,
		},
		"action":     p.action,
		"effect":     p.effect,
		"created_at": p.createdAt.Format(lang.DatetimeFormatterStr.Str()),
	}
}
