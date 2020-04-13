package edge

import (
	"strconv"
	"time"

	"github.com/maritimusj/centrum/util"

	"github.com/maritimusj/centrum/gate/web/app"

	"github.com/maritimusj/centrum/gate/lang"

	"github.com/maritimusj/centrum/gate/web/edge"
	"github.com/maritimusj/centrum/gate/web/helper"
	"github.com/maritimusj/centrum/gate/web/resource"

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

type Perf struct {
	Rate int `json:"rate"`
}

func Feedback(deviceID int64, ctx iris.Context) {
	device, err := app.Store().GetDevice(deviceID)
	if err != nil {
		if err != lang.Error(lang.ErrDeviceNotFound) {
			log.Debugln("[Feedback 1]", err)
		} else {
			edge.Remove(strconv.FormatInt(deviceID, 10))
		}
		return
	}

	var form struct {
		Log     *Log     `json:"log"`
		Status  *Status  `json:"status"`
		Measure *Measure `json:"measure"`
		Alarm   *Alarm   `json:"alarm"`
		Perf    *Perf    `json:"perf"`
	}

	if err := ctx.ReadJSON(&form); err != nil {
		log.Debugln("[Feedback 2]", err)
		return
	}

	if form.Status != nil {
		if form.Status.Index == int(edgeLang.Disconnected) {
			global.UpdateDevicePerf(device, iris.Map{})

			org, _, _ := global.GetDeviceStatus(device)
			if org == int(edgeLang.Connected) {
				device.Logger().Warningln(lang.Error(lang.ErrDeviceDisconnected))
			}
		}
		global.UpdateDeviceStatus(device, form.Status.Index, form.Status.Title)
	}

	if form.Measure != nil {
		measure, err := app.Store().GetMeasureFromTagName(device.GetID(), form.Measure.TagName)
		if err != nil {
			if err != lang.Error(lang.ErrMeasureNotFound) {
				log.Debugln("[Feedback 4]", err)
				return
			}
			kind := resource.ParseMeasureKind(form.Measure.TagName)
			if kind == resource.UnknownKind {
				log.Debugln("[Feedback 5]", lang.Error(lang.ErrMeasureNotFound))
				return
			}
			measure, err = app.Store().CreateMeasure(device.GetID(), form.Measure.Title, form.Measure.TagName, kind)
			if err != nil {
				log.Debugln("[Feedback 6]", err)
				return
			}
		} else {
			measure.SetTitle(form.Measure.Title)
			err = measure.Save()
			if err != nil {
				log.Debugln("[Feedback 7]", err)
			}
		}
	}

	if form.Alarm != nil {
		//保存警报信息
		tag, _ := form.Alarm.Tags["tag"]
		measure, err := app.Store().GetMeasureFromTagName(device.GetID(), tag)
		if err != nil {
			log.Debugln("[Feedback 8] [", tag, "]", err)
			return
		}

		alarm, _, err := app.Store().GetLastUnconfirmedAlarm(helper.Device(device.GetID()), helper.Measure(measure.GetID()))
		if err != nil {
			if err != lang.Error(lang.ErrAlarmNotFound) {
				log.Debugln("[Feedback 9]", err)
				return
			}
			alarm, err = app.Store().CreateAlarm(device, measure.GetID(), map[string]interface{}{
				"name":   form.Alarm.Name,
				"tags":   form.Alarm.Tags,
				"fields": form.Alarm.Fields,
				"time":   form.Alarm.Time,
			})
			if err != nil {
				log.Debugln("[Feedback 10]", err)
			}
		} else {
			alarm.Updated()
			if err = alarm.Save(); err != nil {
				log.Debugln("[Feedback 11]", err)
			}
		}
	}

	if form.Perf != nil {
		rate := uint64((*form.Perf).Rate)
		str := util.FormatFileSize(rate)
		data := iris.Map{
			"rate": str + "/s",
		}
		level := 1
		if rate < 1 {
			level = 1
		} else if rate < 10 {
			level = 2
		} else if rate < 50 {
			level = 3
		} else if rate < 100 {
			level = 4
		} else {
			level = 5
		}
		data["level"] = level
		global.UpdateDevicePerf(device, data)
	}

	if form.Log != nil {
		level, err := log.ParseLevel(form.Log.Level)
		if err != nil {
			log.Debugln("[Feedback 3]", err)
		} else {
			device.Logger().Log(level, form.Log.Message)
		}
	}
}
