package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kuba--/diff"
)

func main() {
	flag.Usage = func() {
		fmt.Printf("%s sig-file new-file delta-file\n", flag.CommandLine.Name())
	}
	flag.Parse()
	args := flag.Args()
	if len(args) != 3 {
		flag.Usage()
		os.Exit(1)
	}

	sigFile, err := os.Open(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer sigFile.Close()

	newFile, err := os.Open(args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer newFile.Close()

	deltaFile, err := os.Create(args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer deltaFile.Close()

	sig, err := diff.ReadSignature(sigFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if err = diff.WriteDelta(sig, newFile, deltaFile); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
