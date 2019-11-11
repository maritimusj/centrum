package mysqlStore

import (
	"errors"
	"time"

	"github.com/maritimusj/centrum/gate/web/alarm"

	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/dirty"
	"github.com/maritimusj/centrum/gate/web/helper"
	"github.com/maritimusj/centrum/gate/web/model"
	"github.com/maritimusj/centrum/gate/web/resource"
	"github.com/maritimusj/centrum/gate/web/status"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type State struct {
	id          int64
	enable      int8
	title       string
	desc        string
	equipmentID int64
	measureID   int64
	extra       []byte
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

func (s *State) OrganizationID() int64 {
	equipment := s.Equipment()
	if equipment != nil {
		return equipment.OrganizationID()
	}
	return 0
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

func (s *State) GetChildrenResources(options ...helper.OptionFN) ([]model.Resource, int64, error) {
	return []model.Resource{}, 0, nil
}

func (s *State) GetID() int64 {
	return s.id
}

func (s *State) Option() map[string]interface{} {
	return gjson.ParseBytes(s.extra).Value().(map[string]interface{})
}

func (s *State) GetOption(path string) gjson.Result {
	return gjson.GetBytes(s.extra, path)
}

func (s *State) SetOption(path string, value interface{}) error {
	data, err := sjson.SetBytes(s.extra, path, value)
	if err != nil {
		return err
	}

	s.extra = data
	s.dirty.Set("extra", func() interface{} {
		return s.extra
	})

	return nil
}

func (s *State) CreatedAt() time.Time {
	return s.createdAt
}

func (s *State) Save() error {
	if s.dirty.Any() {
		err := SaveData(s.store.db, TbStates, s.dirty.Data(true), "id=?", s.id)
		if err != nil {
			return lang.InternalError(err)
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
	if s.enable != status.Enable {
		s.enable = status.Enable
		s.dirty.Set("enable", func() interface{} {
			return s.enable
		})
	}
}

func (s *State) Disable() {
	if s.enable != status.Disable {
		s.enable = status.Disable
		s.dirty.Set("enable", func() interface{} {
			return s.enable
		})
	}
}

func (s *State) IsEnabled() bool {
	return s.enable == status.Enable
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

func (s *State) Measure() model.Measure {
	if s.measureID > 0 {
		measure, _ := s.store.GetMeasure(s.measureID)
		if measure != nil {
			return measure
		}
	}
	return nil
}

func (s *State) Equipment() model.Equipment {
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
	case model.Measure:
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

func (s *State) IsAlarmEnabled() bool {
	return s.GetOption("__alarm.enable").Bool()
}

func (s *State) EnableAlarm() {
	_ = s.SetOption("__alarm.enable", true)
}

func (s *State) DisableAlarm() {
	_ = s.SetOption("__alarm.enable", false)
}

func (s *State) AlarmDeadBand() float32 {
	return float32(s.GetOption("__alarm.deadband").Float())
}

func (s *State) SetAlarmDeadBand(v float32) {
	_ = s.SetOption("__alarm.deadband", v)
}

func (s *State) AlarmDelaySecond() int {
	return int(s.GetOption("__alarm.delay").Int())
}

func (s *State) SetAlarmDelay(seconds int) {
	_ = s.SetOption("__alarm.delay", seconds)
}

func (s *State) GetAlarmEntries() map[string]float32 {
	data := map[string]float32{}
	for _, name := range []string{alarm.HF, alarm.HH, alarm.HI, alarm.LF, alarm.LL, alarm.LO} {
		if v, ok := s.GetAlarmEntry(name); ok {
			data[name] = v
		}
	}
	return data
}

func (s *State) GetAlarmEntry(name string) (float32, bool) {
	entry := s.GetOption("__alarm." + name)
	if entry.Exists() {
		return float32(entry.Get("val").Float()), entry.Get("enable").Bool()
	}
	return 0, false
}

func (s *State) SetAlarmEntry(name string, value float32) {
	entry := s.GetOption("__alarm." + name)
	data := map[string]interface{}{}
	if entry.Exists() {
		data["enable"] = entry.Get("enable").Bool()
	} else {
		data["enable"] = true
	}
	data["val"] = value
	_ = s.SetOption("__alarm."+name, data)
}

func (s *State) EnableAlarmEntry(name string) {
	_ = s.SetOption("__alarm."+name+".enable", true)
}

func (s *State) DisableAlarmEntry(name string) {
	_ = s.SetOption("__alarm."+name+".enable", false)
}

func (s *State) Simple() model.Map {
	if s == nil {
		return model.Map{}
	}
	return model.Map{
		"id":     s.id,
		"enable": s.IsEnabled(),
		"title":  s.title,
	}
}

func (s *State) Brief() model.Map {
	if s == nil {
		return model.Map{}
	}
	return model.Map{
		"id":         s.id,
		"enable":     s.IsEnabled(),
		"title":      s.title,
		"desc":       s.desc,
		"created_at": s.createdAt,
	}
}

func (s *State) Detail() model.Map {
	if s == nil {
		return model.Map{}
	}

	detail := model.Map{
		"id":     s.id,
		"enable": s.IsEnabled(),
		"title":  s.title,
		"desc":   s.desc,
		"alarm": model.Map{
			"enable":   s.IsAlarmEnabled(),
			"deadband": s.AlarmDeadBand(),
			"delay":    s.AlarmDelaySecond(),
			"entries":  s.GetAlarmEntries(),
		},
		"created_at": s.createdAt,
	}

	measure := s.Measure()
	if measure != nil {
		detail["measure"] = measure.Detail()
	}

	return detail
}
