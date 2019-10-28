package statistics

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	lang2 "github.com/maritimusj/centrum/gate/lang"
	app2 "github.com/maritimusj/centrum/gate/web/app"
	resource2 "github.com/maritimusj/centrum/gate/web/resource"
	response2 "github.com/maritimusj/centrum/gate/web/response"
	"time"
)

func Measure(measureID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		var form struct {
			Start    *time.Time    `json:"start"`
			End      *time.Time    `json:"end"`
			Interval time.Duration `json:"interval"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang2.ErrInvalidRequestData
		}

		s := app2.Store()
		measure, err := s.GetMeasure(measureID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app2.Allow(admin, measure, resource2.View) {
			return lang2.ErrNoPermission
		}

		device := measure.Device()
		if device == nil {
			return lang2.ErrDeviceNotFound
		}

		org, _ := device.Organization()

		var start time.Time
		if form.Start != nil {
			start = *form.Start
		} else {
			start = time.Now().Add(-time.Hour * 1)
			form.Interval = 15
		}

		result, err := app2.StatsDB.GetMeasureStats(org.Name(), device.GetID(), measure.TagName(), &start, form.End, form.Interval*time.Second)
		if err != nil {
			return lang2.InternalError(err)
		}

		return result
	})
}

func State(stateID int64, ctx iris.Context) hero.Result {
	return response2.Wrap(func() interface{} {
		var form struct {
			Start    *time.Time    `json:"start"`
			End      *time.Time    `json:"end"`
			Interval time.Duration `json:"interval"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang2.ErrInvalidRequestData
		}

		s := app2.Store()
		state, err := s.GetState(stateID)
		if err != nil {
			return err
		}

		measure := state.Measure()
		if measure == nil {
			return lang2.ErrMeasureNotFound
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app2.Allow(admin, measure, resource2.View) {
			return lang2.ErrNoPermission
		}

		device := measure.Device()
		if device == nil {
			return lang2.ErrDeviceNotFound
		}

		org, _ := device.Organization()

		var start time.Time
		if form.Start != nil {
			start = *form.Start
		} else {
			start = time.Now().Add(-time.Hour * 1)
			form.Interval = 15
		}

		result, err := app2.StatsDB.GetMeasureStats(org.Name(), device.GetID(), measure.TagName(), &start, form.End, form.Interval*time.Second)
		if err != nil {
			return lang2.InternalError(err)
		}

		return result
	})
}
