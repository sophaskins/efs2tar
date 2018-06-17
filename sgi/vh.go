package sgi

import (
	"bytes"
	"encoding/binary"
)

type VolumeHeader struct {
	MagicNumber      uint32
	Root             int16
	Swap             int16
	Bootfile         [16]byte
	BootDeviceParams DeviceParameters
	VolumeDirectory  [15]FileHeader
	Partitions       [16]Partition
	Checksum         int32
	Padding          int32
}

type DeviceParameters struct {
	Skew           uint8
	Gap1           uint8
	Gap2           uint8
	SparesCylinder uint8
	Cylinder       uint16
	Shd0           uint16
	Tracks         uint16
	CtqDepth       uint8
	CylindersShi   uint8
	Unused         uint16
	Sectors        uint16
	SectorBytes    uint16
	Interleave     uint16
	Flags          uint32
	DataRate       uint32
	NumRetries     uint32
	Mspw           uint32
	Xgap1          uint16
	Xsync          uint16
	XRdly          uint16
	Xgap2          uint16
	Xrgate         uint16
	Xwcont         uint16
}

// FileHeader points to a special file stored in the volume header
// These are typically not accessed from user-space - examples
// of uses include the `fx` partitioning tool included on OS
// installation discs
type FileHeader struct {
	Name  [8]byte
	Block int32
	Bytes int32
}

// Partition points to where on disk various partitions are
// located. Partitions _can_ overlap, and they don't _have_
// to be filesystems (eg, swap)
type Partition struct {
	Blocks int32
	First  int32
	Type   PartitionType
}

type PartitionType int32

const (
	VolumeHeaderPartition = PartitionType(0)
	TrackRepl             = PartitionType(1)
	Repl                  = PartitionType(2)
	Raw                   = PartitionType(3)
	BSD                   = PartitionType(4)
	SystemV               = PartitionType(5)
	Volume                = PartitionType(6)
	EFS                   = PartitionType(7)
	LVol                  = PartitionType(8)
	RLVol                 = PartitionType(9)
	XFS                   = PartitionType(10)
	XFSLog                = PartitionType(11)
	XLV                   = PartitionType(12)
	XVM                   = PartitionType(13)
)

func (pt PartitionType) String() string {
	switch pt {
	case VolumeHeaderPartition:
		return "VolumeHeader"
	case TrackRepl:
		return "TrackRepl"
	case Repl:
		return "Repl"
	case Raw:
		return "Raw"
	case BSD:
		return "BSD"
	case SystemV:
		return "SystemV"
	case Volume:
		return "Volume"
	case EFS:
		return "EFS"
	case LVol:
		return "LogicalVolume"
	case RLVol:
		return "RLVolume"
	case XFS:
		return "XFS"
	case XFSLog:
		return "XFSLog"
	case XLV:
		return "XLV"
	case XVM:
		return "XVM"
	}

	return "Unknown"
}

func NewVolumeHeader(raw []byte) VolumeHeader {
	r := bytes.NewReader(raw)
	vh := VolumeHeader{}
	binary.Read(r, binary.BigEndian, &vh)
	return vh
}
