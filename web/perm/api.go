package perm

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/model"
	"github.com/maritimusj/centrum/request"
	"github.com/maritimusj/centrum/store"
	"github.com/maritimusj/centrum/web/response"
)

func Check(ctx iris.Context, store store.Store) hero.Result {
	result := MustAdminOk(ctx, func(admin model.User) interface{} {
		router := ctx.GetCurrentRoute()
		req, err := request.NewApiRequest(store, router.Name(), router.Method())
		if err != nil {
			return err
		}
		return admin.IsAllowed(req)
	})

	if err := result.(error); err != nil {
		ctx.StopExecution()
		return response.Wrap(result)
	}

	ctx.Next()
	return nil
}
