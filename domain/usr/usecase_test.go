package usr

import (
	"context"
	"errors"
	"testing"

	"firebase.google.com/go/auth"
	"github.com/golang/mock/gomock"
	"github.com/guregu/null"
)

func TestUsecase_Create(t *testing.T) {
	terr := errors.New("test")

	cases := []struct {
		name            string
		userID          null.String
		req             CreateRequest
		err             error
		repoCreateCount int
		repoCreateErr   error
	}{
		{
			name:   "normal",
			userID: null.StringFrom("userid"),
			req: CreateRequest{
				Name: "name",
			},
			err:             nil,
			repoCreateCount: 1,
			repoCreateErr:   nil,
		},
		{
			name:   "no userID",
			userID: null.String{},
			req: CreateRequest{
				Name: "name",
			},
			err:             ErrUnexpectedUserID,
			repoCreateCount: 0,
			repoCreateErr:   nil,
		},
		{
			name:   "repoCreate has error",
			userID: null.StringFrom("userid"),
			req: CreateRequest{
				Name: "name",
			},
			err:             terr,
			repoCreateCount: 1,
			repoCreateErr:   terr,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			if c.userID.Valid {
				ctx = SetUserIDToContext(ctx, c.userID.String)
			}

			newUser := User{}

			fbs := NewMockFireBaseAuthClient(ctrl)
			gen := func(attr UserAttribute) (*User, error) {
				return &newUser, nil
			}
			repo := NewMockRepository(ctrl)
			repo.EXPECT().create(
				ctx, &newUser,
			).Return(c.repoCreateErr).Times(c.repoCreateCount)

			sut := NewUsecase(fbs, gen, repo)
			user, err := sut.Create(ctx, c.req)
			if err != c.err {
				t.Fatal(err)
			}
			if err == nil && user != &newUser {
				t.Fatal()
			}
		})
	}
}

func TestUsecase_Find(t *testing.T) {
	terr := errors.New("test")
	newUser := User{}
	user := User{}

	cases := []struct {
		name            string
		userID          string
		user            *User
		err             error
		fbsGetUserCount int
		fbsGetUser      *auth.UserRecord
		fbsGetUserErr   error
		repoFindCount   int
		repoFindUser    *User
		repoFindErr     error
		repoCreateCount int
		repoCreateErr   error
	}{
		{
			name:            "exists user",
			userID:          "userid",
			user:            &user,
			err:             nil,
			repoFindCount:   1,
			repoFindUser:    &user,
			repoFindErr:     nil,
			fbsGetUserCount: 0,
			fbsGetUser:      &auth.UserRecord{},
			fbsGetUserErr:   nil,
			repoCreateCount: 0,
			repoCreateErr:   nil,
		},
		{
			name:            "not exists user",
			userID:          "userid",
			user:            &newUser,
			err:             nil,
			repoFindCount:   1,
			repoFindUser:    nil,
			repoFindErr:     nil,
			fbsGetUserCount: 1,
			fbsGetUser: &auth.UserRecord{
				UserInfo: &auth.UserInfo{
					DisplayName: "display name",
				},
			},
			fbsGetUserErr:   nil,
			repoCreateCount: 1,
			repoCreateErr:   nil,
		},
		{
			name:            "repo find has error",
			userID:          "userid",
			user:            nil,
			err:             terr,
			repoFindCount:   1,
			repoFindUser:    nil,
			repoFindErr:     terr,
			fbsGetUserCount: 0,
			fbsGetUser:      &auth.UserRecord{},
			fbsGetUserErr:   nil,
			repoCreateCount: 0,
			repoCreateErr:   nil,
		},
		{
			name:            "fbs get user has error",
			userID:          "userid",
			user:            nil,
			err:             terr,
			repoFindCount:   1,
			repoFindUser:    nil,
			repoFindErr:     nil,
			fbsGetUserCount: 1,
			fbsGetUser: &auth.UserRecord{
				UserInfo: &auth.UserInfo{
					DisplayName: "display name",
				},
			},
			fbsGetUserErr:   terr,
			repoCreateCount: 0,
			repoCreateErr:   nil,
		},
		{
			name:            "repo create has error",
			userID:          "userid",
			user:            nil,
			err:             terr,
			repoFindCount:   1,
			repoFindUser:    nil,
			repoFindErr:     nil,
			fbsGetUserCount: 1,
			fbsGetUser: &auth.UserRecord{
				UserInfo: &auth.UserInfo{
					DisplayName: "display name",
				},
			},
			fbsGetUserErr:   nil,
			repoCreateCount: 1,
			repoCreateErr:   terr,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			fbs := NewMockFireBaseAuthClient(ctrl)
			fbs.EXPECT().GetUser(ctx, c.userID).Return(
				c.fbsGetUser, c.fbsGetUserErr,
			).Times(c.fbsGetUserCount)

			gen := func(attr UserAttribute) (*User, error) {
				return &newUser, nil
			}

			repo := NewMockRepository(ctrl)
			repo.EXPECT().find(ctx, UserID(c.userID)).Return(
				c.repoFindUser, c.repoFindErr,
			).Times(c.repoFindCount)
			repo.EXPECT().create(
				ctx, &newUser,
			).Return(c.repoCreateErr).Times(c.repoCreateCount)

			sut := NewUsecase(fbs, gen, repo)
			user, err := sut.Find(ctx, c.userID)
			if err != c.err {
				t.Fatal(err)
			}
			if user != c.user {
				t.Fatal()
			}
		})
	}
}
