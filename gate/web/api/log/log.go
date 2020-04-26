package log

import (
	"encoding/json"

	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/logStore"
	"github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/gate/web/response"
	log "github.com/sirupsen/logrus"
)

func Level() hero.Result {
	return response.Wrap(func() interface{} {
		return []iris.Map{
			{
				"id":    "trace",
				"title": lang.LogTrace.Str(),
			},
			{
				"id":    "debug",
				"title": lang.LogDebug.Str(),
			},
			{
				"id":    "info",
				"title": lang.LogInfo.Str(),
			},
			{
				"id":    "warning",
				"title": lang.LogWarning.Str(),
			},
			{
				"id":    "error",
				"title": lang.LogError.Str(),
			},
			{
				"id":    "fatal",
				"title": lang.LogFatal.Str(),
			},
			{
				"id":    "panic",
				"title": lang.LogPanic.Str(),
			},
		}
	})
}

func List(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		admin := app.Store().MustGetUserFromContext(ctx)
		var orgID int64
		if app.IsDefaultAdminUser(admin) && ctx.URLParamExists("org") {
			orgID = ctx.URLParamInt64Default("org", admin.OrganizationID())
		}
		return GetLogList(ctx, orgID, logStore.SystemLog)
	})
}

func Delete(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		admin := app.Store().MustGetUserFromContext(ctx)
		var orgID int64
		if app.IsDefaultAdminUser(admin) && ctx.URLParamExists("org") {
			orgID = ctx.URLParamInt64Default("org", admin.OrganizationID())
		}
		return DeleteLog(ctx, orgID, logStore.SystemLog)
	})
}

func GetLogList(ctx iris.Context, orgID int64, src string) interface{} {
	level := ctx.URLParam("level")
	start := ctx.URLParamInt64Default("start", 0)
	page := ctx.URLParamInt64Default("page", 1)
	pageSize := ctx.URLParamInt64Default("pagesize", app.Config.DefaultPageSize())

	admin := app.Store().MustGetUserFromContext(ctx)

	x := uint64(start)
	logs, total, err := app.LogDBStore.GetList(orgID, src, level, &x, uint64((page-1)*pageSize), uint64(pageSize))
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
		"stats": app.LogDBStore.Stats(admin.OrganizationID()),
		"start": x,
		"total": total,
		"list":  result,
	}
}

func DeleteLog(ctx iris.Context, orgID int64, src string) interface{} {
	admin := app.Store().MustGetUserFromContext(ctx)

	err := app.LogDBStore.Delete(orgID, src)
	if err != nil {
		return err
	}

	log.WithField("src", logStore.SystemLog).Info(lang.LogDeletedByUser.Str(admin.Name()))
	return lang.Ok
}
