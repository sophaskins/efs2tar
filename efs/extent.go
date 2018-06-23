package efs

import (
	"bytes"
	"encoding/binary"
)

type Extent struct {
	Magic  uint8
	Block  uint32
	Length uint8
	Offset uint32
}

func NewExtent(raw []byte) Extent {
	paddedExtent := make([]byte, 10)
	copy(paddedExtent[0:1], raw[0:1])
	copy(paddedExtent[2:6], raw[1:5])
	copy(paddedExtent[7:10], raw[5:8])

	r := bytes.NewReader(paddedExtent)
	e := Extent{}
	binary.Read(r, binary.BigEndian, &e)
	return e
}
