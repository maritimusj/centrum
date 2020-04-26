package enUS

import "github.com/maritimusj/centrum/edge/lang"

func init() {
	lang.Register(lang.EnUS, strMap, errStrMap)
}

var (
	strMap = map[lang.StrIndex]string{
		lang.EdgeUnknownState:    "Unknown",
		lang.AdapterInitializing: "Initializing",
		lang.Connecting:          "Connecting",
		lang.Connected:           "Connected",
		lang.Disconnected:        "Disconnected",
		lang.MalFunctioned:       "MalFunctioned",
		lang.InfluxDBError:       "InfluxDBError",
	}

	errStrMap = map[lang.ErrIndex]string{
		lang.Ok:                    "Ok",
		lang.ErrDeviceNotExists:    "device does not exists!",
		lang.ErrDeviceNotConnected: "device does not connected！",
	}
)
