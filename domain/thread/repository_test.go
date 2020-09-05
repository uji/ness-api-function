package thread

import (
	"context"
	"database/sql"
	"testing"

	"example.com/ness-api-function/infra/db"
	"github.com/google/go-cmp/cmp"
)

func Test_psqlRepositoryGet(t *testing.T) {
	db := db.NewDB()
	repo := NewPsqlRepository(db)

	cases := []struct {
		name    string
		limit   int
		offset  int
		threads []*Thread
	}{
		{"normal", 10, 0, []*Thread{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(*testing.T) {
			tx, err := db.BeginTxx(context.Background(), &sql.TxOptions{})
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback()

			res, err := repo.get(context.Background(), repositoryGetRequest{
				limit:  c.limit,
				offset: c.offset,
			})
			if err != nil {
				t.Fatal(err)
			}

			opt := cmp.AllowUnexported(Thread{})
			if diff := cmp.Diff(c.threads, res, opt); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
