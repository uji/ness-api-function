package elsch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/uji/ness-api-function/domain/thread"
	"github.com/uji/ness-api-function/reqctx"
)

var (
	ErrThreadNotFound = errors.New("Thread not found")
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

var _ thread.ElasticSearch = &Client{}

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
		Index:      string(c.threadIndexName),
		DocumentID: req.ID,
		Body:       strings.NewReader(string(bytes)),
		Refresh:    "true",
	}.Do(ctx, c.client)
	return err
}

func (c *Client) DeleteThread(ctx context.Context, threadID string) error {
	_, err := esapi.DeleteRequest{
		Index:      string(c.threadIndexName),
		DocumentID: threadID,
		Refresh:    "true",
	}.Do(ctx, c.client)
	return err
}

func (c *Client) FindThread(ctx context.Context, id string) (thread.ElasticSearchThreadRow, error) {
	ainfo, err := reqctx.GetAuthenticationInfo(ctx)
	if err != nil {
		return thread.ElasticSearchThreadRow{}, err
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match_phrase": map[string]string{
							"id": id,
						},
					},
					{
						"match_phrase": map[string]string{
							"teamID": ainfo.TeamID(),
						},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return thread.ElasticSearchThreadRow{}, err
	}
	size := 1

	res, err := esapi.SearchRequest{
		Index: []string{string(c.threadIndexName)},
		Body:  &buf,
		Size:  &size,
	}.Do(ctx, c.client)
	if err != nil {
		return thread.ElasticSearchThreadRow{}, err
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return thread.ElasticSearchThreadRow{}, err
	}

	rslt := new(struct {
		Hits struct {
			Hits []struct {
				ID     string `json:"_id"`
				Source struct {
					ID        string    `json:"id"`
					TeamID    string    `json:"teamID"`
					CreatorID string    `json:"creatorID"`
					Title     string    `json:"title"`
					Closed    bool      `json:"closed"`
					CreatedAt time.Time `json:"createdAt"`
					UpdatedAt time.Time `json:"updatedAt"`
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	})

	if err := json.Unmarshal(bytes, rslt); err != nil {
		return thread.ElasticSearchThreadRow{}, err
	}

	if len(rslt.Hits.Hits) == 0 {
		return thread.ElasticSearchThreadRow{}, ErrThreadNotFound
	}

	rows := make([]thread.ElasticSearchThreadRow, len(rslt.Hits.Hits))
	for i, row := range rslt.Hits.Hits {
		rows[i] = thread.ElasticSearchThreadRow{
			Id:        row.ID,
			TeamID:    row.Source.TeamID,
			CreaterID: row.Source.CreatorID,
			Title:     row.Source.Title,
			Closed:    row.Source.Closed,
			CreatedAt: row.Source.CreatedAt,
			UpdatedAt: row.Source.UpdatedAt,
		}
	}
	return rows[0], nil
}

func (c *Client) SearchThreads(
	ctx context.Context,
	req thread.SearchThreadIDsRequest,
	opts ...thread.SearchThreadIDsOption,
) ([]thread.ElasticSearchThreadRow, error) {
	ainfo, err := reqctx.GetAuthenticationInfo(ctx)
	if err != nil {
		return nil, err
	}

	must := make([]map[string]interface{}, 0, 3)
	must = append(must, map[string]interface{}{
		"match_phrase": map[string]string{
			"teamID": ainfo.TeamID(),
		},
	})

	if req.Word != "" {
		must = append(must, map[string]interface{}{
			"match": map[string]string{
				"title": req.Word,
			},
		})
	}

	for _, opt := range opts {
		switch opt {
		case thread.SearchThreadIDsOptionOnlyClosed:
			must = append(must, map[string]interface{}{
				"match_phrase": map[string]bool{
					"closed": true,
				},
			})
		case thread.SearchThreadIDsOptionOnlyOpened:
			must = append(must, map[string]interface{}{
				"match_phrase": map[string]bool{
					"closed": false,
				},
			})
		}
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := esapi.SearchRequest{
		Index: []string{string(c.threadIndexName)},
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
				ID     string `json:"_id"`
				Source struct {
					ID        string    `json:"id"`
					TeamID    string    `json:"teamID"`
					CreatorID string    `json:"creatorID"`
					Title     string    `json:"title"`
					Closed    bool      `json:"closed"`
					CreatedAt time.Time `json:"createdAt"`
					UpdatedAt time.Time `json:"updatedAt"`
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	})

	if err := json.Unmarshal(bytes, rslt); err != nil {
		return nil, err
	}

	rows := make([]thread.ElasticSearchThreadRow, len(rslt.Hits.Hits))
	for i, row := range rslt.Hits.Hits {
		rows[i] = thread.ElasticSearchThreadRow{
			Id:        row.ID,
			TeamID:    row.Source.TeamID,
			CreaterID: row.Source.CreatorID,
			Title:     row.Source.Title,
			Closed:    row.Source.Closed,
			CreatedAt: row.Source.CreatedAt,
			UpdatedAt: row.Source.UpdatedAt,
		}
	}
	return rows, nil
}
