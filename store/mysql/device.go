package mysqlStore

import (
	"errors"
	"fmt"
	"github.com/maritimusj/centrum/dirty"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/status"
	"github.com/maritimusj/centrum/store"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
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

func (d *Device) Enable() error {
	if d.enable != status.Enable {
		d.enable = status.Enable
		d.dirty.Set("enable", func() interface{} {
			return d.enable
		})
	}
	return d.Save()
}

func (d *Device) Disable() error {
	if d.enable != status.Disable {
		d.enable = status.Disable
		d.dirty.Set("enable", func() interface{} {
			return d.enable
		})
	}
	return d.Save()
}

func (d *Device) IsEnabled() bool {
	return d.enable == status.Enable
}

func (d *Device) Title() string {
	return d.title
}

func (d *Device) SetTitle(title string) error {
	if d.title != title {
		d.title = title
		d.dirty.Set("title", func() interface{} {
			return d.title
		})
	}
	return d.Save()
}

func (d *Device) GetOption(key string) gjson.Result {
	return gjson.GetBytes(d.options, key)
}

func (d *Device) SetOption(key string, value interface{}) error {
	data, err := sjson.SetBytes(d.options, key, value)
	if err != nil {
		return err
	}

	d.options = data
	d.dirty.Set("options", func() interface{} {
		return d.options
	})

	return d.Save()
}

func (d *Device) SetGroups(groups ...interface{}) error {
	err := d.store.TransactionDo(func(db store.DB) interface{} {
		err := RemoveData(db, TbDeviceGroups, "device_id", d.id)
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
				panic(errors.New("SetGroups: unknown groups"))
			}
			_, err := d.store.GetGroup(groupID)
			if err != nil {
				return err
			}
			_, err = CreateData(db, TbDeviceGroups, map[string]interface{}{
				"device_id":  d.id,
				"group_id":   groupID,
				"created_at": now,
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

func (d *Device) Groups() ([]model.Group, error) {
	groups, _, err := d.store.GetGroupList(store.Device(d.id))
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (d *Device) GetMeasureList(typ model.MeasureKind, page, pageSize int64) ([]model.Measure, int64, error) {
	panic("implement me")
}

func (d *Device) CreateMeasure(title string, tag string, typ model.MeasureKind) (model.Measure, error) {
	panic("implement me")
}

func (d *Device) CreatedAt() time.Time {
	return d.createdAt
}

func (d *Device) Destroy() error {
	return d.store.RemoveDevice(d.id)
}

func (d *Device) Save() error {
	if d.dirty.Any() {
		err := SaveData(d.store.db, TbDevices, d.dirty.Data(true), "id=?", d.id)
		if err != nil {
			return lang.InternalError(err)
		}
	}

	return nil
}

func (d *Device) Simple() model.Map {
	return model.Map{
		"id":     d.id,
		"enable": d.enable,
		"name":   d.title,
	}
}

func (d *Device) Brief() model.Map {
	return model.Map{
		"id":         d.id,
		"enable":     d.enable,
		"title":      d.title,
		"created_at": d.createdAt,
	}
}

func (d *Device) Detail() model.Map {
	return model.Map{
		"id":         d.id,
		"enable":     d.enable,
		"title":      d.title,
		"created_at": d.createdAt,
	}
}
