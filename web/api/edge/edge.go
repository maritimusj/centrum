package edge

import (
	"github.com/kataras/iris"
	edgeLang "github.com/maritimusj/centrum/edge/lang"
	"github.com/maritimusj/centrum/global"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/web/app"
	"github.com/maritimusj/centrum/web/resource"
	log "github.com/sirupsen/logrus"
	"time"
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
		log.Error("[Feedback]", err)
		return
	}

	var form struct {
		Log     *Log     `json:"log"`
		Status  *Status  `json:"status"`
		Measure *Measure `json:"measure"`
		Alarm   *Alarm   `json:"alarm"`
	}

	if err := ctx.ReadJSON(&form); err != nil {
		log.Error("[Feedback]", err)
		return
	}

	if form.Log != nil {
		level, err := log.ParseLevel(form.Log.Level)
		if err != nil {
			log.Error("[Feedback]", err)
		} else {
			device.Logger().Log(level, form.Log.Message)
		}
	}

	if form.Status != nil {
		if form.Status.Index == int(edgeLang.Disconnected) {
			org, _ := global.GetDeviceStatus(device)
			if org == int(edgeLang.Connected) {

			}
		}
		global.UpdateDeviceStatus(device, form.Status.Index, form.Status.Title)
	}

	if form.Measure != nil {
		measure, err := app.Store().GetMeasureFromTagName(device.GetID(), form.Measure.TagName)
		if err != nil {
			if err != lang.Error(lang.ErrMeasureNotFound) {
				log.Error("[Feedback]", err)
				return
			}
			kind := resource.ParseMeasureKind(form.Measure.TagName)
			if kind == resource.UnknownKind {
				log.Error("[Feedback]", lang.Error(lang.ErrMeasureNotFound))
				return
			}
			measure, err = app.Store().CreateMeasure(device.GetID(), form.Measure.Title, form.Measure.TagName, kind)
			if err != nil {
				log.Error("[Feedback]", err)
				return
			}
		} else {
			measure.SetTitle(form.Measure.Title)
			err = measure.Save()
			if err != nil {
				log.Error("[Feedback]", err)
			}
		}
	}

	if form.Alarm != nil {
		//保存警报信息
		tag, _ := form.Alarm.Tags["tag"]
		measure, err := app.Store().GetMeasureFromTagName(device.GetID(), tag)
		if err != nil {
			log.Error("[Feedback]", err)
			return
		}

		alarm, err := app.Store().GetLastUnconfirmedAlarm(device, measure.GetID())
		if err != nil {
			if err != lang.Error(lang.ErrAlarmNotFound) {
				log.Error("[Feedback]", err)
				return
			}
			alarm, err = app.Store().CreateAlarm(device, measure.GetID(), map[string]interface{}{
				"name":   form.Alarm.Name,
				"tags":   form.Alarm.Tags,
				"fields": form.Alarm.Fields,
				"time":   form.Alarm.Time,
			})
			if err != nil {
				log.Error("[Feedback]", err)
			}
		} else {
			alarm.Updated()
			if err = alarm.Save(); err != nil {
				log.Error("[Feedback]", err)
			}
		}
	}
}
