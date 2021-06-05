package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/uji/ness-api-function/infra/elsch"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of tool\n")
	fmt.Fprintf(os.Stderr, "\tdbtool create # create indexes\n")
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "command arg required\n")
		usage()
		os.Exit(2)
	}

	arg := flag.Arg(0)
	switch arg {
	case "create":
		clt, err := elsch.NewClient(elsch.ThreadIndexName)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
		}
		if err := elsch.CreateIndices(clt); err != nil {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
		}
	case "delete":
		clt, err := elsch.NewClient(elsch.ThreadIndexName)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
		}
		if err := elsch.DeleteIndices(clt); err != nil {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
		}
	default:
		fmt.Fprintf(os.Stderr, "%s not found\n", arg)
		usage()
		os.Exit(2)
	}
}
