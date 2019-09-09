package mysqlStore

import (
	"github.com/maritimusj/centrum/model"
	"time"
)

type Role struct {
	id        int64
	enable    int8
	title     string
	createdAt time.Time

	store *mysqlStore
}

func (r *Role) GetID() int64 {
	panic("implement me")
}

func (r *Role) CreatedAt() time.Time {
	panic("implement me")
}

func (r *Role) Save() error {
	panic("implement me")
}

func (r *Role) Destroy() error {
	panic("implement me")
}

func (r *Role) Enable() error {
	panic("implement me")
}

func (r *Role) Disable() error {
	panic("implement me")
}

func (r *Role) IsEnabled() bool {
	panic("implement me")
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

func (r *Role) Title() string {
	panic("implement me")
}

func (r *Role) SetTitle(title string) error {
	panic("implement me")
}

func (r *Role) SetPolicy(resource interface{}, action model.Action, effect model.Effect) (model.Policy, error) {
	panic("implement me")
}

func (r *Role) GetPolicy(resource interface{}) (map[model.Action]model.Policy, error) {
	panic("implement me")
}

func (r *Role) IsAllowed(request model.Request) error {
	panic("implement me")
}
