package elsch

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestPutThread(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	tid1 := uuid.New()
	tid2 := uuid.New()
	cid1 := uuid.New()
	cid2 := uuid.New()

	cases := []struct {
		name string
		reqs []PutThreadRequest
		expt PutThreadRequest
	}{
		{
			name: "put 2 document",
			reqs: []PutThreadRequest{
				{
					ID:        id.String(),
					TeamID:    tid1.String(),
					CreatorID: cid1.String(),
					Title:     "test1",
					Closed:    false,
					CreatedAt: time.Date(2021, 5, 4, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2021, 5, 4, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        uuid.New().String(),
					TeamID:    tid2.String(),
					CreatorID: cid2.String(),
					Title:     "test2",
					Closed:    true,
					CreatedAt: time.Date(2021, 5, 4, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2021, 5, 4, 12, 0, 0, 0, time.UTC),
				},
			},
			expt: PutThreadRequest{
				ID:        id.String(),
				TeamID:    tid1.String(),
				CreatorID: cid1.String(),
				Title:     "test1",
				Closed:    false,
				CreatedAt: time.Date(2021, 5, 4, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2021, 5, 4, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "put 2 times to same document",
			reqs: []PutThreadRequest{
				{
					ID:        id.String(),
					TeamID:    tid1.String(),
					CreatorID: cid1.String(),
					Title:     "test1",
					Closed:    false,
					CreatedAt: time.Date(2021, 5, 4, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2021, 5, 4, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        id.String(),
					TeamID:    tid2.String(),
					CreatorID: cid2.String(),
					Title:     "test2",
					Closed:    true,
					CreatedAt: time.Date(2021, 5, 4, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2021, 5, 4, 12, 0, 0, 0, time.UTC),
				},
			},
			expt: PutThreadRequest{
				ID:        id.String(),
				TeamID:    tid2.String(),
				CreatorID: cid2.String(),
				Title:     "test2",
				Closed:    true,
				CreatedAt: time.Date(2021, 5, 4, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2021, 5, 4, 12, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			clt, err := NewClient()
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.Background()
			for _, r := range c.reqs {
				if err := clt.PutThread(ctx, r); err != nil {
					t.Fatal(err)
				}
			}

			res, err := esapi.GetRequest{
				Index:      threadIndexName,
				DocumentID: id.String(),
			}.Do(ctx, clt.client)
			if err != nil {
				t.Fatal(err)
			}

			bytes, err := ioutil.ReadAll(res.Body) // TODO: use io.ReadAll
			if err != nil {
				t.Fatal(err)
			}

			v := new(struct {
				Source PutThreadRequest `json:"_source"`
			})
			if err := json.Unmarshal(bytes, v); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(c.expt, v.Source); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestDeleteThread(t *testing.T) {
	t.Parallel()

	id := uuid.New()

	clt, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// create test data
	_, err = esapi.IndexRequest{
		Index:      threadIndexName,
		DocumentID: id.String(),
		Body:       nil,
	}.Do(ctx, clt.client)
	if err != nil {
		t.Fatal(err)
	}

	if err := clt.DeleteThread(ctx, id.String()); err != nil {
		t.Fatal(err)
	}

	res, err := esapi.GetRequest{
		Index:      threadIndexName,
		DocumentID: id.String(),
	}.Do(ctx, clt.client)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 404 {
		t.Fatal(res)
	}
}
