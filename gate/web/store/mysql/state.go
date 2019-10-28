package mysqlStore

import (
	"errors"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	dirty2 "github.com/maritimusj/centrum/gate/web/dirty"
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	status2 "github.com/maritimusj/centrum/gate/web/status"
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

	dirty *dirty2.Dirty
	store *mysqlStore
}

func NewState(s *mysqlStore, id int64) *State {
	return &State{
		id:    id,
		dirty: dirty2.New(),
		store: s,
	}
}

func (s *State) OrganizationID() int64 {
	equipment := s.Equipment()
	if equipment != nil {
		return equipment.OrganizationID()
	}
	return 0
}

func (s *State) ResourceClass() resource2.Class {
	return resource2.State
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

func (s *State) GetChildrenResources(options ...helper2.OptionFN) ([]model2.Resource, int64, error) {
	return []model2.Resource{}, 0, nil
}

func (s *State) GetID() int64 {
	return s.id
}

func (s *State) CreatedAt() time.Time {
	return s.createdAt
}

func (s *State) Save() error {
	if s.dirty.Any() {
		err := SaveData(s.store.db, TbStates, s.dirty.Data(true), "id=?", s.id)
		if err != nil {
			return lang2.InternalError(err)
		}
	}
	return nil
}

func (s *State) Destroy() error {
	policies, _, err := s.store.GetPolicyList(s)
	if err != nil {
		return err
	}

	for _, policy := range policies {
		if err = policy.Destroy(); err != nil {
			return err
		}
	}
	return s.store.RemoveState(s.id)
}

func (s *State) Enable() {
	if s.enable != status2.Enable {
		s.enable = status2.Enable
		s.dirty.Set("enable", func() interface{} {
			return s.enable
		})
	}
}

func (s *State) Disable() {
	if s.enable != status2.Disable {
		s.enable = status2.Disable
		s.dirty.Set("enable", func() interface{} {
			return s.enable
		})
	}
}

func (s *State) IsEnabled() bool {
	return s.enable == status2.Enable
}

func (s *State) Title() string {
	return s.title
}

func (s *State) SetTitle(title string) {
	if s.title != title {
		s.title = title
		s.dirty.Set("title", func() interface{} {
			return s.title
		})
	}
}

func (s *State) Desc() string {
	return s.desc
}

func (s *State) SetDesc(desc string) {
	if s.desc != desc {
		s.desc = desc
		s.dirty.Set("desc", func() interface{} {
			return s.desc
		})
	}
}

func (s *State) Script() string {
	return s.script
}

func (s *State) SetScript(script string) {
	if s.script != script {
		s.script = script
		s.dirty.Set("script", func() interface{} {
			return s.script
		})
	}
}

func (s *State) Measure() model2.Measure {
	if s.measureID > 0 {
		measure, _ := s.store.GetMeasure(s.measureID)
		if measure != nil {
			return measure
		}
	}
	return nil
}

func (s *State) Equipment() model2.Equipment {
	if s.equipmentID > 0 {
		equipment, _ := s.store.GetEquipment(s.equipmentID)
		if equipment != nil {
			return equipment
		}
	}
	return nil
}

func (s *State) SetMeasure(measure interface{}) {
	var measureID int64
	switch v := measure.(type) {
	case int64:
		measureID = v
	case model2.Measure:
		measureID = v.GetID()
	default:
		panic(errors.New("state SetMeasure: unknown measure"))
	}

	if measureID != s.measureID {
		s.measureID = measureID
		s.dirty.Set("measure_id", func() interface{} {
			return s.measureID
		})
	}
}

func (s *State) Simple() model2.Map {
	if s == nil {
		return model2.Map{}
	}
	return model2.Map{
		"id":     s.id,
		"enable": s.IsEnabled(),
		"title":  s.title,
	}
}

func (s *State) Brief() model2.Map {
	if s == nil {
		return model2.Map{}
	}
	return model2.Map{
		"id":         s.id,
		"enable":     s.IsEnabled(),
		"title":      s.title,
		"desc":       s.desc,
		"created_at": s.createdAt,
	}
}

func (s *State) Detail() model2.Map {
	if s == nil {
		return model2.Map{}
	}

	detail := model2.Map{
		"id":         s.id,
		"enable":     s.IsEnabled(),
		"title":      s.title,
		"desc":       s.desc,
		"script":     s.script,
		"created_at": s.createdAt,
	}

	measure := s.Measure()
	if measure != nil {
		detail["measure"] = measure.Detail()
	}

	return detail
}
