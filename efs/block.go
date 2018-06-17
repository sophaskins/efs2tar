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
	Payload    [8]byte // this is a union struct
}

func NewBasicBlock(raw []byte) BasicBlock {
	bb := BasicBlock{}
	copy(bb[:], raw)

	return bb
}

func (bb BasicBlock) ToInode() Inode {
	r := bytes.NewReader(bb[:])
	inode := Inode{}
	binary.Read(r, binary.BigEndian, &inode)
	return inode
}
