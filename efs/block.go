package efs

import (
	"bytes"
	"encoding/binary"
)

const VolumeHeaderMagic = 0x0BE5A941

type BasicBlock [512]byte
type CylinderGroup []BasicBlock

const NormalMagicNumber = 0x00072959
const GrownMagicNumber = 0x0007295A

func NewBasicBlock(raw []byte) BasicBlock {
	bb := BasicBlock{}
	copy(bb[:], raw)

	return bb
}

type SuperBlock struct {
	Size        int32 // filesystem size (in BasicBlocks)
	FirstCG     int32 // BasicBlock offset of the first CylinderGroup
	CGSize      int32 // CylinderGroups per sector
	CGInodeSize int16 // number of indoes per CylinderGroup
	Sectors     int16 // sectors per track
	Heads       int16 // heads per cylinder
	CGCount     int16 // CylinderGroups in the filesystem
	Dirty       int16 // whether an fsck is required
	Pad0        int16
	CTime       int32
	Magic       int32   //
	FSName      [6]byte //
	FSPack      [6]byte //
	BMSize      int32   // size in bytes of bitmap
	FreeBlocks  int32   // count of free blocks
	FreeInodes  int32   // count of free inodes
	BMBlock     int32   // offset of the BitMap if the fs has been grown
	ReplSB      int32   // offset of the replacement superblock
	LastInode   int32   // last unallocated inode
	Spare       [20]int8
	Checksum    int32
}

func NewSuperBlock(bb BasicBlock) SuperBlock {
	r := bytes.NewReader(bb[:])
	sb := SuperBlock{}
	binary.Read(r, binary.BigEndian, &sb)
	return sb
}
