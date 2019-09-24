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
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"time"
)

type Device struct {
	id    int64
	orgID int64

	enable int8

	title     string
	options   []byte
	createdAt time.Time

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

func (d *Device) OrganizationID() int64 {
	return d.orgID
}

func (d *Device) Organization() (model.Organization, error) {
	return d.store.GetOrganization(d.orgID)
}

func (d *Device) LogUID() string {
	return fmt.Sprintf("device:%d", d.id)
}

func (d *Device) Logger() *log.Entry {
	return log.WithField("src", d.LogUID())
}

func (d *Device) ResourceID() int64 {
	return d.id
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

func (d *Device) GetChildrenResources(options ...helper.OptionFN) ([]model.Resource, int64, error) {
	options = append(options, helper.Device(d.GetID()))
	measures, total, err := d.store.GetMeasureList(options...)
	if err != nil {
		return nil, 0, err
	}

	result := make([]model.Resource, 0, len(measures))
	for _, measure := range measures {
		result = append(result, measure)
	}

	return result, total, nil
}

func (d *Device) GetID() int64 {
	return d.id
}

func (d *Device) Enable() {
	if d.enable != status.Enable {
		d.enable = status.Enable
		d.dirty.Set("enable", func() interface{} {
			return d.enable
		})
	}
}

func (d *Device) Disable() {
	if d.enable != status.Disable {
		d.enable = status.Disable
		d.dirty.Set("enable", func() interface{} {
			return d.enable
		})
	}
}

func (d *Device) IsEnabled() bool {
	return d.enable == status.Enable
}

func (d *Device) Title() string {
	return d.title
}

func (d *Device) SetTitle(title string) {
	if d.title != title {
		d.title = title
		d.dirty.Set("title", func() interface{} {
			return d.title
		})
	}
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
	err := RemoveData(d.store.db, TbDeviceGroups, "device_id", d.id)
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
			panic(errors.New("device SetGroups: unknown groups"))
		}
		_, err := d.store.GetGroup(groupID)
		if err != nil {
			return err
		}
		_, err = CreateData(d.store.db, TbDeviceGroups, map[string]interface{}{
			"device_id":  d.id,
			"group_id":   groupID,
			"created_at": now,
		})
		if err != nil {
			return lang.InternalError(err)
		}
	}
	return nil
}

func (d *Device) Groups() ([]model.Group, error) {
	groups, _, err := d.store.GetGroupList(helper.Device(d.id))
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (d *Device) GetMeasureList(options ...helper.OptionFN) ([]model.Measure, int64, error) {
	return d.store.GetMeasureList(options...)
}

func (d *Device) CreateMeasure(title string, tag string, kind resource.MeasureKind) (model.Measure, error) {
	return d.store.CreateMeasure(d.GetID(), title, tag, kind)
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
	if d == nil {
		return model.Map{}
	}
	return model.Map{
		"id":     d.id,
		"enable": d.IsEnabled(),
		"title":  d.title,
	}
}

func (d *Device) Brief() model.Map {
	if d == nil {
		return model.Map{}
	}
	return model.Map{
		"id":         d.id,
		"enable":     d.IsEnabled(),
		"title":      d.title,
		"created_at": d.createdAt,
	}
}

func (d *Device) Detail() model.Map {
	if d == nil {
		return model.Map{}
	}
	return model.Map{
		"id":         d.id,
		"enable":     d.IsEnabled(),
		"title":      d.title,
		"created_at": d.createdAt,
	}
}