package util

import (
	"bytes"
	"encoding/gob"
)

type GobCopier struct {
	encoder *gob.Encoder
	decoder *gob.Decoder
}

func NewCopier() *GobCopier {
	buffer := bytes.Buffer{}
	return &GobCopier{
		encoder: gob.NewEncoder(&buffer),
		decoder: gob.NewDecoder(&buffer),
	}
}

func (c *GobCopier) Encode(v interface{}) error {
	return c.encoder.Encode(v)
}

func (c *GobCopier) Decode(data interface{}) error {
	return c.decoder.Decode(data)
}
