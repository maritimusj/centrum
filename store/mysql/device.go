package mysqlStore

import (
	"fmt"
	"github.com/maritimusj/centrum/dirty"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/status"
	"time"
)

type Device struct {
	id     int64
	enable int8

	title     string
	options   []byte
	createdAt time.Time

	resourceUID *string

	dirty *dirty.Dirty
	store *mysqlStore
}

func NewDDevice(store *mysqlStore, id int64) *Device {
	return &Device{
		id:    id,
		dirty: dirty.New(),
		store: store,
	}
}

func (d *Device) ResourceUID() string {
	if d.resourceUID == nil {
		uid := fmt.Sprintf("%d.%d", resource.Device, d.id)
		d.resourceUID = &uid
	}
	return *d.resourceUID
}

func (d *Device) ResourceClass() resource.Class {
	return resource.Device
}

func (d *Device) ResourceTitle() string {
	return d.title
}

func (d *Device) ResourceDesc() string {
	return d.title
}

func (d *Device) GetID() int64 {
	return d.id
}

func (d *Device) Destroy() error {
	return d.store.RemoveDevice(d.id)
}

func (d *Device) Enable() error {
	d.enable = status.Enable
	return nil
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
