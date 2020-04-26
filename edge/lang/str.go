package lang

type StrIndex int

func (index StrIndex) Str(params ...interface{}) string {
	return Str(index, params...)
}

const (
	EdgeUnknownState StrIndex = iota
	AdapterInitializing
	Connecting
	Connected
	Disconnected
	MalFunctioned
	InfluxDBError
)
