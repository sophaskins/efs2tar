package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("/Users/haski/Downloads/IRIS_Development_Option_5.3.iso")
	if err != nil {
		log.Fatal(err)
	}
	stat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(stat)
	fmt.Println("hello world")
}
