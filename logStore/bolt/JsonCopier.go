package bolt

import (
	"bytes"
	"encoding/json"
)

type JsonCopier struct {
	encoder *json.Encoder
	decoder *json.Decoder
}

func NewJsonCopier() *JsonCopier {
	buffer := bytes.Buffer{}
	return &JsonCopier{
		encoder: json.NewEncoder(&buffer),
		decoder: json.NewDecoder(&buffer),
	}
}

func (c *JsonCopier) encode(v interface{}) error {
	return c.encoder.Encode(v)
}

func (c *JsonCopier) decode(data interface{}) error {
	return c.decoder.Decode(data)
}