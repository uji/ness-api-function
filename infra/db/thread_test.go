package db

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/uji/ness-api-function/reqctx"
)

func TestThreadQueryGet(t *testing.T) {
	type responseThread struct {
		id        string
		teamID    string
		createrID string
		title     string
		closed    bool
		createdAt time.Time
		updatedAt time.Time
	}

	cases := []struct {
		name       string
		data       []threadSchema
		requestIDs []string
		expt       []responseThread
		err        error
	}{
		{
			name: "normal",
			data: []threadSchema{
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
			requestIDs: []string{"Thread#1", "Thread#0"},
			expt: []responseThread{
				{
					id:        "Thread#1",
					teamID:    "Team#0",
					createrID: "UserID#1",
					title:     "Thread1",
					closed:    true,
					createdAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
					updatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					id:        "Thread#0",
					teamID:    "Team#0",
					createrID: "UserID#0",
					title:     "Thread0",
					closed:    false,
					createdAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
					updatedAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			dnmdb := NewDynamoDB()
			tbl := CreateThreadTestTable(dnmdb, t)
			defer DestroyTestTable(&tbl, t)

			sut := NewThreadQuery(dnmdb, tbl.Name())

			for _, d := range c.data {
				if err := tbl.Put(d).Run(); err != nil {
					t.Fatal(err)
				}
			}

			ainfo := reqctx.NewAuthenticationInfo("Team#0", "User#0")
			ctx := reqctx.NewRequestContext(context.Background(), ainfo)
			res, err := sut.GetThreadsByIDs(ctx, c.requestIDs)
			if err != nil {
				t.Fatal(err)
			}

			rslt := make([]responseThread, len(res))
			for i, t := range res {
				rslt[i] = responseThread{
					id:        t.ID(),
					teamID:    string(t.TeamID()),
					createrID: string(t.CreatorID()),
					title:     t.Title(),
					closed:    t.Closed(),
					createdAt: t.CreatedAt(),
					updatedAt: t.UpdatedAt(),
				}
			}

			if diff := cmp.Diff(rslt, c.expt, cmp.AllowUnexported(responseThread{})); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
