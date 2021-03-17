//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package usr

import (
	"context"

	"firebase.google.com/go/auth"
	"github.com/uji/ness-api-function/reqctx"
)

type FireBaseAuthClient interface {
	GetUser(ctx context.Context, uid string) (*auth.UserRecord, error)
}

type Repository interface {
	create(context.Context, *User) error
	find(context.Context, UserID) (*User, error)
}

type Usecase struct {
	fbsauth FireBaseAuthClient
	gen     Generator
	repo    Repository
}

func NewUsecase(
	fbsauth FireBaseAuthClient,
	gen Generator,
	repo Repository,
) *Usecase {
	return &Usecase{fbsauth, gen, repo}
}

type (
	CreateRequest struct {
		Name string
	}
)

func (u *Usecase) Create(ctx context.Context, req CreateRequest) (*User, error) {
	ainfo, err := reqctx.GetAuthenticationInfo(ctx)
	if err != nil {
		return nil, err
	}
	usr, err := u.gen(UserAttribute{
		userID: ainfo.UserID(),
		name:   req.Name,
	})
	if err != nil {
		return nil, err
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

	if usr != nil {
		return usr, nil
	}

	fbsUsr, err := u.fbsauth.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	newUsr, err := u.gen(UserAttribute{
		userID: userID,
		name:   fbsUsr.DisplayName,
	})
	if err != nil {
		return nil, err
	}
	if err := u.repo.create(ctx, newUsr); err != nil {
		return nil, err
	}
	return newUsr, nil
}
