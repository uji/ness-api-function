package thread

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

func TestUsecaseGet(t *testing.T) {
	cases := []struct {
		name      string
		reqLimit  int
		reqOffset int
		limit     int
		offset    int
	}{
		{"normal", 5, 5, 5, 5},
		{"limit too small", -1, 5, 1, 5},
		{"limit too big", 101, 5, 100, 5},
		{"offset too small", 5, -1, 5, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockrepository(ctrl)
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
					limit:  c.limit,
					offset: c.offset,
				},
			).Return(threads, nil)

			uc := NewUsecase(repo)
			res, err := uc.Get(context.Background(), GetRequest{
				Limit:  c.reqLimit,
				Offset: c.reqOffset,
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

func TestUsecaseCreate(t *testing.T) {
	cases := []struct {
		name          string
		title         string
		err           error
		useRepository bool
	}{
		{"normal", "thread1", nil, true},
		{"blank title", "", ErrorCreate01, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockrepository(ctrl)
			thrd := Thread{
				id:     "thread1",
				title:  "thread1",
				closed: false,
			}
			if c.useRepository {
				repo.EXPECT().create(
					context.Background(),
					repositoryCreateRequest{
						title: c.title,
					},
				).Return(thrd, nil)
			}

			uc := NewUsecase(repo)
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
			if diff := cmp.Diff(thrd, res, opts); diff != "" {
				t.Fatal(res)
			}
		})
	}
}