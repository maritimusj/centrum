package zhCN

import "github.com/maritimusj/centrum/edge/lang"

func init() {
	lang.Register(lang.ZhCN, strMap, errStrMap)
}

var (
	strMap = map[lang.StrIndex]string{
		lang.AdapterInitializing: "正在初始化",
		lang.Connecting:          "正在连接",
		lang.Connected:           "已连接",
		lang.Disconnected:        "已断开",
		lang.MalFunctioned:       "故障",
		lang.InfluxDBError:       "数据库错误",
	}

	errStrMap = map[lang.ErrorCode]string{
		lang.Ok:                    "成功！",
		lang.ErrDeviceNotExists:    "设备不存在！",
		lang.ErrDeviceNotConnected: "设备没有连接！",
	}
)
