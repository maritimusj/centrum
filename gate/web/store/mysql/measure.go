package mysqlStore

import (
	lang2 "github.com/maritimusj/centrum/gate/lang"
	dirty2 "github.com/maritimusj/centrum/gate/web/dirty"
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	status2 "github.com/maritimusj/centrum/gate/web/status"
	"time"
)

type Measure struct {
	id     int64
	enable int8

	deviceID  int64
	kind      int
	title     string
	tag       string
	createdAt time.Time

	dirty *dirty2.Dirty
	store *mysqlStore
}

func NewMeasure(s *mysqlStore, id int64) *Measure {
	return &Measure{
		id:    id,
		dirty: dirty2.New(),
		store: s,
	}
}

func (m *Measure) OrganizationID() int64 {
	device := m.Device()
	if device != nil {
		return device.OrganizationID()
	}
	return 0
}

func (m *Measure) ResourceClass() resource2.Class {
	return resource2.Measure
}

func (m *Measure) ResourceID() int64 {
	return m.id
}

func (m *Measure) ResourceTitle() string {
	return m.title
}

func (m *Measure) ResourceDesc() string {
	return m.title
}

func (m *Measure) GetChildrenResources(options ...helper2.OptionFN) ([]model2.Resource, int64, error) {
	return []model2.Resource{}, 0, nil
}

func (m *Measure) GetID() int64 {
	return m.id
}

func (m *Measure) Enable() {
	if m.enable != status2.Enable {
		m.enable = status2.Enable
		m.dirty.Set("enable", func() interface{} {
			return m.enable
		})
	}
}

func (m *Measure) Disable() {
	if m.enable != status2.Disable {
		m.enable = status2.Disable
		m.dirty.Set("enable", func() interface{} {
			return m.enable
		})
	}
}

func (m *Measure) IsEnabled() bool {
	return m.enable == status2.Enable
}

func (m *Measure) Device() model2.Device {
	if m.deviceID > 0 {
		device, _ := m.store.GetDevice(m.deviceID)
		return device
	}
	return nil
}

func (m *Measure) Title() string {
	return m.title
}

func (m *Measure) SetTitle(title string) {
	if m.title != title {
		m.title = title
		m.dirty.Set("title", func() interface{} {
			return m.title
		})
	}
}

func (m *Measure) TagName() string {
	return m.tag
}

func (m *Measure) Kind() resource2.MeasureKind {
	return resource2.MeasureKind(m.kind)
}

func (m *Measure) CreatedAt() time.Time {
	return m.createdAt
}

func (m *Measure) Destroy() error {

	alarms, _, err := m.store.GetAlarmList(nil, nil, helper2.Measure(m.GetID()))
	if err != nil {
		return err
	}

	for _, alarm := range alarms {
		err = alarm.Destroy()
		if err != nil {
			return err
		}
	}

	policies, _, err := m.store.GetPolicyList(m)
	if err != nil {
		return err
	}

	for _, policy := range policies {
		if err = policy.Destroy(); err != nil {
			return err
		}
	}
	return m.store.RemoveMeasure(m.id)
}

func (m *Measure) Save() error {
	if m.dirty.Any() {
		err := SaveData(m.store.db, TbMeasures, m.dirty.Data(true), "id=?", m.id)
		if err != nil {
			return lang2.InternalError(err)
		}
	}
	return nil
}

func (m *Measure) Simple() model2.Map {
	if m == nil {
		return model2.Map{}
	}
	return model2.Map{
		"id":     m.id,
		"enable": m.IsEnabled(),
		"kind":   m.kind,
		"title":  m.title,
	}
}

func (m *Measure) Brief() model2.Map {
	if m == nil {
		return model2.Map{}
	}
	return model2.Map{
		"id":         m.id,
		"enable":     m.IsEnabled(),
		"kind":       m.kind,
		"title":      m.title,
		"tag":        m.tag,
		"created_at": m.createdAt,
	}
}

func (m *Measure) Detail() model2.Map {
	if m == nil {
		return model2.Map{}
	}
	detail := model2.Map{
		"id":         m.id,
		"enable":     m.IsEnabled(),
		"kind":       m.kind,
		"title":      m.title,
		"tag":        m.tag,
		"created_at": m.createdAt,
	}

	device := m.Device()
	if device != nil {
		detail["device"] = device.Brief()
	}
	return detail
}
