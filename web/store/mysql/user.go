package mysqlStore

import (
	"errors"
	"fmt"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/util"
	"github.com/maritimusj/centrum/web/dirty"
	"github.com/maritimusj/centrum/web/helper"
	"github.com/maritimusj/centrum/web/model"
	"github.com/maritimusj/centrum/web/resource"
	"github.com/maritimusj/centrum/web/status"
	log "github.com/sirupsen/logrus"
	"time"
)

type User struct {
	id    int64
	orgID int64

	enable int8

	name      string
	title     string
	password  []byte
	mobile    string
	email     string
	createdAt time.Time

	dirty *dirty.Dirty
	store *mysqlStore
}

func NewUser(s *mysqlStore, id int64) *User {
	return &User{
		id:    id,
		dirty: dirty.New(),
		store: s,
	}
}

func (u *User) LogUID() string {
	if u != nil {
		return fmt.Sprintf("user:%d", u.id)
	}
	return "user:<unknown>"
}

func (u *User) Logger() *log.Entry {
	return log.WithFields(log.Fields{
		"org": u.OrganizationID(),
		"src": u.LogUID(),
	})
}

func (u *User) OrganizationID() int64 {
	if u != nil {
		return u.orgID
	}
	return 0
}

func (u *User) Organization() (model.Organization, error) {
	if u != nil {
		return u.store.GetOrganization(u.orgID)
	}
	return nil, lang.Error(lang.ErrUserNotFound)
}

func (u *User) GetID() int64 {
	if u != nil {
		return u.id
	}
	return 0
}

func (u *User) Enable() {
	if u != nil && u.enable != status.Enable {
		u.enable = status.Enable
		u.dirty.Set("enable", func() interface{} {
			return u.enable
		})
	}
}

func (u *User) Disable() {
	if u != nil && u.enable != status.Disable {
		u.enable = status.Disable
		u.dirty.Set("enable", func() interface{} {
			return u.enable
		})
	}
}

func (u *User) IsEnabled() bool {
	if u != nil {
		return u.enable == status.Enable
	}
	return false
}

func (u *User) Name() string {
	if u != nil {
		return u.name
	}
	return "<unknown>"
}

func (u *User) Title() string {
	if u != nil {
		return u.title
	}
	return "<unknown>"
}

func (u *User) Mobile() string {
	if u != nil {
		return u.mobile
	}
	return "<unknown>"
}

func (u *User) Email() string {
	if u != nil {
		return u.email
	}
	return "<unknown>"
}

func (u *User) CreatedAt() time.Time {
	if u != nil {
		return u.createdAt
	}
	return time.Time{}
}

func (u *User) ResetPassword(password string) {
	if u != nil {
		data, _ := util.HashPassword([]byte(password))

		u.password = data
		u.dirty.Set("password", func() interface{} {
			return u.password
		})
	}
}

func (u *User) CheckPassword(password string) bool {
	if u != nil {
		return util.ComparePassword(u.password, []byte(password))
	}
	return false
}

func (u *User) Update(profile model.Map) {
	if u != nil {
		if enable, ok := profile["enable"].(int8); ok && enable != u.enable {
			u.enable = enable
			u.dirty.Set("enable", func() interface{} {
				return u.enable
			})
		}
		if title, ok := profile["title"].(string); ok && title != u.title {
			u.title = title
			u.dirty.Set("title", func() interface{} {
				return u.title
			})
		}

		if mobile, ok := profile["mobile"].(string); ok && mobile != u.mobile {
			u.mobile = mobile
			u.dirty.Set("mobile", func() interface{} {
				return u.mobile
			})
		}

		if email, ok := profile["email"].(string); ok && email != u.email {
			u.email = email
			u.dirty.Set("email", func() interface{} {
				return u.email
			})
		}
	}
}

func (u *User) SetRoles(roles ...interface{}) error {
	if u == nil {
		return lang.Error(lang.ErrUserNotFound)
	}

	err := RemoveData(u.store.db, TbUserRoles, "user_id=?", u.id)
	if err != nil {
		return err
	}

	now := time.Now()

	for _, r := range roles {
		role, err := u.store.GetRole(r)
		if err != nil {
			return err
		}
		_, err = CreateData(u.store.db, TbUserRoles, map[string]interface{}{
			"user_id":    u.id,
			"role_id":    role.GetID(),
			"created_at": now,
		})
		if err != nil {
			return lang.InternalError(err)
		}
	}
	return nil
}

