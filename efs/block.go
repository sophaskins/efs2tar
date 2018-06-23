package efs

import (
	"bytes"
	"encoding/binary"
)

// BasicBlock is how the disk is divided up in EFS
// Each BasicBlock can be:
//   * a collection of inodes
//   * part of the body of a file
//   * a directory listing
//   * a collection of indirect extents
//   * other things?
type BasicBlock [512]byte

func NewBasicBlock(raw []byte) BasicBlock {
	bb := BasicBlock{}
	copy(bb[:], raw)
	return bb
}

const maxInodesPerBlock = 4

func (bb BasicBlock) ToInodes() []Inode {
	r := bytes.NewReader(bb[:])
	inodes := make([]Inode, maxInodesPerBlock)
	for i := range inodes {
		binary.Read(r, binary.BigEndian, &inodes[i])
	}
	return inodes
}

func (bb BasicBlock) ToDirectory() Directory {
	d := Directory{}
	r := bytes.NewReader(bb[:])
	binary.Read(r, binary.BigEndian, &d)
	return d
}

const extentsPerBlock = 64

func (bb BasicBlock) ToExtents() []Extent {
	// Note that not all of these extents are going to be valid -
	// the inode knows how many extents it should actually use
	// (and then discards the rest, which may contain garbage, who knows)
	extents := make([]Extent, extentsPerBlock)
	for i := range extents {
		extents[i] = NewExtent(bb[8*i : 8*(i+1)])
	}
	return extents
}
