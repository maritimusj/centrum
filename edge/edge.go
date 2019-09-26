package edge

type MeasureData struct {
	Name   string
	Tags   map[string]string
	Fields map[string]interface{}
}

type CtrlData struct {
	Values map[string]interface{}
	Error  chan error
}
