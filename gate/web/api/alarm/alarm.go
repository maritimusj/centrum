package alarm

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/web/app"
	"github.com/maritimusj/centrum/web/helper"
	"github.com/maritimusj/centrum/web/model"
	"github.com/maritimusj/centrum/web/resource"
	"github.com/maritimusj/centrum/web/response"
	"time"
)

func List(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()

		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", app.Config.DefaultPageSize())

		var params = []helper.OptionFN{
			helper.Page(page, pageSize),
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app.IsDefaultAdminUser(admin) {
			params = append(params, helper.DefaultEffect(app.Config.DefaultEffect()))
			params = append(params, helper.User(admin.GetID()))
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
			s, err := time.Parse("2006-01-02_15:04:05", ctx.URLParam("start"))
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

		return alarm.Detail()
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

		if !app.Allow(admin, measure, resource.Ctrl) {
			return lang.ErrNoPermission
		}

		err = alarm.Confirm(map[string]interface{}{
			"admin": admin.Brief(),
			"time":  time.Now(),
			"ip":    ctx.RemoteAddr(),
			"desc":  form.Desc,
		})
		if err != nil {
			return err
		}
		return lang.Ok
	})
}

func Delete(alarmID int64, ctx iris.Context) hero.Result {
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
