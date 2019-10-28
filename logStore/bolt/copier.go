package bolt

import (
	"bytes"
	"encoding/gob"
)

type JsonCopier struct {
	encoder *gob.Encoder
	decoder *gob.Decoder
}

func NewJsonCopier() *JsonCopier {
	buffer := bytes.Buffer{}
	return &JsonCopier{
		encoder: gob.NewEncoder(&buffer),
		decoder: gob.NewDecoder(&buffer),
	}
}

func (c *JsonCopier) encode(v interface{}) error {
	return c.encoder.Encode(v)
}

func (c *JsonCopier) decode(data interface{}) error {
	return c.decoder.Decode(data)
}
