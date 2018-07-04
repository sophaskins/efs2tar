package main

import (
	"fmt"
	"os"

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
	fs.WalkFilesystem(func(in efs.Inode, path string) {
		fmt.Println(path)
	})
}
