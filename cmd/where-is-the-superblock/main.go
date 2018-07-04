package main

import (
	"os"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	file, _ := os.Open("./input.iso")
	b := make([]byte, 1024)
	file.Read(b)
	spew.Dump(b[512:1024])
}
