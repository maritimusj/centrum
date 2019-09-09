package mysqlStore

import (
	"github.com/maritimusj/centrum/model"
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

	store *mysqlStore
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
