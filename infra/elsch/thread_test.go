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
	cases := []struct {
		name string
		req  PutThreadRequest
	}{
		{
			name: "normal",
			req: PutThreadRequest{
				ID:        uuid.New().String(),
				TeamID:    uuid.New().String(),
				CreatorID: uuid.New().String(),
				Title:     "test",
				Closed:    true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
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
			if err := clt.PutThread(ctx, c.req); err != nil {
				t.Fatal(err)
			}

			res, err := esapi.GetRequest{
				Index:      "thread",
				DocumentID: c.req.ID,
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

			if diff := cmp.Diff(c.req, v.Source); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
