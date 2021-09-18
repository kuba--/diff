package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kuba--/diff"
)

func main() {
	flag.Usage = func() {
		fmt.Printf("%s basis-file delta-file recreated-file\n", flag.CommandLine.Name())
	}
	flag.Parse()
	args := flag.Args()
	if len(args) != 3 {
		flag.Usage()
		os.Exit(1)
	}

	basisFile, err := os.Open(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer basisFile.Close()

	deltaFile, err := os.Open(args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer deltaFile.Close()

	recreatedFile, err := os.Create(args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer recreatedFile.Close()

	if err = diff.Patch(basisFile, deltaFile, recreatedFile); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
