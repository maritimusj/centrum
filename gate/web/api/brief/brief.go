package brief

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/gate/web/helper"
	"github.com/maritimusj/centrum/gate/web/response"
	"github.com/maritimusj/centrum/global"
)

func Simple(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		devices, total, err := app.Store().GetDeviceList()
		if err != nil {
			return err
		}

		var statsMap = make(map[int]int)
		for _, device := range devices {
			stats, _ := global.GetDeviceStatus(device)
			if v, ok := statsMap[stats]; ok {
				statsMap[stats] = v + 1
			} else {
				statsMap[stats] = 1
			}
		}

		result := iris.Map{
			"device": iris.Map{
				"total": total,
				"stats": statsMap,
			},
		}

		_, total, err = app.Store().GetAlarmList(nil, nil, helper.Limit(1))
		if err != nil {
			return err
		}

		_, total, err = app.Store().GetLastUnconfirmedAlarm()
		if err != nil {
			if err != lang.Error(lang.ErrAlarmNotFound) {
				return err
			}

		}

		result["alarm"] = iris.Map{
			"total": total,
		}

		return result
	})
}
