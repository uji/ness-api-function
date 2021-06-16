package usr

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/uji/ness-api-function/infra/gcp/dtstr"
)

type testStoreSut struct {
	str *store
	clt *datastore.Client
}

func createTestStoreSut(t *testing.T) testStoreSut {
	clt, err := dtstr.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	userEntityKey := "User" + t.Name()
	userInfoKey := "UserInfo" + t.Name()
	store := NewStore(userEntityKey, userInfoKey)
	return testStoreSut{store, clt}
}

func TestStore_create(t *testing.T) {
	sut := createTestStoreSut(t)
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
	tx, err := sut.clt.NewTransaction(reqctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	if err := sut.str.create(reqctx, tx, user); err != nil {
		t.Fatal(err)
	}
}

func TestStore_find(t *testing.T) {

}

func TestStore_createMembers(t *testing.T) {

}

func TestStore_findMembers(t *testing.T) {

}