func (u *User) GetRoles() ([]model.Role, error) {
	if u == nil {
		return nil, lang.Error(lang.ErrUserNotFound)
	}

	roles, _, err := u.store.GetRoleList(helper.User(u.id))
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (u *User) Is(role interface{}) (bool, error) {
	if u == nil {
		return false, lang.Error(lang.ErrUserNotFound)
	}

	roles, err := u.GetRoles()
	if err != nil {
		return false, err
	}
	var fn func(role model.Role) bool
	switch v := role.(type) {
	case int64:
		fn = func(role model.Role) bool {
			return role.GetID() == v
		}
	case string:
		fn = func(role model.Role) bool {
			return role.Name() == v
		}
	case model.Role:
		fn = func(role model.Role) bool {
			return role.GetID() == v.GetID()
		}
	default:
		panic(errors.New("user.Is() unknown role"))
	}
	for _, role := range roles {
		if fn(role) {
			return true, nil
		}
	}
	return false, nil
}

func (u *User) Destroy() error {
	if u != nil {
		return u.store.RemoveUser(u.id)
	}

	return lang.Error(lang.ErrUserNotFound)
}

func (u *User) Save() error {
	if u == nil {
		return lang.Error(lang.ErrUserNotFound)
	}

	if u.dirty.Any() {
		err := SaveData(u.store.db, TbUsers, u.dirty.Data(true), "id=?", u.id)
		if err != nil {
			return lang.InternalError(err)
		}
	}
	return nil
}

func (u *User) RemovePolicies(res model.Resource) error {
	if u == nil {
		return lang.Error(lang.ErrUserNotFound)
	}

	roles, err := u.GetRoles()
	if err != nil {
		return err
	}

	for _, role := range roles {
		err = role.RemovePolicy(res)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *User) SetDeny(res model.Resource, actions ...resource.Action) error {
	if u == nil {
		return lang.Error(lang.ErrUserNotFound)
	}

	role, err := u.store.GetRole(u.Name())
	if err != nil {
		return err
	}

	for _, action := range actions {
		_, err = role.SetPolicy(res, action, resource.Deny, make(map[model.Resource]struct{}))
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *User) SetAllow(res model.Resource, actions ...resource.Action) error {
	if u == nil {
		return lang.Error(lang.ErrUserNotFound)
	}

	role, err := u.store.GetRole(u.Name())
	if err != nil {
		return err
	}

	for _, action := range actions {
		_, err = role.SetPolicy(res, action, resource.Allow, make(map[model.Resource]struct{}))
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *User) IsAllow(res model.Resource, action resource.Action) (bool, error) {
	if u == nil {
		return false, lang.Error(lang.ErrUserNotFound)
	}

	if res.OrganizationID() > 0 && res.OrganizationID() != u.OrganizationID() {
		return false, lang.Error(lang.ErrOrganizationDifferent)
	}

	roles, err := u.GetRoles()
	if err != nil {
		return false, err
	}

	var denied bool
	for _, role := range roles {
		if allowed, err := role.IsAllow(res, action); allowed {
			return allowed, err
		} else {
			if err == lang.Error(lang.ErrNoPermission) {
				denied = true
			}
		}
	}

	if denied {
		return false, lang.Error(lang.ErrNoPermission)
	}

	return false, lang.Error(lang.ErrPolicyNotFound)
}

func (u *User) Simple() model.Map {
	if u == nil {
		return model.Map{}
	}
	return model.Map{
		"id":     u.id,
		"enable": u.IsEnabled(),
		"name":   u.name,
	}
}

func (u *User) Brief() model.Map {
	if u == nil {
		return model.Map{}
	}
	return model.Map{
		"id":         u.id,
		"enable":     u.IsEnabled(),
		"name":       u.name,
		"title":      u.title,
		"mobile":     u.mobile,
		"email":      u.email,
		"created_at": u.createdAt,
	}
}

func (u *User) Detail() model.Map {
	if u == nil {
		return model.Map{}
	}
	detail := model.Map{
		"id":         u.id,
		"enable":     u.IsEnabled(),
		"name":       u.name,
		"title":      u.title,
		"mobile":     u.mobile,
		"email":      u.email,
		"created_at": u.createdAt,
	}
	rolesData := make(map[string]model.Map, 0)
	roles, _ := u.GetRoles()
	for _, role := range roles {
		if role.Name() != u.Name() {
			rolesData[role.Name()] = model.Map{
				"id":     role.GetID(),
				"enable": role.IsEnabled(),
				"title":  role.Title(),
			}
		}
	}
	detail["roles"] = rolesData
	return detail
}
