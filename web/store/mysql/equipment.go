package mysqlStore

import (
	"errors"
	"fmt"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/web/dirty"
	"github.com/maritimusj/centrum/web/helper"
	"github.com/maritimusj/centrum/web/model"
	"github.com/maritimusj/centrum/web/resource"
	"github.com/maritimusj/centrum/web/status"
	log "github.com/sirupsen/logrus"
	"time"
)

type Equipment struct {
	id    int64
	orgID int64

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

func (e *Equipment) OrganizationID() int64 {
	return e.orgID
}

func (e *Equipment) Organization() (model.Organization, error) {
	return e.store.GetOrganization(e.orgID)
}

func (e *Equipment) LogUID() string {
	return fmt.Sprintf("equipment:%d", e.id)
}

func (e *Equipment) Logger() *log.Entry {
	return log.WithFields(log.Fields{
		"org": e.OrganizationID(),
		"src": e.LogUID(),
	})
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

func (e *Equipment) GetChildrenResources(options ...helper.OptionFN) ([]model.Resource, int64, error) {
	options = append(options, helper.Equipment(e.GetID()))

	states, total, err := e.store.GetStateList(options...)
	if err != nil {
		return nil, 0, err
	}

	result := make([]model.Resource, 0, len(states))
	for _, state := range states {
		result = append(result, state)
	}

	return result, total, nil
}

func (e *Equipment) GetID() int64 {
	return e.id
}

func (e *Equipment) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Equipment) Save() error {
	if e != nil {
		if e.dirty.Any() {
			err := SaveData(e.store.db, TbEquipments, e.dirty.Data(true), "id=?", e.id)
			if err != nil {
				return lang.InternalError(err)
			}
		}
		return nil
	}
	return lang.Error(lang.ErrEquipmentNotFound)
}

func (e *Equipment) Destroy() error {
	states, _, err := e.store.GetStateList(helper.Equipment(e.GetID()))
	if err != nil {
		return err
	}

	for _, state := range states {
		if err = state.Destroy(); err != nil {
			return err
		}
	}

	err = e.SetGroups(nil)
	if err != nil {
		return err
	}

	policies, _, err := e.store.GetPolicyList(e)
	if err != nil {
		return err
	}

	for _, policy := range policies {
		if err = policy.Destroy(); err != nil {
			return err
		}
	}

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
	err := RemoveData(e.store.db, TbEquipmentGroups, "equipment_id=?", e.id)
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
		_, err = CreateData(e.store.db, TbEquipmentGroups, map[string]interface{}{
			"equipment_id": e.id,
			"group_id":     groupID,
			"created_at":   now,
		})
		if err != nil {
			return lang.InternalError(err)
		}
	}
	return nil
}

func (e *Equipment) Groups() ([]model.Group, error) {
	groups, err := e.store.GetEquipmentGroups(e.GetID())
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
		"desc":       e.desc,
		"created_at": e.createdAt,
	}
}

func (e *Equipment) Detail() model.Map {
	if e == nil {
		return model.Map{}
	}
	detail := model.Map{
		"id":         e.id,
		"enable":     e.IsEnabled(),
		"title":      e.title,
		"desc":       e.desc,
		"created_at": e.createdAt,
	}

	groups, _ := e.Groups()
	if len(groups) > 0 {
		groupsProfile := make([]model.Map, 0, len(groups))
		for _, g := range groups {
			groupsProfile = append(groupsProfile, g.Simple())
		}
		detail["groups"] = groupsProfile
	}
	return detail
}
