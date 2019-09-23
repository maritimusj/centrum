package my

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/web/app"
	"github.com/maritimusj/centrum/web/model"
	"github.com/maritimusj/centrum/web/resource"
	"github.com/maritimusj/centrum/web/response"
	"strconv"
)

func Detail(ctx iris.Context) hero.Result {
	s := app.Store()
	my := s.MustGetUserFromContext(ctx)

	return response.Wrap(func() interface{} {
		return my.Detail()
	})
}

func Update(ctx iris.Context) hero.Result {
	s := app.Store()
	my := s.MustGetUserFromContext(ctx)

	return response.Wrap(func() interface{} {
		var form struct {
			Password *string `json:"password"`
			Title    *string `json:"title"`
			Mobile   *string `json:"mobile"`
			Email    *string `json:"email"`
		}

		err := ctx.ReadJSON(&form)
		if err != nil {
			return lang.ErrInvalidRequestData
		}

		if form.Password != nil && *form.Password != "" {
			my.ResetPassword(*form.Password)
		}

		var data = model.Map{}
		if form.Title != nil {
			data["title"] = *form.Title
		}
		if form.Mobile != nil {
			data["mobile"] = *form.Mobile
		}
		if form.Email != nil {
			data["email"] = *form.Email
		}

		if len(data) > 0 {
			my.Update(data)
		}

		err = my.Save()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}

func Perm(class string, ctx iris.Context) hero.Result {
	s := app.Store()
	my := s.MustGetUserFromContext(ctx)

	return response.Wrap(func() interface{} {
		var res model.Resource
		var err error
		switch class {
		case "api":
			id := ctx.URLParam("id")
			if id != "" {
				res, err = s.GetApiResource(id)
			} else {
				err = lang.Error(lang.ErrInvalidRequestData)
			}
		case "group":
			res, err = s.GetResource(resource.Group, ctx.URLParamInt64Default("id", 0))
		case "device":
			res, err = s.GetResource(resource.Device, ctx.URLParamInt64Default("id", 0))
		case "measure":
			res, err = s.GetResource(resource.Measure, ctx.URLParamInt64Default("id", 0))
		case "equipment":
			res, err = s.GetResource(resource.Equipment, ctx.URLParamInt64Default("id", 0))
		case "state":
			res, err = s.GetResource(resource.State, ctx.URLParamInt64Default("id", 0))
		default:
			err = lang.Error(lang.ErrInvalidRequestData)
		}

		if err != nil {
			return err
		}

		if class == "api" {
			return iris.Map{
				"invoke": app.Allow(my, res, resource.Invoke),
			}
		}
		return iris.Map{
			"view": app.Allow(my, res, resource.View),
			"ctrl": app.Allow(my, res, resource.Ctrl),
		}
	})
}

func MultiPerm(class string, ctx iris.Context) hero.Result {
	s := app.Store()
	my := s.MustGetUserFromContext(ctx)

	return response.Wrap(func() interface{} {
		var form struct {
			Names []string `json:"names"`
			IDs   []int64  `json:"res"`
		}

		if err := ctx.ReadJSON(&form); err != nil || resource.IsValidClass(class) {
			return lang.ErrInvalidRequestData
		}

		var perms = iris.Map{}
		switch class {
		case "api":
			for _, name := range form.Names {
				res, err := s.GetApiResource(name)
				if err != nil {
					return err
				}
				perms[name] = app.Allow(my, res, resource.Invoke)
			}
			return perms
		case "group", "device", "measure", "equipment", "state":
			for _, id := range form.IDs {
				res, err := s.GetResource(resource.ParseClass(class), id)
				if err != nil {
					return err
				}
				perms[strconv.FormatInt(id, 10)] = iris.Map{
					"view": app.Allow(my, res, resource.View),
					"ctrl": app.Allow(my, res, resource.Ctrl),
				}
			}
		default:
			return lang.ErrInvalidRequestData
		}
		return perms
	})
}
