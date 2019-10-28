package mysqlStore

import (
	lang2 "github.com/maritimusj/centrum/gate/lang"
	dirty2 "github.com/maritimusj/centrum/gate/web/dirty"
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	status2 "github.com/maritimusj/centrum/gate/web/status"
	"time"
)

type Role struct {
	id    int64
	orgID int64

	enable int8

	name      string
	title     string
	desc      string
	createdAt time.Time

	dirty *dirty2.Dirty
	store *mysqlStore
}

func NewRole(s *mysqlStore, id int64) *Role {
	return &Role{
		id:    id,
		dirty: dirty2.New(),
		store: s,
	}
}

func (r *Role) OrganizationID() int64 {
	return r.orgID
}

func (r *Role) Organization() (model2.Organization, error) {
	return r.store.GetOrganization(r.orgID)
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
			return lang2.InternalError(err)
		}
	}

	return nil
}

func (r *Role) Destroy() error {
	policies, _, err := r.store.GetPolicyList(nil, helper2.Role(r.GetID()))
	if err != nil {
		return err
	}
	for _, p := range policies {
		err = p.Destroy()
		if err != nil {
			return err
		}
	}

	return r.store.RemoveRole(r.id)
}

func (r *Role) Enable() {
	if r.enable != status2.Enable {
		r.enable = status2.Enable
		r.dirty.Set("enable", func() interface{} {
			return r.enable
		})
	}
}

func (r *Role) Disable() {
	if r.enable != status2.Disable {
		r.enable = status2.Disable
		r.dirty.Set("enable", func() interface{} {
			return r.enable
		})
	}
}

func (r *Role) IsEnabled() bool {
	return r.enable == status2.Enable
}

func (r *Role) Name() string {
	return r.name
}

func (r *Role) Title() string {
	return r.title
}

func (r *Role) Desc() string {
	return r.desc
}

func (r *Role) SetTitle(title string) {
	if r.title != title {
		r.title = title
		r.dirty.Set("title", func() interface{} {
			return r.title
		})
	}
}

func (r *Role) SetDesc(desc string) {
	if r.desc != desc {
		r.desc = desc
		r.dirty.Set("desc", func() interface{} {
			return r.desc
		})
	}
}

func (r *Role) GetUserList(options ...helper2.OptionFN) ([]model2.User, int64, error) {
	return r.store.GetUserList(helper2.Role(r.GetID()))
}

func (r *Role) SetPolicy(res model2.Resource, action resource2.Action, effect resource2.Effect, recursiveMap map[model2.Resource]struct{}) (model2.Policy, error) {
	if recursiveMap != nil {
		if _, ok := recursiveMap[res]; ok {
			return nil, lang2.Error(lang2.ErrRecursiveDetected)
		}
		recursiveMap[res] = struct{}{}
	}

	policy, err := r.store.GetPolicyFrom(r.id, res, action)
	if err != nil {
		if err != lang2.Error(lang2.ErrPolicyNotFound) {
			return nil, err
		}

		policy, err = r.store.CreatePolicy(r.id, res, action, effect)
		if err != nil {
			return nil, err
		}
	}

	if recursiveMap != nil {
		//递归设置所有子资源的权限
		children, _, err := res.GetChildrenResources()
		if err != nil {
			return nil, err
		}

		for _, res := range children {
			_, err = r.SetPolicy(res, action, effect, recursiveMap)
			if err != nil {
				return nil, err
			}
		}
	}

	policy.SetEffect(effect)
	err = policy.Save()
	if err != nil {
		return nil, err
	}

	return policy, nil
}

func (r *Role) RemovePolicy(res model2.Resource) error {
	policies, _, err := r.store.GetPolicyList(res, helper2.Role(r.id))
	if err != nil {
		return err
	}
	for _, policy := range policies {
		if err = policy.Destroy(); err != nil {
			return err
		}
	}
	return nil
}

func (r *Role) GetPolicy(res model2.Resource) (map[resource2.Action]model2.Policy, error) {
	policies, _, err := r.store.GetPolicyList(res, helper2.Role(r.id))
	if err != nil {
		return nil, err
	}
	result := make(map[resource2.Action]model2.Policy)
	for _, policy := range policies {
		result[policy.Action()] = policy
	}
	return result, nil
}

func (r *Role) IsAllow(res model2.Resource, action resource2.Action) (bool, error) {
	pm, err := r.GetPolicy(res)
	if err != nil {
		return false, err
	}

	if v, ok := pm[action]; ok {
		if v.Effect() == resource2.Allow {
			return true, nil
		}
		return false, lang2.Error(lang2.ErrNoPermission)
	}

	return false, lang2.Error(lang2.ErrPolicyNotFound)
}

func (r *Role) Simple() model2.Map {
	if r == nil {
		return model2.Map{}
	}
	return model2.Map{
		"id":     r.id,
		"enable": r.IsEnabled(),
		"name":   r.name,
		"title":  r.title,
	}
}

func (r *Role) Brief() model2.Map {
	if r == nil {
		return model2.Map{}
	}
	return model2.Map{
		"id":         r.id,
		"enable":     r.IsEnabled(),
		"name":       r.name,
		"title":      r.title,
		"desc":       r.desc,
		"created_at": r.createdAt,
	}
}

func (r *Role) Detail() model2.Map {
	if r == nil {
		return model2.Map{}
	}
	return model2.Map{
		"id":         r.id,
		"enable":     r.IsEnabled(),
		"name":       r.name,
		"title":      r.title,
		"desc":       r.desc,
		"created_at": r.createdAt,
	}
}
