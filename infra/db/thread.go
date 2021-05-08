package db

import (
	"strings"
	"testing"
	"time"

	"github.com/guregu/dynamo"
	"github.com/uji/ness-api-function/domain/thread"
	"github.com/uji/ness-api-function/reqctx"
	"golang.org/x/net/context"
)

var (
	ThreadTableName = "Thread"
)

type threadTable struct {
	PK        string    `dynamo:",hash"`
	SK        string    `dynamo:",range"`
	CreatedAt time.Time `localIndex:"PK-CreatedAt-index,range"`
	Closed    string    `localIndex:"PK-Closed-index,range"`
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

func DestroyThreadTable(db *dynamo.DB, name string) error {
	tbl := db.Table(name)
	return tbl.DeleteTable().Run()
}

type threadSchema struct {
	PK        string    // Hash key
	SK        string    // Range key
	CreatorID string    `dynamo:"CreatorID"`
	Content   string    `dynamo:"Content"`
	Closed    string    `dynamo:"Closed"`
	CreatedAt time.Time `dynamo:"CreatedAt"`
	UpdatedAt time.Time `dynamo:"UpdatedAt"`
}

func (t threadSchema) toThread() thread.Thread {
	clsd := false
	if t.Closed == "true" {
		clsd = true
	}

	return thread.NewThread(
		t.SK,
		thread.TeamID(t.PK),
		thread.UserID(t.CreatorID),
		t.Content,
		clsd,
		t.CreatedAt,
		t.UpdatedAt,
	)
}

type threadQuery struct {
	db  *dynamo.DB
	tbl *dynamo.Table
}

func NewThreadQuery(
	db *dynamo.DB,
	tableName string,
) *threadQuery {
	tbl := db.Table(tableName)
	return &threadQuery{db, &tbl}
}

func (t *threadQuery) GetThreadsByIDs(ctx context.Context, ids []string) ([]thread.Thread, error) {
	ainfo, err := reqctx.GetAuthenticationInfo(ctx)
	if err != nil {
		return nil, err
	}

	keys := make([]dynamo.Keyed, len(ids))
	for i, id := range ids {
		keys[i] = dynamo.Keys{ainfo.TeamID(), id}
	}
	b := t.tbl.Batch("PK", "SK")
	var rslt []threadSchema
	if err := b.Get(keys...).All(&rslt); err != nil {
		return nil, err
	}

	res := make([]thread.Thread, len(rslt))
	for i, r := range rslt {
		res[i] = r.toThread()
	}
	return res, nil
}
