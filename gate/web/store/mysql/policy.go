package mysqlStore

import (
	lang2 "github.com/maritimusj/centrum/gate/lang"
	dirty2 "github.com/maritimusj/centrum/gate/web/dirty"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	"time"
)

type Policy struct {
	id            int64
	roleID        int64
	resourceClass resource2.Class
	resourceID    int64
	action        resource2.Action
	effect        resource2.Effect
	createdAt     time.Time

	dirty *dirty2.Dirty
	store *mysqlStore
}

func NewPolicy(s *mysqlStore, id int64) *Policy {
	return &Policy{
		id:    id,
		dirty: dirty2.New(),
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
			return lang2.InternalError(err)
		}
	}
	return nil
}

func (p *Policy) Role() model2.Role {
	role, _ := p.store.GetRole(p.roleID)
	return role
}

func (p *Policy) Resource() model2.Resource {
	res, _ := p.store.GetResource(p.resourceClass, p.resourceID)
	return res
}

func (p *Policy) Action() resource2.Action {
	return p.action
}

func (p *Policy) SetEffect(effect resource2.Effect) {
	if p.effect != effect {
		p.effect = effect
		p.dirty.Set("effect", func() interface{} {
			return p.effect
		})
	}
}

func (p *Policy) Effect() resource2.Effect {
	return p.effect
}

func (p *Policy) Simple() model2.Map {
	if p == nil {
		return model2.Map{}
	}
	return model2.Map{
		"id": p.id,
		"resource": model2.Map{
			"class": p.resourceClass,
			"id":    p.resourceID,
		},
		"action": p.action,
		"effect": p.effect,
	}
}

func (p *Policy) Brief() model2.Map {
	if p == nil {
		return model2.Map{}
	}
	return model2.Map{
		"id": p.id,
		"resource": model2.Map{
			"class": p.resourceClass,
			"id":    p.resourceID,
		},
		"action":     p.action,
		"effect":     p.effect,
		"created_at": p.createdAt,
	}
}

func (p *Policy) Detail() model2.Map {
	if p == nil {
		return model2.Map{}
	}
	return model2.Map{
		"id": p.id,
		"resource": model2.Map{
			"class": p.resourceClass,
			"id":    p.resourceID,
		},
		"action":     p.action,
		"effect":     p.effect,
		"created_at": p.createdAt,
	}
}
