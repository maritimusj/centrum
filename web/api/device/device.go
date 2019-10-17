package device

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/event"
	"github.com/maritimusj/centrum/global"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/web/app"
	"github.com/maritimusj/centrum/web/edge"
	"github.com/maritimusj/centrum/web/helper"
	"github.com/maritimusj/centrum/web/model"
	"github.com/maritimusj/centrum/web/resource"
	"github.com/maritimusj/centrum/web/response"
	"github.com/maritimusj/centrum/web/store"
	"gopkg.in/go-playground/validator.v9"
	"strconv"
	"strings"
)

func List(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var params []helper.OptionFN
		var orgID int64

		s := app.Store()
		admin := s.MustGetUserFromContext(ctx)
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

		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", app.Config.DefaultPageSize())
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

		var result = make([]model.Map, 0, len(devices))
		for _, device := range devices {
			brief := device.Brief()
			brief["perm"] = iris.Map{
				"view": true,
				"ctrl": app.Allow(admin, device, resource.Ctrl),
			}
			if baseInfo, err := edge.GetBaseInfo(strconv.FormatInt(device.GetID(), 10)); err != nil {
				index, title := global.GetDeviceStatus(device)
				brief["edge"] = iris.Map{
					"status": iris.Map{
						"index": index,
						"title": title,
					},
				}
			} else {
				brief["edge"] = baseInfo
			}

			result = append(result, brief)
		}

		return iris.Map{
			"total": total,
			"list":  result,
		}
	})
}

func Create(ctx iris.Context, validate *validator.Validate) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			OrgID    int64   `json:"org"`
			Title    string  `json:"title" validate:"required"`
			Groups   []int64 `json:"groups"`
			ConnStr  string  `json:"params.connStr" validate:"required"`
			Interval int64   `json:"params.interval" validate:"required"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if err := validate.Struct(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		if govalidator.IsIPv4(form.ConnStr) {
			form.ConnStr += ":502"
		} else if govalidator.IsIPv6(form.ConnStr) {
			form.ConnStr = fmt.Sprintf("[%s]:502", form.ConnStr)
		} else if govalidator.IsMAC(form.ConnStr) {
			form.ConnStr = strings.ToLower(form.ConnStr)
		}

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

		return device.Detail()
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
