package statistics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/kataras/iris/hero"
	"github.com/maritimusj/centrum/gate/web/response"

	"github.com/kataras/iris"
	"github.com/maritimusj/centrum/gate/lang"
	"github.com/maritimusj/centrum/gate/web/app"
	"github.com/maritimusj/centrum/gate/web/model"
	"github.com/maritimusj/centrum/gate/web/resource"
	"github.com/maritimusj/centrum/util"
)

var (
	exportFiles = sync.Map{}
)

type exportStats struct {
	Msg            string
	LocalFileName  string
	ExportFileName string
	IsOk           bool
	Error          error
}

func ExportStats(uid string) hero.Result {
	return response.Wrap(func() interface{} {
		if v, ok := exportFiles.Load(uid); ok {
			stats := v.(*exportStats)
			if stats.Error != nil {
				return stats.Error
			}
			return iris.Map{
				"ok":   stats.IsOk,
				"text": stats.Msg,
			}
		}
		return lang.ErrExportNotExists
	})
}

func ExportDownload(uid string, ctx iris.Context) {
	if v, ok := exportFiles.Load(uid); ok {
		stats := v.(*exportStats)
		if stats.IsOk && stats.Error == nil {
			_ = ctx.SendFile(stats.LocalFileName, stats.ExportFileName)
			_ = os.Remove(stats.LocalFileName)
			return
		}
	}

	ctx.StatusCode(iris.StatusNotFound)
}

func Export(ctx iris.Context) hero.Result {
	return response.Wrap(func() interface{} {
		var form struct {
			MeasureIDs []int64    `json:"measures"`
			StatesIDs  []int64    `json:"states"`
			Start      time.Time  `json:"start"`
			End        *time.Time `json:"end"`
			Interval   *string    `json:"interval"`
		}

		if err := ctx.ReadJSON(&form); err != nil {
			return lang.ErrInvalidRequestData
		}

		csvFile, err := ioutil.TempFile("", "tempFile")
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			return lang.InternalError(err)
		}

		uid := util.RandStr(16, util.RandNum)
		stats := &exportStats{
			Msg:            lang.Str(lang.ExportInitialized),
			LocalFileName:  csvFile.Name(),
			ExportFileName: time.Now().Format("2006-01-02_15_04_05") + ".csv",
		}

		exportFiles.Store(uid, stats)

		go func() {
			//写入UTF-8 BOM，防止中文乱码
			_, _ = csvFile.WriteString("\xEF\xBB\xBF")

			var (
				s     = app.Store()
				admin = s.MustGetUserFromContext(ctx)
			)

			rangeMeasuresFN := func(fn func(device model.Device, measure model.Measure) error) error {
				for _, measureID := range form.MeasureIDs {
					measure, err := s.GetMeasure(measureID)
					if err != nil {
						//return err
						continue
					}

					if !app.Allow(admin, measure, resource.View) {
						//return lang.Error(lang.ErrNoPermission)
						continue
					}

					device := measure.Device()
					if device == nil {
						//return lang.Error(lang.ErrDeviceNotFound)
						continue
					}

					if fn != nil {
						err = fn(device, measure)
						if err != nil {
							//return err
							continue
						}
					}
				}
				return nil
			}

			rangeStatesFN := func(fn func(equipment model.Equipment, state model.State) error) error {
				for _, stateID := range form.StatesIDs {
					state, err := s.GetState(stateID)
					if err != nil {
						//return err
						continue
					}

					if !app.Allow(admin, state, resource.View) {
						//return lang.Error(lang.ErrNoPermission)
						continue
					}

					equipment := state.Equipment()
					if equipment == nil {
						//return lang.Error(lang.ErrEquipmentNotFound)
						continue
					}

					if fn != nil {
						err = fn(equipment, state)
						if err != nil {
							//return err
							continue
						}
					}
				}
				return nil
			}

			var (
				measureValues = make(map[int][]string)
				timeValues    = make([]int, 0)
			)

			getMeasureDataFN := func(device model.Device, measure model.Measure) error {
				stats.Msg = lang.Str(lang.ExportingData, device.Title(), measure.Title())

				org, err := device.Organization()
				if err != nil {
					return err
				}

				rows, err := app.StatsDB.GetMeasureStats(org.Name(), device.GetID(), measure.TagName(), &form.Start, form.End, form.Interval)
				if err != nil {
					return err
				}

				for _, data := range rows.Values {
					sec, _ := data[0].(json.Number).Int64()
					index := int(sec)
					var val string
					switch v := data[1].(type) {
					case json.Number:
						val = v.String()
					case bool:
						val = util.If(v, "1", "0").(string)
					default:
						val = "<unknown>"
					}
					if _, ok := measureValues[index]; !ok {
						measureValues[index] = []string{val}
						timeValues = append(timeValues, index)
					} else {
						measureValues[index] = append(measureValues[index], val)
					}
				}
				return nil
			}

			header := []string{""}
			stats.Error = rangeMeasuresFN(func(device model.Device, measure model.Measure) error {
				header = append(header, fmt.Sprintf("%s_%s", device.Title(), measure.Title()))
				return nil
			})
			if err != nil {
				stats.IsOk = true
				return
			}

			stats.Error = rangeStatesFN(func(equipment model.Equipment, state model.State) error {
				header = append(header, fmt.Sprintf("%s_%s", equipment.Title(), state.Title()))
				return nil
			})

			if err != nil {
				stats.IsOk = true
				return
			}

			_, err = csvFile.WriteString(strings.Join(header, ",") + "\r\n")
			if err != nil {
				stats.Error = lang.InternalError(err)
				stats.IsOk = true
				return
			}

			stats.Error = rangeMeasuresFN(func(device model.Device, measure model.Measure) error {
				return getMeasureDataFN(device, measure)
			})
			if err != nil {
				stats.IsOk = true
				return
			}

			err = rangeStatesFN(func(equipment model.Equipment, state model.State) error {
				measure := state.Measure()
				if measure == nil {
					return lang.Error(lang.ErrMeasureNotFound)
				}

				device := measure.Device()
				if device == nil {
					return lang.Error(lang.ErrDeviceNotFound)
				}

				return getMeasureDataFN(device, measure)
			})

			stats.Msg = lang.Str(lang.ArrangingData)
			sort.Ints(timeValues)

			total := len(timeValues)
			if total > 0 {
				for _, index := range timeValues {
					stats.Msg = lang.Str(lang.WritingData, int((float32(index)+1)/float32(total)*100))
					values := measureValues[index]
					ts := time.Unix(int64(index), 0)
					_, err = csvFile.WriteString(ts.Format("2006-01-02 15:04:05") + "," + strings.Join(values, ",") + "\r\n")
					if err != nil {
						continue
					}
				}
			}

			_ = csvFile.Close()
			stats.IsOk = true
		}()

		return iris.Map{
			"uid": uid,
		}
	})
}
