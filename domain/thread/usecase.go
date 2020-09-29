package thread

import (
	"context"
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
	}

	repositoryGetRequest struct {
		limit           int64
		lastEvaluatedID *string
	}

	repositoryCreateRequest struct {
		thread Thread
	}
)

func NewUsecase(gen Generator, repo Repository) *Usecase {
	return &Usecase{gen, repo}
}

type GetRequest struct {
	Limit           *int
	LastEvaluatedID *string
}

func (u *Usecase) Get(ctx context.Context, req GetRequest) ([]Thread, error) {
	var l int64
	if req.Limit == nil {
		l = defaultLimit
	} else if *req.Limit < 1 {
		l = 1
	} else if *req.Limit > maxLimit {
		l = maxLimit
	} else {
		l = int64(*req.Limit)
	}

	return u.repo.get(ctx, repositoryGetRequest{
		limit:           l,
		lastEvaluatedID: req.LastEvaluatedID,
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
