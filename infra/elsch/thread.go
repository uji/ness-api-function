package elsch

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

var (
	threadIndexName = "thread"
)

type PutThreadRequest struct {
	ID        string    `json:"id"`
	TeamID    string    `json:"teamID"`
	CreatorID string    `json:"creatorID"`
	Title     string    `json:"title"`
	Closed    bool      `json:"closed"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (c *Client) PutThread(ctx context.Context, req PutThreadRequest) error {
	bytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	_, err = esapi.IndexRequest{
		Index:      threadIndexName,
		DocumentID: req.ID,
		Body:       strings.NewReader(string(bytes)),
		Refresh:    "true",
	}.Do(ctx, c.client)
	return err
}
