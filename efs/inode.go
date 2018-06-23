package efs

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

const DirectExtentsLimit = 12

func (in Inode) Extents(fs *Filesystem) []Extent {
	extents := make([]Extent, in.NumExtents)
	if in.NumExtents <= DirectExtentsLimit {
		// if all of the extents fit inside of Payload (aka "direct extents")
		// we have a much simpler time reading the extents
		for i := range extents {
			extents[i] = NewExtent(in.Payload[8*i : 8*(i+1)])
		}
	} else {
		// if we have more extents than will fit in Payload, then the extents
		// in Payload (aka "indirect extents") point to ranges of blocks that
		// themselves contain the actual extents.

		// the NumIndirectExtents field on the first indirect extent tells us how many
		// indirect extents there are inside Payload.
		firstIndirectExtent := NewExtent(in.Payload[0:8])
		indirectExtents := make([]Extent, firstIndirectExtent.NumIndirectExtents)
		indirectExtents[0] = firstIndirectExtent
		for i := 1; i < int(indirectExtents[0].NumIndirectExtents); i++ {
			indirectExtents[i] = NewExtent(in.Payload[8*i : 8*(i+1)])
		}

		extentsFetched := 0
		for _, indirectExtent := range indirectExtents {
			extentBBs := fs.BlocksAt(int32(indirectExtent.StartBlock), int32(indirectExtent.Length))
			for _, extentBB := range extentBBs {
				// copy respecting the length of extents saves us from
				// accidentally including the garbage extents at the end
				// of the last block (beyond NumExtents)
				copy(extents[extentsFetched:], extentBB.ToExtents())
				extentsFetched += extentsPerBlock
			}
		}
	}

	return extents
}

func (in Inode) FormatMode() string {
	modeString := ""
	switch in.Type() {
	case FileTypeFIFO:
		modeString += "p"
	case FileTypeCharacterDevice:
		modeString += "c"
	case FileTypeDirectory:
		modeString += "d"
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

func (in Inode) FileContents(fs *Filesystem) []byte {
	extents := in.Extents(fs)
	fileBytes := make([]byte, in.Size)
	blockIndex := 0
	for _, extent := range extents {
		blocks := fs.BlocksAt(int32(extent.StartBlock), int32(extent.Length))
		for _, block := range blocks {
			copy(fileBytes[BlockSize*blockIndex:], block[:])
			blockIndex++
		}
	}

	return fileBytes
}
