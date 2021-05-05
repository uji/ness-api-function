package elsch

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/uji/ness-api-function/domain/thread"
)

var (
	threadIndexName = "thread"
)

type putThreadRequest struct {
	ID        string    `json:"id"`
	TeamID    string    `json:"teamID"`
	CreatorID string    `json:"creatorID"`
	Title     string    `json:"title"`
	Closed    bool      `json:"closed"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (c *Client) PutThread(ctx context.Context, thread thread.Thread) error {
	req := putThreadRequest{
		ID:        thread.ID(),
		TeamID:    string(thread.TeamID()),
		CreatorID: string(thread.CreatorID()),
		Title:     thread.Title(),
		Closed:    thread.Closed(),
		CreatedAt: thread.CreatedAt(),
		UpdatedAt: thread.UpdatedAt(),
	}
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

func (c *Client) DeleteThread(ctx context.Context, threadID string) error {
	_, err := esapi.DeleteRequest{
		Index:      threadIndexName,
		DocumentID: threadID,
		Refresh:    "true",
	}.Do(ctx, c.client)
	return err
}
