package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kuba--/diff"
)

func main() {
	flag.Usage = func() {
		fmt.Printf("%s old-file delta-file new-file\n", flag.CommandLine.Name())
	}
	flag.Parse()
	args := flag.Args()
	if len(args) != 3 {
		flag.Usage()
		os.Exit(1)
	}

	oldFile, err := os.Open(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer oldFile.Close()

	deltaFile, err := os.Open(args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer deltaFile.Close()

	newFile, err := os.Create(args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer newFile.Close()

	if err = diff.Patch(oldFile, deltaFile, newFile); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
