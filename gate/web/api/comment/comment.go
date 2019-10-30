package comment

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/gate/web/helper"
	"github.com/maritimusj/centrum/gate/web/model"
	"github.com/maritimusj/centrum/gate/web/response"
)

func List(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var (
			s = app.Store()

			alarmID  = ctx.Params().GetInt64Default("alarm", 0)
			page     = ctx.URLParamInt64Default("page", 1)
			pageSize = ctx.URLParamInt64Default("pagesize", app.Config.DefaultPageSize())

			params = []helper.OptionFN{
				helper.Page(page, pageSize),
			}
		)

		alarm, err := s.GetAlarm(alarmID)
		if err != nil {
			return err
		}

		comments, total, err := s.GetCommentList(alarm.GetID(), params...)
		if err != nil {
			return err
		}
		var result = make([]model.Map, 0, len(comments))
		for _, comment := range comments {
			brief := comment.Brief()
			result = append(result, brief)
		}

		return iris.Map{
			"total": total,
			"list":  result,
		}
	})
}

func Create(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			AlarmID  int64  `json:"alarm"`
			ParentID *int64 `json:"reply_to"`
			Content  string `json:"content"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		var (
			s        = app.Store()
			admin    = s.MustGetUserFromContext(ctx)
			parentID int64
		)

		if form.ParentID != nil {
			parentID = *form.ParentID
		}

		comment, err := s.CreateComment(admin.GetID(), form.AlarmID, parentID, iris.Map{
			"content": form.Content,
			"ip":      ctx.RemoteAddr(),
		})
		if err != nil {
			return err
		}
		return comment.Brief()
	})
}

func Detail(commentID int64) hero.Result {
	return response.Wrap(func() interface{} {
		comment, err := app.Store().GetComment(commentID)
		if err != nil {
			return err
		}
		return comment.Detail()
	})
}

func Delete(commentID int64) hero.Result {
	return response.Wrap(func() interface{} {
		comment, err := app.Store().GetComment(commentID)
		if err != nil {
			return err
		}
		err = comment.Destroy()
		if err != nil {
			return err
		}
		return lang.Ok
	})
}
