package thread

import (
	"context"

	"github.com/guregu/dynamo"
)

type (
	dynamoRepository struct {
		tbl *dynamo.Table
	}

	item struct {
		Team   string // Hash key
		Title  string // Range key
		Closed string `dynamo:"Closed"`
	}
)

var _ repository = &dynamoRepository{}

func NewDynamoRepository(db *dynamo.DB) *dynamoRepository {
	tbl := db.Table("Thread")
	return &dynamoRepository{&tbl}
}

func (d *dynamoRepository) get(ctx context.Context, req repositoryGetRequest) ([]*Thread, error) {
	var items []item
	if err := d.tbl.Scan().All(&items); err != nil {
		return nil, err
	}

	rslts := make([]*Thread, len(items))
	for i, item := range items {
		clsd := false
		if item.Closed == "true" {
			clsd = true
		}

		rslts[i] = &Thread{
			id:     item.Team,
			title:  item.Title,
			closed: clsd,
		}
	}
	return rslts, nil
}

func (d *dynamoRepository) create(ctx context.Context, req repositoryCreateRequest) (Thread, error) {
	return Thread{}, nil
}
