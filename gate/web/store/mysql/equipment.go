package mysqlStore

import (
	"errors"
	"fmt"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	dirty2 "github.com/maritimusj/centrum/gate/web/dirty"
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	status2 "github.com/maritimusj/centrum/gate/web/status"
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

	dirty *dirty2.Dirty
	store *mysqlStore
}

func NewEquipment(s *mysqlStore, id int64) *Equipment {
	return &Equipment{
		id:    id,
		dirty: dirty2.New(),
		store: s,
	}
}

func (e *Equipment) OrganizationID() int64 {
	return e.orgID
}

func (e *Equipment) Organization() (model2.Organization, error) {
	return e.store.GetOrganization(e.orgID)
}

func (e *Equipment) UID() string {
	return fmt.Sprintf("equipment:%d", e.id)
}

func (e *Equipment) Logger() *log.Entry {
	return log.WithFields(log.Fields{
		"org": e.OrganizationID(),
		"src": e.UID(),
	})
}

func (e *Equipment) ResourceClass() resource2.Class {
	return resource2.Equipment
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

func (e *Equipment) GetChildrenResources(options ...helper2.OptionFN) ([]model2.Resource, int64, error) {
	options = append(options, helper2.Equipment(e.GetID()))

	states, total, err := e.store.GetStateList(options...)
	if err != nil {
		return nil, 0, err
	}

	result := make([]model2.Resource, 0, len(states))
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
				return lang2.InternalError(err)
			}
		}
		return nil
	}
	return lang2.Error(lang2.ErrEquipmentNotFound)
}

func (e *Equipment) Destroy() error {
	states, _, err := e.store.GetStateList(helper2.Equipment(e.GetID()))
	if err != nil {
		return err
	}

	for _, state := range states {
		if err = state.Destroy(); err != nil {
			return err
		}
	}

	err = e.SetGroups()
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
	if e.enable != status2.Enable {
		e.enable = status2.Enable
		e.dirty.Set("enable", func() interface{} {
			return e.enable
		})
	}
}

func (e *Equipment) Disable() {
	if e.enable != status2.Disable {
		e.enable = status2.Disable
		e.dirty.Set("enable", func() interface{} {
			return e.enable
		})
	}
}

func (e *Equipment) IsEnabled() bool {
	return e.enable == status2.Enable
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
		case model2.Group:
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
			return lang2.InternalError(err)
		}
	}
	return nil
}

func (e *Equipment) Groups() ([]model2.Group, error) {
	groups, err := e.store.GetEquipmentGroups(e.GetID())
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (e *Equipment) GetStateList(options ...helper2.OptionFN) ([]model2.State, int64, error) {
	return e.store.GetStateList(options...)
}

func (e *Equipment) CreateState(title, desc string, measure interface{}, script string) (model2.State, error) {
	var measureID int64
	switch v := measure.(type) {
	case int64:
		measureID = v
	case model2.Measure:
		measureID = v.GetID()
	default:
		panic(errors.New("equipment CreateState: unknown measure"))
	}
	return e.store.CreateState(e.GetID(), measureID, title, desc, script)
}

func (e *Equipment) Simple() model2.Map {
	if e == nil {
		return model2.Map{}
	}
	return model2.Map{
		"id":     e.id,
		"enable": e.IsEnabled(),
		"title":  e.title,
	}
}

func (e *Equipment) Brief() model2.Map {
	if e == nil {
		return model2.Map{}
	}
	return model2.Map{
		"id":         e.id,
		"enable":     e.IsEnabled(),
		"title":      e.title,
		"desc":       e.desc,
		"created_at": e.createdAt,
	}
}

func (e *Equipment) Detail() model2.Map {
	if e == nil {
		return model2.Map{}
	}
	detail := model2.Map{
		"id":         e.id,
		"enable":     e.IsEnabled(),
		"title":      e.title,
		"desc":       e.desc,
		"created_at": e.createdAt,
	}

	groups, _ := e.Groups()
	if len(groups) > 0 {
		groupsProfile := make([]model2.Map, 0, len(groups))
		for _, g := range groups {
			groupsProfile = append(groupsProfile, g.Simple())
		}
		detail["groups"] = groupsProfile
	}
	return detail
}
