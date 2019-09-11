package mysqlStore

import (
	"github.com/maritimusj/centrum/dirty"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/status"
	"time"
)

type Role struct {
	id        int64
	enable    int8
	title     string
	createdAt time.Time

	dirty *dirty.Dirty
	store *mysqlStore
}

func NewRole(s *mysqlStore, id int64) *Role {
	return &Role{
		id:    id,
		dirty: dirty.New(),
		store: s,
	}
}

func (r *Role) GetID() int64 {
	return r.id
}

func (r *Role) CreatedAt() time.Time {
	return r.createdAt
}

func (r *Role) Save() error {
	if r.dirty.Any() {
		err := SaveData(r.store.db, TbRoles, r.dirty.Data(true), "id=?", r.id)
		if err != nil {
			return lang.InternalError(err)
		}
	}

	return nil
}

func (r *Role) Destroy() error {
	return r.store.RemoveRole(r.id)
}

func (r *Role) Enable() error {
	if r.enable != status.Enable {
		r.enable = status.Enable
		r.dirty.Set("enable", func() interface{} {
			return r.enable
		})
	}
	return r.Save()
}

func (r *Role) Disable() error {
	if r.enable != status.Disable {
		r.enable = status.Disable
		r.dirty.Set("enable", func() interface{} {
			return r.enable
		})
	}
	return r.Save()
}

func (r *Role) IsEnabled() bool {
	return r.enable == status.Enable
}

func (r *Role) Title() string {
	return r.title
}

func (r *Role) SetTitle(title string) error {
	if r.title != title {
		r.title = title
		r.dirty.Set("title", func() interface{} {
			return r.title
		})
	}
	return r.Save()
}

func (r *Role) SetPolicy(res resource.Resource, action resource.Action, effect resource.Effect) (model.Policy, error) {
	policy, err := r.store.CreatePolicyIsNotExists(r.id, res.ResourceUID(), action)
	if err != nil {
		return nil, err
	}
	err = policy.SetEffect(effect)
	if err != nil {
		return nil, err
	}
	return policy, nil
}

func (r *Role) GetPolicy(res resource.Resource) (map[resource.Action]model.Policy, error) {
	policies, _, err := r.store.GetPolicyList(r.id, res.ResourceUID())
	if err != nil {
		return nil, err
	}
	result := make(map[resource.Action]model.Policy)
	for _, policy := range policies {
		result[policy.Action()] = policy
	}
	return result, nil
}

func (r *Role) IsAllowed(res resource.Resource, action resource.Action) (bool, error) {
	pm, err := r.GetPolicy(res)
	if err != nil {
		return false, err
	}
	if v, ok := pm[action]; ok {
		if v.Effect() == resource.Allow {
			return true, nil
		}
		return false, lang.Error(lang.ErrNoPermission)
	}

	return false, lang.Error(lang.ErrPolicyNotFound)
}

func (r *Role) Simple() model.Map {
	panic("implement me")
}

func (r *Role) Brief() model.Map {
	panic("implement me")
}

func (r *Role) Detail() model.Map {
	panic("implement me")
}
