package thread

import (
	"context"
	"testing"

	"example.com/ness-api-function/infra/db"
	"example.com/ness-api-function/tools/dbtool"
	"github.com/google/go-cmp/cmp"
)

func TestRepoGet(t *testing.T) {
	db := db.NewDynamoDB()
	tbl := dbtool.CreateThreadTestTable(db, t)
	defer tbl.DeleteTable().Run()

	sut := NewDynamoRepository(db, tbl.Name())

	data := []item{
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
	}

	expt := []*Thread{
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
	}

	for _, d := range data {
		if err := tbl.Put(d).Run(); err != nil {
			t.Fatal(err)
		}
	}

	res, err := sut.get(context.Background(), repositoryGetRequest{
		limit:  0,
		offset: 0,
	})
	if err != nil {
		t.Fatal(err)
	}

	opt := cmp.AllowUnexported(Thread{})
	if diff := cmp.Diff(expt, res, opt); diff != "" {
		t.Fatal(diff)
	}
}

func TestRepoCreate(t *testing.T) {
	db := db.NewDynamoDB()
	tbl := dbtool.CreateThreadTestTable(db, t)
	defer tbl.DeleteTable().Run()

	sut := NewDynamoRepository(db, tbl.Name())

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
