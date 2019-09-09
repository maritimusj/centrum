package mysqlStore

import (
	"github.com/maritimusj/centrum/model"
	"time"
)

type State struct {
	id          int64
	enable      int8
	title       string
	equipmentID int64
	measureID   int64
	script      string
	createdAt   time.Time

	store *mysqlStore
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

func (s *State) Script() string {
	panic("implement me")
}

func (s *State) SetScript(string) error {
	panic("implement me")
}
