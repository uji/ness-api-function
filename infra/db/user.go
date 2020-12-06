package db

import (
	"strings"
	"testing"
	"time"

	"github.com/guregu/dynamo"
)

var (
	UserTableName = "User"
)

type userTable struct {
	PK        string    `dynamo:",hash"`
	SK        string    `dynamo:",range"`
	Name      string    `localIndex:"PK-Name-index,range"`
	CreatedAt time.Time `localIndex:"PK-CreatedAt-index,range"`
}

func CreateUserTable(db *dynamo.DB, name string) (dynamo.Table, error) {
	ctbl := db.CreateTable(name, userTable{})
	if err := ctbl.Run(); err != nil {
		return dynamo.Table{}, err
	}
	return db.Table(name), nil
}

func CreateUserTestTable(db *dynamo.DB, t *testing.T) dynamo.Table {
	tName := strings.ReplaceAll(t.Name(), "/", "-")
	tbl, err := CreateUserTable(db, "User-"+tName)
	if err != nil {
		t.Fatal("create User table", err)
	}
	return tbl
}

func DestroyUserTable(db *dynamo.DB, name string) error {
	tbl := db.Table(name)
	return tbl.DeleteTable().Run()
}
