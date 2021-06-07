package usr

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/go-cmp/cmp"
	"github.com/uji/ness-api-function/infra/gcp/dtstr"
)

func TestStore_create(t *testing.T) {
	clt, err := dtstr.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	sut := NewStore()
	reqctx := context.Background()
	user := &User{
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
	}
	tx, err := clt.NewTransaction(reqctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	if err := sut.create(reqctx, tx, user); err != nil {
		t.Fatal(err)
	}

	var rslt []userRow
	q := datastore.NewQuery("User").Transaction(tx).Limit(1).Filter(
		"__key__ =",
		string(user.userID),
	)
	if _, err := clt.GetAll(reqctx, q, &rslt); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(rslt, []userRow{NewUserRow(user)}); diff != "" {
		t.Fatal(diff)
	}
}

func TestStore_find(t *testing.T) {

}

func TestStore_createMembers(t *testing.T) {

}

func TestStore_findMembers(t *testing.T) {

}
