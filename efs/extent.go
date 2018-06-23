package efs

import (
	"bytes"
	"encoding/binary"
)

// Extents are basically "pointers" to ranges of blocks that
// contain the body of their parent entity (a file, directory, etc)
type Extent struct {
	Magic              uint8
	StartBlock         uint32
	Length             uint8
	NumIndirectExtents uint32
}

func NewExtent(raw []byte) Extent {
	// You can't just binary.Read this struct directly because
	// StartBlock and NumIndirectExtents fields are only given
	// 24 bits in the on-disk struct and I'm not aware of a way
	// to hint that information to binary.Read w/o just adding
	// more padding like this
	paddedExtent := make([]byte, 10)
	copy(paddedExtent[0:1], raw[0:1])
	copy(paddedExtent[2:6], raw[1:5])
	copy(paddedExtent[7:10], raw[5:8])

	r := bytes.NewReader(paddedExtent)
	e := Extent{}
	binary.Read(r, binary.BigEndian, &e)
	return e
}
