package db

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

func NewDynamoDB() *dynamo.DB {
	ep := os.Getenv("DB_ENDPOINT")
	ssn, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	return dynamo.New(ssn,
		&aws.Config{
			Endpoint: &ep,
		},
	)
}

func DestroyTestTable(tbl *dynamo.Table, t *testing.T) {
	if err := tbl.DeleteTable().Run(); err != nil {
		t.Fatal("create Thread table", err)
	}
}
