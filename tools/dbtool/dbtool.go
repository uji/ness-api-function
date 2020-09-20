package dbtool

import (
	"testing"

	"github.com/guregu/dynamo"
)

type threadTable struct {
	PK string `dynamo:",hash"`
	SK string `dynamo:",range"`
}

func CreateThreadTable(db *dynamo.DB, name string) (dynamo.Table, error) {
	ctbl := db.CreateTable(name, threadTable{})
	if err := ctbl.Run(); err != nil {
		return dynamo.Table{}, err
	}
	return db.Table(name), nil
}

func CreateThreadTestTable(db *dynamo.DB, t *testing.T) dynamo.Table {
	tbl, err := CreateThreadTable(db, "Thread-"+t.Name())
	if err != nil {
		t.Fatal("create Thread table")
	}
	return tbl
}
