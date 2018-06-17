package efs

import (
	"bytes"
	"encoding/binary"
)

type Directory struct {
	Magic     uint16
	FirstUsed uint8
	Slots     uint8
	Data      [508]byte
}

const DirectoryHeaderSize = 4

func (d Directory) Entries() []DirectoryEntry {
	entries := make([]DirectoryEntry, d.Slots)

	for i := 0; i < int(d.Slots); i++ {
		offset := (int(d.Data[i]) << 1) - DirectoryHeaderSize
		r := bytes.NewReader(d.Data[offset:])
		var inodeIndex uint32
		var nameLength uint8
		binary.Read(r, binary.BigEndian, &inodeIndex)
		binary.Read(r, binary.BigEndian, &nameLength)

		name := make([]byte, nameLength)
		binary.Read(r, binary.BigEndian, &name)
		entries[i] = DirectoryEntry{
			InodeIndex: inodeIndex,
			Name:       string(name),
		}
	}

	return entries
}

type DirectoryEntry struct {
	InodeIndex uint32
	Name       string
}
