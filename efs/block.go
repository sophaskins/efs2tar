package efs

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/davecgh/go-spew/spew"
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
const IndirectExtentsPerBlock = 64

func (in Inode) Extents(fs Filesystem) []Extent {
	extents := make([]Extent, in.NumExtents)
	if in.NumExtents <= DirectExtentsLimit {
		for i := 0; i < int(in.NumExtents); i++ {
			extents[i] = NewExtent(in.Payload[8*i : 8*(i+1)])
		}
	} else {
		// the first indirect extent tells us how many indirect extents total
		firstIndirectExtent := NewExtent(in.Payload[0:8])
		spew.Dump(firstIndirectExtent)
		indirectExtents := make([]Extent, firstIndirectExtent.Offset)
		indirectExtents[0] = firstIndirectExtent
		for i := 1; i < int(indirectExtents[0].Offset); i++ {
			indirectExtents[i] = NewExtent(in.Payload[8*i : 8*(i+1)])
		}

		// this is definitely not working properly
		extentsFetched := 0
		for _, indirectExtent := range indirectExtents {
			indirectBBs := fs.BlocksAt(int32(indirectExtent.Block), int32(indirectExtent.Length))
			for _, indirectBB := range indirectBBs {
				for i := 0; i < IndirectExtentsPerBlock && extentsFetched < int(in.NumExtents); i++ {
					extents[extentsFetched] = NewExtent(indirectBB[8*i : 8*(i+1)])
					extentsFetched++
				}
			}

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

func (in Inode) FormatMode() string {
	modeString := ""
	switch in.Type() {
	case FileTypeFIFO:
		modeString += "p" // FIFO
	case FileTypeCharacterDevice:
		modeString += "c" // character device
	case FileTypeDirectory:
		modeString += "d" // directory
	case FileTypeBlockDevice:
		modeString += "b"
	case FileTypeRegular:
		modeString += "-"
	case FileTypeSymlink:
		modeString += "l"
	case FileTypeSocket:
		modeString += "s"
	}

	return modeString

}

func (in Inode) Type() uint16 {
	return in.Mode >> 9
}

const (
	FileTypeFIFO            = 010
	FileTypeCharacterDevice = 020
	FileTypeDirectory       = 040
	FileTypeBlockDevice     = 060
	FileTypeRegular         = 0100
	FileTypeSymlink         = 0120
	FileTypeSocket          = 0140
)

func (in Inode) ToRegularFile(fs Filesystem) []byte {
	extents := in.Extents(fs)
	fileBytes := make([]byte, in.Size)
	blockIndex := 0
	for _, extent := range extents {
		blocks := fs.BlocksAt(int32(extent.Block), int32(extent.Length))
		for _, block := range blocks {
			copy(fileBytes[BlockSize*blockIndex:], block[:])
			blockIndex++
		}
	}

	if blockIndex == 0 && in.Size != 0 {
		fmt.Println("somehow have no bytes")
		spew.Dump(in.Payload)
	}

	return fileBytes
}
