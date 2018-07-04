package main

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/sophaskins/efs2tar/efs"
	"github.com/sophaskins/efs2tar/sgi"
)

func main() {
	file, _ := os.Open("./input.iso")
	b := make([]byte, 51200)
	_, _ = file.Read(b)

	vh := sgi.NewVolumeHeader(b)
	p := vh.Partitions[7]
	fs := efs.NewFilesystem(file, p.Blocks, p.First)
	extents := fs.RootInode().PayloadExtents()
	blocks := fs.ExtentToBlocks(extents[0])
	spew.Dump(blocks[0].ToDirectory().Entries())
}
