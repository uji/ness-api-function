package usr

import "context"

type Repository interface {
	find(context.Context, UserID) (User, error)
}

type Usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) *Usecase {
	return &Usecase{repo}
}

func (u *Usecase) GetUser(ctx context.Context, userID string) (User, error) {
	usr, err := u.repo.find(ctx, UserID(userID))
	if err != nil {
		return User{}, err
	}

	return usr, err
}
