package mysqlStore

import (
	"fmt"
	"github.com/maritimusj/centrum/dirty"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"time"
)

type Equipment struct {
	id        int64
	enable    int8
	title     string
	desc      string
	createdAt time.Time

	resourceUID *string

	dirty *dirty.Dirty
	store *mysqlStore
}

func NewEquipment(s *mysqlStore, id int64) *Equipment {
	return &Equipment{
		id:    id,
		dirty: dirty.New(),
		store: s,
	}
}

func (e *Equipment) ResourceClass() resource.Class {
	return resource.Equipment
}

func (e *Equipment) ResourceUID() string {
	if e.resourceUID == nil {
		uid := fmt.Sprintf("%d.%d", resource.Equipment, e.id)
		e.resourceUID = &uid
	}
	return *e.resourceUID
}

func (e *Equipment) ResourceTitle() string {
	return e.title
}

func (e *Equipment) ResourceDesc() string {
	return e.desc
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
