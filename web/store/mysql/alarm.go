package mysqlStore

import (
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/web/dirty"
	"github.com/maritimusj/centrum/web/model"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"time"
)

type Alarm struct {
	id int64

	orgID     int64
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

	return alarm.Save()
}
