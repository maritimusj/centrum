package statistics

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/gate/web/resource"
	"github.com/maritimusj/centrum/gate/web/response"
	"time"
)

func Measure(measureID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			Start    *time.Time    `json:"start"`
			End      *time.Time    `json:"end"`
			Interval time.Duration `json:"interval"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		s := app.Store()
		measure, err := s.GetMeasure(measureID)
		if err != nil {
			return err
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app.Allow(admin, measure, resource.View) {
			return lang.ErrNoPermission
		}

		device := measure.Device()
		if device == nil {
			return lang.ErrDeviceNotFound
		}

		org, _ := device.Organization()

		var start time.Time
		if form.Start != nil {
			start = *form.Start
		} else {
			start = time.Now().Add(-time.Hour * 1)
			form.Interval = 15
		}

		result, err := app.StatsDB.GetMeasureStats(org.Name(), device.GetID(), measure.TagName(), &start, form.End, form.Interval*time.Second)
		if err != nil {
			return lang.InternalError(err)
		}

		return result
	})
}

func State(stateID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			Start    *time.Time    `json:"start"`
			End      *time.Time    `json:"end"`
			Interval time.Duration `json:"interval"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		s := app.Store()
		state, err := s.GetState(stateID)
		if err != nil {
			return err
		}

		measure := state.Measure()
		if measure == nil {
			return lang.ErrMeasureNotFound
		}

		admin := s.MustGetUserFromContext(ctx)
		if !app.Allow(admin, measure, resource.View) {
			return lang.ErrNoPermission
		}

		device := measure.Device()
		if device == nil {
			return lang.ErrDeviceNotFound
		}

		org, _ := device.Organization()

		var start time.Time
		if form.Start != nil {
			start = *form.Start
		} else {
			start = time.Now().Add(-time.Hour * 1)
			form.Interval = 15
		}

		result, err := app.StatsDB.GetMeasureStats(org.Name(), device.GetID(), measure.TagName(), &start, form.End, form.Interval*time.Second)
		if err != nil {
			return lang.InternalError(err)
		}

		return result
	})
}
