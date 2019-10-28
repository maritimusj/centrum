package edge

import (
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/gate/web/resource"
	"time"

	"github.com/kataras/iris"
	edgeLang "github.com/maritimusj/centrum/edge/lang"
	"github.com/maritimusj/centrum/global"
	log "github.com/sirupsen/logrus"
)

type Log struct {
	Level     string    `json:"level"`
	Message   string    `json:"msg"`
	CreatedAt time.Time `json:"time"`
}

type Status struct {
	Index int    `json:"index"`
	Title string `json:"title"`
}

type Measure struct {
	TagName string `json:"tag"`
	Title   string `json:"title"`
}

type Alarm struct {
	Name   string                 `json:"name"`
	Tags   map[string]string      `json:"tags"`
	Fields map[string]interface{} `json:"fields"`
	Time   time.Time              `json:"time"`
}

func Feedback(deviceID int64, ctx iris.Context) {
	device, err := app.Store().GetDevice(deviceID)
	if err != nil {
		log.Error("[Feedback 1]", err)
		return
	}

	var form struct {
		Log     *Log     `json:"log"`
		Status  *Status  `json:"status"`
		Measure *Measure `json:"measure"`
		Alarm   *Alarm   `json:"alarm"`
	}

	if err := ctx.ReadJSON(&form); err != nil {
		log.Error("[Feedback 2]", err)
		return
	}

	if form.Log != nil {
		level, err := log.ParseLevel(form.Log.Level)
		if err != nil {
			log.Error("[Feedback 3]", err)
		} else {
			device.Logger().Log(level, form.Log.Message)
		}
	}

	if form.Status != nil {
		if form.Status.Index == int(edgeLang.Disconnected) {
			org, _ := global.GetDeviceStatus(device)
			if org == int(edgeLang.Connected) {
				//断线警报？
			}
		}
		global.UpdateDeviceStatus(device, form.Status.Index, form.Status.Title)
	}

	if form.Measure != nil {
		measure, err := app.Store().GetMeasureFromTagName(device.GetID(), form.Measure.TagName)
		if err != nil {
			if err != lang.Error(lang.ErrMeasureNotFound) {
				log.Error("[Feedback 4]", err)
				return
			}
			kind := resource.ParseMeasureKind(form.Measure.TagName)
			if kind == resource.UnknownKind {
				log.Error("[Feedback 5]", lang.Error(lang.ErrMeasureNotFound))
				return
			}
			measure, err = app.Store().CreateMeasure(device.GetID(), form.Measure.Title, form.Measure.TagName, kind)
			if err != nil {
				log.Error("[Feedback 6]", err)
				return
			}
		} else {
			measure.SetTitle(form.Measure.Title)
			err = measure.Save()
			if err != nil {
				log.Error("[Feedback 7]", err)
			}
		}
	}

	if form.Alarm != nil {
		//保存警报信息
		tag, _ := form.Alarm.Tags["tag"]
		measure, err := app.Store().GetMeasureFromTagName(device.GetID(), tag)
		if err != nil {
			log.Error("[Feedback 8]", err)
			return
		}

		alarm, err := app.Store().GetLastUnconfirmedAlarm(device, measure.GetID())
		if err != nil {
			if err != lang.Error(lang.ErrAlarmNotFound) {
				log.Error("[Feedback 9]", err)
				return
			}
			alarm, err = app.Store().CreateAlarm(device, measure.GetID(), map[string]interface{}{
				"name":   form.Alarm.Name,
				"tags":   form.Alarm.Tags,
				"fields": form.Alarm.Fields,
				"time":   form.Alarm.Time,
			})
			if err != nil {
				log.Error("[Feedback 10]", err)
			}
		} else {
			alarm.Updated()
			if err = alarm.Save(); err != nil {
				log.Error("[Feedback 11]", err)
			}
		}
	}
}
