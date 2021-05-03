package elsch

import (
	"os"

	"github.com/elastic/go-elasticsearch/v7"
)

type Client struct {
	client *elasticsearch.Client
}

func NewClient() (*Client, error) {
	addresses := []string{
		"http://elasticsearch:9200",
		"http://elasticsearch:9300",
	}

	adrs1 := os.Getenv("ELASTICSEARCH_ADDRESS_1")
	adrs2 := os.Getenv("ELASTICSEARCH_ADDRESS_2")

	if adrs1 != "" && adrs2 != "" {
		addresses = []string{adrs1, adrs2}
	}

	cfg := elasticsearch.Config{
		Addresses:             addresses,
		Username:              os.Getenv("ELASTICSEARCH_USERNAME"),
		Password:              os.Getenv("ELASTICSEARCH_PASSWORD"),
		CloudID:               "",
		APIKey:                "",
		Header:                map[string][]string{},
		CACert:                nil,
		RetryOnStatus:         []int{},
		DisableRetry:          false,
		EnableRetryOnTimeout:  false,
		MaxRetries:            0,
		DiscoverNodesOnStart:  false,
		DiscoverNodesInterval: 0,
		EnableMetrics:         false,
		EnableDebugLogger:     false,
		DisableMetaHeader:     false,
		RetryBackoff:          nil,
		Transport:             nil,
		Logger:                nil,
		Selector:              nil,
		ConnectionPoolFunc:    nil,
	}
	clt, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: clt,
	}, nil
}
