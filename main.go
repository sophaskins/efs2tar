package main

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/sophaskins/efs2tar/efs"
	"github.com/sophaskins/efs2tar/sgi"
)

func main() {
	path := os.Args[1]
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, 51200)
	_, err = file.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	scs := spew.ConfigState{
		Indent: "\t",
	}
	vh := sgi.NewVolumeHeader(b)

	p := vh.Partitions[7]
	fs := efs.NewFilesystem(file, p.Blocks, p.First)

	rootInode := fs.RootInode()
	scs.Dump(rootInode)
	rootInodeExtents := rootInode.Extents()
	scs.Dump(rootInodeExtents)
	blocks := fs.BlocksAt(int32(rootInodeExtents[0].Block), int32(rootInodeExtents[0].Length))
	for _, b := range blocks {
		scs.Dump(b.ToDirectory().Entries())
	}
}
