package mysqlStore

import (
	"github.com/maritimusj/centrum/model"
	"time"
)

type User struct {
	id     int64
	enable int8

	name      string
	title     string
	password  string
	mobile    string
	email     string
	createdAt time.Time

	store *mysqlStore
}

func (u *User) GetID() int64 {
	panic("implement me")
}

func (u *User) Enable() error {
	panic("implement me")
}

func (u *User) Disable() error {
	panic("implement me")
}

func (u *User) IsEnabled() bool {
	panic("implement me")
}

func (u *User) Simple() model.Map {
	panic("implement me")
}

func (u *User) Brief() model.Map {
	panic("implement me")
}

func (u *User) Detail() model.Map {
	panic("implement me")
}

func (u *User) Name() string {
	panic("implement me")
}

func (u *User) Title() string {
	panic("implement me")
}

func (u *User) Mobile() string {
	panic("implement me")
}

func (u *User) Email() string {
	panic("implement me")
}

func (u *User) ResetPassword(password string) error {
	panic("implement me")
}

func (u *User) CheckPassword(password string) error {
	panic("implement me")
}

func (u *User) Update(profile model.Map) error {
	panic("implement me")
}

func (u *User) SetRoles(roles ...interface{}) error {
	panic("implement me")
}

func (u *User) GetRoles() ([]model.Role, error) {
	panic("implement me")
}

func (u *User) CreatedAt() time.Time {
	panic("implement me")
}

func (u *User) Destroy() error {
	panic("implement me")
}

func (u *User) IsAllowed(request model.Request) error {
	panic("implement me")
}

func (u *User) Save() error {
	panic("implement me")
}
