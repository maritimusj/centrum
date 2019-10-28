package log

import (
	"encoding/json"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/logStore"
	app2 "github.com/maritimusj/centrum/gate/web/app"
	response2 "github.com/maritimusj/centrum/gate/web/response"
	log "github.com/sirupsen/logrus"
)

func Level() hero.Result {
	return response2.Wrap(func() interface{} {
		return []iris.Map{
			{
				"id":    "trace",
				"title": lang2.Str(lang2.LogTrace),
			},
			{
				"id":    "debug",
				"title": lang2.Str(lang2.LogDebug),
			},
			{
				"id":    "info",
				"title": lang2.Str(lang2.LogInfo),
			},
			{
				"id":    "warn",
				"title": lang2.Str(lang2.LogWarn),
			},
			{
				"id":    "error",
				"title": lang2.Str(lang2.LogError),
			},
			{
				"id":    "fatal",
				"title": lang2.Str(lang2.LogFatal),
			},
			{
				"id":    "panic",
				"title": lang2.Str(lang2.LogPanic),
			},
		}
	})
}

func List(ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		admin := app2.Store().MustGetUserFromContext(ctx)
		var orgID int64
		if app2.IsDefaultAdminUser(admin) && ctx.URLParamExists("org") {
			orgID = ctx.URLParamInt64Default("org", admin.OrganizationID())
		}
		return GetLogList(ctx, orgID, logStore.SystemLog)
	})
}

func Delete(ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		admin := app2.Store().MustGetUserFromContext(ctx)
		var orgID int64
		if app2.IsDefaultAdminUser(admin) && ctx.URLParamExists("org") {
			orgID = ctx.URLParamInt64Default("org", admin.OrganizationID())
		}
		return DeleteLog(ctx, orgID, logStore.SystemLog)
	})
}

func GetLogList(ctx iris.Context, orgID int64, src string) interface{} {
	level := ctx.URLParam("level")
	start := ctx.URLParamInt64Default("start", 0)
	page := ctx.URLParamInt64Default("page", 1)
	pageSize := ctx.URLParamInt64Default("pagesize", app2.Config.DefaultPageSize())

	admin := app2.Store().MustGetUserFromContext(ctx)

	x := uint64(start)
	logs, total, err := app2.LogDBStore.GetList(orgID, src, level, &x, uint64((page-1)*pageSize), uint64(pageSize))
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
		"stats": app2.LogDBStore.Stats(admin.OrganizationID()),
		"start": x,
		"total": total,
		"list":  result,
	}
}

func DeleteLog(ctx iris.Context, orgID int64, src string) interface{} {
	admin := app2.Store().MustGetUserFromContext(ctx)

	err := app2.LogDBStore.Delete(orgID, src)
	if err != nil {
		return err
	}

	log.Info(lang2.Str(lang2.LogDeletedByUser, admin.Name()))
	return lang2.Ok
}
