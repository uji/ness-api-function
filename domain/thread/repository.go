package thread

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
)

type repository struct {
	db  *dynamo.DB
	tbl *dynamo.Table
}

var _ Repository = &repository{}

func repositoryError(err error) error {
	return fmt.Errorf("repository: %w", err)
}

func NewDynamoRepository(
	db *dynamo.DB,
	tableName string,
) *repository {
	tbl := db.Table(tableName)
	return &repository{db, &tbl}
}

func (d *repository) get(ctx context.Context, req repositoryGetRequest) ([]Thread, error) {
	var items []item
	teamID := "Team#0"

	qr := d.tbl.Get("PK", teamID).Index("PK-CreatedAt-index").Order(false).Limit(req.limit)
	if req.offsetTime.Valid {
		qr = qr.Range("CreatedAt", dynamo.Less, req.offsetTime.Time)
	}

	err := qr.All(&items)
	if err != nil {
		return nil, repositoryError(err)
	}

	rslts := make([]Thread, len(items))
	for i, item := range items {
		rslts[i] = item.toThread()
	}
	return rslts, nil
}

func (d *repository) create(ctx context.Context, req repositoryCreateRequest) (Thread, error) {
	condition := "attribute_not_exists(PK) AND attribute_not_exists(SK)"
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

func (d *repository) update(ctx context.Context, req repositoryUpdateRequest) (Thread, error) {
	condition := "attribute_exists(PK) AND attribute_exists(SK)"
	if err := d.tbl.Put(newItem(req.thread)).If(condition).Run(); err != nil {
		return nil, err
	}
	return req.thread, nil
}

func (d *repository) close(ctx context.Context, req repositoryCloseRequest) (Thread, error) {
	var itm item
	if err := d.tbl.Get("PK", "Team#0").Range("SK", dynamo.Equal, req.threadID).One(&itm); err != nil {
		return nil, err
	}

	th := itm.toThread()
	th.Close()

	return d.update(ctx, repositoryUpdateRequest{th})
}

type item struct {
	PK        string    // Hash key
	SK        string    // Range key
	Content   string    `dynamo:"Content"`
	Closed    string    `dynamo:"Closed"`
	CreatedAt time.Time `dynamo:"CreatedAt"`
	UpdatedAt time.Time `dynamo:"UpdatedAt"`
}

func newItem(thread Thread) *item {
	clsd := "false"
	if thread.Closed() {
		clsd = "true"
	}

	return &item{
		PK:        "Team#0",
		SK:        thread.ID(),
		Content:   thread.Title(),
		Closed:    clsd,
		CreatedAt: thread.CreatedAt(),
		UpdatedAt: thread.UpdatedAt(),
	}
}

func (i *item) toThread() Thread {
	clsd := false
	if i.Closed == "true" {
		clsd = true
	}

	return &thread{
		id:        i.SK,
		title:     i.Content,
		closed:    clsd,
		createdAt: i.CreatedAt,
		updatedAt: i.UpdatedAt,
	}
}
