package thread

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
	"github.com/uji/ness-api-function/reqctx"
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
	ainfo, err := reqctx.GetAuthenticationInfo(ctx)
	if err != nil {
		return nil, err
	}

	qr := d.tbl.Get("PK", ainfo.TeamID()).Index("PK-CreatedAt-index").Order(false)
	if req.offsetTime.Valid {
		qr = qr.Range("CreatedAt", dynamo.Less, req.offsetTime.Time)
	}
	if req.closed.Valid {
		clsd := "false"
		if req.closed.Bool {
			clsd = "true"
		}
		qr = qr.Filter("Closed = ?", clsd)
	}

	if err := qr.All(&items); err != nil {
		return nil, repositoryError(err)
	}

	rslts := make([]Thread, len(items))
	for i, item := range items {
		rslts[i] = item.toThread()
	}
	return rslts, nil
}

func (d *repository) find(ctx context.Context, req repositoryFindRequest) (Thread, error) {
	var itm item
	ainfo, err := reqctx.GetAuthenticationInfo(ctx)
	if err != nil {
		return nil, err
	}

	if err := d.tbl.Get("PK", ainfo.TeamID()).Range("SK", dynamo.Equal, req.threadID).One(&itm); err != nil {
		return nil, err
	}
	return itm.toThread(), nil
}

func (d *repository) create(ctx context.Context, req repositoryCreateRequest) error {
	condition := "attribute_not_exists(PK) AND attribute_not_exists(SK)"
	itm := newItem(req.thread)

	err := d.tbl.Put(&itm).If(condition).Run()
	if err != nil {
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

func (d *repository) update(ctx context.Context, req repositoryUpdateRequest) error {
	condition := "attribute_exists(PK) AND attribute_exists(SK)"
	return d.tbl.Put(newItem(req.thread)).If(condition).Run()
}

type item struct {
	PK        string    // Hash key
	SK        string    // Range key
	CreatorID string    `dynamo:"CreatorID"`
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
		PK:        string(thread.TeamID()),
		SK:        thread.ID(),
		CreatorID: string(thread.CreatorID()),
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
		teamID:    TeamID(i.PK),
		createrID: UserID(i.CreatorID),
		title:     i.Content,
		closed:    clsd,
		createdAt: i.CreatedAt,
		updatedAt: i.UpdatedAt,
	}
}
