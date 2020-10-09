package thread

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/guregu/null"
	"github.com/uji/ness-api-function/infra/db"
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
					Content:   "Thread0",
					Closed:    "false",
					CreatedAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
				},
				{
					PK:        "Team#0",
					SK:        "Thread#1",
					Content:   "Thread1",
					Closed:    "true",
					CreatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			offsetTime: null.Time{},
			expt: []Thread{
				&thread{
					id:        "Thread#1",
					title:     "Thread1",
					closed:    true,
					createdAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
				},
				&thread{
					id:        "Thread#0",
					title:     "Thread0",
					closed:    false,
					createdAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
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
				},
				{
					PK:        "Team#0",
					SK:        "Thread#1",
					Content:   "Thread1",
					Closed:    "true",
					CreatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			offsetTime: null.TimeFrom(time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC)),
			expt: []Thread{
				&thread{
					id:        "Thread#0",
					title:     "Thread0",
					closed:    false,
					createdAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
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
				},
				{
					PK:        "Team#0",
					SK:        "Thread#1",
					Content:   "Thread1",
					Closed:    "true",
					CreatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			closed: null.NewBool(false, true),
			expt: []Thread{
				&thread{
					id:        "Thread#0",
					title:     "Thread0",
					closed:    false,
					createdAt: time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC),
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
				},
				{
					PK:        "Team#0",
					SK:        "Thread#1",
					Content:   "Thread1",
					Closed:    "true",
					CreatedAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			closed: null.NewBool(true, true),
			expt: []Thread{
				&thread{
					id:        "Thread#1",
					title:     "Thread1",
					closed:    true,
					createdAt: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dnmdb := db.NewDynamoDB()
			tbl := db.CreateThreadTestTable(dnmdb, t)
			defer db.DestroyTestTable(&tbl, t)

			sut := NewDynamoRepository(dnmdb, tbl.Name())

			for _, d := range c.items {
				if err := tbl.Put(d).Run(); err != nil {
					t.Fatal(err)
				}
			}

			res, err := sut.get(context.Background(), repositoryGetRequest{
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

func TestRepoCreate(t *testing.T) {
	dnmdb := db.NewDynamoDB()
	tbl := db.CreateThreadTestTable(dnmdb, t)
	defer db.DestroyTestTable(&tbl, t)

	sut := NewDynamoRepository(dnmdb, tbl.Name())

	thrd := thread{
		id:     "thread1",
		title:  "thread1",
		closed: false,
	}
	res, err := sut.create(
		context.Background(),
		repositoryCreateRequest{
			thread: &thrd,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	opt := cmp.AllowUnexported(thread{})
	if diff := cmp.Diff(&thrd, res, opt); diff != "" {
		t.Fatal(diff)
	}
}

func testRepo_update(
	t *testing.T,
	items []item,
	req repositoryUpdateRequest,
	errMsg string,
) {
	dnmdb := db.NewDynamoDB()
	tbl := db.CreateThreadTestTable(dnmdb, t)
	defer db.DestroyTestTable(&tbl, t)

	sut := NewDynamoRepository(dnmdb, tbl.Name())

	for _, itm := range items {
		if err := tbl.Put(itm).Run(); err != nil {
			t.Fatal(err)
		}
	}

	res, err := sut.update(
		context.Background(),
		req,
	)
	if errMsg == "" {
		if err != nil {
			t.Fatal(err)
		}
	} else {
		if err.Error() != errMsg {
			t.Fatal(err)
		}
		return
	}

	opt := cmp.AllowUnexported(thread{})
	if diff := cmp.Diff(req.thread, res, opt); diff != "" {
		t.Fatal(diff)
	}
}

func TestRepo_update(t *testing.T) {
	cases := []struct {
		name   string
		items  []item
		req    repositoryUpdateRequest
		errMsg string
	}{
		{
			name: "normal",
			items: []item{
				{
					PK:        "Team#0",
					SK:        "Thread#0",
					Content:   "thread0",
					Closed:    "false",
					CreatedAt: time.Date(2020, 10, 2, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 10, 2, 0, 0, 0, 0, time.UTC),
				},
				{
					PK:        "Team#0",
					SK:        "Thread#1",
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
					closed:    true,
					createdAt: time.Date(2020, 10, 2, 0, 0, 0, 0, time.UTC),
					updatedAt: time.Date(2020, 10, 2, 12, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "dont create when not found",
			items: []item{
				{
					PK:        "Team#0",
					SK:        "Thread#1",
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
					closed:    true,
					createdAt: time.Date(2020, 10, 2, 0, 0, 0, 0, time.UTC),
					updatedAt: time.Date(2020, 10, 2, 12, 0, 0, 0, time.UTC),
				},
			},
			errMsg: "ConditionalCheckFailedException: The conditional request failed",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			testRepo_update(t, c.items, c.req, c.errMsg)
		})
	}
}
