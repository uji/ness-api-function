package dtstr

import (
	"context"
	"os"

	"cloud.google.com/go/datastore"
)

func NewClient() (*datastore.Client, error) {
	return datastore.NewClient(context.Background(), os.Getenv("DATASTORE_PROJECT_ID"))
}
