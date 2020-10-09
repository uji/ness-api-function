package thread

import (
	"context"
	"time"

	"github.com/guregu/null"
)

type (
	repositoryGetRequest struct {
		offsetTime null.Time
		closed     null.Bool
	}

	repositoryCreateRequest struct {
		thread Thread
	}

	repositoryUpdateRequest struct {
		thread Thread
	}

	repositoryOpenRequest struct {
		threadID string
	}

	repositoryCloseRequest struct {
		threadID string
	}

	Repository interface {
		get(context.Context, repositoryGetRequest) ([]Thread, error)
		create(context.Context, repositoryCreateRequest) (Thread, error)
		update(context.Context, repositoryUpdateRequest) (Thread, error)
		open(context.Context, repositoryOpenRequest) (Thread, error)
		close(context.Context, repositoryCloseRequest) (Thread, error)
	}

	Usecase struct {
		gen  Generator
		repo Repository
	}
)

func NewUsecase(gen Generator, repo Repository) *Usecase {
	return &Usecase{gen, repo}
}

type (
	GetRequest struct {
		OffsetTime null.String
		Closed     null.Bool
	}
	CreateRequest struct {
		Title string
	}
	OpenRequest struct {
		ThreadID string
	}
	CloseRequest struct {
		ThreadID string
	}
)

func (u *Usecase) Get(ctx context.Context, req GetRequest) ([]Thread, error) {
	var ofst null.Time
	if req.OffsetTime.Valid {
		t, err := time.Parse(time.RFC3339, req.OffsetTime.String)
		if err != nil {
			return nil, ErrorTimeFormatInValid
		}
		ofst = null.TimeFrom(t)
	}

	return u.repo.get(ctx, repositoryGetRequest{
		closed:     req.Closed,
		offsetTime: ofst,
	})
}

func (u *Usecase) Create(ctx context.Context, req CreateRequest) (Thread, error) {
	if req.Title == "" {
		return nil, ErrorTitleIsRequired
	}
	th, err := u.gen(ThreadAttribute(req))
	if err != nil {
		return nil, err
	}
	return u.repo.create(ctx, repositoryCreateRequest{
		thread: th,
	})
}

func (u *Usecase) Open(ctx context.Context, req OpenRequest) (Thread, error) {
	res, err := u.repo.open(ctx, repositoryOpenRequest{
		threadID: req.ThreadID,
	})
	return res, err
}

func (u *Usecase) Close(ctx context.Context, req CloseRequest) (Thread, error) {
	res, err := u.repo.close(ctx, repositoryCloseRequest{
		threadID: req.ThreadID,
	})
	return res, err
}
