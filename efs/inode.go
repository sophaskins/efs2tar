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
	// Payload is a union struct - sometimes it contains extents, but
	// it also can contain other stuff (like link targets and device
	// descriptors, which are not implemented here)
	Payload [96]byte
}

const DirectExtentsLimit = 12

const (
	FileTypeFIFO            = 010
	FileTypeCharacterDevice = 020
	FileTypeDirectory       = 040
	FileTypeBlockDevice     = 060
	FileTypeRegular         = 0100
	FileTypeSymlink         = 0120
	FileTypeSocket          = 0140
)

// FormatMode is the beginnings of something that could format
// an Inode like `/bin/ls` does
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

func (in Inode) usesDirectExtents() bool {
	return in.NumExtents <= DirectExtentsLimit
}

// payloadExtents unpacks the extents stored in the
// Payload field. The number of these varies based on
// whether or not we're using Direct Extents or not
func (in Inode) payloadExtents() []Extent {
	var extents []Extent

	if in.usesDirectExtents() {
		// the number of payload extents is just "all the
		// extents" if this inode uses Direct Extents
		extents = make([]Extent, in.NumExtents)
		for i := range extents {
			extents[i] = NewExtent(in.Payload[8*i : 8*(i+1)])
		}
	} else {
		// the number of payload extents is contained in the
		// NumIndirectExtents field of the first payload extent
		// if we're using indirect extents
		firstExtent := NewExtent(in.Payload[0:8])
		extents = make([]Extent, firstExtent.NumIndirectExtents)
		extents[0] = firstExtent
		for i := 1; i < int(extents[0].NumIndirectExtents); i++ {
			extents[i] = NewExtent(in.Payload[8*i : 8*(i+1)])
		}
	}

	return extents
}
