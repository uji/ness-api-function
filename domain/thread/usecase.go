package thread

import (
	"context"
)

const (
	maxlimit = 100
)

type (
	Usecase struct {
		repo Repository
	}

	Repository interface {
		get(context.Context, repositoryGetRequest) ([]*Thread, error)
		create(context.Context, repositoryCreateRequest) (*Thread, error)
	}

	repositoryGetRequest struct {
		limit  int
		offset int
	}

	repositoryCreateRequest struct {
		thread *Thread
	}
)

func NewUsecase(repo Repository) *Usecase {
	return &Usecase{repo}
}

type GetRequest struct {
	Limit  int
	Offset int
}

func (u *Usecase) Get(ctx context.Context, req GetRequest) ([]*Thread, error) {
	l := req.Limit
	if req.Limit < 1 {
		l = 1
	} else if req.Limit > maxlimit {
		l = maxlimit
	}

	o := req.Offset
	if req.Offset < 0 {
		o = 0
	}

	return u.repo.get(ctx, repositoryGetRequest{
		limit:  l,
		offset: o,
	})
}

type (
	CreateRequest struct {
		Title string
	}
)

func (u *Usecase) Create(ctx context.Context, req CreateRequest) (*Thread, error) {
	if req.Title == "" {
		return nil, ErrorCreate01
	}
	th, err := NewThread(req.Title)
	if err != nil {
		return nil, err
	}
	return u.repo.create(ctx, repositoryCreateRequest{
		thread: th,
	})
}
