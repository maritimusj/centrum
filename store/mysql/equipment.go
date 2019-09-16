package mysqlStore

import (
	"errors"
	"github.com/maritimusj/centrum/dirty"
	"github.com/maritimusj/centrum/helper"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/status"
	"time"
)

type Equipment struct {
	id        int64
	enable    int8
	title     string
	desc      string
	createdAt time.Time

	dirty *dirty.Dirty
	store *mysqlStore
}

func NewEquipment(s *mysqlStore, id int64) *Equipment {
	return &Equipment{
		id:    id,
		dirty: dirty.New(),
		store: s,
	}
}

func (e *Equipment) ResourceClass() resource.Class {
	return resource.Equipment
}

func (e *Equipment) ResourceID() int64 {
	return e.id
}

func (e *Equipment) ResourceTitle() string {
	return e.title
}

func (e *Equipment) ResourceDesc() string {
	return e.desc
}

func (e *Equipment) GetID() int64 {
	return e.id
}

func (e *Equipment) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Equipment) Save() error {
	panic("implement me")
}

func (e *Equipment) Destroy() error {
	return e.store.RemoveEquipment(e.id)
}

func (e *Equipment) Enable() {
	if e.enable != status.Enable {
		e.enable = status.Enable
		e.dirty.Set("enable", func() interface{} {
			return e.enable
		})
	}
}

func (e *Equipment) Disable() {
	if e.enable != status.Disable {
		e.enable = status.Disable
		e.dirty.Set("enable", func() interface{} {
			return e.enable
		})
	}
}

func (e *Equipment) IsEnabled() bool {
	return e.enable == status.Enable
}

func (e *Equipment) Title() string {
	return e.title
}

func (e *Equipment) SetTitle(title string) {
	if e.title != title {
		e.title = title
		e.dirty.Set("title", func() interface{} {
			return e.title
		})
	}
}

func (e *Equipment) Desc() string {
	return e.desc
}

func (e *Equipment) SetDesc(desc string) {
	if e.desc != desc {
		e.desc = desc
		e.dirty.Set("desc", func() interface{} {
			return e.desc
		})
	}
}

func (e *Equipment) SetGroups(groups ...interface{}) error {
	err := e.store.TransactionDo(func(db helper.DB) interface{} {
		err := RemoveData(db, TbEquipmentGroups, "equipment_id", e.id)
		if err != nil {
			return err
		}
		now := time.Now()
		for _, group := range groups {
			var groupID int64
			switch v := group.(type) {
			case int64:
				groupID = v
			case model.Group:
				groupID = v.GetID()
			default:
				panic(errors.New("equipment SetGroups: unknown groups"))
			}
			_, err := e.store.GetGroup(groupID)
			if err != nil {
				return err
			}
			_, err = CreateData(db, TbEquipmentGroups, map[string]interface{}{
				"equipment_id": e.id,
				"group_id":     groupID,
				"created_at":   now,
			})
			if err != nil {
				return lang.InternalError(err)
			}
		}
		return nil
	})
	if err != nil {
		return err.(error)
	}
	return nil
}

func (e *Equipment) Groups() ([]model.Group, error) {
	groups, _, err := e.store.GetGroupList(helper.Device(e.id))
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (e *Equipment) GetStateList(options ...helper.OptionFN) ([]model.State, int64, error) {
	return e.store.GetStateList(options...)
}

func (e *Equipment) CreateState(title, desc string, measure interface{}, script string) (model.State, error) {
	var measureID int64
	switch v := measure.(type) {
	case int64:
		measureID = v
	case model.Measure:
		measureID = v.GetID()
	default:
		panic(errors.New("equipment CreateState: unknown measure"))
	}
	return e.store.CreateState(e.GetID(), measureID, title, desc, script)
}

func (e *Equipment) Simple() model.Map {
	if e == nil {
		return model.Map{}
	}
	return model.Map{
		"id":     e.id,
		"enable": e.IsEnabled(),
		"title":  e.title,
	}
}

func (e *Equipment) Brief() model.Map {
	if e == nil {
		return model.Map{}
	}
	return model.Map{
		"id":         e.id,
		"enable":     e.IsEnabled(),
		"title":      e.title,
		"created_at": e.createdAt,
	}
}

func (e *Equipment) Detail() model.Map {
	if e == nil {
		return model.Map{}
	}
	return model.Map{
		"id":         e.id,
		"enable":     e.IsEnabled(),
		"title":      e.title,
		"created_at": e.createdAt,
	}
}
