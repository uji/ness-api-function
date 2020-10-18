package usr

import "context"

type Repository interface {
	create(context.Context, *User) error
	find(context.Context, UserID) (*User, error)
}

type Usecase struct {
	gen  Generator
	repo Repository
}

func NewUsecase(
	gen Generator,
	repo Repository,
) *Usecase {
	return &Usecase{gen, repo}
}

type (
	CreateRequest struct {
		UserID string
		Name   string
	}
)

func (u *Usecase) Create(ctx context.Context, req CreateRequest) (*User, error) {
	usr, err := u.gen(userAttribute{
		userID: req.UserID,
		name:   req.Name,
	})
	if err != nil {
		return nil, nil
	}
	if err := u.repo.create(ctx, usr); err != nil {
		return nil, err
	}
	return usr, nil
}

func (u *Usecase) Find(ctx context.Context, userID string) (*User, error) {
	usr, err := u.repo.find(ctx, UserID(userID))
	if err != nil {
		return nil, err
	}
	return usr, err
}
