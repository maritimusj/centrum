package statistics

import (
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/gate/web/resource"
	"github.com/maritimusj/centrum/gate/web/response"
)

func Measure(measureID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			Start *time.Time `json:"start"`
			End   *time.Time `json:"end"`
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

		var (
			start time.Time
			end   time.Time
		)

		if form.Start != nil {
			start = *form.Start
		} else {
			start = time.Now().Add(-time.Hour * 1)
		}

		if form.End != nil {
			end = *form.End
		} else {
			end = time.Now()
		}

		//最多取10000个点位数据
		interval := int64(end.Sub(start).Seconds() / 10000)
		if interval < 1 {
			interval = 1
		}

		result, err := app.StatsDB.GetMeasureStats(org.Name(), device.GetID(), measure.TagName(), &start, form.End, time.Duration(interval)*time.Second)
		if err != nil {
			return iris.Map{}
		}

		return result
	})
}

func State(stateID int64, ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			Start *time.Time `json:"start"`
			End   *time.Time `json:"end"`
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
		var end time.Time
		if form.Start != nil {
			start = *form.Start
		} else {
			start = time.Now().Add(-time.Hour * 1)
		}
		if form.End != nil {
			end = *form.End
		} else {
			end = time.Now()
		}

		//最多取10000个点位数据
		interval := int64(end.Sub(start).Seconds() / 10000)
		if interval < 1 {
			interval = 1
		}

		result, err := app.StatsDB.GetMeasureStats(org.Name(), device.GetID(), measure.TagName(), &start, form.End, time.Duration(interval)*time.Second)
		if err != nil {
			return lang.InternalError(err)
		}

		return result
	})
}

func Alarm(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			Start     *time.Time `json:"start"`
			End       *time.Time `json:"end"`
			Device    *int64     `json:"device"`
			Equipment *int64     `json:"equipment"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		panic("not finished")
	})
}
