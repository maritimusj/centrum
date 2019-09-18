package mysqlStore

import (
	"database/sql"
	"errors"
	"github.com/maritimusj/centrum/dirty"
	"github.com/maritimusj/centrum/helper"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"time"
)

type Group struct {
	id    int64
	orgID int64

	parentID  int64
	title     string
	desc      string
	createdAt time.Time

	dirty *dirty.Dirty
	store *mysqlStore
}

func NewGroup(s *mysqlStore, id int64) *Group {
	return &Group{
		id:    id,
		dirty: dirty.New(),
		store: s,
	}
}

func (g *Group) ResourceClass() resource.Class {
	return resource.Group
}

func (g *Group) ResourceID() int64 {
	return g.id
}

func (g *Group) ResourceTitle() string {
	return g.title
}

func (g *Group) ResourceDesc() string {
	return g.desc
}

func (g *Group) OrganizationID() int64 {
	return g.orgID
}

func (g *Group) Organization() (model.Organization, error) {
	return g.store.GetOrganization(g.orgID)
}

func (g *Group) GetID() int64 {
	return g.id
}

func (g *Group) CreatedAt() time.Time {
	return g.createdAt
}

func (g *Group) Save() error {
	if g.dirty.Any() {
		err := SaveData(g.store.db, TbGroups, g.dirty.Data(true), "id=?", g.id)
		if err != nil {
			return lang.InternalError(err)
		}
	}
	return nil
}

func (g *Group) Destroy() error {
	return g.store.RemoveGroup(g.id)
}

func (g *Group) Parent() model.Group {
	if g.parentID > 0 {
		group, err := g.store.GetGroup(g.parentID)
		if err == nil {
			return group
		}
	}
	return nil
}

func (g *Group) Title() string {
	return g.title
}

func (g *Group) SetTitle(title string) {
	if g.title != title {
		g.title = title
		g.dirty.Set("title", func() interface{} {
			return g.title
		})
	}
}

func (g *Group) Desc() string {
	return g.desc
}

func (g *Group) SetDesc(desc string) {
	if g.desc != desc {
		g.desc = desc
		g.dirty.Set("desc", func() interface{} {
			return g.desc
		})
	}
}

func (g *Group) SetParent(parent interface{}) {
	var parentID int64
	switch v := parent.(type) {
	case int64:
		parentID = v
	case model.Group:
		parentID = v.GetID()
	default:
		panic(errors.New("group SetParent: unknown group"))
	}

	if g.parentID != parentID {
		g.parentID = parentID
		g.dirty.Set("parent_id", func() interface{} {
			return g.parentID
		})
	}
}

func (g *Group) AddDevice(devices ...interface{}) error {
	now := time.Now()

	for _, device := range devices {
		var deviceID int64
		switch v := device.(type) {
		case int64:
			deviceID = v
		case model.Device:
			deviceID = v.GetID()
		default:
			panic(errors.New("AddDevice: unknown device"))
		}

		_, err := g.store.GetDevice(deviceID)
		if err != nil {
			return err
		}

		if exists, err := IsDataExists(g.store.db, TbDeviceGroups, "group_id=? AND device_id=?", g.id, deviceID); err != nil {
			return lang.InternalError(err)
		} else if exists {
			continue
		}

		_, err = CreateData(g.store.db, TbDeviceGroups, map[string]interface{}{
			"group_id":   g.id,
			"device_id":  deviceID,
			"created_at": now,
		})
		if err != nil {
			return lang.InternalError(err)
		}
	}
	return nil
}

func (g *Group) RemoveDevice(devices ...interface{}) error {
	for _, device := range devices {
		var deviceID int64
		switch v := device.(type) {
		case int64:
			deviceID = v
		case model.Device:
			deviceID = v.GetID()
		default:
			panic(errors.New("RemoveDevice: unknown device"))
		}

		_, err := g.store.GetDevice(deviceID)
		if err != nil {
			return err
		}

		if err := RemoveData(g.store.db, TbDeviceGroups, "group_id=? AND device_id=?", g.id, deviceID); err != nil {
			if err != sql.ErrNoRows {
				return lang.InternalError(err)
			}
		}
	}
	return nil
}

func (g *Group) GetDeviceList(options ...helper.OptionFN) ([]model.Device, int64, error) {
	return g.store.GetDeviceList(options...)
}

func (g *Group) AddEquipment(equipments ...interface{}) error {
	now := time.Now()

	for _, equipment := range equipments {
		var equipmentID int64
		switch v := equipment.(type) {
		case int64:
			equipmentID = v
		case model.Equipment:
			equipmentID = v.GetID()
		default:
			panic(errors.New("AddEquipment: unknown equipment"))
		}

		_, err := g.store.GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		if exists, err := IsDataExists(g.store.db, TbEquipmentGroups, "group_id=? AND equipment_id=?", g.id, equipmentID); err != nil {
			return lang.InternalError(err)
		} else if exists {
			continue
		}

		_, err = CreateData(g.store.db, TbEquipmentGroups, map[string]interface{}{
			"group_id":     g.id,
			"equipment_id": equipmentID,
			"created_at":   now,
		})
		if err != nil {
			return lang.InternalError(err)
		}
	}
	return nil
}

func (g *Group) RemoveEquipment(equipments ...interface{}) error {
	for _, equipment := range equipments {
		var equipmentID int64
		switch v := equipment.(type) {
		case int64:
			equipmentID = v
		case model.Equipment:
			equipmentID = v.GetID()
		default:
			panic(errors.New("RemoveEquipment: unknown equipment"))
		}

		_, err := g.store.GetEquipment(equipmentID)
		if err != nil {
			return err
		}

		if err := RemoveData(g.store.db, TbDeviceGroups, "group_id=? AND equipment_id=?", g.id, equipmentID); err != nil {
			if err != sql.ErrNoRows {
				return lang.InternalError(err)
			}
		}
	}
	return nil
}

func (g *Group) GetEquipmentList(options ...helper.OptionFN) ([]model.Equipment, int64, error) {
	return g.store.GetEquipmentList(options...)
}

func (g *Group) Simple() model.Map {
	if g == nil {
		return model.Map{}
	}
	return model.Map{
		"id":    g.id,
		"title": g.title,
	}
}

func (g *Group) Brief() model.Map {
	if g == nil {
		return model.Map{}
	}
	return model.Map{
		"id":         g.id,
		"title":      g.title,
		"desc":       g.desc,
		"created_at": g.createdAt,
	}
}

func (g *Group) Detail() model.Map {
	if g == nil {
		return model.Map{}
	}
	detail := model.Map{
		"id":         g.id,
		"title":      g.title,
		"desc":       g.desc,
		"created_at": g.createdAt,
	}
	parent := g.Parent()
	if parent != nil {
		detail["parent"] = parent.Simple()
	}
	return detail
}
