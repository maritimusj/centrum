package devices

import "time"

type MeasureData struct {
	Name   string
	Tags   map[string]string
	Fields map[string]interface{}
	Time   time.Time
}

func NewMeasureData(name string) *MeasureData {
	return &MeasureData{
		Name:   name,
		Tags: map[string]string{},
		Fields: map[string]interface{}{},
		Time:   time.Now(),
	}
}

func (measure *MeasureData) AddTag(name, val string) *MeasureData {
	measure.Tags[name] = val
	return measure
}

func (measure *MeasureData) AddField(name string, val interface{}) *MeasureData {
	measure.Fields[name] = val
	return measure
}
