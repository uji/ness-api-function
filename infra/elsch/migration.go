package elsch

import (
	"context"
	"log"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func Migrate(client *Client) error {
	req := esapi.IndicesCreateRequest{
		Index:               "thread",
		Body:                nil,
		IncludeTypeName:     new(bool),
		MasterTimeout:       0,
		Timeout:             0,
		WaitForActiveShards: "",
		Pretty:              false,
		Human:               false,
		ErrorTrace:          false,
		FilterPath:          []string{},
		Header:              map[string][]string{},
	}
	res, err := req.Do(context.Background(), client.client)
	if err != nil {
		panic(err)
	}
	log.Println(res)
	return err
}
