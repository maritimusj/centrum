package log

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/web/response"
)

func Level() hero.Result{
	return response.Wrap(func() interface {}{
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

func List() {

}

func Delete() {

}