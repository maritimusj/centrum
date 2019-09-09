package mysqlStore

import (
	"github.com/maritimusj/centrum/model"
	"time"
)

type Device struct {
	id     int64
	enable int8

	title     string
	options   []byte
	createdAt time.Time

	store *mysqlStore
}

func (d *Device) GetID() int64 {
	panic("implement me")
}

func (d *Device) Destroy() error {
	panic("implement me")
}

func (d *Device) Enable() error {
	panic("implement me")
}

func (d *Device) Disable() error {
	panic("implement me")
}

func (d *Device) IsEnabled() bool {
	panic("implement me")
}

func (d *Device) Simple() model.Map {
	panic("implement me")
}

func (d *Device) Brief() model.Map {
	panic("implement me")
}

func (d *Device) Detail() model.Map {
	panic("implement me")
}

func (d *Device) Title() string {
	panic("implement me")
}

func (d *Device) SetTitle(title string) error {
	panic("implement me")
}

func (d *Device) Option() model.Map {
	panic("implement me")
}

func (d *Device) SetOption(option model.Map) error {
	panic("implement me")
}

func (d *Device) SetGroups(groups ...model.Group) error {
	panic("implement me")
}

func (d *Device) Groups() ([]model.Group, error) {
	panic("implement me")
}

func (d *Device) GetMeasureList(typ model.MeasureKind, page, pageSize int64) ([]model.Measure, int64, error) {
	panic("implement me")
}

func (d *Device) CreateMeasure(title string, tag string, typ model.MeasureKind) (model.Measure, error) {
	panic("implement me")
}

func (d *Device) CreatedAt() time.Time {
	panic("implement me")
}

func (d *Device) Save() error {
	panic("implement me")
}
