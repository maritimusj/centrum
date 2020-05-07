package response

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/lang"
)

type responseData struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}

func (data *responseData) Dispatch(c iris.Context) {
	_, _ = c.JSON(data)
}

func Wrap(data interface{}) hero.Result {
	switch v := data.(type) {
	case lang.ErrIndex:
		return &responseData{
			Status: v == lang.Ok,
			Data: map[string]interface{}{
				"msg": lang.ErrorStr(v),
			},
		}
	case error:
		return &responseData{
			Status: false,
			Data: map[string]interface{}{
				"msg": v.Error(),
			},
		}
	case string:
		return &responseData{
			Status: true,
			Data: map[string]interface{}{
				"msg": v,
			},
		}
	case func() interface{}:
		return Wrap(v())
	default:
		return &responseData{
			Status: true,
			Data:   v,
		}
	}
}
