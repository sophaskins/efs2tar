package efs

import (
	"bytes"
	"encoding/binary"
	"os"
)

const BlockSize = 512

type Filesystem struct {
	device *os.File
	size   int32 // size in blocks
	offset int32 // offset in blocks
	sb     *SuperBlock
}

type SuperBlock struct {
	Size         int32    // filesystem size (in BasicBlocks)
	FirstCG      int32    // BasicBlock offset of the first CylinderGroup
	CGSize       int32    // CylinderGroup size (in BasicBlocks)
	CGInodeSize  int16    // Number of BasicBlocks per CylinderGroup that are Inodes
	Sectors      int16    // sectors per track
	Heads        int16    // heads per cylinder
	CGCount      int16    // CylinderGroups in the filesystem
	Dirty        int16    // whether an fsck is required
	_            int16    // padding
	CTime        int32    // last SuperBlock updated time
	Magic        int32    // filesystem magic number
	FSName       [6]byte  // name of the filesystem
	FSPack       [6]byte  // fs "pack" name
	BMSize       int32    // size in bytes of bitmap
	FreeBlocks   int32    // count of free blocks
	FreeInodes   int32    // count of free inodes
	BMBlock      int32    // offset of the bitmap
	ReplicatedSB int32    // offset of the replicated superblock
	LastInode    int32    // last unallocated inode
	_            [20]int8 // padding
	Checksum     int32
}

func NewFilesystem(device *os.File, size int32, offset int32) *Filesystem {
	fs := &Filesystem{
		device: device,
		size:   size,
		offset: offset,
	}
	fs.initSuperBlock()
	return fs
}

func (fs *Filesystem) WalkFilesystem(callback func(Inode, string)) {
	fs.walkTree(fs.RootInode(), "", callback)
}

func (fs *Filesystem) ExtentToBlocks(e Extent) []BasicBlock {
	blocks := make([]BasicBlock, e.Length)
	for i := range blocks {
		blocks[i] = fs.blockAt(int32(e.StartBlock) + int32(i))
	}
	return blocks
}

func (fs *Filesystem) FileContents(in Inode) []byte {
	extents := fs.extents(in)
	fileBytes := make([]byte, in.Size)
	blockIndex := 0
	for _, extent := range extents {
		for _, block := range fs.ExtentToBlocks(extent) {
			copy(fileBytes[BlockSize*blockIndex:], block[:])
			blockIndex++
		}
	}

	return fileBytes
}

func (fs *Filesystem) extents(in Inode) []Extent {
	payloadExtents := in.PayloadExtents()
	if in.usesDirectExtents() {
		// if all of the extents fit inside of Payload (aka "direct extents")
		// we have a much simpler time reading the extents
		return payloadExtents
	}

	// if we have more extents than will fit in Payload, then the extents
	// in Payload (aka "indirect extents") point to ranges of blocks that
	// themselves contain the actual extents.
	extents := make([]Extent, in.NumExtents)
	extentsFetched := 0
	for _, indirectExtent := range payloadExtents {
		for _, extentBB := range fs.ExtentToBlocks(indirectExtent) {
			// copy respecting the length of extents saves us from
			// accidentally including the garbage extents at the end
			// of the last block (beyond NumExtents)
			copy(extents[extentsFetched:], extentBB.ToExtents())
			extentsFetched += extentsPerBlock
		}
	}

	return extents
}

func (fs *Filesystem) walkTree(in Inode, prefix string, callback func(Inode, string)) {
	switch in.Type() {
	case FileTypeDirectory:
		dirExtents := fs.extents(in)
		// I don't _believe_ it's possible for a directory to take up more than one block
		// (and so also to only have one extent) but...if I'm wrong, this is where it would
		// make a huge problem
		blocks := fs.ExtentToBlocks(dirExtents[0])
		for _, b := range blocks {
			for _, entry := range b.ToDirectory().Entries() {
				if entry.Name != "." && entry.Name != ".." {
					fs.walkTree(fs.inodeForIndex(int32(entry.InodeIndex)), prefix+"/"+entry.Name, callback)
				}
			}
		}
		fallthrough
	default:
		callback(in, prefix)
	}
}

func (fs *Filesystem) initSuperBlock() {
	bb := BasicBlock{}
	fs.device.ReadAt(bb[:], int64((fs.offset+1)*BlockSize))
	r := bytes.NewReader(bb[:])
	sb := SuperBlock{}
	binary.Read(r, binary.BigEndian, &sb)
	fs.sb = &sb
}

// InodeForIndex returns the Inode struct for the given inodeIndex
// (which is to say, indexed in to the list of all inodes, _not_ a
// block offset)
func (fs *Filesystem) inodeForIndex(inodeIndex int32) Inode {
	inodeBlocksPerCG := int32(fs.sb.CGInodeSize)
	inodeCGIndex := inodeIndex / (inodeBlocksPerCG * maxInodesPerBlock)
	inodeBBinCG := inodeIndex % (inodeBlocksPerCG * maxInodesPerBlock) / maxInodesPerBlock
	bbIndex := fs.sb.FirstCG + inodeCGIndex*fs.sb.CGSize + inodeBBinCG
	bb := fs.blockAt(bbIndex)

	offsetInBB := inodeIndex & (maxInodesPerBlock - 1)
	return bb.ToInodes()[offsetInBB]
}

func (fs *Filesystem) RootInode() Inode {
	return fs.inodeForIndex(2)
}

func (fs *Filesystem) blockAt(index int32) BasicBlock {
	rawOffset := int64((fs.offset + index) * BlockSize)
	bb := BasicBlock{}
	fs.device.ReadAt(bb[:], rawOffset)

	return bb
}
