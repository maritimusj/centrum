package statistics

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/influxdata/influxdb1-client/models"
	"github.com/kataras/iris"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/gate/web/resource"
	"time"

	log "github.com/sirupsen/logrus"
)

type AxisMap map[string]byte

func (m AxisMap) Set(name string, start byte) {
	m[name] = start
}

func (m AxisMap) Get(name string) byte {
	if v, ok := m[name]; ok {
		return v
	}
	return 'A'
}
func (m AxisMap) Next(name string) byte {
	if v, ok := m[name]; ok {
		m[name] = v + 1
		return v + 1
	} else {
		m[name] = 'A'
		return m.Next(name)
	}
}

func Export(ctx iris.Context) {
	res := func() interface{} {
		var form struct {
			MeasureIDs []int64    `json:"measures"`
			StatesIDs  []int64    `json:"states"`
			Start      time.Time  `json:"start"`
			End        *time.Time `json:"end"`
		}
		if err := ctx.ReadJSON(&form); err != nil {
			return err
		}

		var (
			axisMap = AxisMap{}

			s             = app.Store()
			admin         = s.MustGetUserFromContext(ctx)
			excel         = excelize.NewFile()
			alarmStyle, _ = excel.NewStyle(`{"font":{"color":"#f44336"}}`)
		)

		exportMeasureFN := func(sheetName string, rows *models.Row) {
			excel.SetActiveSheet(excel.NewSheet(sheetName))
			excel.SetColWidth(sheetName, "A", "A", 20)

			col := axisMap.Next(sheetName)
			excel.SetCellValue(sheetName, fmt.Sprintf("%c1", col), rows.Name)

			for i, data := range rows.Values {
				cell := fmt.Sprintf("%c%d", col, i+2)
				sec, _ := data[0].(json.Number).Int64()
				excel.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), time.Unix(sec, 0))
				switch v := data[1].(type) {
				case json.Number:
					val, _ := v.Float64()
					excel.SetCellValue(sheetName, cell, val)
				case bool:
					excel.SetCellBool(sheetName, cell, v)
				}

				if data[2] != nil {
					excel.SetCellStyle(sheetName, cell, cell, alarmStyle)
				}
			}
		}

		for _, measureID := range form.MeasureIDs {
			measure, err := s.GetMeasure(measureID)
			if err != nil {
				return err
			}

			if !app.Allow(admin, measure, resource.View) {
				return lang.Error(lang.ErrNoPermission)
			}

			device := measure.Device()
			if device == nil {
				return lang.Error(lang.ErrDeviceNotFound)
			}

			org, err := device.Organization()
			if err != nil {
				return err
			}

			rows, err := app.StatsDB.GetMeasureStats(org.Name(), device.GetID(), measure.TagName(), &form.Start, form.End, 0)
			if err != nil {
				return err
			}

			sheetName := device.Title()
			exportMeasureFN(sheetName, rows)
		}

		for _, stateID := range form.StatesIDs {
			state, err := s.GetState(stateID)
			if err != nil {
				return err
			}

			if !app.Allow(admin, state, resource.View) {
				return lang.Error(lang.ErrNoPermission)
			}

			equipment := state.Equipment()
			if equipment == nil {
				return lang.Error(lang.ErrEquipmentNotFound)
			}

			measure := state.Measure()
			if measure == nil {
				return lang.Error(lang.ErrMeasureNotFound)
			}

			device := measure.Device()
			if device == nil {
				return lang.Error(lang.ErrDeviceNotFound)
			}

			org, err := equipment.Organization()
			if err != nil {
				return err
			}

			rows, err := app.StatsDB.GetMeasureStats(org.Name(), device.GetID(), measure.TagName(), &form.Start, form.End, 0)
			if err != nil {
				return err
			}

			sheetName := equipment.Title()
			exportMeasureFN(sheetName, rows)
		}

		excel.DeleteSheet("Sheet1")
		err := excel.Write(ctx)
		if err != nil {
			log.Error(err)
			return err
		}

		return nil
	}()

	if err, ok := res.(error); ok {
		excel := excelize.NewFile()
		excel.SetCellValue("Sheet1", "A1", err.Error())
		_ = excel.Write(ctx)
	}
}
