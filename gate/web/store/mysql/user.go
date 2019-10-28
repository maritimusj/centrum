package mysqlStore

import (
	"errors"
	"fmt"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	dirty2 "github.com/maritimusj/centrum/gate/web/dirty"
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	status2 "github.com/maritimusj/centrum/gate/web/status"
	"github.com/maritimusj/centrum/util"
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

	dirty *dirty2.Dirty
	store *mysqlStore
}

func NewUser(s *mysqlStore, id int64) *User {
	return &User{
		id:    id,
		dirty: dirty2.New(),
		store: s,
	}
}

func (u *User) UID() string {
	if u != nil {
		return fmt.Sprintf("user:%d", u.id)
	}
	return "user:<unknown>"
}

func (u *User) Logger() *log.Entry {
	return log.WithFields(log.Fields{
		"org": u.OrganizationID(),
		"src": u.UID(),
	})
}

func (u *User) OrganizationID() int64 {
	if u != nil {
		return u.orgID
	}
	return 0
}

func (u *User) Organization() (model2.Organization, error) {
	if u != nil {
		return u.store.GetOrganization(u.orgID)
	}
	return nil, lang2.Error(lang2.ErrUserNotFound)
}

func (u *User) GetID() int64 {
	if u != nil {
		return u.id
	}
	return 0
}

func (u *User) Enable() {
	if u != nil && u.enable != status2.Enable {
		u.enable = status2.Enable
		u.dirty.Set("enable", func() interface{} {
			return u.enable
		})
	}
}

func (u *User) Disable() {
	if u != nil && u.enable != status2.Disable {
		u.enable = status2.Disable
		u.dirty.Set("enable", func() interface{} {
			return u.enable
		})
	}
}

func (u *User) IsEnabled() bool {
	if u != nil {
		return u.enable == status2.Enable
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

func (u *User) Update(profile model2.Map) {
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
		return lang2.Error(lang2.ErrUserNotFound)
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
			return lang2.InternalError(err)
		}
	}
	return nil
}

func (u *User) GetRoles() ([]model2.Role, error) {
	if u == nil {
		return nil, lang2.Error(lang2.ErrUserNotFound)
	}

	roles, _, err := u.store.GetRoleList(helper2.User(u.id))
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (u *User) Is(role interface{}) (bool, error) {
	if u == nil {
		return false, lang2.Error(lang2.ErrUserNotFound)
	}

	roles, err := u.GetRoles()
	if err != nil {
		return false, err
	}
	var fn func(role model2.Role) bool
	switch v := role.(type) {
	case int64:
		fn = func(role model2.Role) bool {
			return role.GetID() == v
		}
	case string:
		fn = func(role model2.Role) bool {
			return role.Name() == v
		}
	case model2.Role:
		fn = func(role model2.Role) bool {
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

	return lang2.Error(lang2.ErrUserNotFound)
}

func (u *User) Save() error {
	if u == nil {
		return lang2.Error(lang2.ErrUserNotFound)
	}

	if u.dirty.Any() {
		err := SaveData(u.store.db, TbUsers, u.dirty.Data(true), "id=?", u.id)
		if err != nil {
			return lang2.InternalError(err)
		}
	}
	return nil
}

func (u *User) RemovePolicies(res model2.Resource) error {
	if u == nil {
		return lang2.Error(lang2.ErrUserNotFound)
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

func (u *User) SetDeny(res model2.Resource, actions ...resource2.Action) error {
	if u == nil {
		return lang2.Error(lang2.ErrUserNotFound)
	}

	role, err := u.store.GetRole(u.Name())
	if err != nil {
		return err
	}

	for _, action := range actions {
		_, err = role.SetPolicy(res, action, resource2.Deny, make(map[model2.Resource]struct{}))
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *User) SetAllow(res model2.Resource, actions ...resource2.Action) error {
	if u == nil {
		return lang2.Error(lang2.ErrUserNotFound)
	}

	role, err := u.store.GetRole(u.Name())
	if err != nil {
		return err
	}

	for _, action := range actions {
		_, err = role.SetPolicy(res, action, resource2.Allow, make(map[model2.Resource]struct{}))
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *User) IsAllow(res model2.Resource, action resource2.Action) (bool, error) {
	if u == nil {
		return false, lang2.Error(lang2.ErrUserNotFound)
	}

	if res.OrganizationID() > 0 && res.OrganizationID() != u.OrganizationID() {
		return false, lang2.Error(lang2.ErrOrganizationDifferent)
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
			if err == lang2.Error(lang2.ErrNoPermission) {
				denied = true
			}
		}
	}

	if denied {
		return false, lang2.Error(lang2.ErrNoPermission)
	}

	return false, lang2.Error(lang2.ErrPolicyNotFound)
}

func (u *User) Simple() model2.Map {
	if u == nil {
		return model2.Map{}
	}
	return model2.Map{
		"id":     u.id,
		"enable": u.IsEnabled(),
		"name":   u.name,
	}
}

func (u *User) Brief() model2.Map {
	if u == nil {
		return model2.Map{}
	}
	return model2.Map{
		"id":         u.id,
		"enable":     u.IsEnabled(),
		"name":       u.name,
		"title":      u.title,
		"mobile":     u.mobile,
		"email":      u.email,
		"created_at": u.createdAt,
	}
}

func (u *User) Detail() model2.Map {
	if u == nil {
		return model2.Map{}
	}
	detail := model2.Map{
		"id":         u.id,
		"enable":     u.IsEnabled(),
		"name":       u.name,
		"title":      u.title,
		"mobile":     u.mobile,
		"email":      u.email,
		"created_at": u.createdAt,
	}
	rolesData := make(map[string]model2.Map, 0)
	roles, _ := u.GetRoles()
	for _, role := range roles {
		if role.Name() != u.Name() {
			rolesData[role.Name()] = model2.Map{
				"id":     role.GetID(),
				"enable": role.IsEnabled(),
				"title":  role.Title(),
			}
		}
	}
	detail["roles"] = rolesData
	return detail
}
