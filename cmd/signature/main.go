package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"os"

	"github.com/kuba--/diff"
)

const (
	// 64MB
	maxBlockSize = 64 * 1024 * 1024
)

var (
	blockSize  int
	strongSize int
)

func main() {
	flag.IntVar(&blockSize, "b", 0, "block size")
	flag.IntVar(&strongSize, "s", 0, "strong size")
	flag.Usage = func() {
		fmt.Printf("%s [-b block size (<= %d)] [-s strong size] basis-file sig-file\n", flag.CommandLine.Name(), maxBlockSize)
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

	switch {
	case blockSize < 0:
		fmt.Printf("block size must be > 0 <= %d\n", maxBlockSize)
		os.Exit(2)
	case blockSize == 0:
		oldInfo, err := basisFile.Stat()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		if blockSize = int(oldInfo.Size() / 10); blockSize > maxBlockSize {
			blockSize = maxBlockSize
		}
	}

	switch {
	case strongSize < 0:
		fmt.Printf("strong size must be in range (0, %d]\n", md5.Size)
		os.Exit(2)
	case strongSize == 0:
		strongSize = md5.Size / 2
	case strongSize > md5.Size:
		strongSize = md5.Size
	}

	sigFile, err := os.Create(args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer sigFile.Close()

	if _, err = diff.WriteSignature(basisFile, sigFile, uint32(blockSize), byte(strongSize)); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
