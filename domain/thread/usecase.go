package thread

import (
	"context"
	"errors"
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
		create(context.Context, repositoryCreateRequest) (Thread, error)
	}

	repositoryGetRequest struct {
		limit  int
		offset int
	}

	repositoryCreateRequest struct {
		title string
	}
)

func NewUsecase(repo repository) *Usecase {
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

func (u *Usecase) Create(ctx context.Context, req CreateRequest) (Thread, error) {
	if req.Title == "" {
		return Thread{}, errors.New("title is blank")
	}
	return u.repo.create(ctx, repositoryCreateRequest{
		title: req.Title,
	})
}
