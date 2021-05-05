package elsch

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
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

type GetThreadsRequest struct {
	Size int
	From int
}

type GetThreadsOptions interface {
}

func (c *Client) GetThreadIDs(ctx context.Context, req GetThreadsRequest, opts ...GetThreadsOptions) ([]string, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := esapi.SearchRequest{
		Index: []string{threadIndexName},
		From:  &req.From,
		Body:  &buf,
		Size:  &req.Size,
	}.Do(ctx, c.client)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(res.Body) // TODO: use io.ReadAll
	if err != nil {
		return nil, err
	}

	rslt := new(struct {
		Hits struct {
			Hits []struct {
				ID string `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
	})

	if err := json.Unmarshal(bytes, rslt); err != nil {
		return nil, err
	}

	ids := make([]string, len(rslt.Hits.Hits))
	for i, h := range rslt.Hits.Hits {
		ids[i] = h.ID
	}
	return ids, nil
}
