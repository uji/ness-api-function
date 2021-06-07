package usr

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/go-cmp/cmp"
	"github.com/uji/ness-api-function/infra/gcp/dtstr"
)

func TestRepository_create(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		requser *User
	}{
		{
			name: "normal",
			requser: &User{
				userID:          "User0",
				onCheckInTeamID: "Team0",
				name:            "test",
				createdAt:       time.Now(),
				updatedAt:       time.Now(),
				members: []*Member{
					{
						teamID:    "Team0",
						name:      "test",
						createdAt: time.Now(),
						updatedAt: time.Now(),
					},
				},
			},
		},
		{
			name: "user has no member",
			requser: &User{
				userID:          "User0",
				onCheckInTeamID: "Team0",
				name:            "test",
				createdAt:       time.Now(),
				updatedAt:       time.Now(),
				members:         []*Member{},
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			reqctx := context.Background()
			dtstr, err := dtstr.NewClient()
			if err != nil {
				t.Fatal(err)
			}
			cmpopts := cmp.Options{
				cmp.AllowUnexported(User{}),
				cmp.AllowUnexported(Member{}),
			}
			cmd := StoreMock{
				createFunc: func(ctx context.Context, tx *datastore.Transaction, user *User) error {
					if !cmp.Equal(ctx, reqctx) {
						t.Fatal(ctx)
					}
					if diff := cmp.Diff(user, c.requser, cmpopts); diff != "" {
						t.Fatal(ctx)
					}
					return nil
				},
				createMembersFunc: func(ctx context.Context, tx *datastore.Transaction, user *User) error {
					if !cmp.Equal(ctx, reqctx) {
						t.Fatal(ctx)
					}
					if diff := cmp.Diff(user, c.requser, cmpopts); diff != "" {
						t.Fatal(ctx)
					}
					return nil
				},
			}
			sut := NewRepository(dtstr, &cmd)
			if err := sut.create(reqctx, c.requser); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestRepository_find(t *testing.T) {
	cases := []struct {
		name string
	}{
		{
			name: "normal",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			reqctx := context.Background()
			reqUserID := UserID("User0")
			dtstr, err := dtstr.NewClient()
			if err != nil {
				t.Fatal(err)
			}
			user := &User{
				userID:          reqUserID,
				onCheckInTeamID: "Team0",
				name:            "test",
				createdAt:       time.Now(),
				updatedAt:       time.Now(),
				members:         nil,
			}
			members := []*Member{
				{
					teamID:    "Team0",
					name:      "test",
					createdAt: time.Now(),
					updatedAt: time.Now(),
				},
			}
			expt := &User{
				userID:          reqUserID,
				onCheckInTeamID: "Team0",
				name:            "test",
				createdAt:       user.createdAt,
				updatedAt:       user.updatedAt,
				members:         members,
			}
			cmd := StoreMock{
				findFunc: func(ctx context.Context, tx *datastore.Transaction, userID UserID) (*User, error) {
					if !cmp.Equal(ctx, reqctx) {
						t.Fatal(ctx)
					}
					if userID != reqUserID {
						t.Fatal(userID)
					}
					return user, nil
				},
				findMembersFunc: func(ctx context.Context, tx *datastore.Transaction, userID UserID) ([]*Member, error) {
					if !cmp.Equal(ctx, reqctx) {
						t.Fatal(ctx)
					}
					if userID != reqUserID {
						t.Fatal(userID)
					}
					return members, nil
				},
			}
			sut := NewRepository(dtstr, &cmd)
			res, err := sut.find(reqctx, reqUserID)
			if err != nil {
				t.Fatal(err)
			}
			opts := cmp.Options{
				cmp.AllowUnexported(User{}),
				cmp.AllowUnexported(Member{}),
			}
			if !cmp.Equal(res, expt, opts) {
				t.Fatal(res)
			}
		})
	}
}
