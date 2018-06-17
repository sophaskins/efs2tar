package efs

import (
	"bytes"
	"encoding/binary"
	"os"
)

const BlockSize = 512
const NormalMagicNumber = 0x00072959
const GrownMagicNumber = 0x0007295A

type Filesystem struct {
	device *os.File
	size   int32 // size in blocks
	offset int32 // offset in blocks
	sb     *SuperBlock
}

type SuperBlock struct {
	Size         int32 // filesystem size (in BasicBlocks)
	FirstCG      int32 // BasicBlock offset of the first CylinderGroup
	CGSize       int32 // CylinderGroup size (in BasicBlocks)
	CGInodeSize  int16 // Number of BasicBlocks per CylinderGroup that are Inodes
	Sectors      int16 // sectors per track
	Heads        int16 // heads per cylinder
	CGCount      int16 // CylinderGroups in the filesystem
	Dirty        int16 // whether an fsck is required
	Pad0         int16
	CTime        int32   // last SuperBlock updated time
	Magic        int32   // filesystem magic number
	FSName       [6]byte // name of the filesystem
	FSPack       [6]byte // fs "pack" name
	BMSize       int32   // size in bytes of bitmap
	FreeBlocks   int32   // count of free blocks
	FreeInodes   int32   // count of free inodes
	BMBlock      int32   // offset of the bitmap
	ReplicatedSB int32   // offset of the replicated superblock
	LastInode    int32   // last unallocated inode
	Spare        [20]int8
	Checksum     int32
}

func NewFilesystem(device *os.File, size int32, offset int32) Filesystem {
	return Filesystem{
		device: device,
		size:   size,
		offset: offset,
	}
}

func NewSuperBlock(bb BasicBlock) SuperBlock {
	r := bytes.NewReader(bb[:])
	sb := SuperBlock{}
	binary.Read(r, binary.BigEndian, &sb)
	return sb
}

func (fs Filesystem) SuperBlock() *SuperBlock {
	if fs.sb == nil {
		bb := BasicBlock{}
		fs.device.ReadAt(bb[:], int64((fs.offset+1)*BlockSize))
		sb := NewSuperBlock(bb)
		fs.sb = &sb
	}
	return fs.sb
}

func (fs Filesystem) FirstCG() CylinderGroup {
	sb := fs.SuperBlock()
	blocks := make([]BasicBlock, sb.CGSize)

	for i := 0; i < int(sb.CGSize); i++ {
		bb := BasicBlock{}
		offset := int64((fs.offset + sb.FirstCG + int32(i)) * BlockSize)
		fs.device.ReadAt(bb[:], offset)
		blocks[i] = bb
	}

	return fs.NewCylinderGroup(blocks)
}

func (fs Filesystem) InodeForIndex(inodeIndex int32) Inode {
	sb := fs.SuperBlock()
	inodesPerBB := int32(4)
	inodeBlocksPerCG := int32(sb.CGInodeSize)
	inodeCGIndex := inodeIndex / (inodeBlocksPerCG * inodesPerBB)
	inodeBBinCG := inodeIndex % (inodeBlocksPerCG * inodesPerBB) / inodesPerBB
	bbIndex := sb.FirstCG + inodeCGIndex*sb.CGSize + inodeBBinCG
	bb := fs.BlockAt(bbIndex)

	offsetInBB := inodeIndex & (inodesPerBB - 1)
	return bb.ToInodes()[offsetInBB]
}

func (fs Filesystem) RootInode() Inode {
	return fs.InodeForIndex(2)
}

func (fs Filesystem) NewCylinderGroup(blocks []BasicBlock) CylinderGroup {
	return CylinderGroup{
		blocks: blocks,
		fs:     &fs,
	}
}

func (fs Filesystem) BlockAt(index int32) BasicBlock {
	rawOffset := int64((fs.offset + index) * BlockSize)
	bb := BasicBlock{}
	fs.device.ReadAt(bb[:], rawOffset)

	return bb
}

func (fs Filesystem) BlocksAt(index int32, length int32) []BasicBlock {
	blocks := make([]BasicBlock, length)
	for i := 0; i < int(length); i++ {
		blocks[i] = fs.BlockAt(index + int32(i))
	}

	return blocks
}
