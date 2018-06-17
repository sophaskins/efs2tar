package efs

import (
	"bytes"
	"encoding/binary"
)

const VolumeHeaderMagic = 0x0BE5A941

type BasicBlock [512]byte

type Inode struct {
	Mode       uint16
	NumLinks   int16
	UID        uint16
	GID        uint16
	Size       int32
	ATime      uint32
	MTime      uint32
	CTime      uint32
	Generation int32
	NumExtents int16
	Version    uint8
	Spare      uint8
	Payload    [96]byte // this is a union struct
}

func NewBasicBlock(raw []byte) BasicBlock {
	bb := BasicBlock{}
	copy(bb[:], raw)

	return bb
}

func (bb BasicBlock) ToInodes() []Inode {
	r := bytes.NewReader(bb[:])
	inodes := make([]Inode, 4)
	for i := 0; i < 4; i++ {
		binary.Read(r, binary.BigEndian, &inodes[i])
	}
	return inodes
}

const DirectExtentsLimit = 12

func (in Inode) Extents() []Extent {
	extents := make([]Extent, in.NumExtents)
	if in.NumExtents <= DirectExtentsLimit {
		for i := 0; i < int(in.NumExtents); i++ {
			extents[i] = NewExtent(in.Payload[8*i : 8*(i+1)])
		}
	}

	return extents
}

func (bb BasicBlock) ToDirectory() Directory {
	d := Directory{}
	r := bytes.NewReader(bb[:])
	binary.Read(r, binary.BigEndian, &d)

	return d
}
