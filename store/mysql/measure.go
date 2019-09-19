package mysqlStore

import (
	"github.com/maritimusj/centrum/dirty"
	"github.com/maritimusj/centrum/helper"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/status"
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

	dirty *dirty.Dirty
	store *mysqlStore
}

func NewMeasure(s *mysqlStore, id int64) *Measure {
	return &Measure{
		id:    id,
		dirty: dirty.New(),
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

func (m *Measure) ResourceClass() resource.Class {
	return resource.Measure
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

func (m *Measure) GetChildrenResources(options ...helper.OptionFN)([]model.Resource, int64, error) {
	panic("implement me")
}

func (m *Measure) GetID() int64 {
	return m.id
}

func (m *Measure) Enable() {
	if m.enable != status.Enable {
		m.enable = status.Enable
		m.dirty.Set("enable", func() interface{} {
			return m.enable
		})
	}
}

func (m *Measure) Disable() {
	if m.enable != status.Disable {
		m.enable = status.Disable
		m.dirty.Set("enable", func() interface{} {
			return m.enable
		})
	}
}

func (m *Measure) IsEnabled() bool {
	return m.enable == status.Enable
}

func (m *Measure) Device() model.Device {
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

func (m *Measure) Tag() string {
	return m.tag
}

func (m *Measure) Kind() resource.MeasureKind {
	return resource.MeasureKind(m.kind)
}

func (m *Measure) CreatedAt() time.Time {
	return m.createdAt
}

func (m *Measure) Destroy() error {
	return m.store.RemoveMeasure(m.id)
}

func (m *Measure) Save() error {
	if m.dirty.Any() {
		err := SaveData(m.store.db, TbMeasures, m.dirty.Data(true), "id=?", m.id)
		if err != nil {
			return lang.InternalError(err)
		}
	}
	return nil
}

func (m *Measure) Simple() model.Map {
	if m == nil {
		return model.Map{}
	}
	return model.Map{
		"id":     m.id,
		"enable": m.IsEnabled(),
		"kind":   m.kind,
		"title":  m.title,
	}
}

func (m *Measure) Brief() model.Map {
	if m == nil {
		return model.Map{}
	}
	return model.Map{
		"id":         m.id,
		"enable":     m.IsEnabled(),
		"kind":       m.kind,
		"title":      m.title,
		"tag":        m.tag,
		"created_at": m.createdAt,
	}
}

func (m *Measure) Detail() model.Map {
	if m == nil {
		return model.Map{}
	}
	detail := model.Map{
		"id":         m.id,
		"enable":     m.IsEnabled(),
		"kind":       m.kind,
		"title":      m.title,
		"tag":        m.tag,
		"created_at": m.createdAt,
	}

	device := m.Device()
	if device != nil {
		detail["device"] = device.Simple()
	}
	return detail
}
