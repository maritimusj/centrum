package edge

import (
	"fmt"
	model2 "github.com/maritimusj/centrum/gate/web/model"
	"github.com/maritimusj/centrum/global"
	"github.com/maritimusj/centrum/json_rpc"
	"strconv"
	"time"
)

func ActiveDevice(device model2.Device) error {
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
		LogLevel:         "error",
	}

	return Active(conf)
}

func ResetConfig(device model2.Device) {
	Reset(strconv.FormatInt(device.GetID(), 10))
}

func GetStatus(device model2.Device) (map[string]interface{}, error) {
	return GetBaseInfo(strconv.FormatInt(device.GetID(), 10))
}

func GetData(device model2.Device) ([]interface{}, error) {
	return GetRealtimeData(strconv.FormatInt(device.GetID(), 10))
}

func SetCHValue(device model2.Device, chTagName string, v interface{}) error {
	return SetValue(strconv.FormatInt(device.GetID(), 10), chTagName, v)
}

func GetCHValue(device model2.Device, chTagName string) (map[string]interface{}, error) {
	return GetValue(strconv.FormatInt(device.GetID(), 10), chTagName)
}
