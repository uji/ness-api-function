package usr

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/uji/ness-api-function/infra/db"
)

func testRepository_create(
	t *testing.T,
	user *User,
	terr error,
	itms []*userItem,
) {
	dnmdb := db.NewDynamoDB()
	tbl := db.CreateUserTestTable(dnmdb, t)
	defer db.DestroyTestTable(&tbl, t)

	sut := NewDynamoRepository(dnmdb, tbl.Name())
	if err := sut.create(context.Background(), user); err != terr {
		t.Fatal(err)
	}

	createdItems := make([]*userItem, 0, len(itms))
	if err := tbl.Scan().All(&createdItems); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(itms, createdItems); diff != "" {
		t.Fatal(diff)
	}
}
func TestRepository_Create(t *testing.T) {
	cases := []struct {
		name         string
		user         *User
		err          error
		createdItems []*userItem
	}{
		{
			name: "normal",
			user: &User{
				userID:          "User#0",
				onCheckInTeamID: "Team#0",
				name:            "user name",
				createdAt:       time.Date(2020, 10, 18, 0, 0, 0, 0, time.UTC),
				updatedAt:       time.Date(2020, 10, 19, 0, 0, 0, 0, time.UTC),
				members: []*Member{
					{
						teamID:    "Team#0",
						name:      "member name1",
						createdAt: time.Date(2020, 10, 20, 0, 0, 0, 0, time.UTC),
						updatedAt: time.Date(2020, 10, 21, 0, 0, 0, 0, time.UTC),
					},
					{
						teamID:    "Team#1",
						name:      "member name2",
						createdAt: time.Date(2020, 10, 22, 0, 0, 0, 0, time.UTC),
						updatedAt: time.Date(2020, 10, 23, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			createdItems: []*userItem{
				{
					PK:              "User#0",
					SK:              "0",
					Name:            "user name",
					OnCheckInTeamID: "Team#0",
					CreatedAt:       time.Date(2020, 10, 18, 0, 0, 0, 0, time.UTC),
					UpdatedAt:       time.Date(2020, 10, 19, 0, 0, 0, 0, time.UTC),
				},
				{
					PK:        "User#0",
					SK:        "Team#0",
					Name:      "member name1",
					CreatedAt: time.Date(2020, 10, 20, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 10, 21, 0, 0, 0, 0, time.UTC),
				},
				{
					PK:        "User#0",
					SK:        "Team#1",
					Name:      "member name2",
					CreatedAt: time.Date(2020, 10, 22, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 10, 23, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "user has no member",
			user: &User{
				userID:          "User#0",
				onCheckInTeamID: "",
				name:            "user name",
				createdAt:       time.Date(2020, 10, 18, 0, 0, 0, 0, time.UTC),
				updatedAt:       time.Date(2020, 10, 19, 0, 0, 0, 0, time.UTC),
				members:         nil,
			},
			createdItems: []*userItem{
				{
					PK:              "User#0",
					SK:              "0",
					Name:            "user name",
					OnCheckInTeamID: "",
					CreatedAt:       time.Date(2020, 10, 18, 0, 0, 0, 0, time.UTC),
					UpdatedAt:       time.Date(2020, 10, 19, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			testRepository_create(
				t,
				c.user,
				c.err,
				c.createdItems,
			)
		})
	}
}

func testRepository_find(
	t *testing.T,
	itms []*userItem,
	userID UserID,
	expt *User,
	terr error,
) {
	dnmdb := db.NewDynamoDB()
	tbl := db.CreateUserTestTable(dnmdb, t)
	defer db.DestroyTestTable(&tbl, t)

	sut := NewDynamoRepository(dnmdb, tbl.Name())
	for _, itm := range itms {
		if err := tbl.Put(itm).Run(); err != nil {
			t.Fatal(err)
		}
	}
	res, err := sut.find(context.Background(), userID)
	if err != terr {
		t.Fatal(err)
	}
	opts := cmp.Options{
		cmp.AllowUnexported(User{}),
		cmp.AllowUnexported(Member{}),
	}
	if diff := cmp.Diff(expt, res, opts); diff != "" {
		t.Fatal(diff)
	}
}
func TestRepository_Find(t *testing.T) {
	cases := []struct {
		name   string
		items  []*userItem
		userID UserID
		expt   *User
		err    error
	}{
		{
			name: "normal",
			items: []*userItem{
				{
					PK:              "User#0",
					SK:              "0",
					Name:            "user name",
					OnCheckInTeamID: "Team#0",
					CreatedAt:       time.Date(2020, 10, 18, 0, 0, 0, 0, time.UTC),
					UpdatedAt:       time.Date(2020, 10, 19, 0, 0, 0, 0, time.UTC),
				},
				{
					PK:        "User#0",
					SK:        "Team#0",
					Name:      "member name1",
					CreatedAt: time.Date(2020, 10, 20, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 10, 21, 0, 0, 0, 0, time.UTC),
				},
				{
					PK:        "User#0",
					SK:        "Team#1",
					Name:      "member name2",
					CreatedAt: time.Date(2020, 10, 22, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2020, 10, 23, 0, 0, 0, 0, time.UTC),
				},
			},
			userID: "User#0",
			expt: &User{
				userID:          "User#0",
				onCheckInTeamID: "Team#0",
				name:            "user name",
				createdAt:       time.Date(2020, 10, 18, 0, 0, 0, 0, time.UTC),
				updatedAt:       time.Date(2020, 10, 19, 0, 0, 0, 0, time.UTC),
				members: []*Member{
					{
						teamID:    "Team#0",
						name:      "member name1",
						createdAt: time.Date(2020, 10, 20, 0, 0, 0, 0, time.UTC),
						updatedAt: time.Date(2020, 10, 21, 0, 0, 0, 0, time.UTC),
					},
					{
						teamID:    "Team#1",
						name:      "member name2",
						createdAt: time.Date(2020, 10, 22, 0, 0, 0, 0, time.UTC),
						updatedAt: time.Date(2020, 10, 23, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},
		{
			name:   "not found",
			items:  []*userItem{},
			userID: "User#0",
			expt:   nil,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			testRepository_find(
				t,
				c.items,
				c.userID,
				c.expt,
				c.err,
			)
		})
	}
}
