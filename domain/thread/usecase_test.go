package thread

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/guregu/null"
)

func TestUsecaseGet(t *testing.T) {
	lst := time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC)
	lststr := lst.Format(time.RFC3339)

	cases := []struct {
		name     string
		req      GetRequest
		repoReq  repositoryGetRequest
		callRepo bool
		err      error
	}{
		{
			name:     "normal",
			req:      GetRequest{null.IntFrom(5), null.StringFrom(lststr)},
			repoReq:  repositoryGetRequest{5, null.TimeFrom(lst)},
			callRepo: true,
		},
		{
			name:     "limit too small",
			req:      GetRequest{null.IntFrom(-1), null.String{}},
			repoReq:  repositoryGetRequest{1, null.Time{}},
			callRepo: true,
		},
		{
			name:     "limit too big",
			req:      GetRequest{null.IntFrom(101), null.String{}},
			repoReq:  repositoryGetRequest{100, null.Time{}},
			callRepo: true,
		},
		{
			name:     "limit too big",
			req:      GetRequest{null.Int{}, null.String{}},
			repoReq:  repositoryGetRequest{30, null.Time{}},
			callRepo: true,
		},
		{
			name:     "last evaluated time format invalid",
			req:      GetRequest{null.IntFrom(5), null.StringFrom("test")},
			callRepo: false,
			err:      ErrorTimeFormatInValid,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gen := DefaultGenerator
			repo := NewMockRepository(ctrl)
			threads := []Thread{
				&thread{
					id:     "thread1",
					title:  "thread1",
					closed: false,
				},
				&thread{
					id:     "thread2",
					title:  "thread2",
					closed: true,
				},
			}

			if c.callRepo {
				repo.EXPECT().get(
					context.Background(),
					c.repoReq,
				).Return(threads, nil)
			}

			uc := NewUsecase(gen, repo)

			res, err := uc.Get(context.Background(), c.req)
			if err != c.err {
				t.Fatal(err)
			}
			if c.err != nil {
				return
			}

			opts := cmp.Options{
				cmp.AllowUnexported(thread{}),
			}
			if diff := cmp.Diff(threads, res, opts); diff != "" {
				t.Fatal(res)
			}
		})
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

			thrd := thread{
				id:     "thread1",
				title:  c.title,
				closed: false,
			}

			gen := func(attr ThreadAttribute) (Thread, error) {
				return &thrd, nil
			}

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
				cmp.AllowUnexported(thread{}),
			}
			if diff := cmp.Diff(&thrd, res, opts); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
