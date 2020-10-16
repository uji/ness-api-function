package thread

import (
	"context"
	"time"

	"github.com/guregu/null"
	"github.com/uji/ness-api-function/domain/nessauth"
)

type (
	repositoryGetRequest struct {
		offsetTime null.Time
		closed     null.Bool
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

	repositoryOpenRequest struct {
		threadID string
	}

	repositoryCloseRequest struct {
		threadID string
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
	uid, err := nessauth.GetUserIDToContext(ctx)
	tid, err := nessauth.GetTeamIDToContext(ctx)
	th, err := u.gen(threadAttribute{
		Title:     req.Title,
		TeamID:    TeamID(tid),
		CreatorID: UserID(uid),
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
	tid, err := nessauth.GetTeamIDToContext(ctx)
	if err != nil {
		return nil, err
	}
	th, err := u.repo.find(ctx, repositoryFindRequest{
		teamID:   TeamID(tid),
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
	tid, err := nessauth.GetTeamIDToContext(ctx)
	if err != nil {
		return nil, err
	}
	th, err := u.repo.find(ctx, repositoryFindRequest{
		teamID:   TeamID(tid),
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
