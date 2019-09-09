package mysqlStore

import (
	"github.com/maritimusj/centrum/model"
	"time"
)

type Equipment struct {
	id        int64
	enable    int8
	title     string
	desc      string
	createdAt time.Time

	store *mysqlStore
}

func (e *Equipment) GetID() int64 {
	panic("implement me")
}

func (e *Equipment) CreatedAt() time.Time {
	panic("implement me")
}

func (e *Equipment) Save() error {
	panic("implement me")
}

func (e *Equipment) Destroy() error {
	panic("implement me")
}

func (e *Equipment) Enable() error {
	panic("implement me")
}

func (e *Equipment) Disable() error {
	panic("implement me")
}

func (e *Equipment) IsEnabled() bool {
	panic("implement me")
}

func (e *Equipment) Simple() model.Map {
	panic("implement me")
}

func (e *Equipment) Brief() model.Map {
	panic("implement me")
}

func (e *Equipment) Detail() model.Map {
	panic("implement me")
}

func (e *Equipment) Title() string {
	panic("implement me")
}

func (e *Equipment) SetTitle(title string) error {
	panic("implement me")
}

func (e *Equipment) Desc() string {
	panic("implement me")
}

func (e *Equipment) SetDesc(string) error {
	panic("implement me")
}

func (e *Equipment) SetGroups(groups ...model.Group) error {
	panic("implement me")
}

func (e *Equipment) Groups() ([]model.Group, error) {
	panic("implement me")
}

func (e *Equipment) GetStateList(keyword string, page, pageSize int64) ([]model.State, int64, error) {
	panic("implement me")
}

func (e *Equipment) CreateState(title string, measure interface{}, script string) (model.State, error) {
	panic("implement me")
}
