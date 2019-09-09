package mysqlStore

import (
	"github.com/maritimusj/centrum/model"
	"time"
)

type Group struct {
	id        int64
	parentID  int64
	title     string
	createdAt time.Time

	store *mysqlStore
}

func (g *Group) GetID() int64 {
	panic("implement me")
}

func (g *Group) CreatedAt() time.Time {
	panic("implement me")
}

func (g *Group) Save() error {
	panic("implement me")
}

func (g *Group) Destroy() error {
	panic("implement me")
}

func (g *Group) Simple() model.Map {
	panic("implement me")
}

func (g *Group) Brief() model.Map {
	panic("implement me")
}

func (g *Group) Detail() model.Map {
	panic("implement me")
}

func (g *Group) Parent() model.Group {
	panic("implement me")
}

func (g *Group) Title() string {
	panic("implement me")
}

func (g *Group) SetTitle(title string) error {
	panic("implement me")
}

func (g *Group) SetParent(group model.Group) error {
	panic("implement me")
}

func (g *Group) AddDevice(devices ...interface{}) error {
	panic("implement me")
}

func (g *Group) RemoveDevice(devices ...interface{}) error {
	panic("implement me")
}

func (g *Group) GetDeviceList(keyword string, page, PageSize int64) ([]model.Device, int64, error) {
	panic("implement me")
}

func (g *Group) AddEquipment(equipment ...interface{}) error {
	panic("implement me")
}

func (g *Group) RemoveEquipment(equipment ...interface{}) error {
	panic("implement me")
}

func (g *Group) GetEquipmentList(keyword string, page, PageSize int64) ([]model.Equipment, int64, error) {
	panic("implement me")
}
