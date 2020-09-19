package thread

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
)

type item struct {
	PK      string // Hash key
	SK      string // Range key
	Content string `dynamo:"Content"`
	Closed  string `dynamo:"Closed"`
}

func newItem(thread *Thread) *item {
	clsd := "false"
	if thread.closed {
		clsd = "true"
	}

	return &item{
		PK:      "Team#0",
		SK:      thread.id,
		Content: thread.title,
		Closed:  clsd,
	}
}

func (i *item) toThread() *Thread {
	clsd := false
	if i.Closed == "true" {
		clsd = true
	}

	return &Thread{
		id:     i.SK,
		title:  i.Content,
		closed: clsd,
	}
}

type repository struct {
	tbl *dynamo.Table
}

var _ Repository = &repository{}

func repositoryError(err error) error {
	return fmt.Errorf("repository: %w", err)
}

func NewDynamoRepository(db *dynamo.DB) *repository {
	tbl := db.Table("Thread")
	return &repository{&tbl}
}

func (d *repository) get(ctx context.Context, req repositoryGetRequest) ([]*Thread, error) {
	var items []item
	if err := d.tbl.Get("PK", "Team#0").All(&items); err != nil {
		return nil, repositoryError(err)
	}

	rslts := make([]*Thread, len(items))
	for i, item := range items {
		rslts[i] = item.toThread()
	}
	return rslts, nil
}

func (d *repository) create(ctx context.Context, req repositoryCreateRequest) (*Thread, error) {
	condition := "attribute_not_exists(PK) AND attribute_not_exists(SK) "
	itm := newItem(req.thread)

	err := d.tbl.Put(&itm).If(condition).Run()
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return nil, err
			default:
				return nil, err
			}
		}
	}

	return req.thread, nil
}
