package alarm

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/maritimusj/centrum/gate/web/app"

	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/lang"

	"github.com/maritimusj/centrum/gate/web/helper"
	"github.com/maritimusj/centrum/gate/web/model"
	"github.com/maritimusj/centrum/gate/web/resource"
	"github.com/maritimusj/centrum/gate/web/response"
)

func List(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var (
			s = app.Store()

			page     = ctx.URLParamInt64Default("page", 1)
			pageSize = ctx.URLParamInt64Default("pagesize", app.Config.DefaultPageSize())

			params = []helper.OptionFN{
				helper.Page(page, pageSize),
			}

			admin = s.MustGetUserFromContext(ctx)
		)

		if !app.IsDefaultAdminUser(admin) {
			params = append(params, helper.DefaultEffect(app.Config.DefaultEffect()))
			params = append(params, helper.User(admin.GetID()))
		}

		if ctx.URLParamExists("device") {
			deviceID, err := strconv.ParseInt(ctx.URLParam("device"), 10, 0)
			if err != nil {
				return lang.ErrInvalidRequestData
			}
			device, err := s.GetDevice(deviceID)
			if err != nil {
				return err
			}
			if !app.Allow(admin, device, resource.View) {
				return lang.ErrNoPermission
			}
			params = append(params, helper.Device(deviceID))
		}

		if ctx.URLParamExists("equipment") {
			equipmentID, err := strconv.ParseInt(ctx.URLParam("equipment"), 10, 0)
			if err != nil {
				return lang.ErrInvalidRequestData
			}
			equipment, err := s.GetEquipment(equipmentID)
			if err != nil {
				return err
			}
			if !app.Allow(admin, equipment, resource.View) {
				return lang.ErrNoPermission
			}
			params = append(params, helper.Equipment(equipmentID))
		}

		var (
			start *time.Time
			end   *time.Time
		)

		if ctx.URLParamExists("start") {
			s, err := time.Parse("2006-01-02_15:04:05", ctx.URLParam("start"))
			if err != nil {
				return lang.ErrInvalidRequestData
			}
			start = &s
		}

		if ctx.URLParamExists("end") {
			s, err := time.Parse("2006-01-02_15:04:05", ctx.URLParam("end"))
			if err != nil {
				return lang.ErrInvalidRequestData
			}
			end = &s
		}

		alarms, total, err := s.GetAlarmList(start, end, params...)
		if err != nil {
			return err
		}

		var result = make([]model.Map, 0, len(alarms))
		for _, alarm := range alarms {
			brief := alarm.Brief()
			measure, err := alarm.Measure()
			if err != nil {
				brief["perm"] = iris.Map{
					"err": err,
				}
			} else {
				brief["perm"] = iris.Map{
					"view": true,
					"ctrl": app.Allow(admin, measure, resource.Ctrl),
				}
			}
			lastID := alarm.GetOption(fmt.Sprintf("read.%d", admin.GetID())).Int()
			_, total, err := s.GetCommentList(alarm, lastID, helper.Limit(1))
			if err != nil {
				return err
			}
			brief["comment"] = iris.Map{
				"total": total,
			}
			result = append(result, brief)
		}

		return iris.Map{
			"total": total,
			"list":  result,
		}
	})
}

func Detail(alarmID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()
		admin := s.MustGetUserFromContext(ctx)

		alarm, err := s.GetAlarm(alarmID)
		if err != nil {
			return err
		}

		measure, err := alarm.Measure()
		if err != nil {
			return err
		}

		if !app.Allow(admin, measure, resource.View) {
			return lang.ErrNoPermission
		}

		detail := alarm.Detail()

		lastID := alarm.GetOption(fmt.Sprintf("read.%d", admin.GetID())).Int()
		_, total, err := s.GetCommentList(alarm, lastID, helper.Limit(1))
		if err != nil {
			return err
		}
		detail["comment"] = iris.Map{
			"total": total,
		}
		return detail
	})
}

