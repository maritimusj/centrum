package edge

import (
	"fmt"
	"strconv"
	"time"

	"github.com/maritimusj/centrum/gate/config"

	"github.com/maritimusj/centrum/gate/web/model"
	"github.com/maritimusj/centrum/global"
	"github.com/maritimusj/centrum/json_rpc"
)

func ActiveDevice(device model.Device, config *config.Config) error {
	org, err := device.Organization()
	if err != nil {
		return err
	}

	influxDBConfig := config.InfluxDBConfig()
	conf := &json_rpc.Conf{
		UID:              strconv.FormatInt(device.GetID(), 10),
		Address:          device.GetOption("params.connStr").Str,
		Interval:         time.Second * time.Duration(device.GetOption("params.interval").Int()),
		DB:               org.Title(),
		InfluxDBUrl:      influxDBConfig["url"],
		InfluxDBUserName: influxDBConfig["username"],
		InfluxDBPassword: influxDBConfig["password"],
		CallbackURL:      fmt.Sprintf("%s/%d", global.Params.MustGet("callbackURL"), device.GetID()),
		LogLevel:         "error",
	}

	return Active(conf)
}

func ResetConfig(device model.Device) {
	Reset(strconv.FormatInt(device.GetID(), 10))
}

func GetStatus(device model.Device) (map[string]interface{}, error) {
	return GetBaseInfo(strconv.FormatInt(device.GetID(), 10))
}

func GetRealTimeData(device model.Device) (interface{}, error) {
	return GetRealtimeData(strconv.FormatInt(device.GetID(), 10))
}

func SetCHValue(device model.Device, chTagName string, v interface{}) error {
	return SetValue(strconv.FormatInt(device.GetID(), 10), chTagName, v)
}

func GetCHValue(device model.Device, chTagName string) (map[string]interface{}, error) {
	return GetValue(strconv.FormatInt(device.GetID(), 10), chTagName)
}
