package brief

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/gate/web/helper"
	"github.com/maritimusj/centrum/gate/web/response"
	"github.com/maritimusj/centrum/global"
	"github.com/maritimusj/centrum/version"
)

func Simple(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		devices, total, err := app.Store().GetDeviceList()
		if err != nil {
			return err
		}

		var statsMap = make(map[int]int)
		for _, device := range devices {
			stats, _, _ := global.GetDeviceStatus(device)
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

		opts := []helper.OptionFN{
			helper.Limit(1),
		}

		var (
			s     = app.Store()
			admin = s.MustGetUserFromContext(ctx)
		)

		if !app.IsDefaultAdminUser(admin) {
			opts = append(opts, helper.User(admin.GetID()))
		}

		_, total, err = app.Store().GetAlarmList(nil, nil, opts...)
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

		result["version"] = iris.Map{
			"edge": iris.Map{
				"ver":   version.EdgeVersion,
				"build": version.EdgeBuildDate,
			},
			"gate": iris.Map{
				"ver":   version.GateVersion,
				"build": version.GeteBuildDate,
			},
		}
		return result
	})
}
