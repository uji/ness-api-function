package thread

import (
	"context"

	"github.com/guregu/null"
	"github.com/uji/ness-api-function/reqctx"
)

type (
	repositoryGetRequest struct {
		closed null.Bool
		size   int
		from   int
		word   string
	}

	repositoryFindRequest struct {
		teamID   TeamID
		threadID string
	}

	repositoryCreateRequest struct {
		thread Thread
	}

	repositoryUpdateRequest struct {
		thread Thread
	}

	Repository interface {
		get(context.Context, repositoryGetRequest) ([]Thread, error)
		find(context.Context, repositoryFindRequest) (Thread, error)
		create(context.Context, repositoryCreateRequest) error
		update(context.Context, repositoryUpdateRequest) error
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
		Size       int
		From       int
		Word       string
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
	return u.repo.get(ctx, repositoryGetRequest{
		closed: req.Closed,
		size:   req.Size,
		from:   req.From,
		word:   req.Word,
	})
}

func (u *Usecase) Create(ctx context.Context, req CreateRequest) (Thread, error) {
	if req.Title == "" {
		return nil, ErrorTitleIsRequired
	}
	ainfo, err := reqctx.GetAuthenticationInfo(ctx)
	if err != nil {
		return nil, err
	}
	th, err := u.gen(threadAttribute{
		Title:     req.Title,
		TeamID:    TeamID(ainfo.TeamID()),
		CreatorID: UserID(ainfo.UserID()),
	})
	if err != nil {
		return nil, err
	}
	if err := u.repo.create(ctx, repositoryCreateRequest{
		thread: th,
	}); err != nil {
		return nil, err
	}
	return th, nil
}

func (u *Usecase) Open(ctx context.Context, req OpenRequest) (Thread, error) {
	ainfo, err := reqctx.GetAuthenticationInfo(ctx)
	if err != nil {
		return nil, err
	}
	th, err := u.repo.find(ctx, repositoryFindRequest{
		teamID:   TeamID(ainfo.TeamID()),
		threadID: req.ThreadID,
	})
	if err != nil {
		return nil, err
	}
	th.Open()
	if err := u.repo.update(ctx, repositoryUpdateRequest{
		thread: th,
	}); err != nil {
		return nil, err
	}
	return th, nil
}

func (u *Usecase) Close(ctx context.Context, req CloseRequest) (Thread, error) {
	ainfo, err := reqctx.GetAuthenticationInfo(ctx)
	if err != nil {
		return nil, err
	}
	th, err := u.repo.find(ctx, repositoryFindRequest{
		teamID:   TeamID(ainfo.TeamID()),
		threadID: req.ThreadID,
	})
	if err != nil {
		return nil, err
	}
	th.Close()
	if err := u.repo.update(ctx, repositoryUpdateRequest{
		thread: th,
	}); err != nil {
		return nil, err
	}
	return th, nil
}
