package thread

import (
	"context"
	"fmt"

	"github.com/guregu/dynamo"
)

type (
	dynamoRepository struct {
		tbl *dynamo.Table
	}

	item struct {
		ThreadID string // Hash key
		Key      string // Range key
		Content  string `dynamo:"Content"`
		Closed   string `dynamo:"Closed"`
	}
)

var _ repository = &dynamoRepository{}

func repositoryError(err error) error {
	return fmt.Errorf("repository: %w", err)
}

func NewDynamoRepository(db *dynamo.DB) *dynamoRepository {
	tbl := db.Table("Thread")
	return &dynamoRepository{&tbl}
}

func (d *dynamoRepository) get(ctx context.Context, req repositoryGetRequest) ([]*Thread, error) {
	var items []item
	if err := d.tbl.Get("ThreadID", "0").All(&items); err != nil {
		return nil, repositoryError(err)
	}

	rslts := make([]*Thread, len(items))
	for i, item := range items {
		clsd := false
		if item.Closed == "true" {
			clsd = true
		}

		rslts[i] = &Thread{
			id:     item.ThreadID,
			title:  item.Content,
			closed: clsd,
		}
	}
	return rslts, nil
}

func (d *dynamoRepository) create(ctx context.Context, req repositoryCreateRequest) (Thread, error) {
	return Thread{}, nil
}
