package main

import (
	"archive/tar"
	"log"
	"os"

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

	outputFile, err := os.OpenFile("/Users/haski/Downloads/out.tar", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	tw := tar.NewWriter(outputFile)

	vh := sgi.NewVolumeHeader(b)
	p := vh.Partitions[7]
	fs := efs.NewFilesystem(file, p.Blocks, p.First)
	rootInode := fs.RootInode()

	fs.WalkTree(rootInode, "", buildTarCallback(tw, fs))
	tw.Close()
}

func buildTarCallback(tw *tar.Writer, fs *efs.Filesystem) func(efs.Inode, string) {
	return func(in efs.Inode, path string) {
		if path == "" {
			return
		}

		if in.Type() == efs.FileTypeDirectory {
			hdr := &tar.Header{
				Name:     path,
				Mode:     0755,
				Typeflag: tar.TypeDir,
			}
			if err := tw.WriteHeader(hdr); err != nil {
				log.Fatal(err)
			}
		} else if in.Type() == efs.FileTypeRegular {

			contents := in.FileContents(fs)
			hdr := &tar.Header{
				Name: path,
				Mode: 0755,
				Size: int64(len(contents)),
			}
			if err := tw.WriteHeader(hdr); err != nil {
				log.Fatal(err)
			}
			if _, err := tw.Write([]byte(contents)); err != nil {
				log.Fatal(err)
			}
		}
	}
}
