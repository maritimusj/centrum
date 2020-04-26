package mysqlStore

import (
	"time"

	"github.com/kataras/iris"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/dirty"
	"github.com/maritimusj/centrum/gate/web/model"
	"github.com/maritimusj/centrum/gate/web/status"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Alarm struct {
	id int64

	orgID     int64
	status    int
	deviceID  int64
	measureID int64

	extra     []byte
	createdAt time.Time
	updatedAt time.Time

	dirty *dirty.Dirty
	store *mysqlStore
}

func NewAlarm(store *mysqlStore, id int64) *Alarm {
	return &Alarm{
		id:    id,
		dirty: dirty.New(),
		store: store,
	}
}

func (alarm *Alarm) GetID() int64 {
	return alarm.id
}

func (alarm *Alarm) DeviceID() int64 {
	return alarm.deviceID
}

func (alarm *Alarm) MeasureID() int64 {
	return alarm.measureID
}

func (alarm *Alarm) Status() (int, string) {
	return alarm.status, lang.AlarmStatusDesc(alarm.status)
}

func (alarm *Alarm) Confirm(data map[string]interface{}) error {
	if err := alarm.SetOption("confirm", data); err != nil {
		return err
	}
	alarm.status = status.Confirmed
	alarm.dirty.Set("status", func() interface{} {
		return alarm.status
	})
	return alarm.Save()
}

func (alarm *Alarm) Device() (model.Device, error) {
	return alarm.store.GetDevice(alarm.deviceID)
}

func (alarm *Alarm) Measure() (model.Measure, error) {
	return alarm.store.GetMeasure(alarm.measureID)
}

func (alarm *Alarm) CreatedAt() time.Time {
	return alarm.createdAt
}

func (alarm *Alarm) UpdatedAt() time.Time {
	return alarm.updatedAt
}

func (alarm *Alarm) Updated() {
	alarm.updatedAt = time.Now()
	alarm.dirty.Set("updated_at", func() interface{} {
		return alarm.updatedAt
	})
}

func (alarm *Alarm) Save() error {
	if alarm.dirty.Any() {
		err := SaveData(alarm.store.db, TbAlarms, alarm.dirty.Data(true), "id=?", alarm.id)
		if err != nil {
			return lang.InternalError(err)
		}
	}
	return nil
}

func (alarm *Alarm) Destroy() error {
	return alarm.store.RemoveAlarm(alarm.id)
}

func (alarm *Alarm) Option() map[string]interface{} {
	return gjson.ParseBytes(alarm.extra).Value().(map[string]interface{})
}

func (alarm *Alarm) GetOption(path string) gjson.Result {
	return gjson.GetBytes(alarm.extra, path)
}

func (alarm *Alarm) SetOption(path string, value interface{}) error {
	data, err := sjson.SetBytes(alarm.extra, path, value)
	if err != nil {
		return err
	}

	alarm.extra = data
	alarm.dirty.Set("extra", func() interface{} {
		return alarm.extra
	})

	return nil
}

func (alarm *Alarm) Simple() model.Map {
	device, _ := alarm.Device()
	measure, _ := alarm.Measure()
	return model.Map{
		"id":          alarm.GetID(),
		"status":      alarm.status,
		"status_desc": lang.AlarmStatusDesc(alarm.status),
		"device":      device.Brief(),
		"measure":     measure.Brief(),
		"raw": iris.Map{
			"alarm": alarm.GetOption("tags.alarm").String(),
			"val":   alarm.GetOption("fields.val").String(),
		},
		"updated_at": alarm.updatedAt.Format(lang.DatetimeFormatterStr.Str()),
	}
}

func (alarm *Alarm) Brief() model.Map {
	device, _ := alarm.Device()
	measure, _ := alarm.Measure()
	return model.Map{
		"id":          alarm.GetID(),
		"status":      alarm.status,
		"status_desc": lang.AlarmStatusDesc(alarm.status),
		"device":      device.Brief(),
		"measure":     measure.Brief(),
		"raw":         alarm.Option(),
		"created_at":  alarm.createdAt.Format(lang.DatetimeFormatterStr.Str()),
		"updated_at":  alarm.updatedAt.Format(lang.DatetimeFormatterStr.Str()),
	}
}

func (alarm *Alarm) Detail() model.Map {
	device, _ := alarm.Device()
	measure, _ := alarm.Measure()
	return model.Map{
		"id":          alarm.GetID(),
		"status":      alarm.status,
		"status_desc": lang.AlarmStatusDesc(alarm.status),
		"device":      device.Brief(),
		"measure":     measure.Brief(),
		"raw":         alarm.Option(),
		"created_at":  alarm.createdAt.Format(lang.DatetimeFormatterStr.Str()),
		"updated_at":  alarm.updatedAt.Format(lang.DatetimeFormatterStr.Str()),
	}
}
