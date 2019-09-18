package mysqlStore

import (
	"errors"
	"fmt"
	"github.com/maritimusj/centrum/dirty"
	"github.com/maritimusj/centrum/helper"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/status"
	"github.com/maritimusj/centrum/util"

	log "github.com/sirupsen/logrus"

	"time"
)

type User struct {
	id     int64
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
	return fmt.Sprintf("user:%d", u.id)
}

func (u *User) Logger() *log.Entry {
	return log.WithField("src", u.LogUID())
}

func (u *User) GetID() int64 {
	return u.id
}

func (u *User) Enable() {
	if u.enable != status.Enable {
		u.enable = status.Enable
		u.dirty.Set("enable", func() interface{} {
			return u.enable
		})
	}
}

func (u *User) Disable() {
	if u.enable != status.Disable {
		u.enable = status.Disable
		u.dirty.Set("enable", func() interface{} {
			return u.enable
		})
	}
}

func (u *User) IsEnabled() bool {
	return u.enable == status.Enable
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Title() string {
	return u.title
}

func (u *User) Mobile() string {
	return u.mobile
}

func (u *User) Email() string {
	return u.email
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) ResetPassword(password string) {
	data, _ := util.HashPassword([]byte(password))

	u.password = data
	u.dirty.Set("password", func() interface{} {
		return u.password
	})

}

func (u *User) CheckPassword(password string) bool {
	return util.ComparePassword(u.password, []byte(password))
}

func (u *User) Update(profile model.Map) {
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

func (u *User) SetRoles(roles ...interface{}) error {
	err := RemoveData(u.store.db, TbUserRoles, "user_id=?", u.id)
	if err != nil {
		return err
	}

	now := time.Now()

	for _, role := range roles {
		var roleID int64
		switch v := role.(type) {
		case int64:
			roleID = v
		case model.Role:
			roleID = v.GetID()
		default:
			panic(errors.New("SetRoles: unknown roles"))
		}
		_, err := u.store.GetRole(roleID)
		if err != nil {
			return err
		}
		_, err = CreateData(u.store.db, TbUserRoles, map[string]interface{}{
			"user_id":    u.id,
			"role_id":    roleID,
			"created_at": now,
		})
		if err != nil {
			return lang.InternalError(err)
		}
	}
	return nil
}

func (u *User) GetRoles() ([]model.Role, error) {
	roles, _, err := u.store.GetRoleList(helper.User(u.id))
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (u *User) Destroy() error {
	return u.store.RemoveUser(u.id)
}

func (u *User) Save() error {
	if u.dirty.Any() {
		err := SaveData(u.store.db, TbUsers, u.dirty.Data(true), "id=?", u.id)
		if err != nil {
			return lang.InternalError(err)
		}
	}
	return nil
}

func (u *User) IsAllow(res resource.Resource, action resource.Action) (bool, error) {
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
