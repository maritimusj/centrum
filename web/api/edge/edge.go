package edge

import (
	"github.com/kataras/iris"
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
}
