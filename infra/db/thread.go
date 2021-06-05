package db

import (
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

func newThreadSchema(thread thread.Thread) threadSchema {
	clsd := "false"
	if thread.Closed() {
		clsd = "true"
	}

	return threadSchema{
		PK:        string(thread.TeamID()),
		SK:        thread.ID(),
		CreatorID: string(thread.CreatorID()),
		Content:   thread.Title(),
		Closed:    clsd,
		CreatedAt: thread.CreatedAt(),
		UpdatedAt: thread.UpdatedAt(),
	}
}

// func (t threadSchema) toThread() thread.Thread {
// 	clsd := false
// 	if t.Closed == "true" {
// 		clsd = true
// 	}
//
// 	return thread.NewThread(
// 		t.SK,
// 		thread.TeamID(t.PK),
// 		thread.UserID(t.CreatorID),
// 		t.Content,
// 		clsd,
// 		t.CreatedAt,
// 		t.UpdatedAt,
// 	)
// }

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

var _ thread.DynamoDB = &threadQuery{}

func (t *threadQuery) GetThreadsByIDs(ctx context.Context, ids []string) (map[string]thread.DynamoDBThreadRow, error) {
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

	res := make(map[string]thread.DynamoDBThreadRow, len(rslt))
	for _, r := range rslt {
		clsd := false
		if r.Closed == "true" {
			clsd = true
		}
		res[r.SK] = thread.DynamoDBThreadRow{
			Id:        r.SK,
			TeamID:    r.PK,
			CreaterID: r.CreatorID,
			Title:     r.Content,
			Closed:    clsd,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		}
	}
	return res, nil
}

func (t *threadQuery) Find(ctx context.Context, id string) (thread.DynamoDBThreadRow, error) {
	var thrd threadSchema
	ainfo, err := reqctx.GetAuthenticationInfo(ctx)
	if err != nil {
		return thread.DynamoDBThreadRow{}, err
	}
	if err := t.tbl.Get("PK", ainfo.TeamID()).Range("SK", dynamo.Equal, id).One(&thrd); err != nil {
		return thread.DynamoDBThreadRow{}, err
	}
	return thread.DynamoDBThreadRow{
		Id:        thrd.SK,
		TeamID:    thrd.PK,
		CreaterID: thrd.CreatorID,
		Title:     thrd.Content,
		Closed:    false,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}, err
}

func (t *threadQuery) Create(ctx context.Context, thread thread.Thread) error {
	condition := "attribute_not_exists(PK) AND attribute_not_exists(SK)"
	thrd := newThreadSchema(thread)

	if err := t.tbl.Put(&thrd).If(condition).Run(); err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return err
			default:
				return err
			}
		}
	}
	return nil
}

func (t *threadQuery) Update(ctx context.Context, thread thread.Thread) error {
	condition := "attribute_exists(PK) AND attribute_exists(SK)"
	thrd := newThreadSchema(thread)
	return t.tbl.Put(&thrd).If(condition).Run()
}
