package thread

import (
	"context"
)

const (
	maxlimit = 100
)

type (
	Usecase struct {
		repo repository
	}

	repository interface {
		get(context.Context, repositoryGetRequest) ([]*Thread, error)
	}

	repositoryGetRequest struct {
		limit  int
		offset int
	}
)

func NewUsecase(repo repository) *Usecase {
	return &Usecase{repo}
}

type UsecaseGetRequest struct {
	Limit  int
	Offset int
}

func (u *Usecase) Get(ctx context.Context, req UsecaseGetRequest) ([]*Thread, error) {
	l := req.Limit
	if req.Limit < 1 {
		l = 1
	} else if req.Limit > maxlimit {
		l = maxlimit
	}

	o := req.Offset
	if req.Offset < 1 {
		o = 1
	}

	return u.repo.get(ctx, repositoryGetRequest{
		limit:  l,
		offset: o,
	})
}
