package dbtool

import (
	"strings"
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
	tName := strings.ReplaceAll(t.Name(), "/", "-")
	tbl, err := CreateThreadTable(db, "Thread-"+tName)
	if err != nil {
		t.Fatal("create Thread table", err)
	}
	return tbl
}
