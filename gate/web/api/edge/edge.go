package edge

import (
	"fmt"
	"strconv"
	"time"

	"github.com/maritimusj/centrum/gate/web/model"

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
	Delay int `json:"delay"`
}

func createMsg(res model.Resource, data interface{}) {
	store := app.Store()
	global.AddMessage(data, func(uid string, userId int64) bool {
		user, err := store.GetUser(userId)
		if err != nil {
			global.Close(uid, userId)
			return false
		}

		return app.Allow(user, res, resource.View)
	})
}

//将新创建的点位授权给对设备有相应权限的用户
//暂时采用循环遍历用户来判断
func AssignPermissionsToUsers(device model.Device, measure model.Measure) error {
	if app.Config.DefaultEffect() == resource.Allow {
		return nil
	}

	users, _, err := app.Store().GetUserList()
	if err != nil {
		return err
	}

	for _, user := range users {
		if app.IsDefaultAdminUser(user) {
			continue
		}

		if app.Allow(user, device, resource.Ctrl) {
			_ = app.SetAllow(user, measure, resource.View, resource.Ctrl)
		} else if app.Allow(user, device, resource.View) {
			_ = app.SetAllow(user, measure, resource.View)
		}
	}

	return nil
}

func Feedback(deviceID int64, ctx iris.Context) {
	device, err := app.Store().GetDevice(deviceID)
	if err != nil {
		if err != lang.ErrDeviceNotFound.Error() {
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
		org, _, _ := global.GetDeviceStatus(device)

		if form.Status.Index == int(edgeLang.Disconnected) {
			global.UpdateDevicePerf(device, iris.Map{})
			if org == int(edgeLang.Connected) {
				createMsg(device, iris.Map{
					"level":      "warning",
					"title":      device.Title(),
					"message":    lang.ErrDeviceDisconnected.Str(),
					"created_at": time.Now().Format(lang.DatetimeFormatterStr.Str()),
				})
				device.Logger().Warningln(lang.ErrDeviceDisconnected.Str())
			}
		} else if form.Status.Index == int(edgeLang.Connected) {
			if org != int(edgeLang.Connected) {
				createMsg(device, iris.Map{
					"level":      "success",
					"title":      device.Title(),
					"message":    lang.DeviceConnected.Str(),
					"created_at": time.Now().Format(lang.DatetimeFormatterStr.Str()),
				})
			}
		}
		global.UpdateDeviceStatus(device, form.Status.Index, form.Status.Title)
	}

	if form.Measure != nil {
		measure, err := app.Store().GetMeasureFromTagName(device.GetID(), form.Measure.TagName)
		if err != nil {
			if err != lang.ErrMeasureNotFound.Error() {
				log.Debugln("[Feedback 4]", err)
				return
			}

			kind := resource.ParseMeasureKind(form.Measure.TagName)
			if kind == resource.UnknownKind {
				log.Debugln("[Feedback 5]", lang.ErrMeasureNotFound.Str())
				return
			}

			measure, err = app.Store().CreateMeasure(device.GetID(), form.Measure.Title, form.Measure.TagName, kind)
			if err != nil {
				log.Debugln("[Feedback 6]", err)
				return
			}

			//将新创建的点位授权给对设备有相应权限的用户
			err = AssignPermissionsToUsers(device, measure)
			if err != nil {
				log.Debugln("[Feedback 7]", err)
			}

		} else {
			measure.SetTitle(form.Measure.Title)
			err = measure.Save()
			if err != nil {
				log.Debugln("[Feedback 8]", err)
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
			if err != lang.ErrAlarmNotFound.Error() {
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
		delay := time.Duration((*form.Perf).Delay) / time.Millisecond
		data := iris.Map{
			"rate": fmt.Sprintf("%dms", delay),
		}
		level := 1
		if delay == -1 {
			level = 1
		} else if delay < 20 {
			level = 5
		} else if delay < 100 {
			level = 4
		} else if delay < 300 {
			level = 3
		} else if delay < 600 {
			level = 2
		} else {
			level = 1
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
