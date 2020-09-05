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
			res, err := uc.Get(context.Background(), UsecaseGetRequest{
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
