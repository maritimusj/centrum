package lang

type StrIndex int

const (
	_ StrIndex = iota
	AdapterInitializing
	Connecting
	Connected
	Disconnected
)
