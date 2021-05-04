package elsch

import (
	"context"
	"log"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func CreateIndices(client *Client) error {
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

func DeleteIndices(client *Client) error {
	req := esapi.IndicesDeleteRequest{
		Index:             []string{"thread"},
		AllowNoIndices:    new(bool),
		ExpandWildcards:   "",
		IgnoreUnavailable: new(bool),
		MasterTimeout:     0,
		Timeout:           0,
		Pretty:            false,
		Human:             false,
		ErrorTrace:        false,
		FilterPath:        []string{},
		Header:            map[string][]string{},
	}
	res, err := req.Do(context.Background(), client.client)
	if err != nil {
		panic(err)
	}
	log.Println(res)
	return err
}
