package my

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	app2 "github.com/maritimusj/centrum/gate/web/app"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	response2 "github.com/maritimusj/centrum/gate/web/response"
	store2 "github.com/maritimusj/centrum/gate/web/store"
	"strconv"
)

func Detail(ctx iris.Context) hero.Result {
	my := app2.Store().MustGetUserFromContext(ctx)
	return response2.Wrap(func() interface{} {
		return my.Detail()
	})
}

func Update(ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		var form struct {
			Password *string `json:"password"`
			Title    *string `json:"title"`
			Mobile   *string `json:"mobile"`
			Email    *string `json:"email"`
		}

		err := ctx.ReadJSON(&form)
		if err != nil {
			return lang2.ErrInvalidRequestData
		}

		return app2.TransactionDo(func(s store2.Store) interface{} {
			my := s.MustGetUserFromContext(ctx)
			if form.Password != nil && *form.Password != "" {
				my.ResetPassword(*form.Password)
			}

			var data = model2.Map{}
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
			return lang2.Ok
		})
	})
}

func Perm(class string, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		s := app2.Store()

		var res model2.Resource
		var err error
		switch class {
		case "api":
			id := ctx.URLParam("id")
			if id != "" {
				res, err = s.GetApiResource(id)
			} else {
				err = lang2.Error(lang2.ErrInvalidRequestData)
			}
		case "group":
			res, err = s.GetResource(resource2.Group, ctx.URLParamInt64Default("id", 0))
		case "device":
			res, err = s.GetResource(resource2.Device, ctx.URLParamInt64Default("id", 0))
		case "measure":
			res, err = s.GetResource(resource2.Measure, ctx.URLParamInt64Default("id", 0))
		case "equipment":
			res, err = s.GetResource(resource2.Equipment, ctx.URLParamInt64Default("id", 0))
		case "state":
			res, err = s.GetResource(resource2.State, ctx.URLParamInt64Default("id", 0))
		default:
			err = lang2.Error(lang2.ErrInvalidRequestData)
		}

		if err != nil {
			return err
		}

		my := s.MustGetUserFromContext(ctx)

		if class == "api" {
			return iris.Map{
				"invoke": app2.Allow(my, res, resource2.Invoke),
			}
		}
		return iris.Map{
			"view": app2.Allow(my, res, resource2.View),
			"ctrl": app2.Allow(my, res, resource2.Ctrl),
		}
	})
}

func MultiPerm(class string, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		var form struct {
			Names []string `json:"names"`
			IDs   []int64  `json:"res"`
		}

		if err := ctx.ReadJSON(&form); err != nil || resource2.IsValidClass(class) {
			return lang2.ErrInvalidRequestData
		}

		s := app2.Store()
		my := s.MustGetUserFromContext(ctx)

		var perms = iris.Map{}
		switch class {
		case "api":
			for _, name := range form.Names {
				res, err := s.GetApiResource(name)
				if err != nil {
					return err
				}
				perms[name] = app2.Allow(my, res, resource2.Invoke)
			}
			return perms
		case "group", "device", "measure", "equipment", "state":
			for _, id := range form.IDs {
				res, err := s.GetResource(resource2.ParseClass(class), id)
				if err != nil {
					return err
				}
				perms[strconv.FormatInt(id, 10)] = iris.Map{
					"view": app2.Allow(my, res, resource2.View),
					"ctrl": app2.Allow(my, res, resource2.Ctrl),
				}
			}
		default:
			return lang2.ErrInvalidRequestData
		}
		return perms
	})
}
