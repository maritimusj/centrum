package log

import (
	"encoding/json"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/config"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/logStore"
	"github.com/maritimusj/centrum/web/perm"
	"github.com/maritimusj/centrum/web/response"
	log "github.com/sirupsen/logrus"
)

func Level() hero.Result {
	return response.Wrap(func() interface{} {
		return []iris.Map{
			{
				"id":    "trace",
				"title": lang.Str(lang.LogTrace),
			},
			{
				"id":    "debug",
				"title": lang.Str(lang.LogDebug),
			},
			{
				"id":    "info",
				"title": lang.Str(lang.LogInfo),
			},
			{
				"id":    "warn",
				"title": lang.Str(lang.LogWarn),
			},
			{
				"id":    "error",
				"title": lang.Str(lang.LogError),
			},
			{
				"id":    "fatal",
				"title": lang.Str(lang.LogFatal),
			},
			{
				"id":    "panic",
				"title": lang.Str(lang.LogPanic),
			},
		}
	})
}

func List(src string, ctx iris.Context, store logStore.Store, cfg config.Config) hero.Result {
	return response.Wrap(func() interface{} {
		level := ctx.URLParam("level")
		start := ctx.URLParamInt64Default("start", 0)
		page := ctx.URLParamInt64Default("page", 1)
		pageSize := ctx.URLParamInt64Default("pagesize", cfg.DefaultPageSize())

		x := uint64(start)
		logs, total, err := store.Get(src, level, &x, uint64((page-1)*pageSize), uint64(pageSize))
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
			"stats": store.Stats(),
			"start": x,
			"total": total,
			"list":  result,
		}
	})
}

func Delete(src string, ctx iris.Context, store logStore.Store) hero.Result {
	return response.Wrap(func() interface {}{
		err := store.Delete(src)
		if err != nil {
			return err
		}

		log.Info(lang.Str(lang.LogDeletedByUser, perm.AdminUser(ctx).Name()))
		return lang.Ok
	})
}
