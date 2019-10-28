package mysqlStore

import (
	lang2 "github.com/maritimusj/centrum/gate/lang"
	dirty2 "github.com/maritimusj/centrum/gate/web/dirty"
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	status2 "github.com/maritimusj/centrum/gate/web/status"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"time"
)

type Organization struct {
	id     int64
	enable int8

	name      string
	title     string
	extra     []byte
	createdAt time.Time

	dirty *dirty2.Dirty
	store *mysqlStore
}

func NewOrganization(store *mysqlStore, id int64) *Organization {
	return &Organization{
		id:    id,
		dirty: dirty2.New(),
		store: store,
	}
}

func (o *Organization) GetID() int64 {
	return o.id
}

func (o *Organization) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Organization) Enable() {
	if o.enable != status2.Enable {
		o.enable = status2.Enable
		o.dirty.Set("enable", func() interface{} {
			return o.enable
		})
	}
}

func (o *Organization) Disable() {
	if o.enable != status2.Disable {
		o.enable = status2.Disable
		o.dirty.Set("enable", func() interface{} {
			return o.enable
		})
	}
}

func (o *Organization) IsEnabled() bool {
	return o.enable == status2.Enable
}

func (o *Organization) Option() map[string]interface{} {
	return gjson.ParseBytes(o.extra).Value().(map[string]interface{})
}

func (o *Organization) GetOption(path string) gjson.Result {
	return gjson.GetBytes(o.extra, path)
}

func (o *Organization) SetOption(path string, value interface{}) error {
	data, err := sjson.SetBytes(o.extra, path, value)
	if err != nil {
		return err
	}

	o.extra = data
	o.dirty.Set("extra", func() interface{} {
		return o.extra
	})

	return nil
}

func (o *Organization) Name() string {
	return o.name
}

func (o *Organization) Title() string {
	return o.title
}

func (o *Organization) SetTitle(title string) {
	if o.title != title {
		o.title = title
		o.dirty.Set("title", func() interface{} {
			return o.title
		})
	}
}

func (o *Organization) Save() error {
	if o.dirty.Any() {
		err := SaveData(o.store.db, TbOrganization, o.dirty.Data(true), "id=?", o.id)
		if err != nil {
			return lang2.InternalError(err)
		}
	}
	return nil
}

func (o *Organization) Destroy() error {
	users, _, err := o.store.GetUserList(helper2.Organization(o.GetID()))
	if err != nil {
		return err
	}
	for _, user := range users {
		if err = user.Destroy(); err != nil {
			return err
		}
	}

	devices, _, err := o.store.GetDeviceList(helper2.Organization(o.GetID()))
	if err != nil {
		return err
	}
	for _, device := range devices {
		if err = device.Destroy(); err != nil {
			return err
		}
	}

	equipments, _, err := o.store.GetEquipmentList(helper2.Organization(o.GetID()))
	if err != nil {
		return err
	}
	for _, equipment := range equipments {
		if err = equipment.Destroy(); err != nil {
			return err
		}
	}

	groups, _, err := o.store.GetGroupList(helper2.Organization(o.GetID()))
	if err != nil {
		return err
	}

	for _, group := range groups {
		if err = group.Destroy(); err != nil {
			return err
		}
	}

	return o.store.RemoveOrganization(o.id)
}

func (o *Organization) Simple() model2.Map {
	if o == nil {
		return model2.Map{}
	}
	return model2.Map{
		"id":    o.id,
		"name":  o.name,
		"title": o.title,
	}
}

func (o *Organization) Brief() model2.Map {
	if o == nil {
		return model2.Map{}
	}
	return model2.Map{
		"id":         o.id,
		"name":       o.name,
		"title":      o.title,
		"enable":     o.IsEnabled(),
		"created_at": o.createdAt,
	}
}

func (o *Organization) Detail() model2.Map {
	if o == nil {
		return model2.Map{}
	}
	return model2.Map{
		"id":         o.id,
		"name":       o.name,
		"title":      o.title,
		"enable":     o.IsEnabled(),
		"created_at": o.createdAt,
	}
}
