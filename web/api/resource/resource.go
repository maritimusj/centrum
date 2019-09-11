package resource

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/config"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/store"
	"github.com/maritimusj/centrum/web/response"
)

func GroupList(store store.Store) hero.Result {
	return response.Wrap(func() interface{} {
		return store.GetResourceGroupList()
	})
}

func List(classID int, ctx iris.Context, s store.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", cfg.DefaultPageSize())
		keyword := ctx.URLParam("keyword")

		var params = []store.OptionFN{
			store.Page(page, pageSize),
		}

		if keyword != "" {
			params = append(params, store.Keyword(keyword))
		}

		class := resource.Class(classID)
		if class == resource.Api {
			sub := ctx.URLParam("sub")
			if sub != "" {
				params = append(params, store.Name(sub))
			}
		} else {
			sub := ctx.URLParamInt64Default("sub", -1)
			if sub != -1 {
				switch class {
				case resource.Group:
					params = append(params, store.Parent(sub))
				case resource.Device:
					params = append(params, store.Group(sub))
				case resource.Measure:
					params = append(params, store.Device(sub))
				case resource.Equipment:
					params = append(params, store.Group(sub))
				case resource.State:
					params = append(params, store.Equipment(sub))
				}
			}
		}

		resources, total, err := s.GetResourceList(resource.Class(classID), params...)
		if err != nil {
			return err
		}

		var result = make([]model.Map, 0, len(resources))
		for _, res := range resources {
			result = append(result, model.Map{
				"class":       classID,
				"class_title": lang.ResourceClassTitle(class),
				"uid":         res.ResourceUID(),
				"title":       res.ResourceTitle(),
				"desc":        res.ResourceDesc(),
			})
		}

		return iris.Map{
			"total": total,
			"list":  result,
		}
	})
}
