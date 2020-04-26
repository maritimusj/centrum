package statistics

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
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
			exportFiles.Delete(uid)
			return
		}
	}

	ctx.StatusCode(iris.StatusNotFound)
}

type Int64Slice []int64

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

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

		//写入UTF-8 BOM，防止中文乱码
		_, _ = csvFile.Write([]byte{0xEF, 0xBB, 0xBF})

		uid := util.RandStr(16, util.RandNum)
		stats := &exportStats{
			Msg:            lang.ExportInitialized.Str(),
			LocalFileName:  csvFile.Name(),
			ExportFileName: time.Now().Format("2006-01-02_15_04_05") + ".csv",
		}

		exportFiles.Store(uid, stats)

		go func() {

			var (
				s     = app.Store()
				admin = s.MustGetUserFromContext(ctx)
			)

			rangeMeasuresFN := func(fn func(device model.Device, measure model.Measure) error) error {
				if fn != nil {
					var (
						device  model.Device
						measure model.Measure
					)

					for _, measureID := range form.MeasureIDs {
						device = nil
						measure = nil

						measure, _ = s.GetMeasure(measureID)
						if measure != nil {
							if app.Allow(admin, measure, resource.View) {
								device = measure.Device()
							}
						}

						err = fn(device, measure)
						if err != nil {
							//忽略错误
							//return err
						}
					}
				}

				return nil
			}

			rangeStatesFN := func(fn func(equipment model.Equipment, state model.State) error) error {
				var (
					state     model.State
					equipment model.Equipment
				)

				for _, stateID := range form.StatesIDs {
					equipment = nil
					state = nil

					state, _ = s.GetState(stateID)
					if state != nil {
						if app.Allow(admin, state, resource.View) {
							equipment = state.Equipment()
						}
					}

					err = fn(equipment, state)
					if err != nil {
						return err
					}
				}
				return nil
			}

			var (
				measureValues = make(map[int64]map[string]string)
				timeValues    = make(Int64Slice, 0)
			)

			getMeasureDataFN := func(device model.Device, measure model.Measure, title string) error {
				if device != nil && measure != nil {
					stats.Msg = lang.ExportingData.Str(device.Title(), measure.Title())

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
						timestamp := sec
						var val string
						switch v := data[1].(type) {
						case json.Number:
							val = v.String()
						case bool:
							val = util.If(v, "1", "0").(string)
						case nil:
							val = ""
						default:
							val = fmt.Sprintf("%#v", data[1])
						}
						if _, ok := measureValues[timestamp]; !ok {
							measureValues[timestamp] = map[string]string{}
							measureValues[timestamp][title] = val
							timeValues = append(timeValues, timestamp)
						} else {
							measureValues[timestamp][title] = val
						}
					}
				}

				return nil
			}

			var headers []string
			getDeviceHeaderTitle := func(device model.Device, measure model.Measure) string {
				var (
					deviceTitle  string
					measureTitle string
					idStr        string
				)
				if device == nil {
					deviceTitle = "unknown"
					idStr = "0:"
				} else {
					deviceTitle = device.Title()
					idStr = strconv.FormatInt(device.GetID(), 10) + ":"
				}
				if measure == nil {
					measureTitle = "unknown"
					idStr += "0"
				} else {
					measureTitle = measure.Title()
					idStr += strconv.FormatInt(measure.GetID(), 10)
				}

				return deviceTitle + "[" + measureTitle + "]" + "#!" + idStr
			}

			getEquipmentHeaderTitle := func(equipment model.Equipment, state model.State) string {
				var (
					equipmentTitle string
					stateTitle     string
					idStr          string
				)
				if equipment == nil {
					equipmentTitle = "unknown"
					idStr = "0:"
				} else {
					equipmentTitle = equipment.Title()
					idStr = strconv.FormatInt(equipment.GetID(), 10) + ":"
				}
				if state == nil {
					stateTitle = "unknown"
					idStr = "0:"
				} else {
					stateTitle = state.Title()
					idStr += strconv.FormatInt(state.GetID(), 10)
				}
				return equipmentTitle + "(" + stateTitle + ")" + "#!" + idStr
			}

			_ = rangeMeasuresFN(func(device model.Device, measure model.Measure) error {
				headers = append(headers, getDeviceHeaderTitle(device, measure))
				return nil
			})

			_ = rangeStatesFN(func(equipment model.Equipment, state model.State) error {
				headers = append(headers, getEquipmentHeaderTitle(equipment, state))
				return nil
			})

			csvWriter := csv.NewWriter(csvFile)

			formattedHeaders := []string{"#"}
			for _, header := range headers {
				arr := strings.SplitN(header, "#!", 2)
				formattedHeaders = append(formattedHeaders, arr[0])
			}

			_ = csvWriter.Write(formattedHeaders)

			_ = rangeMeasuresFN(func(device model.Device, measure model.Measure) error {
				return getMeasureDataFN(device, measure, getDeviceHeaderTitle(device, measure))
			})

			_ = rangeStatesFN(func(equipment model.Equipment, state model.State) error {
				measure := state.Measure()
				if measure != nil {
					device := measure.Device()
					if device != nil {
						return getMeasureDataFN(device, measure, getEquipmentHeaderTitle(equipment, state))
					}
				}

				//忽略错误
				return nil
			})

			stats.Msg = lang.ArrangingData.Str()
			sort.Sort(timeValues)

			total := len(timeValues)
			if total > 0 {
				for i, index := range timeValues {
					stats.Msg = lang.WritingData.Str(int((float32(i+1) / float32(total)) * 100))

					valuesMap := measureValues[index]
					ts := time.Unix(index, 0)
					valueSlice := []string{ts.Format(lang.DatetimeFormatterStr.Str())}

					for _, header := range headers {
						if v, ok := valuesMap[header]; ok {
							valueSlice = append(valueSlice, v)
						} else {
							valueSlice = append(valueSlice, "")
						}
					}

					err = csvWriter.Write(valueSlice)
					if err != nil {
						continue
					}
				}
			}

			csvWriter.Flush()
			_ = csvFile.Close()

			stats.Msg = lang.ExportReady.Str()
			stats.IsOk = true
		}()

		return iris.Map{
			"uid": uid,
		}
	})
}
