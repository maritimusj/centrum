package global

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	edgeLang "github.com/maritimusj/centrum/edge/lang"
	_ "github.com/maritimusj/centrum/edge/lang/enUS"
	_ "github.com/maritimusj/centrum/edge/lang/zhCN"

	"github.com/maritimusj/centrum/gate/web/model"
)

func AddMessage(data interface{}, fn func(string, int64) bool) {
	println("add message")
	Messages.Add(data, func(uid string) bool {
		if fn != nil {
			arr := strings.Split(uid, ":")
			if len(arr) != 2 {
				return false
			}
			userId, err := strconv.ParseInt(arr[1], 10, 0)
			if err != nil {
				return false
			}
			return fn(arr[0], userId)
		}

		return true
	})
}

func GetAllMessage(uid string, userId int64) []*msg {
	println("get message:", uid, userId)
	return Messages.GetAll(fmt.Sprintf("%s:%d", uid, userId))
}

func Create(uid string, userId int64) {
	println("create: ", uid, userId)
	Messages.Create(fmt.Sprintf("%s:%d", uid, userId))
}

func Close(uid string, userId int64) {
	println("close: ", uid, userId)
	Messages.Close(fmt.Sprintf("%s:%d", uid, userId))
}

func UpdateDeviceStatus(device model.Device, index int, title string) {
	path := fmt.Sprintf("device.%d.stats", device.GetID())
	data := map[string]interface{}{
		"index": index,
		"title": title,
	}
	orgIndex, _, from := GetDeviceStatus(device)
	if orgIndex != int(edgeLang.Connected) && index == int(edgeLang.Connected) {
		data["from"] = time.Now().Format(time.RFC3339)
	} else {
		data["from"] = from
	}

	if index != int(edgeLang.Connected) {
		UpdateDevicePerf(device, map[string]interface{}{})
	}
	_ = Stats.Set(path, data)
}

func GetDeviceStatus(device model.Device) (int, string, time.Time) {
	key := fmt.Sprintf("device.%d.stats", device.GetID())
	if v, ok := Stats.Get(key); ok {
		if vv, ok := v.(map[string]interface{}); ok {
			from, _ := time.Parse(time.RFC3339, vv["from"].(string))
			return int(vv["index"].(float64)), vv["title"].(string), from
		}
	}
	return int(edgeLang.EdgeUnknownState), edgeLang.Str(edgeLang.EdgeUnknownState), time.Now()
}

func UpdateDevicePerf(device model.Device, data map[string]interface{}) {
	path := fmt.Sprintf("device.%d.perf", device.GetID())
	_ = Stats.Set(path, data)
}

func GetDevicePerf(device model.Device) map[string]interface{} {
	key := fmt.Sprintf("device.%d.perf", device.GetID())
	if v, ok := Stats.Get(key); ok {
		if vv, ok := v.(map[string]interface{}); ok {
			return vv
		}
	}
	return map[string]interface{}{}
}
