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
	BootDeviceParams DeviceParams
	VolDirs          [15]VolDir
	Partitions       [16]Partition
	Checksum         int32
	Padding          int32
}

type DeviceParams struct {
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

type VolDir struct {
	Name  [8]byte
	Block int32
	Bytes int32
}
type Partition struct {
	Blocks int32
	First  int32
	Type   int32
}

func NewVolumeHeader(raw []byte) VolumeHeader {
	r := bytes.NewReader(raw)
	vh := VolumeHeader{}
	binary.Read(r, binary.BigEndian, &vh)
	return vh
}
