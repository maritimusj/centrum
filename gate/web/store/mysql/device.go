package mysqlStore

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/dirty"
	"github.com/maritimusj/centrum/gate/web/helper"
	"github.com/maritimusj/centrum/gate/web/model"
	"github.com/maritimusj/centrum/gate/web/resource"
	"github.com/maritimusj/centrum/gate/web/status"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
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
	if d != nil {
		return d.orgID
	}
	return 0
}

func (d *Device) Organization() (model.Organization, error) {
	if d != nil {
		return d.store.GetOrganization(d.orgID)
	}
	return nil, lang.ErrDeviceNotFound.Error()
}

func (d *Device) UID() string {
	if d != nil {
		return fmt.Sprintf("device:%d", d.id)
	}
	return "device:<unknown>"
}

func (d *Device) Logger() *log.Entry {
	return log.WithFields(log.Fields{
		"org": d.OrganizationID(),
		"src": d.UID(),
	})
}

func (d *Device) ResourceID() int64 {
	if d != nil {
		return d.id
	}
	return 0
}

func (d *Device) ResourceClass() resource.Class {
	return resource.Device
}

func (d *Device) ResourceTitle() string {
	if d != nil {
		return d.title
	}
	return "<unknown>"
}

func (d *Device) ResourceDesc() string {
	if d != nil {
		return d.title
	}
	return "<unknown>"
}

func (d *Device) GetChildrenResources(options ...helper.OptionFN) ([]model.Resource, int64, error) {
	if d != nil {
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

	return nil, 0, lang.ErrDeviceNotFound.Error()
}

func (d *Device) GetID() int64 {
	if d != nil {
		return d.id
	}
	return 0
}

func (d *Device) Enable() {
	if d != nil && d.enable != status.Enable {
		d.enable = status.Enable
		d.dirty.Set("enable", func() interface{} {
			return d.enable
		})
	}
}

func (d *Device) Disable() {
	if d != nil && d.enable != status.Disable {
		d.enable = status.Disable
		d.dirty.Set("enable", func() interface{} {
			return d.enable
		})
	}
}

func (d *Device) IsEnabled() bool {
	return d != nil && d.enable == status.Enable
}

func (d *Device) Title() string {
	if d != nil {
		return d.title
	}
	return "<unknown>"
}

func (d *Device) SetTitle(title string) {
	if d != nil && d.title != title {
		d.title = title
		d.dirty.Set("title", func() interface{} {
			return d.title
		})
	}
}

func (d *Device) Option() map[string]interface{} {
	return gjson.ParseBytes(d.options).Value().(map[string]interface{})
}

func (d *Device) GetOption(key string) gjson.Result {
	if d != nil {
		return gjson.GetBytes(d.options, key)
	}
	return gjson.Result{}
}

func (d *Device) SetOption(key string, value interface{}) error {
	if d != nil {
		data, err := sjson.SetBytes(d.options, key, value)
		if err != nil {
			return err
		}

		d.options = data
		d.dirty.Set("options", func() interface{} {
			return d.options
		})

		return nil
	}
	return lang.ErrDeviceNotFound.Error()
}

func (d *Device) SetGroups(groups ...interface{}) error {
	if d != nil {
		err := RemoveData(d.store.db, TbDeviceGroups, "device_id=?", d.id)
		if err != nil {
			if err != sql.ErrNoRows {
				return lang.InternalError(err)
			}
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
	return lang.ErrDeviceNotFound.Error()
}

func (d *Device) Groups() ([]model.Group, error) {
	if d != nil {
		groups, err := d.store.GetDeviceGroups(d.GetID())
		if err != nil {
			return nil, err
		}
		return groups, nil
	}
	return nil, lang.ErrDeviceNotFound.Error()
}

func (d *Device) GetMeasureList(options ...helper.OptionFN) ([]model.Measure, int64, error) {
	if d != nil {
		return d.store.GetMeasureList(options...)
	}
	return nil, 0, lang.ErrDeviceNotFound.Error()
}

func (d *Device) CreateMeasure(title string, tag string, kind resource.MeasureKind) (model.Measure, error) {
	if d != nil {
		return d.store.CreateMeasure(d.GetID(), title, tag, kind)
	}
	return nil, lang.ErrDeviceNotFound.Error()
}

func (d *Device) CreatedAt() time.Time {
	if d != nil {
		return d.createdAt
	}
	return time.Time{}
}

func (d *Device) Destroy() error {
	if d == nil {
		return lang.ErrDeviceNotFound.Error()
	}

	alarms, _, err := d.store.GetAlarmList(nil, nil, helper.Device(d.GetID()))
	if err != nil {
		return err
	}

	for _, alarm := range alarms {
		err = alarm.Destroy()
		if err != nil {
			return err
		}
	}

	measures, _, err := d.store.GetMeasureList(helper.Device(d.GetID()))
	if err != nil {
		return err
	}

	for _, measure := range measures {
		if err = measure.Destroy(); err != nil {
			return err
		}
	}

	err = d.SetGroups()
	if err != nil {
		return err
	}

	policies, _, err := d.store.GetPolicyList(d)
	if err != nil {
		return err
	}

	for _, policy := range policies {
		if err = policy.Destroy(); err != nil {
			return err
		}
	}

	return d.store.RemoveDevice(d.id)
}

func (d *Device) Save() error {
	if d != nil {
		if d.dirty.Any() {
			err := SaveData(d.store.db, TbDevices, d.dirty.Data(true), "id=?", d.id)
			if err != nil {
				return lang.InternalError(err)
			}
		}
		return nil
	}
	return lang.ErrDeviceNotFound.Error()
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
		"id":             d.id,
		"enable":         d.IsEnabled(),
		"title":          d.title,
		"params.connStr": d.GetOption("params.connStr").Str,
		"created_at":     d.createdAt.Format(lang.DatetimeFormatterStr.Str()),
	}
}

func (d *Device) Detail() model.Map {
	if d == nil {
		return model.Map{}
	}

	detail := model.Map{
		"id":              d.id,
		"enable":          d.IsEnabled(),
		"title":           d.title,
		"params.connStr":  d.GetOption("params.connStr").String(),
		"params.interval": d.GetOption("params.interval").Int(),
		"created_at":      d.createdAt.Format(lang.DatetimeFormatterStr.Str()),
	}

	groups, _ := d.Groups()
	if len(groups) > 0 {
		groupsProfile := make([]model.Map, 0, len(groups))
		for _, g := range groups {
			groupsProfile = append(groupsProfile, g.Simple())
		}
		detail["groups"] = groupsProfile
	}

	return detail
}
