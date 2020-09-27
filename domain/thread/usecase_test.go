package thread

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

func TestUsecaseGet(t *testing.T) {
	cases := []struct {
		name     string
		reqLimit sql.NullInt64
		limit    int64
	}{
		{"normal", nInt64(5), 5},
		{"limit too small", nInt64(-1), 1},
		{"limit too big", nInt64(101), 100},
		{"limit is null", sql.NullInt64{}, 30},
	}

	lastEvaluatedID := "lastEvaluatedID"

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gen := NewGeneratorConfigured()
			repo := NewMockRepository(ctrl)
			threads := []*Thread{
				{
					id:     "thread1",
					title:  "thread1",
					closed: false,
				},
				{
					id:     "thread2",
					title:  "thread2",
					closed: true,
				},
			}
			repo.EXPECT().get(
				context.Background(),
				repositoryGetRequest{
					limit:           c.limit,
					lastEvaluatedID: &lastEvaluatedID,
				},
			).Return(threads, nil)

			uc := NewUsecase(gen, repo)

			l := new(int)
			if c.reqLimit.Valid {
				*l = int(c.reqLimit.Int64)
			} else {
				l = nil
			}
			res, err := uc.Get(context.Background(), GetRequest{
				Limit:           l,
				LastEvaluatedID: &lastEvaluatedID,
			})
			if err != nil {
				t.Fatal(err)
			}

			opts := cmp.Options{
				cmp.AllowUnexported(Thread{}),
			}
			if diff := cmp.Diff(threads, res, opts); diff != "" {
				t.Fatal(res)
			}
		})
	}
}

func nInt64(n int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: n,
		Valid: true,
	}
}

func TestUsecaseCreate(t *testing.T) {
	cases := []struct {
		name          string
		title         string
		err           error
		useGenerator  bool
		useRepository bool
	}{
		{"normal", "thread1", nil, true, true},
		{"blank title", "", ErrorTitleIsRequired, false, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			thrd := Thread{
				id:     "thread1",
				title:  c.title,
				closed: false,
			}

			gen := NewGenerator(func(title string) (*Thread, error) {
				return &thrd, nil
			})

			repo := NewMockRepository(ctrl)
			if c.useRepository {
				repo.EXPECT().create(
					context.Background(),
					repositoryCreateRequest{
						thread: &thrd,
					},
				).Return(&thrd, nil)
			}

			uc := NewUsecase(gen, repo)
			res, err := uc.Create(context.Background(), CreateRequest{
				Title: c.title,
			})
			if err != c.err {
				t.Fatal(err)
			}
			if err != nil {
				return
			}

			opts := cmp.Options{
				cmp.AllowUnexported(Thread{}),
			}
			if diff := cmp.Diff(&thrd, res, opts); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
