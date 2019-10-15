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
		Inverse:          false,
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
