package elsch

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/uji/ness-api-function/domain/thread"
)

func TestPutThread(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	tid1 := uuid.New()
	tid2 := uuid.New()
	cid1 := uuid.New()
	cid2 := uuid.New()

	type putThread struct {
		id        string
		teamID    string
		creatorID string
		title     string
		closed    bool
		createdAt time.Time
		updatedAt time.Time
	}

	cases := []struct {
		name string
		reqs []putThread
		expt putThreadRequest
	}{
		{
			name: "put 2 document",
			reqs: []putThread{
				{
					id:        id.String(),
					teamID:    tid1.String(),
					creatorID: cid1.String(),
					title:     "test1",
					closed:    false,
					createdAt: time.Date(2021, 5, 4, 0, 0, 0, 0, time.UTC),
					updatedAt: time.Date(2021, 5, 4, 0, 0, 0, 0, time.UTC),
				},
				{
					id:        uuid.New().String(),
					teamID:    tid2.String(),
					creatorID: cid2.String(),
					title:     "test2",
					closed:    true,
					createdAt: time.Date(2021, 5, 4, 0, 0, 0, 0, time.UTC),
					updatedAt: time.Date(2021, 5, 4, 12, 0, 0, 0, time.UTC),
				},
			},
			expt: putThreadRequest{
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
			reqs: []putThread{
				{
					id:        id.String(),
					teamID:    tid1.String(),
					creatorID: cid1.String(),
					title:     "test1",
					closed:    false,
					createdAt: time.Date(2021, 5, 4, 0, 0, 0, 0, time.UTC),
					updatedAt: time.Date(2021, 5, 4, 0, 0, 0, 0, time.UTC),
				},
				{
					id:        id.String(),
					teamID:    tid2.String(),
					creatorID: cid2.String(),
					title:     "test2",
					closed:    true,
					createdAt: time.Date(2021, 5, 4, 0, 0, 0, 0, time.UTC),
					updatedAt: time.Date(2021, 5, 4, 12, 0, 0, 0, time.UTC),
				},
			},
			expt: putThreadRequest{
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
			ctrl := gomock.NewController(t)

			clt, err := NewClient()
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.Background()
			for _, r := range c.reqs {
				th := thread.NewMockThread(ctrl)
				th.EXPECT().ID().Return(r.id)
				th.EXPECT().TeamID().Return(thread.TeamID(r.teamID))
				th.EXPECT().CreatorID().Return(thread.UserID(r.creatorID))
				th.EXPECT().Closed().Return(r.closed)
				th.EXPECT().Title().Return(r.title)
				th.EXPECT().CreatedAt().Return(r.createdAt)
				th.EXPECT().UpdatedAt().Return(r.updatedAt)
				if err := clt.PutThread(ctx, th); err != nil {
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
				Source putThreadRequest `json:"_source"`
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
