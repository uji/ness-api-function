package elsch

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/uji/ness-api-function/domain/thread"
	"github.com/uji/ness-api-function/reqctx"
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
			defer ctrl.Finish()

			clt, err := NewClient(IndexName(uuid.New().String()))
			if err != nil {
				t.Fatal(err)
			}
			CreateIndexForTest(t, clt)
			defer DeleteIndexForTest(t, clt)

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
				Index:      string(clt.threadIndexName),
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

	clt, err := NewClient(IndexName(uuid.New().String()))
	if err != nil {
		t.Fatal(err)
	}
	CreateIndexForTest(t, clt)
	defer DeleteIndexForTest(t, clt)

	ctx := context.Background()

	// create test data
	_, err = esapi.IndexRequest{
		Index:      string(clt.threadIndexName),
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
		Index:      string(clt.threadIndexName),
		DocumentID: id.String(),
	}.Do(ctx, clt.client)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 404 {
		t.Fatal(res)
	}
}

func TestGetThreads(t *testing.T) {
	type threadSchema struct {
		ID        string    `json:"id"`
		TeamID    string    `json:"teamID"`
		CreatorID string    `json:"creatorID"`
		Title     string    `json:"title"`
		Closed    bool      `json:"closed"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	id1 := uuid.New().String()
	id2 := uuid.New().String()
	id3 := uuid.New().String()
	myTeamID := uuid.New().String()
	tid2 := uuid.New().String()
	myUserID := uuid.New().String()
	uid2 := uuid.New().String()

	cases := []struct {
		name string
		data []threadSchema
		req  thread.SearchThreadIDsRequest
		opts []thread.SearchThreadIDsOption
		res  []string
	}{
		{
			name: "normal",
			data: []threadSchema{
				{
					ID:        id1,
					TeamID:    myTeamID,
					CreatorID: myUserID,
					Title:     "test",
					Closed:    true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			req: thread.SearchThreadIDsRequest{
				Size: 10,
				From: 0,
			},
			res: []string{id1},
		},
		{
			name: "specify word",
			data: []threadSchema{
				{
					ID:        id1,
					TeamID:    myTeamID,
					CreatorID: myUserID,
					Title:     "test",
					Closed:    true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        id2,
					TeamID:    myTeamID,
					CreatorID: myUserID,
					Title:     "test thread",
					Closed:    true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        id3,
					TeamID:    myTeamID,
					CreatorID: myUserID,
					Title:     "xxx",
					Closed:    true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			req: thread.SearchThreadIDsRequest{
				Size: 10,
				From: 0,
				Word: "test",
			},
			res: []string{id1, id2},
		},
		{
			name: "use closed only option",
			data: []threadSchema{
				{
					ID:        id1,
					TeamID:    myTeamID,
					CreatorID: myUserID,
					Title:     "test1",
					Closed:    true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        id2,
					TeamID:    myTeamID,
					CreatorID: myUserID,
					Title:     "test2",
					Closed:    true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        id3,
					TeamID:    myTeamID,
					CreatorID: myUserID,
					Title:     "test3",
					Closed:    false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			req: thread.SearchThreadIDsRequest{
				Size: 10,
				From: 0,
				Word: "",
			},
			opts: []thread.SearchThreadIDsOption{
				thread.SearchThreadIDsOptionOnlyClosed,
			},
			res: []string{id1, id2},
		},
		{
			name: "use opened only option",
			data: []threadSchema{
				{
					ID:        id1,
					TeamID:    myTeamID,
					CreatorID: myUserID,
					Title:     "test1",
					Closed:    true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        id2,
					TeamID:    myTeamID,
					CreatorID: myUserID,
					Title:     "test2",
					Closed:    true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        id3,
					TeamID:    myTeamID,
					CreatorID: myUserID,
					Title:     "test3",
					Closed:    false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			req: thread.SearchThreadIDsRequest{
				Size: 10,
				From: 0,
				Word: "",
			},
			opts: []thread.SearchThreadIDsOption{
				thread.SearchThreadIDsOptionOnlyOpened,
			},
			res: []string{id3},
		},
		{
			name: "size is smaller than total",
			data: []threadSchema{
				{
					ID:        id1,
					TeamID:    myTeamID,
					CreatorID: myUserID,
					Title:     "test1",
					Closed:    true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        id2,
					TeamID:    myTeamID,
					CreatorID: myUserID,
					Title:     "test2",
					Closed:    false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			req: thread.SearchThreadIDsRequest{
				Size: 1,
				From: 1,
			},
			res: []string{id2},
		},
		{
			name: "include other users thread",
			data: []threadSchema{
				{
					ID:        id1,
					TeamID:    myTeamID,
					CreatorID: myUserID,
					Title:     "test1",
					Closed:    true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        id2,
					TeamID:    myTeamID,
					CreatorID: uid2,
					Title:     "test2",
					Closed:    true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			req: thread.SearchThreadIDsRequest{
				Size: 10,
				From: 0,
			},
			res: []string{id1, id2},
		},
		{
			name: "include other teams thread",
			data: []threadSchema{
				{
					ID:        id1,
					TeamID:    myTeamID,
					CreatorID: myUserID,
					Title:     "test1",
					Closed:    true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        id2,
					TeamID:    tid2,
					CreatorID: uid2,
					Title:     "test2",
					Closed:    true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			req: thread.SearchThreadIDsRequest{
				Size: 10,
				From: 0,
			},
			res: []string{id1},
		},
		{
			name: "no data",
			data: []threadSchema{},
			req: thread.SearchThreadIDsRequest{
				Size: 10,
				From: 0,
			},
			res: []string{},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			clt, err := NewClient(IndexName(uuid.New().String()))
			if err != nil {
				t.Fatal(err)
			}
			CreateIndexForTest(t, clt)
			defer DeleteIndexForTest(t, clt)

			ctx := reqctx.NewRequestContext(
				context.Background(),
				reqctx.NewAuthenticationInfo(myTeamID, myUserID),
			)

			// create test data
			for _, d := range c.data {
				bytes, err := json.Marshal(d)
				if err != nil {
					t.Fatal(err)
				}

				_, err = esapi.IndexRequest{
					Index:      string(clt.threadIndexName),
					DocumentID: d.ID,
					Body:       strings.NewReader(string(bytes)),
					Refresh:    "true",
				}.Do(ctx, clt.client)
				if err != nil {
					t.Fatal(err)
				}
			}

			res, err := clt.GetThreadIDs(ctx, c.req, c.opts...)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(c.res, res); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
