package mysqlStore

import (
	"github.com/maritimusj/centrum/dirty"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"time"
)

type State struct {
	id          int64
	enable      int8
	title       string
	desc        string
	equipmentID int64
	measureID   int64
	script      string
	createdAt   time.Time

	dirty *dirty.Dirty
	store *mysqlStore
}

func NewState(s *mysqlStore, id int64) *State {
	return &State{
		id:    id,
		dirty: dirty.New(),
		store: s,
	}
}

func (s *State) ResourceClass() resource.Class {
	return resource.State
}

func (s *State) ResourceID() int64 {
	return s.id
}

func (s *State) ResourceTitle() string {
	return s.title
}

func (s *State) ResourceDesc() string {
	return s.desc
}

func (s *State) GetID() int64 {
	panic("implement me")
}

func (s *State) CreatedAt() time.Time {
	panic("implement me")
}

func (s *State) Save() error {
	panic("implement me")
}

func (s *State) Destroy() error {
	panic("implement me")
}

func (s *State) Enable() error {
	panic("implement me")
}

func (s *State) Disable() error {
	panic("implement me")
}

func (s *State) IsEnabled() bool {
	panic("implement me")
}

func (s *State) Simple() model.Map {
	panic("implement me")
}

func (s *State) Brief() model.Map {
	panic("implement me")
}

func (s *State) Detail() model.Map {
	panic("implement me")
}

func (s *State) Measure() model.Measure {
	panic("implement me")
}

func (s *State) Equipment() model.Equipment {
	panic("implement me")
}

func (s *State) SetMeasure(measure interface{}) error {
	panic("implement me")
}

func (s *State) Title() string {
	panic("implement me")
}

func (s *State) SetTitle(string) error {
	panic("implement me")
}

func (s *State) Desc() string {
	panic("implement me")
}

func (s *State) SetDesc(desc string) error {
	panic("implement me")
}

func (s *State) Script() string {
	panic("implement me")
}

func (s *State) SetScript(string) error {
	panic("implement me")
}
