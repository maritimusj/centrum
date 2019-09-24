package web

import (
	"encoding/json"
	"github.com/kataras/iris"
	log "github.com/sirupsen/logrus"

	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/web/app"
)

func GetLogList(ctx iris.Context, src string) interface{} {
	level := ctx.URLParam("level")
	start := ctx.URLParamInt64Default("start", 0)
	page := ctx.URLParamInt64Default("page", 1)
	pageSize := ctx.URLParamInt64Default("pagesize", app.Config.DefaultPageSize())

	x := uint64(start)
	logs, total, err := app.LogDBStore.Get(src, level, &x, uint64((page-1)*pageSize), uint64(pageSize))
	if err != nil {
		return err
	}

	var result = make([]iris.Map, 0, len(logs))
	for _, v := range logs {
		r := map[string]interface{}{}
		err := json.Unmarshal(v.Content, &r)

		if err != nil {
			result = append(result, iris.Map{
				"id":  v.ID,
				"raw": string(v.Content),
			})
		} else {
			result = append(result, iris.Map{
				"id":      v.ID,
				"content": r,
			})
		}
	}
	return iris.Map{
		"stats": app.LogDBStore.Stats(),
		"start": x,
		"total": total,
		"list":  result,
	}
}

func DeleteLog(ctx iris.Context, src string) interface{} {
	err := app.LogDBStore.Delete(src)
	if err != nil {
		return err
	}

	admin := app.Store().MustGetUserFromContext(ctx)
	log.Info(lang.Str(lang.LogDeletedByUser, admin.Name()))
	return lang.Ok
}
