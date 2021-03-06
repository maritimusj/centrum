package mysqlStore

import (
	"time"

	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/dirty"
	"github.com/maritimusj/centrum/gate/web/helper"
	"github.com/maritimusj/centrum/gate/web/model"
	"github.com/maritimusj/centrum/gate/web/resource"
	"github.com/maritimusj/centrum/gate/web/status"
)

type Role struct {
	id    int64
	orgID int64

	enable int8

	name      string
	title     string
	desc      string
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

func (r *Role) OrganizationID() int64 {
	return r.orgID
}

func (r *Role) Organization() (model.Organization, error) {
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
			return lang.InternalError(err)
		}
	}

	return nil
}

func (r *Role) Destroy() error {
	policies, _, err := r.store.GetPolicyList(nil, helper.Role(r.GetID()))
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
	if r.enable != status.Enable {
		r.enable = status.Enable
		r.dirty.Set("enable", func() interface{} {
			return r.enable
		})
	}
}

func (r *Role) Disable() {
	if r.enable != status.Disable {
		r.enable = status.Disable
		r.dirty.Set("enable", func() interface{} {
			return r.enable
		})
	}
}

func (r *Role) IsEnabled() bool {
	return r.enable == status.Enable
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

func (r *Role) GetUserList(options ...helper.OptionFN) ([]model.User, int64, error) {
	return r.store.GetUserList(helper.Role(r.GetID()))
}

func (r *Role) SetPolicy(res model.Resource, action resource.Action, effect resource.Effect, recursiveMap map[model.Resource]struct{}) (model.Policy, error) {
	if recursiveMap != nil {
		if _, ok := recursiveMap[res]; ok {
			return nil, lang.ErrRecursiveDetected.Error()
		}
		recursiveMap[res] = struct{}{}
	}

	policy, err := r.store.GetPolicyFrom(r.id, res, action)
	if err != nil {
		if err != lang.ErrPolicyNotFound.Error() {
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

func (r *Role) RemovePolicy(res model.Resource) error {
	policies, _, err := r.store.GetPolicyList(res, helper.Role(r.id))
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

func (r *Role) GetPolicy(res model.Resource) (map[resource.Action]model.Policy, error) {
	policies, _, err := r.store.GetPolicyList(res, helper.Role(r.id))
	if err != nil {
		return nil, err
	}
	result := make(map[resource.Action]model.Policy)
	for _, policy := range policies {
		result[policy.Action()] = policy
	}
	return result, nil
}

func (r *Role) IsAllow(res model.Resource, action resource.Action) (bool, error) {
	pm, err := r.GetPolicy(res)
	if err != nil {
		return false, err
	}

	if v, ok := pm[action]; ok {
		if v.Effect() == resource.Allow {
			return true, nil
		}
		return false, lang.ErrNoPermission.Error()
	}

	return false, lang.ErrPolicyNotFound.Error()
}

func (r *Role) Simple() model.Map {
	if r == nil {
		return model.Map{}
	}
	return model.Map{
		"id":     r.id,
		"enable": r.IsEnabled(),
		"name":   r.name,
		"title":  r.title,
	}
}

func (r *Role) Brief() model.Map {
	if r == nil {
		return model.Map{}
	}
	return model.Map{
		"id":         r.id,
		"enable":     r.IsEnabled(),
		"name":       r.name,
		"title":      r.title,
		"desc":       r.desc,
		"created_at": r.createdAt.Format(lang.DatetimeFormatterStr.Str()),
	}
}

func (r *Role) Detail() model.Map {
	if r == nil {
		return model.Map{}
	}
	return model.Map{
		"id":         r.id,
		"enable":     r.IsEnabled(),
		"name":       r.name,
		"title":      r.title,
		"desc":       r.desc,
		"created_at": r.createdAt.Format(lang.DatetimeFormatterStr.Str()),
	}
}
