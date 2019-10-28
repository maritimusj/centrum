package mysqlStore

import (
	"database/sql"
	"errors"
	"fmt"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	dirty2 "github.com/maritimusj/centrum/gate/web/dirty"
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	status2 "github.com/maritimusj/centrum/gate/web/status"
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

	dirty *dirty2.Dirty
	store *mysqlStore
}

func NewDDevice(store *mysqlStore, id int64) *Device {
	return &Device{
		id:    id,
		dirty: dirty2.New(),
		store: store,
	}
}

func (d *Device) OrganizationID() int64 {
	if d != nil {
		return d.orgID
	}
	return 0
}

func (d *Device) Organization() (model2.Organization, error) {
	if d != nil {
		return d.store.GetOrganization(d.orgID)
	}
	return nil, lang2.Error(lang2.ErrDeviceNotFound)
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

func (d *Device) ResourceClass() resource2.Class {
	return resource2.Device
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

func (d *Device) GetChildrenResources(options ...helper2.OptionFN) ([]model2.Resource, int64, error) {
	if d != nil {
		options = append(options, helper2.Device(d.GetID()))
		measures, total, err := d.store.GetMeasureList(options...)
		if err != nil {
			return nil, 0, err
		}

		result := make([]model2.Resource, 0, len(measures))
		for _, measure := range measures {
			result = append(result, measure)
		}

		return result, total, nil
	}

	return nil, 0, lang2.Error(lang2.ErrDeviceNotFound)
}

func (d *Device) GetID() int64 {
	if d != nil {
		return d.id
	}
	return 0
}

func (d *Device) Enable() {
	if d != nil && d.enable != status2.Enable {
		d.enable = status2.Enable
		d.dirty.Set("enable", func() interface{} {
			return d.enable
		})
	}
}

func (d *Device) Disable() {
	if d != nil && d.enable != status2.Disable {
		d.enable = status2.Disable
		d.dirty.Set("enable", func() interface{} {
			return d.enable
		})
	}
}

func (d *Device) IsEnabled() bool {
	return d != nil && d.enable == status2.Enable
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
	return lang2.Error(lang2.ErrDeviceNotFound)
}

func (d *Device) SetGroups(groups ...interface{}) error {
	if d != nil {
		err := RemoveData(d.store.db, TbDeviceGroups, "device_id=?", d.id)
		if err != nil {
			if err != sql.ErrNoRows {
				return lang2.InternalError(err)
			}
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
				return lang2.InternalError(err)
			}
		}
		return nil
	}
	return lang2.Error(lang2.ErrDeviceNotFound)
}

func (d *Device) Groups() ([]model2.Group, error) {
	if d != nil {
		groups, err := d.store.GetDeviceGroups(d.GetID())
		if err != nil {
			return nil, err
		}
		return groups, nil
	}
	return nil, lang2.Error(lang2.ErrDeviceNotFound)
}

func (d *Device) GetMeasureList(options ...helper2.OptionFN) ([]model2.Measure, int64, error) {
	if d != nil {
		return d.store.GetMeasureList(options...)
	}
	return nil, 0, lang2.Error(lang2.ErrDeviceNotFound)
}

func (d *Device) CreateMeasure(title string, tag string, kind resource2.MeasureKind) (model2.Measure, error) {
	if d != nil {
		return d.store.CreateMeasure(d.GetID(), title, tag, kind)
	}
	return nil, lang2.Error(lang2.ErrDeviceNotFound)
}

func (d *Device) CreatedAt() time.Time {
	if d != nil {
		return d.createdAt
	}
	return time.Time{}
}

func (d *Device) Destroy() error {
	if d == nil {
		return lang2.Error(lang2.ErrDeviceNotFound)
	}

	alarms, _, err := d.store.GetAlarmList(nil, nil, helper2.Device(d.GetID()))
	if err != nil {
		return err
	}

	for _, alarm := range alarms {
		err = alarm.Destroy()
		if err != nil {
			return err
		}
	}

	measures, _, err := d.store.GetMeasureList(helper2.Device(d.GetID()))
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
				return lang2.InternalError(err)
			}
		}
		return nil
	}
	return lang2.Error(lang2.ErrDeviceNotFound)
}

func (d *Device) Simple() model2.Map {
	if d == nil {
		return model2.Map{}
	}
	return model2.Map{
		"id":     d.id,
		"enable": d.IsEnabled(),
		"title":  d.title,
	}
}

func (d *Device) Brief() model2.Map {
	if d == nil {
		return model2.Map{}
	}
	return model2.Map{
		"id":             d.id,
		"enable":         d.IsEnabled(),
		"title":          d.title,
		"params.connStr": d.GetOption("params.connStr").Str,
		"created_at":     d.createdAt,
	}
}

func (d *Device) Detail() model2.Map {
	if d == nil {
		return model2.Map{}
	}

	detail := model2.Map{
		"id":              d.id,
		"enable":          d.IsEnabled(),
		"title":           d.title,
		"params.connStr":  d.GetOption("params.connStr").String(),
		"params.interval": d.GetOption("params.interval").Int(),
		"created_at":      d.createdAt,
	}

	groups, _ := d.Groups()
	if len(groups) > 0 {
		groupsProfile := make([]model2.Map, 0, len(groups))
		for _, g := range groups {
			groupsProfile = append(groupsProfile, g.Simple())
		}
		detail["groups"] = groupsProfile
	}

	return detail
}
