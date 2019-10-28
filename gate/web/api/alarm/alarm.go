package alarm

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	app2 "github.com/maritimusj/centrum/gate/web/app"
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	response2 "github.com/maritimusj/centrum/gate/web/response"
	"time"
)

func List(ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		s := app2.Store()

		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", app2.Config.DefaultPageSize())

		var params = []helper2.OptionFN{
			helper2.Page(page, pageSize),
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app2.IsDefaultAdminUser(admin) {
			params = append(params, helper2.DefaultEffect(app2.Config.DefaultEffect()))
			params = append(params, helper2.User(admin.GetID()))
		}

		var (
			start *time.Time
			end   *time.Time
		)
		if ctx.URLParamExists("start") {
			s, err := time.Parse("2006-01-02_15:04:05", ctx.URLParam("start"))
			if err != nil {
				return lang2.ErrInvalidRequestData
			}
			start = &s
		}

		if ctx.URLParamExists("end") {
			s, err := time.Parse("2006-01-02_15:04:05", ctx.URLParam("start"))
			if err != nil {
				return lang2.ErrInvalidRequestData
			}
			end = &s
		}

		alarms, total, err := s.GetAlarmList(start, end, params...)
		if err != nil {
			return err
		}

		var result = make([]model2.Map, 0, len(alarms))
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
					"ctrl": app2.Allow(admin, measure, resource2.Ctrl),
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
	return response2.Wrap(func() interface{} {
		s := app2.Store()
		admin := s.MustGetUserFromContext(ctx)

		alarm, err := s.GetAlarm(alarmID)
		if err != nil {
			return err
		}

		measure, err := alarm.Measure()
		if err != nil {
			return err
		}

		if !app2.Allow(admin, measure, resource2.View) {
			return lang2.ErrNoPermission
		}

		return alarm.Detail()
	})
}

func Confirm(alarmID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		var form struct {
			Desc string `json:"desc"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang2.ErrInvalidRequestData
		}

		s := app2.Store()
		admin := s.MustGetUserFromContext(ctx)

		alarm, err := s.GetAlarm(alarmID)
		if err != nil {
			return err
		}

		measure, err := alarm.Measure()
		if err != nil {
			return err
		}

		if !app2.Allow(admin, measure, resource2.Ctrl) {
			return lang2.ErrNoPermission
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
		return lang2.Ok
	})
}

func Delete(alarmID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		s := app2.Store()
		admin := s.MustGetUserFromContext(ctx)

		alarm, err := s.GetAlarm(alarmID)
		if err != nil {
			return err
		}

		measure, err := alarm.Measure()
		if err != nil {
			return err
		}

		if !app2.Allow(admin, measure, resource2.Ctrl) {
			return lang2.ErrNoPermission
		}

		err = alarm.Destroy()
		if err != nil {
			return err
		}
		return lang2.Ok
	})
}
