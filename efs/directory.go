package efs

import (
	"bytes"
	"encoding/binary"
)

// Directory is the format of a block that stores
// a directory listing. NetBSD's sys/fs/efs/efs_dir.h has
// some detail as to how the Data field is organized
type Directory struct {
	Magic     uint16
	FirstUsed uint8
	Slots     uint8
	Data      [508]byte
}

type DirectoryEntry struct {
	InodeIndex uint32
	Name       string
}

const DirectoryHeaderOffset = 4

func (d Directory) Entries() []DirectoryEntry {
	entries := make([]DirectoryEntry, d.Slots)

	for i := range entries {
		// The "slots" at the low indexes of Data tell us where the
		// "entries" in the high indexes are
		offset := (int(d.Data[i]) << 1) - DirectoryHeaderOffset

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
