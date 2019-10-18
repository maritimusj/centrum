package edge

import (
	"fmt"
	"github.com/maritimusj/centrum/global"
	"github.com/maritimusj/centrum/json_rpc"
	"github.com/maritimusj/centrum/web/model"
	"strconv"
	"time"
)

func ActiveDevice(device model.Device) error {
	org, err := device.Organization()
	if err != nil {
		return err
	}

	conf := &json_rpc.Conf{
		UID:              strconv.FormatInt(device.GetID(), 10),
		Address:          device.GetOption("params.connStr").Str,
		Interval:         time.Second * time.Duration(device.GetOption("params.interval").Int()),
		DB:               org.Title(),
		InfluxDBAddress:  "http://localhost:8086",
		InfluxDBUserName: "",
		InfluxDBPassword: "",
		CallbackURL:      fmt.Sprintf("%s/%d", global.Params.MustGet("callbackURL"), device.GetID()),
		LogLevel:         "trace",
	}

	return Active(conf)
}

func ResetConfig(device model.Device) {
	Reset(strconv.FormatInt(device.GetID(), 10))
}

func GetStatus(device model.Device) (map[string]interface{}, error) {
	return GetBaseInfo(strconv.FormatInt(device.GetID(), 10))
}

func GetData(device model.Device) ([]interface{}, error) {
	return GetRealtimeData(strconv.FormatInt(device.GetID(), 10))
}

func SetCHValue(device model.Device, chTagName string, v interface{}) error {
	return SetValue(strconv.FormatInt(device.GetID(), 10), chTagName, v)
}

func GetCHValue(device model.Device, chTagName string) (map[string]interface{}, error) {
	return GetValue(strconv.FormatInt(device.GetID(), 10), chTagName)
}
