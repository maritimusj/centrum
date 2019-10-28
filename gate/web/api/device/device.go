package device

import (
	"fmt"
	"github.com/maritimusj/centrum/gate/event"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	app2 "github.com/maritimusj/centrum/gate/web/app"
	helper2 "github.com/maritimusj/centrum/gate/web/helper"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	response2 "github.com/maritimusj/centrum/gate/web/response"
	store2 "github.com/maritimusj/centrum/gate/web/store"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/global"
)

func List(ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		var params []helper2.OptionFN
		var orgID int64

		s := app2.Store()
		admin := s.MustGetUserFromContext(ctx)
		if app2.IsDefaultAdminUser(admin) {
			if ctx.URLParamExists("org") {
				orgID = ctx.URLParamInt64Default("org", 0)
			}
		} else {
			orgID = admin.OrganizationID()
		}
		if orgID > 0 {
			params = append(params, helper2.Organization(orgID))
		}

		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", app2.Config.DefaultPageSize())
		params = append(params, helper2.Page(page, pageSize))

		keyword := ctx.URLParam("keyword")
		if keyword != "" {
			params = append(params, helper2.Keyword(keyword))
		}

		if ctx.URLParamExists("group") {
			groupID := ctx.URLParamInt64Default("group", 0)
			params = append(params, helper2.Group(groupID))
		}

		if !app2.IsDefaultAdminUser(admin) {
			params = append(params, helper2.User(admin.GetID()))
			params = append(params, helper2.DefaultEffect(app2.Config.DefaultEffect()))
		}

		devices, total, err := s.GetDeviceList(params...)
		if err != nil {
			return err
		}

		var (
			result = make([]model2.Map, 0, len(devices))
		)

		for _, device := range devices {
			brief := device.Brief()

			brief["perm"] = iris.Map{
				"view": true,
				"ctrl": app2.Allow(admin, device, resource2.Ctrl),
			}

			index, title := global.GetDeviceStatus(device)
			brief["edge"] = iris.Map{
				"status": iris.Map{
					"index": index,
					"title": title,
				},
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
	var form struct {
		OrgID    int64   `json:"org"`
		Title    string  `json:"title" valid:"required"`
		Groups   []int64 `json:"groups"`
		ConnStr  string  `json:"params.connStr" valid:"required"`
		Interval int64   `json:"params.interval"`
	}

	if err := ctx.ReadJSON(&form); err != nil {
		return response2.Wrap(lang2.ErrInvalidRequestData)
	}

	if _, err := govalidator.ValidateStruct(&form); err != nil {
		return response2.Wrap(lang2.ErrInvalidRequestData)
	}

	if govalidator.IsIPv4(form.ConnStr) {
		form.ConnStr += ":502"
	} else if govalidator.IsIPv6(form.ConnStr) {
		form.ConnStr = fmt.Sprintf("[%s]:502", form.ConnStr)
	} else if govalidator.IsMAC(form.ConnStr) {
		form.ConnStr = strings.ToLower(form.ConnStr)
	}

	fn := func() interface{} {
		result := app2.TransactionDo(func(s store2.Store) interface{} {
			var org interface{}

			admin := s.MustGetUserFromContext(ctx)
			if app2.IsDefaultAdminUser(admin) {
				if form.OrgID > 0 {
					org = form.OrgID
				} else {
					org = app2.Config.DefaultOrganization()
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

			err = app2.SetAllow(admin, device, resource2.View, resource2.Ctrl)
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
			app2.Event.Publish(event.DeviceCreated, data.Get("userID"), data.Get("deviceID"))
			return data.Pop("result")
		}

		return result
	}

	return response2.Wrap(fn())
}

func Detail(deviceID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		device, err := app2.Store().GetDevice(deviceID)
		if err != nil {
			return err
		}

		admin := app2.Store().MustGetUserFromContext(ctx)
		if !app2.Allow(admin, device, resource2.View) {
			return lang2.ErrNoPermission
		}

		return device.Detail()
	})
}

func Update(deviceID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		var form struct {
			Title    *string  `json:"title"`
			ConnStr  *string  `json:"params.connStr"`
			Interval *int64   `json:"params.interval"`
			Groups   *[]int64 `json:"groups"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang2.ErrInvalidRequestData
		}

		result := app2.TransactionDo(func(s store2.Store) interface{} {
			device, err := s.GetDevice(deviceID)
			if err != nil {
				return err
			}

			admin := s.MustGetUserFromContext(ctx)
			if !app2.Allow(admin, device, resource2.Ctrl) {
				return lang2.ErrNoPermission
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
			app2.Event.Publish(event.DeviceUpdated, data.GetMulti("userID", "deviceID")...)
			return lang2.Ok
		}

		return result
	})
}

func Delete(deviceID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		result := app2.TransactionDo(func(s store2.Store) interface{} {
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
			if app2.IsDefaultAdminUser(admin) {
				err = device.Destroy()
				if err != nil {
					return err
				}
			} else {
				if !app2.Allow(admin, device, resource2.Ctrl) {
					return lang2.ErrNoPermission
				}
				err = app2.SetDeny(admin, device, resource2.View, resource2.Ctrl)
				if err != nil {
					return err
				}
			}
			return data
		})

		if data, ok := result.(event.Data); ok {
			app2.Event.Publish(event.DeviceDeleted, data.GetMulti("userID", "id", "uid", "title")...)
			return lang2.Ok
		}
		return result
	})
}
