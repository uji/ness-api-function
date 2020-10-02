package thread

import (
	"context"
	"time"

	"github.com/guregu/null"
)

const (
	defaultLimit = 30
	maxLimit     = 100
)

type (
	Usecase struct {
		gen  Generator
		repo Repository
	}

	Repository interface {
		get(context.Context, repositoryGetRequest) ([]Thread, error)
		create(context.Context, repositoryCreateRequest) (Thread, error)
		update(context.Context, repositoryUpdateRequest) (Thread, error)
		close(context.Context, repositoryCloseRequest) (Thread, error)
	}

	repositoryGetRequest struct {
		limit             int64
		lastEvaluatedTime null.Time
	}

	repositoryCreateRequest struct {
		thread Thread
	}

	repositoryUpdateRequest struct {
		thread Thread
	}

	repositoryCloseRequest struct {
		threadID string
	}
)

func NewUsecase(gen Generator, repo Repository) *Usecase {
	return &Usecase{gen, repo}
}

type GetRequest struct {
	Limit             null.Int
	LastEvaluatedTime null.String
}

func (u *Usecase) Get(ctx context.Context, req GetRequest) ([]Thread, error) {
	l := req.Limit.Int64
	if !req.Limit.Valid {
		l = defaultLimit
	} else if l < 1 {
		l = 1
	} else if l > maxLimit {
		l = maxLimit
	}

	var lst null.Time
	if req.LastEvaluatedTime.Valid {
		t, err := time.Parse(time.RFC3339, req.LastEvaluatedTime.String)
		if err != nil {
			return nil, ErrorTimeFormatInValid
		}
		lst = null.TimeFrom(t)
	}

	return u.repo.get(ctx, repositoryGetRequest{
		limit:             l,
		lastEvaluatedTime: lst,
	})
}

type (
	CreateRequest struct {
		Title string
	}
)

func (u *Usecase) Create(ctx context.Context, req CreateRequest) (Thread, error) {
	if req.Title == "" {
		return nil, ErrorTitleIsRequired
	}
	th, err := u.gen(ThreadAttribute{
		Title: req.Title,
	})
	if err != nil {
		return nil, err
	}
	return u.repo.create(ctx, repositoryCreateRequest{
		thread: th,
	})
}

type CloseRequest struct {
	ThreadID string
}

func (u *Usecase) Close(ctx context.Context, req CloseRequest) (Thread, error) {
	res, err := u.repo.close(ctx, repositoryCloseRequest{
		threadID: req.ThreadID,
	})
	return res, err
}
