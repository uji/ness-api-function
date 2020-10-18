package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/uji/ness-api-function/infra/db"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of dbtool\n")
	fmt.Fprintf(os.Stderr, "\tdbtool create # create tables & indexes\n")
	fmt.Fprintf(os.Stderr, "\tdbtool destroy # destroy tables & indexes\n")
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
		dnmdb := db.NewDynamoDB()
		if _, err := db.CreateThreadTable(dnmdb, db.ThreadTableName); err != nil {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
		}
		if _, err := db.CreateUserTable(dnmdb, db.UserTableName); err != nil {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
		}
	case "destroy":
		dnmdb := db.NewDynamoDB()
		if err := db.DestroyThreadTable(dnmdb, db.ThreadTableName); err != nil {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
		}
		if err := db.DestroyUserTable(dnmdb, db.UserTableName); err != nil {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
		}
	default:
		fmt.Fprintf(os.Stderr, "%s not found\n", arg)
		usage()
		os.Exit(2)
	}
}