func Confirm(alarmID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			Desc string `json:"desc"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		var (
			s     = app.Store()
			admin = s.MustGetUserFromContext(ctx)
		)

		alarm, err := s.GetAlarm(alarmID)
		if err != nil {
			return err
		}

		measure, err := alarm.Measure()
		if err != nil {
			return err
		}

		if app.Allow(admin, measure, resource.View) || app.Allow(admin, measure, resource.Ctrl) {
			err = alarm.Confirm(map[string]interface{}{
				"admin": admin.Brief(),
				"time":  time.Now().Format(lang.DatetimeFormatterStr.Str()),
				"ip":    ctx.RemoteAddr(),
				"desc":  form.Desc,
			})
			if err != nil {
				return err
			}

			return lang.Ok
		}

		return lang.ErrNoPermission
	})
}

func Delete(alarmID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var (
			s     = app.Store()
			admin = s.MustGetUserFromContext(ctx)
		)

		alarm, err := s.GetAlarm(alarmID)
		if err != nil {
			return err
		}

		measure, err := alarm.Measure()
		if err != nil {
			return err
		}

		if !app.Allow(admin, measure, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		err = alarm.Destroy()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}

func Statistics(alarmID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var (
			s     = app.Store()
			admin = s.MustGetUserFromContext(ctx)
		)

		alarm, err := s.GetAlarm(alarmID)
		if err != nil {
			return err
		}

		measure, err := alarm.Measure()
		if err != nil {
			return err
		}

		device := measure.Device()
		if device == nil {
			return lang.ErrDeviceNotFound
		}

		org, _ := device.Organization()

		if !app.Allow(admin, measure, resource.View) {
			return lang.ErrNoPermission
		}

		start := alarm.CreatedAt()
		end := alarm.UpdatedAt()

		//最多取10000个点位数据
		interval := int64(end.Sub(start).Seconds() / 10000)
		if interval < 1 {
			interval = 1
		}

		//前后各增加100个点位需要的时间
		start = start.Add(-1 * time.Duration(interval*100) * time.Second)
		end = end.Add(time.Duration(interval*100) * time.Second)
		if end.After(time.Now()) {
			end = time.Now()
		}

		result, err := app.StatsDB.GetMeasureStats(org.Name(), device.GetID(), measure.TagName(), &start, &end, time.Duration(interval)*time.Second)
		if err != nil {
			return lang.InternalError(err)
		}

		return result

	})
}

func Export(ctx iris.Context) {
	res := response.Wrap(func() interface{} {
		csvFile, err := ioutil.TempFile("", "tempFile")
		if err != nil {
			return err
		}

		defer func() {
			_ = os.Remove(csvFile.Name())
		}()

		//写入UTF-8 BOM，防止中文乱码
		_, err = csvFile.WriteString("\xEF\xBB\xBF")
		if err != nil {
			return err
		}

		var (
			s      = app.Store()
			params = make([]helper.OptionFN, 0)
			admin  = s.MustGetUserFromContext(ctx)
		)

		if !app.IsDefaultAdminUser(admin) {
			params = append(params, helper.DefaultEffect(app.Config.DefaultEffect()))
			params = append(params, helper.User(admin.GetID()))
		}

		if ctx.URLParamExists("device") {
			deviceID, err := strconv.ParseInt(ctx.URLParam("device"), 10, 0)
			if err != nil {
				return lang.ErrInvalidRequestData
			}
			device, err := s.GetDevice(deviceID)
			if err != nil {
				return err
			}
			if !app.Allow(admin, device, resource.View) {
				return lang.ErrNoPermission
			}
			params = append(params, helper.Device(deviceID))
		}

		if ctx.URLParamExists("equipment") {
			equipmentID, err := strconv.ParseInt(ctx.URLParam("equipment"), 10, 0)
			if err != nil {
				return lang.ErrInvalidRequestData
			}
			equipment, err := s.GetEquipment(equipmentID)
			if err != nil {
				return err
			}
			if !app.Allow(admin, equipment, resource.View) {
				return lang.ErrNoPermission
			}
			params = append(params, helper.Equipment(equipmentID))
		}

		var (
			start *time.Time
			end   *time.Time
		)

		if ctx.URLParamExists("start") {
			s, err := time.Parse("2006-01-02_15:04:05", ctx.URLParam("start"))
			if err != nil {
				return lang.ErrInvalidRequestData
			}
			start = &s
		}

		if ctx.URLParamExists("end") {
			s, err := time.Parse("2006-01-02_15:04:05", ctx.URLParam("end"))
			if err != nil {
				return lang.ErrInvalidRequestData
			}
			end = &s
		}

		alarms, _, err := s.GetAlarmList(start, end, params...)
		if err != nil {
			return err
		}

		seg := make([]string, 0, 10)

		_, _ = csvFile.WriteString(strings.Join([]string{
			"#",
			lang.CVSHeaderDevice.Str(),
			lang.CVSHeaderPoint.Str(),
			lang.CVSHeaderVal.Str(),
			lang.CVSHeaderThreshold.Str(),
			lang.CVSHeaderAlarm.Str(),
			lang.CVSHeaderCreatedAt.Str(),
			lang.CVSHeaderUpdatedAt.Str(),
			lang.CVSHeaderUser.Str(),
			lang.CVSHeaderConfirmedBy.Str() + "\r\n",
		}, ","))

		for index, alarm := range alarms {
			seg = append(seg, strconv.FormatInt(int64(index+1), 10))

			device, err := alarm.Device()
			if err != nil {
				seg = append(seg, "<unknown>")
			} else {
				seg = append(seg, device.Title())
			}

			seg = append(seg,
				alarm.GetOption("name").String(),
				alarm.GetOption("fields.val").String(),
				alarm.GetOption("fields.threshold").String(),
				alarm.GetOption("tags.alarm").String(),
				alarm.CreatedAt().Format(lang.DatetimeFormatterStr.Str()),
				alarm.UpdatedAt().Format(lang.DatetimeFormatterStr.Str()),
				alarm.GetOption("confirm.admin.name").String(),
				alarm.GetOption("confirm.time").String())

			_, _ = csvFile.WriteString(strings.Join(seg, ",") + "\r\n")

			seg = seg[0:0]
		}

		_ = csvFile.Close()
		_ = ctx.SendFile(csvFile.Name(), time.Now().Format("2006-01-02_15_04_05")+".csv")

		return nil
	})

	if _, ok := res.(error); ok {
		ctx.StatusCode(iris.StatusNotFound)
	}
}
