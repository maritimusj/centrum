package bolt

import (
	"bytes"
	"encoding/gob"
)

type gobCopier struct {
	encoder *gob.Encoder
	decoder *gob.Decoder
}

func NewCopier() *gobCopier {
	buffer := bytes.Buffer{}
	return &gobCopier{
		encoder: gob.NewEncoder(&buffer),
		decoder: gob.NewDecoder(&buffer),
	}
}

func (c *gobCopier) encode(v interface{}) error {
	return c.encoder.Encode(v)
}

func (c *gobCopier) decode(data interface{}) error {
	return c.decoder.Decode(data)
}
