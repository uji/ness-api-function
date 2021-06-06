package thread

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/guregu/null"
)

func Test_get(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		req    repositoryGetRequest
		esReq  SearchThreadIDsRequest
		esOpts []SearchThreadIDsOption
	}{
		{
			name: "normal",
			req: repositoryGetRequest{
				size: 5,
				from: 3,
				word: "test",
			},
			esReq: SearchThreadIDsRequest{
				Size: 5,
				From: 3,
				Word: "test",
			},
			esOpts: []SearchThreadIDsOption{},
		},
		{
			name: "size over limit",
			req: repositoryGetRequest{
				size: 101,
				from: 3,
				word: "test",
			},
			esReq: SearchThreadIDsRequest{
				Size: 100,
				From: 3,
				Word: "test",
			},
			esOpts: []SearchThreadIDsOption{},
		},
		{
			name: "closed true",
			req: repositoryGetRequest{
				closed: null.Bool{
					NullBool: sql.NullBool{
						Bool:  true,
						Valid: true,
					},
				},
				size: 5,
				from: 3,
				word: "test",
			},
			esReq: SearchThreadIDsRequest{
				Size: 5,
				From: 3,
				Word: "test",
			},
			esOpts: []SearchThreadIDsOption{
				SearchThreadIDsOptionOnlyClosed,
			},
		},
		{
			name: "closed true",
			req: repositoryGetRequest{
				closed: null.Bool{
					NullBool: sql.NullBool{
						Bool:  false,
						Valid: true,
					},
				},
				size: 5,
				from: 3,
				word: "test",
			},
			esReq: SearchThreadIDsRequest{
				Size: 5,
				From: 3,
				Word: "test",
			},
			esOpts: []SearchThreadIDsOption{
				SearchThreadIDsOptionOnlyOpened,
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			es := NewMockElasticSearch(ctrl)
			sut := NewRepository(es)

			ctx := context.Background()
			rows := []ElasticSearchThreadRow{
				{
					Id:        "Thread#0",
					TeamID:    "Team#0",
					CreaterID: "User#0",
					Title:     "Title",
					Closed:    true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}
			es.EXPECT().SearchThreads(ctx, c.esReq, c.esOpts).Return(rows, nil)
			res, err := sut.get(ctx, c.req)
			if err != nil {
				t.Fatal(err)
			}

			for i, r := range res {
				if diff := cmp.Diff(r, rows[i].toThread(), cmp.AllowUnexported(thread{})); diff != "" {
					t.Fatal(diff)
				}
			}
		})
	}

	t.Run("ids not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		es := NewMockElasticSearch(ctrl)
		sut := NewRepository(es)

		ctx := context.Background()
		rows := []ElasticSearchThreadRow{}
		req := SearchThreadIDsRequest{}
		es.EXPECT().SearchThreads(ctx, req).Return(rows, nil)
		res, err := sut.get(ctx, repositoryGetRequest{})
		if err != nil {
			t.Fatal(err)
		}

		if len(res) != 0 {
			t.Fatal(res)
		}
	})
}

func TestRepo_find(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	es := NewMockElasticSearch(ctrl)
	sut := NewRepository(es)

	ctx := context.Background()
	req := repositoryFindRequest{
		teamID:   "Team#0",
		threadID: "Thread#0",
	}
	row := ElasticSearchThreadRow{
		Id:        "Thread#0",
		TeamID:    "Team#0",
		CreaterID: "User#0",
		Title:     "Title",
		Closed:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	es.EXPECT().FindThread(ctx, "Thread#0").Return(row, nil)
	res, err := sut.find(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(res, row.toThread(), cmp.AllowUnexported(thread{})); diff != "" {
		t.Fatal(diff)
	}
}

func TestRepo_create(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	es := NewMockElasticSearch(ctrl)
	sut := NewRepository(es)

	ctx := context.Background()
	thrd := NewMockThread(ctrl)
	req := repositoryCreateRequest{
		thread: thrd,
	}
	es.EXPECT().PutThread(ctx, thrd).Return(nil)
	if err := sut.create(ctx, req); err != nil {
		t.Fatal(err)
	}
}

func TestRepo_update(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	es := NewMockElasticSearch(ctrl)
	sut := NewRepository(es)

	ctx := context.Background()
	thrd := NewMockThread(ctrl)
	req := repositoryUpdateRequest{
		thread: thrd,
	}
	es.EXPECT().PutThread(ctx, thrd).Return(nil)
	if err := sut.update(ctx, req); err != nil {
		t.Fatal(err)
	}
}
