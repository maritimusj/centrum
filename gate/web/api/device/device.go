package device

import (
	"fmt"
	"strings"
	"time"

	"github.com/maritimusj/durafmt"

	edgeLang "github.com/maritimusj/centrum/edge/lang"

	"github.com/maritimusj/centrum/gate/event"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/gate/web/helper"
	"github.com/maritimusj/centrum/gate/web/model"
	"github.com/maritimusj/centrum/gate/web/resource"
	"github.com/maritimusj/centrum/gate/web/response"
	"github.com/maritimusj/centrum/gate/web/store"

	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/global"
)

func List(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var (
			params []helper.OptionFN
			orgID  int64
			s      = app.Store()
			admin  = s.MustGetUserFromContext(ctx)
		)

		if app.IsDefaultAdminUser(admin) {
			if ctx.URLParamExists("org") {
				orgID = ctx.URLParamInt64Default("org", 0)
			}
		} else {
			orgID = admin.OrganizationID()
		}

		if orgID > 0 {
			params = append(params, helper.Organization(orgID))
		}

		var (
			page     = ctx.URLParamInt64Default("page", 1)
			pageSize = ctx.URLParamInt64Default("pagesize", app.Config.DefaultPageSize())
		)

		params = append(params, helper.Page(page, pageSize))

		keyword := ctx.URLParam("keyword")
		if keyword != "" {
			params = append(params, helper.Keyword(keyword))
		}

		if ctx.URLParamExists("group") {
			groupID := ctx.URLParamInt64Default("group", 0)
			params = append(params, helper.Group(groupID))
		}

		if !app.IsDefaultAdminUser(admin) {
			params = append(params, helper.User(admin.GetID()))
			params = append(params, helper.DefaultEffect(app.Config.DefaultEffect()))
		}

		devices, total, err := s.GetDeviceList(params...)
		if err != nil {
			return err
		}

		var (
			result = make([]model.Map, 0, len(devices))
		)

		for _, device := range devices {
			brief := device.Brief()

			groups, err := device.Groups()
			if err != nil {
				brief["groups"] = make([]interface{}, 0)
			} else {
				groupTitles := make([]string, 0)
				for _, group := range groups {
					groupTitles = append(groupTitles, group.Title())
				}
				brief["groups"] = groupTitles
			}

			brief["perm"] = iris.Map{
				"view": true,
				"ctrl": app.Allow(admin, device, resource.Ctrl),
			}

			index, title, from := global.GetDeviceStatus(device)
			status := iris.Map{
				"index": index,
				"title": title,
			}
			if index == int(edgeLang.Connected) {
				status["from"] = from.Format(lang.DatetimeFormatterStr.Str())
				status["duration"] = strings.ReplaceAll(durafmt.Parse(time.Now().Sub(from)).LimitFirstN(2).String(), " ", "")
			}
			brief["edge"] = iris.Map{
				"status": status,
			}

			var params = []helper.OptionFN{helper.Device(device.GetID()), helper.Limit(1)}
			if !app.IsDefaultAdminUser(admin) {
				params = append(params, helper.User(admin.GetID()))
			}

			_, total, err := s.GetMeasureList(params...)
			if err != nil {
				return err
			}
			brief["measure"] = iris.Map{
				"total": total,
			}

			_, total, err = s.GetLastUnconfirmedAlarm(params...)
			if err != nil {
				if err != lang.ErrAlarmNotFound.Error() {
					return err
				}
			}
			brief["alarm"] = iris.Map{
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

func Create(ctx iris.Context) hero.Result {
	if !app.IsRegistered() {
		return response.Wrap(lang.ErrRegFirst)
	}

	var form struct {
		OrgID    int64   `json:"org"`
		Title    string  `json:"title" valid:"required"`
		Groups   []int64 `json:"groups"`
		ConnStr  string  `json:"params.connStr" valid:"required"`
		Interval int64   `json:"params.interval"`
	}

	if err := ctx.ReadJSON(&form); err != nil {
		return response.Wrap(lang.ErrInvalidRequestData)
	}

	if _, err := govalidator.ValidateStruct(&form); err != nil {
		return response.Wrap(lang.ErrInvalidRequestData)
	}

	if govalidator.IsIPv4(form.ConnStr) {
		form.ConnStr += ":502"
	} else if govalidator.IsIPv6(form.ConnStr) {
		form.ConnStr = fmt.Sprintf("[%s]:502", form.ConnStr)
	} else if govalidator.IsMAC(form.ConnStr) {
		form.ConnStr = strings.ToLower(form.ConnStr)
	}

	fn := func() interface{} {
		result := app.TransactionDo(func(s store.Store) interface{} {
			var org interface{}

			admin := s.MustGetUserFromContext(ctx)
			if app.IsDefaultAdminUser(admin) {
				if form.OrgID > 0 {
					org = form.OrgID
				} else {
					org = app.Config.DefaultOrganization()
				}
			} else {
				org = admin.OrganizationID()
			}

			device, err := s.CreateDevice(org, form.Title, map[string]interface{}{
				"params": map[string]interface{}{
					"connStr":  form.ConnStr,
					"interval": form.Interval,
				},
			})

			if err != nil {
				return err
			}

			if len(form.Groups) > 0 {
				var groups []interface{}
				for _, g := range form.Groups {
					groups = append(groups, g)
				}
				err = device.SetGroups(groups...)
				if err != nil {
					return err
				}
			}

			err = app.SetAllow(admin, device, resource.View, resource.Ctrl)
			if err != nil {
				return err
			}

			return event.Data{
				"deviceID": device.GetID(),
				"userID":   admin.GetID(),
				"result":   device.Simple(),
			}
		})

		if data, ok := result.(event.Data); ok {
			app.Event.Publish(event.DeviceCreated, data.Get("userID"), data.Get("deviceID"))
			return data.Pop("result")
		}

		return result
	}

	return response.Wrap(fn())
}

func MultiStatus(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			DeviceIDs []int64 `json:"devices"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		var (
			s      = app.Store()
			admin  = s.MustGetUserFromContext(ctx)
			result = make([]iris.Map, 0)
		)

		for _, id := range form.DeviceIDs {
			device, err := s.GetDevice(id)
			if err != nil {
				result = append(result, iris.Map{
					"id":    id,
					"error": err.Error(),
				})
			} else {
				if !app.Allow(admin, device, resource.View) {
					result = append(result, iris.Map{
						"id":    id,
						"error": lang.ErrNoPermission.Error(),
					})
				} else {
					index, title, from := global.GetDeviceStatus(device)
					status := iris.Map{
						"id":    id,
						"index": index,
						"title": title,
					}
					if index == int(edgeLang.Connected) {
						status["from"] = from.Format(lang.DatetimeFormatterStr.Str())
						status["duration"] = strings.ReplaceAll(durafmt.Parse(time.Now().Sub(from)).LimitFirstN(2).String(), " ", "")
					}
					status["perf"] = global.GetDevicePerf(device)
					result = append(result, status)
				}
			}
		}

		return result
	})
}

func Detail(deviceID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		device, err := app.Store().GetDevice(deviceID)
		if err != nil {
			return err
		}

		admin := app.Store().MustGetUserFromContext(ctx)
		if !app.Allow(admin, device, resource.View) {
			return lang.ErrNoPermission
		}

		detail := device.Detail()
		detail["perm"] = iris.Map{
			"view": true,
			"ctrl": app.Allow(admin, device, resource.Ctrl),
		}
		return detail
	})
}

func Update(deviceID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			Title    *string  `json:"title"`
			ConnStr  *string  `json:"params.connStr"`
			Interval *int64   `json:"params.interval"`
			Groups   *[]int64 `json:"groups"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		result := app.TransactionDo(func(s store.Store) interface{} {
			device, err := s.GetDevice(deviceID)
			if err != nil {
				return err
			}

			admin := s.MustGetUserFromContext(ctx)
			if !app.Allow(admin, device, resource.Ctrl) {
				return lang.ErrNoPermission
			}

			logFields := make(map[string]interface{})

			if form.Title != nil {
				device.SetTitle(*form.Title)
				logFields["title"] = form.Title
			}

			if form.ConnStr != nil {
				if govalidator.IsIPv4(*form.ConnStr) {
					*form.ConnStr += ":502"
				} else if govalidator.IsIPv6(*form.ConnStr) {
					*form.ConnStr = fmt.Sprintf("[%s]:502", *form.ConnStr)
				} else if govalidator.IsMAC(*form.ConnStr) {
					*form.ConnStr = strings.ToLower(*form.ConnStr)
				}

				err = device.SetOption("params.connStr", form.ConnStr)
				if err != nil {
					return err
				}
				logFields["connStr"] = form.ConnStr
			}

			if form.Interval != nil {
				err = device.SetOption("params.interval", form.Interval)
				if err != nil {
					return err
				}
				logFields["Interval"] = form.Interval
			}

			if form.Groups != nil {
				var groups []interface{}
				for _, g := range *form.Groups {
					groups = append(groups, g)
				}
				err = device.SetGroups(groups...)
				if err != nil {
					return err
				}
			}

			err = device.Save()
			if err != nil {
				return err
			}

			data := event.Data{
				"userID":   admin.GetID(),
				"deviceID": device.GetID(),
			}

			return data
		})

		if data, ok := result.(event.Data); ok {
			app.Event.Publish(event.DeviceUpdated, data.GetMulti("userID", "deviceID")...)
			return lang.Ok
		}

		return result
	})
}

func Delete(deviceID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		result := app.TransactionDo(func(s store.Store) interface{} {
			device, err := s.GetDevice(deviceID)
			if err != nil {
				return err
			}

			admin := s.MustGetUserFromContext(ctx)
			data := event.Data{
				"id":     device.GetID(),
				"uid":    device.UID(),
				"title":  device.Title(),
				"userID": admin.GetID(),
			}
			if app.IsDefaultAdminUser(admin) {
				err = device.Destroy()
				if err != nil {
					return err
				}
			} else {
				if !app.Allow(admin, device, resource.Ctrl) {
					return lang.ErrNoPermission
				}
				err = app.SetDeny(admin, device, resource.View, resource.Ctrl)
				if err != nil {
					return err
				}
			}
			return data
		})

		if data, ok := result.(event.Data); ok {
			app.Event.Publish(event.DeviceDeleted, data.GetMulti("userID", "id", "uid", "title")...)
			return lang.Ok
		}
		return result
	})
}

func GetLastAlarm(_ int64, measureID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		s := app.Store()
		admin := s.MustGetUserFromContext(ctx)

		measure, err := app.Store().GetMeasure(measureID)
		if err != nil {
			return err
		}

		if !app.Allow(admin, measure, resource.View) {
			return lang.ErrNoPermission
		}

		alarm, _, err := s.GetLastAlarm(helper.Measure(measure.GetID()), helper.OrderBy("id DESC"))
		if err != nil {
			return err
		}

		return alarm.Detail()
	})
}
