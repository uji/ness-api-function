package thread

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"github.com/guregu/dynamo"
)

type (
	dynamoRepository struct {
		tbl *dynamo.Table
	}

	item struct {
		PK      string // Hash key
		SK      string // Range key
		Content string `dynamo:"Content"`
		Closed  string `dynamo:"Closed"`
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
	if err := d.tbl.Get("PK", "Team#0").All(&items); err != nil {
		return nil, repositoryError(err)
	}

	rslts := make([]*Thread, len(items))
	for i, item := range items {
		clsd := false
		if item.Closed == "true" {
			clsd = true
		}

		rslts[i] = &Thread{
			id:     item.SK,
			title:  item.Content,
			closed: clsd,
		}
	}
	return rslts, nil
}

func (d *dynamoRepository) create(ctx context.Context, req repositoryCreateRequest) (Thread, error) {
	itm := item{
		PK:      "Team#0",
		SK:      "Thread#" + uuid.New().String(),
		Content: req.title,
		Closed:  "false",
	}

	err := d.tbl.Put(&itm).If("attribute_not_exists(SK)").Run()
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return Thread{}, err
			default:
				return Thread{}, err
			}
		}
	}

	return Thread{
		id:     itm.SK,
		title:  req.title,
		closed: false,
	}, nil
}
