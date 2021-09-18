package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"os"

	"github.com/kuba--/diff"
)

var (
	blockSize  int
	strongSize byte = byte(md5.Size / 2)
)

func main() {
	flag.IntVar(&blockSize, "b", 0, "block size")
	flag.Usage = func() {
		fmt.Printf("%s [-b block size] basis-file sig-file\n", flag.CommandLine.Name())
	}
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	basisFile, err := os.Open(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer basisFile.Close()

	if blockSize == 0 {
		oldInfo, err := basisFile.Stat()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		blockSize = int(oldInfo.Size() / 10)
	}

	sigFile, err := os.Create(args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer sigFile.Close()

	if _, err = diff.WriteSignature(basisFile, sigFile, uint32(blockSize), strongSize); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
