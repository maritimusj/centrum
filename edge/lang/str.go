package lang

type StrIndex int

const (
	EdgeUnknownState StrIndex = iota
	AdapterInitializing
	Connecting
	Connected
	Disconnected
	MalFunctioned
	InfluxDBError
)
