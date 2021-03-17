package thread

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/guregu/null"
	"github.com/uji/ness-api-function/reqctx"
)

func TestUsecaseGet(t *testing.T) {
	offsetTime := time.Date(2020, 9, 30, 0, 0, 0, 0, time.UTC)
	offsetTimeStr := offsetTime.Format(time.RFC3339)

	cases := []struct {
		name     string
		req      GetRequest
		repoReq  repositoryGetRequest
		callRepo bool
		err      error
	}{
		{
			name:     "normal",
			req:      GetRequest{null.StringFrom(offsetTimeStr), null.NewBool(true, true)},
			repoReq:  repositoryGetRequest{null.NewTime(offsetTime, true), null.NewBool(true, true)},
			callRepo: true,
		},
		{
			name:     "last evaluated time format invalid",
			req:      GetRequest{null.StringFrom("test"), null.Bool{}},
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

			gen := func(attr threadAttribute) (Thread, error) {
				return &thrd, nil
			}

			ainfo := reqctx.NewAuthenticationInfo("Team#0", "User#0")
			ctx := reqctx.NewRequestContext(context.Background(), ainfo)
			repo := NewMockRepository(ctrl)
			if c.useRepository {
				repo.EXPECT().create(
					ctx,
					repositoryCreateRequest{
						thread: &thrd,
					},
				).Return(nil)
			}

			uc := NewUsecase(gen, repo)
			res, err := uc.Create(ctx, CreateRequest{
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

func TestUsecase_Open(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockRepository(ctrl)
	uc := NewUsecase(DefaultGenerator, repo)

	thrd := thread{
		id:     "thread",
		closed: true,
	}
	ainfo := reqctx.NewAuthenticationInfo("Team#0", "User#0")
	ctx := reqctx.NewRequestContext(context.Background(), ainfo)
	repo.EXPECT().find(ctx, repositoryFindRequest{
		threadID: "thread",
		teamID:   "Team#0",
	}).Return(&thrd, nil)
	repo.EXPECT().update(ctx, repositoryUpdateRequest{
		thread: &thrd,
	}).Return(nil)
	res, err := uc.Open(ctx, OpenRequest{
		ThreadID: "thread",
	})
	if err != nil {
		t.Fatal(err)
	}
	opts := cmp.Options{
		cmp.AllowUnexported(thread{}),
	}
	if diff := cmp.Diff(&thrd, res, opts); diff != "" {
		t.Fatal(diff)
	}
}
func TestUsecase_OpenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockRepository(ctrl)
	uc := NewUsecase(DefaultGenerator, repo)

	thrd := thread{
		id: "thread",
	}
	terr := errors.New("test")
	ainfo := reqctx.NewAuthenticationInfo("Team#0", "")
	ctx := reqctx.NewRequestContext(context.Background(), ainfo)
	repo.EXPECT().find(ctx, repositoryFindRequest{
		threadID: "thread",
		teamID:   "Team#0",
	}).Return(&thrd, terr)
	if _, err := uc.Open(ctx, OpenRequest{
		ThreadID: "thread",
	}); err != terr {
		t.Fatal(err)
	}
}

func TestUsecase_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockRepository(ctrl)
	uc := NewUsecase(DefaultGenerator, repo)

	thrd := thread{
		id:     "thread",
		closed: false,
	}
	ainfo := reqctx.NewAuthenticationInfo("Team#0", "")
	ctx := reqctx.NewRequestContext(context.Background(), ainfo)
	repo.EXPECT().find(ctx, repositoryFindRequest{
		threadID: "thread",
		teamID:   "Team#0",
	}).Return(&thrd, nil)
	repo.EXPECT().update(ctx, repositoryUpdateRequest{thread: &thrd}).Return(nil)
	res, err := uc.Close(ctx, CloseRequest{
		ThreadID: "thread",
	})
	if err != nil {
		t.Fatal(err)
	}
	opts := cmp.Options{
		cmp.AllowUnexported(thread{}),
	}
	if diff := cmp.Diff(&thrd, res, opts); diff != "" {
		t.Fatal(diff)
	}
}
func TestUsecase_CloseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockRepository(ctrl)
	uc := NewUsecase(DefaultGenerator, repo)

	thrd := thread{
		id: "thread",
	}
	terr := errors.New("test")
	ainfo := reqctx.NewAuthenticationInfo("Team#0", "")
	ctx := reqctx.NewRequestContext(context.Background(), ainfo)
	repo.EXPECT().find(ctx, repositoryFindRequest{
		threadID: "thread",
		teamID:   "Team#0",
	}).Return(&thrd, terr)
	if _, err := uc.Close(ctx, CloseRequest{
		ThreadID: "thread",
	}); err != terr {
		t.Fatal(err)
	}
}
