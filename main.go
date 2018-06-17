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

	//	scs.Dump(fs.SuperBlock())

	scs.Dump(fs.RootInode())

	//	offset := 64
	//	blocks := make([]efs.BasicBlock, 4)
	//	for i := 0; i < 4; i++ {
	//		blocks[i] = efs.NewBasicBlock(b[512*(i+offset) : 512*(i+offset)+512])
	//	}
	//
	//	sb := efs.NewSuperBlock(blocks[1])
	//	scs.Dump(sb.FSName)
}
