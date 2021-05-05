package thread

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/guregu/dynamo"
	"github.com/guregu/null"
	"github.com/uji/ness-api-function/infra/db"
	"github.com/uji/ness-api-function/reqctx"
)

func TestRepoGet(t *testing.T) {
	cases := []struct {
		name       string
		items      []item
		offsetTime null.Time
		closed     null.Bool
		expt       []Thread
		err        error
	}{
		{
			name: "normal",
			items: []item{
				{
					PK:        "Team#0",
					SK:        "Thread#0",
					CreatorID: "UserID#0",
					Content:   "Thread0",
					Closed:    "false",
					CreatedAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
				},
				{
					PK:        "Team#0",
					SK:        "Thread#1",
					CreatorID: "UserID#1",
					Content:   "Thread1",
					Closed:    "true",
					CreatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					PK:        "Team#1",
					SK:        "Thread#3",
					CreatorID: "UserID#3",
					Content:   "Thread3",
					Closed:    "true",
					CreatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			offsetTime: null.Time{},
			expt: []Thread{
				&thread{
					id:        "Thread#1",
					teamID:    TeamID("Team#0"),
					createrID: UserID("UserID#1"),
					title:     "Thread1",
					closed:    true,
					createdAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
					updatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
				},
				&thread{
					id:        "Thread#0",
					teamID:    TeamID("Team#0"),
					createrID: UserID("UserID#0"),
					title:     "Thread0",
					closed:    false,
					createdAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
					updatedAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "limited and set lastEvaluatedID",
			items: []item{
				{
					PK:        "Team#0",
					SK:        "Thread#0",
					Content:   "Thread0",
					Closed:    "false",
					CreatedAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
				},
				{
					PK:        "Team#0",
					SK:        "Thread#1",
					Content:   "Thread1",
					Closed:    "true",
					CreatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			offsetTime: null.TimeFrom(time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC)),
			expt: []Thread{
				&thread{
					id:        "Thread#0",
					teamID:    TeamID("Team#0"),
					title:     "Thread0",
					closed:    false,
					createdAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
					updatedAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "filtered opened",
			items: []item{
				{
					PK:        "Team#0",
					SK:        "Thread#0",
					Content:   "Thread0",
					Closed:    "false",
					CreatedAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
				},
				{
					PK:        "Team#0",
					SK:        "Thread#1",
					Content:   "Thread1",
					Closed:    "true",
					CreatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			closed: null.NewBool(false, true),
			expt: []Thread{
				&thread{
					id:        "Thread#0",
					teamID:    TeamID("Team#0"),
					title:     "Thread0",
					closed:    false,
					createdAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
					updatedAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "filtered closed",
			items: []item{
				{
					PK:        "Team#0",
					SK:        "Thread#0",
					Content:   "Thread0",
					Closed:    "false",
					CreatedAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
				},
				{
					PK:        "Team#0",
					SK:        "Thread#1",
					Content:   "Thread1",
					Closed:    "true",
					CreatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			closed: null.NewBool(true, true),
			expt: []Thread{
				&thread{
					id:        "Thread#1",
					teamID:    TeamID("Team#0"),
					title:     "Thread1",
					closed:    true,
					createdAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
					updatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			dnmdb := db.NewDynamoDB()
			tbl := db.CreateThreadTestTable(dnmdb, t)
			defer db.DestroyTestTable(&tbl, t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			es := NewMockelasticsearch(ctrl)
			sut := NewDynamoRepository(dnmdb, tbl.Name(), es)

			for _, d := range c.items {
				if err := tbl.Put(d).Run(); err != nil {
					t.Fatal(err)
				}
			}

			ainfo := reqctx.NewAuthenticationInfo("Team#0", "")
			ctx := reqctx.NewRequestContext(context.Background(), ainfo)
			res, err := sut.get(ctx, repositoryGetRequest{
				offsetTime: c.offsetTime,
				closed:     c.closed,
			})
			if err != c.err {
				t.Fatal(err)
			}

			opt := cmp.AllowUnexported(thread{})
			if diff := cmp.Diff(c.expt, res, opt); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestRepo_find(t *testing.T) {
	cases := []struct {
		name  string
		items []item
		req   repositoryFindRequest
		expt  Thread
		err   error
	}{
		{
			name: "normal",
			items: []item{
				{
					PK:        "Team#0",
					SK:        "Thread#0",
					CreatorID: "UserID#0",
					Content:   "Thread0",
					Closed:    "false",
					CreatedAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
				},
				{
					PK:        "Team#0",
					SK:        "Thread#1",
					CreatorID: "UserID#1",
					Content:   "Thread1",
					Closed:    "true",
					CreatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					PK:        "Team#1",
					SK:        "Thread#3",
					CreatorID: "UserID#3",
					Content:   "Thread3",
					Closed:    "true",
					CreatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			req: repositoryFindRequest{
				teamID:   "Team#0",
				threadID: "Thread#1",
			},
			expt: &thread{
				id:        "Thread#1",
				teamID:    TeamID("Team#0"),
				createrID: UserID("UserID#1"),
				title:     "Thread1",
				closed:    true,
				createdAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
				updatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dnmdb := db.NewDynamoDB()
			tbl := db.CreateThreadTestTable(dnmdb, t)
			defer db.DestroyTestTable(&tbl, t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			es := NewMockelasticsearch(ctrl)
			sut := NewDynamoRepository(dnmdb, tbl.Name(), es)

			for _, d := range c.items {
				if err := tbl.Put(d).Run(); err != nil {
					t.Fatal(err)
				}
			}

			ainfo := reqctx.NewAuthenticationInfo("Team#0", "")
			ctx := reqctx.NewRequestContext(context.Background(), ainfo)
			res, err := sut.find(ctx, c.req)
			if err != c.err {
				t.Fatal(err)
			}

			opt := cmp.AllowUnexported(thread{})
			if diff := cmp.Diff(c.expt, res, opt); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestRepoCreate(t *testing.T) {
	dnmdb := db.NewDynamoDB()
	tbl := db.CreateThreadTestTable(dnmdb, t)
	defer db.DestroyTestTable(&tbl, t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	es := NewMockelasticsearch(ctrl)
	sut := NewDynamoRepository(dnmdb, tbl.Name(), es)

	thrd := thread{
		id:        "Thread#0",
		teamID:    "Team#0",
		createrID: "User#0",
		title:     "thread0",
		closed:    false,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}

	ctx := context.Background()
	es.EXPECT().PutThread(ctx, &thrd).Return(nil)

	if err := sut.create(
		ctx,
		repositoryCreateRequest{
			thread: &thrd,
		},
	); err != nil {
		t.Fatal(err)
	}

	var itm item
	if err := tbl.Get("PK", "Team#0").Range("SK", dynamo.Equal, "Thread#0").One(&itm); err != nil {
		t.Fatal(err)
	}
	opt := cmp.AllowUnexported(thread{})
	if diff := cmp.Diff(&thrd, itm.toThread(), opt); diff != "" {
		t.Fatal(diff)
	}
}

func TestRepo_update(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		cases := []struct {
			name  string
			items []item
			req   repositoryUpdateRequest
		}{
			{
				name: "1",
				items: []item{
					{
						PK:        "Team#0",
						SK:        "Thread#0",
						CreatorID: "User#0",
						Content:   "thread0",
						Closed:    "false",
						CreatedAt: time.Date(2020, 10, 2, 0, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2020, 10, 2, 0, 0, 0, 0, time.UTC),
					},
					{
						PK:        "Team#0",
						SK:        "Thread#1",
						CreatorID: "User#1",
						Content:   "thread1",
						Closed:    "true",
						CreatedAt: time.Date(2020, 10, 3, 0, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2020, 10, 3, 12, 0, 0, 0, time.UTC),
					},
				},
				req: repositoryUpdateRequest{
					thread: &thread{
						id:        "Thread#0",
						title:     "thread0",
						teamID:    "Team#0",
						createrID: "User#0",
						closed:    true,
						createdAt: time.Date(2020, 10, 2, 0, 0, 0, 0, time.UTC),
						updatedAt: time.Date(2020, 10, 2, 12, 0, 0, 0, time.UTC),
					},
				},
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				dnmdb := db.NewDynamoDB()
				tbl := db.CreateThreadTestTable(dnmdb, t)
				defer db.DestroyTestTable(&tbl, t)
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				es := NewMockelasticsearch(ctrl)
				sut := NewDynamoRepository(dnmdb, tbl.Name(), es)

				for _, itm := range c.items {
					if err := tbl.Put(itm).Run(); err != nil {
						t.Fatal(err)
					}
				}

				ctx := context.Background()
				es.EXPECT().PutThread(ctx, c.req.thread)

				err := sut.update(
					ctx,
					c.req,
				)
				if err != nil {
					t.Fatal(err)
				}

				var itm item
				if err := tbl.Get("PK", c.req.thread.TeamID()).Range("SK", dynamo.Equal, c.req.thread.ID()).One(&itm); err != nil {
					t.Fatal(err)
				}
				opt := cmp.AllowUnexported(thread{})
				if diff := cmp.Diff(c.req.thread, itm.toThread(), opt); diff != "" {
					t.Fatal(diff)
				}
			})
		}
	})

	t.Run("dont create when not found", func(t *testing.T) {
		dnmdb := db.NewDynamoDB()
		tbl := db.CreateThreadTestTable(dnmdb, t)
		defer db.DestroyTestTable(&tbl, t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		es := NewMockelasticsearch(ctrl)
		sut := NewDynamoRepository(dnmdb, tbl.Name(), es)

		ctx := context.Background()
		err := sut.update(
			ctx,
			repositoryUpdateRequest{
				thread: &thread{
					id:        "Thread#999",
					teamID:    "Team#999",
					createrID: "User#999",
					title:     "Title",
					closed:    true,
					createdAt: time.Now(),
					updatedAt: time.Now(),
				},
			},
		)
		if err.Error() != "ConditionalCheckFailedException: The conditional request failed" {
			t.Fatal(err)
		}
	})
}
