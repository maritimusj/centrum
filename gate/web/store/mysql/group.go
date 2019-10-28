package mysqlStore

import (
	"database/sql"
	"errors"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/dirty"
	"github.com/maritimusj/centrum/gate/web/helper"
	"github.com/maritimusj/centrum/gate/web/model"
	"github.com/maritimusj/centrum/gate/web/resource"
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

func (g *Group) OrganizationID() int64 {
	return g.orgID
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

func (g *Group) GetChildrenResources(options ...helper.OptionFN) ([]model.Resource, int64, error) {
	var result []model.Resource

	groups, groupTotal, err := g.store.GetGroupList(append(options, helper.Parent(g.GetID()))...)
	if err != nil {
		return nil, 0, err
	}

	for _, group := range groups {
		result = append(result, group)
	}

	devices, deviceTotal, err := g.store.GetDeviceList(append(options, helper.Group(g.GetID()))...)
	if err != nil {
		return nil, 0, err
	}
	for _, device := range devices {
		result = append(result, device)
	}
	equipments, equipmentTotal, err := g.store.GetEquipmentList(append(options, helper.Group(g.GetID()))...)
	if err != nil {
		return nil, 0, err
	}
	for _, equipment := range equipments {
		result = append(result, equipment)
	}

	return result, groupTotal + deviceTotal + equipmentTotal, nil
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
	//处理子分组
	groups, _, err := g.store.GetGroupList(helper.Parent(g.GetID()))
	if err != nil {
		return err
	}

	for _, g := range groups {
		g.SetParent(nil)
		if err = g.Save(); err != nil {
			return err
		}
	}

	//处理分组下设备
	if res, _, err := g.GetDeviceList(); err != nil {
		return err
	} else if len(res) > 0 {
		var devices []interface{}
		for _, device := range res {
			devices = append(devices, device)
		}
		err = g.RemoveDevice(devices...)
	}

	//处理分组下自定义设备
	if res, _, err := g.GetEquipmentList(); err != nil {
		return err
	} else if len(res) > 0 {
		var equipments []interface{}
		for _, e := range res {
			equipments = append(equipments, e)
		}
		err = g.RemoveDevice(equipments...)
	}

	policies, _, err := g.store.GetPolicyList(g)
	if err != nil {
		return err
	}

	for _, policy := range policies {
		if err = policy.Destroy(); err != nil {
			return err
		}
	}

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
	case nil:
		parentID = 0
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
	return g.store.GetDeviceList(append(options, helper.Group(g.GetID()))...)
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
	return g.store.GetEquipmentList(append(options, helper.Group(g.GetID()))...)
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
