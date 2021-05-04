package elsch

import (
	"context"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func CreateIndices(client *Client) error {
	req := esapi.IndicesCreateRequest{
		Index:               threadIndexName,
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
	_, err := req.Do(context.Background(), client.client)
	if err != nil {
		panic(err)
	}
	return err
}

func DeleteIndices(client *Client) error {
	req := esapi.IndicesDeleteRequest{
		Index:             []string{threadIndexName},
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
	_, err := req.Do(context.Background(), client.client)
	if err != nil {
		panic(err)
	}
	return err
}
