package mysqlStore

import (
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/web/dirty"
	"github.com/maritimusj/centrum/web/helper"
	"github.com/maritimusj/centrum/web/model"
	"github.com/maritimusj/centrum/web/status"
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

	dirty *dirty.Dirty
	store *mysqlStore
}

func NewOrganization(store *mysqlStore, id int64) *Organization {
	return &Organization{
		id:    id,
		dirty: dirty.New(),
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
	if o.enable != status.Enable {
		o.enable = status.Enable
		o.dirty.Set("enable", func() interface{} {
			return o.enable
		})
	}
}

func (o *Organization) Disable() {
	if o.enable != status.Disable {
		o.enable = status.Disable
		o.dirty.Set("enable", func() interface{} {
			return o.enable
		})
	}
}

func (o *Organization) IsEnabled() bool {
	return o.enable == status.Enable
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
			return lang.InternalError(err)
		}
	}
	return nil
}

func (o *Organization) Destroy() error {
	users, _, err := o.store.GetUserList(helper.Organization(o.GetID()))
	if err != nil {
		return err
	}
	for _, user := range users {
		if err = user.Destroy(); err != nil {
			return err
		}
	}

	devices, _, err := o.store.GetDeviceList(helper.Organization(o.GetID()))
	if err != nil {
		return err
	}
	for _, device := range devices {
		if err = device.Destroy(); err != nil {
			return err
		}
	}

	equipments, _, err := o.store.GetEquipmentList(helper.Organization(o.GetID()))
	if err != nil {
		return err
	}
	for _, equipment := range equipments {
		if err = equipment.Destroy(); err != nil {
			return err
		}
	}

	groups, _, err := o.store.GetGroupList(helper.Organization(o.GetID()))
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

func (o *Organization) Simple() model.Map {
	if o == nil {
		return model.Map{}
	}
	return model.Map{
		"id":    o.id,
		"name":  o.name,
		"title": o.title,
	}
}

func (o *Organization) Brief() model.Map {
	if o == nil {
		return model.Map{}
	}
	return model.Map{
		"id":         o.id,
		"name":       o.name,
		"title":      o.title,
		"enable":     o.IsEnabled(),
		"created_at": o.createdAt,
	}
}

func (o *Organization) Detail() model.Map {
	if o == nil {
		return model.Map{}
	}
	return model.Map{
		"id":         o.id,
		"name":       o.name,
		"title":      o.title,
		"enable":     o.IsEnabled(),
		"created_at": o.createdAt,
	}
}
