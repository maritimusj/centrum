package resource

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/config"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/resource"
	"github.com/maritimusj/centrum/store"
	"github.com/maritimusj/centrum/web/response"
)

func GroupList() hero.Result {
	return response.Wrap(func() interface{} {
		return resource.GetGroupList()
	})
}

func List(groupID int, ctx iris.Context, s store.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", cfg.DefaultPageSize())

		resources, total, err := s.GetResourceList(model.ResourceClass(groupID), store.Page(page, pageSize))
		if err != nil {
			return err
		}

		var result = make([]model.Map, 0, len(resources))
		for _, res := range resources {
			group, id := res.GetResourceID()
			result = append(result, model.Map{
				"id":          id,
				"group":       group,
				"group_title": resource.ClassTitle(group),
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
