package thread

import (
	"context"
	"testing"

	"example.com/ness-api-function/infra/db"
	"github.com/google/go-cmp/cmp"
)

func TestRepoGet(t *testing.T) {
	cases := []struct {
		name   string
		items  []item
		limit  int
		offset int
		expt   []*Thread
	}{
		{
			name: "normal",
			items: []item{
				{
					PK:      "Team#0",
					SK:      "Thread#0",
					Content: "Thread0",
					Closed:  "false",
				},
				{
					PK:      "Team#0",
					SK:      "Thread#1",
					Content: "Thread1",
					Closed:  "true",
				},
			},
			limit:  5,
			offset: 0,
			expt: []*Thread{
				{
					id:     "Thread#0",
					title:  "Thread0",
					closed: false,
				},
				{
					id:     "Thread#1",
					title:  "Thread1",
					closed: true,
				},
			},
		},
		{
			name: "limited",
			items: []item{
				{
					PK:      "Team#0",
					SK:      "Thread#0",
					Content: "Thread0",
					Closed:  "false",
				},
				{
					PK:      "Team#0",
					SK:      "Thread#1",
					Content: "Thread1",
					Closed:  "true",
				},
			},
			limit:  1,
			offset: 0,
			expt: []*Thread{
				{
					id:     "Thread#0",
					title:  "Thread0",
					closed: false,
				},
			},
		},
		// {
		// 	name: "limited and set offset",
		// 	items: []item{
		// 		{
		// 			PK:      "Team#0",
		// 			SK:      "Thread#0",
		// 			Content: "Thread0",
		// 			Closed:  "false",
		// 		},
		// 		{
		// 			PK:      "Team#0",
		// 			SK:      "Thread#1",
		// 			Content: "Thread1",
		// 			Closed:  "true",
		// 		},
		// 	},
		// 	limit:  1,
		// 	offset: 1,
		// 	expt: []*Thread{
		// 		{
		// 			id:     "Thread#1",
		// 			title:  "Thread1",
		// 			closed: true,
		// 		},
		// 	},
		// },
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dnmdb := db.NewDynamoDB()
			tbl := db.CreateThreadTestTable(dnmdb, t)
			defer tbl.DeleteTable().Run()

			sut := NewDynamoRepository(dnmdb, tbl.Name())

			for _, d := range c.items {
				if err := tbl.Put(d).Run(); err != nil {
					t.Fatal(err)
				}
			}

			res, err := sut.get(context.Background(), repositoryGetRequest{
				limit:  c.limit,
				offset: c.offset,
			})
			if err != nil {
				t.Fatal(err)
			}

			opt := cmp.AllowUnexported(Thread{})
			if diff := cmp.Diff(c.expt, res, opt); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestRepoCreate(t *testing.T) {
	dnmdb := db.NewDynamoDB()
	tbl := db.CreateThreadTestTable(dnmdb, t)
	defer tbl.DeleteTable().Run()

	sut := NewDynamoRepository(dnmdb, tbl.Name())

	thrd := Thread{
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

	opt := cmp.AllowUnexported(Thread{})
	if diff := cmp.Diff(&thrd, res, opt); diff != "" {
		t.Fatal(diff)
	}
}
