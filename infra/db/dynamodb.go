package db

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

func NewDynamoDB() *dynamo.DB {
	ep := os.Getenv("DB_ENDPOINT")
	return dynamo.New(session.New(),
		&aws.Config{
			Endpoint: &ep,
		},
	)
}
