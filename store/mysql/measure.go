package mysqlStore

import (
	"fmt"
	"github.com/maritimusj/centrum/dirty"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
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

	resourceUID *string

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

func (m *Measure) ResourceClass() resource.Class {
	return resource.Measure
}

func (m *Measure) ResourceUID() string {
	if m.resourceUID == nil {
		uid := fmt.Sprintf("%d.%d", resource.Measure, m.id)
		m.resourceUID = &uid
	}
	return *m.resourceUID
}

func (m *Measure) ResourceTitle() string {
	return m.title
}

func (m *Measure) ResourceDesc() string {
	return m.title
}

func (m *Measure) GetID() int64 {
	panic("implement me")
}

func (m *Measure) Disable() error {
	panic("implement me")
}

func (m *Measure) IsEnabled() bool {
	panic("implement me")
}

func (m *Measure) Simple() model.Map {
	panic("implement me")
}

func (m *Measure) Brief() model.Map {
	panic("implement me")
}

func (m *Measure) Detail() model.Map {
	panic("implement me")
}

func (m *Measure) Device() model.Device {
	panic("implement me")
}

func (m *Measure) Title() string {
	panic("implement me")
}

func (m *Measure) SetTitle(title string) error {
	panic("implement me")
}

func (m *Measure) Tag() string {
	panic("implement me")
}

func (m *Measure) Kind() model.MeasureKind {
	panic("implement me")
}

func (m *Measure) CreatedAt() time.Time {
	panic("implement me")
}

func (m *Measure) Destroy() error {
	panic("implement me")
}

func (m *Measure) Save() error {
	panic("implement me")
}

func (m *Measure) Enable() error {
	panic("implement me")
}
