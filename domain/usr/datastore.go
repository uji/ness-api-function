package usr

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
)

type store struct{}

var _ Store = &store{}

func NewStore() *store {
	return &store{}
}

type userRow struct {
	UserID          string
	OnCheckInTeamID string
	Name            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewUserRow(user *User) userRow {
	return userRow{
		UserID:          string(user.userID),
		OnCheckInTeamID: string(user.onCheckInTeamID),
		Name:            user.name,
		CreatedAt:       user.createdAt,
		UpdatedAt:       user.updatedAt,
	}
}

func (d *store) create(ctx context.Context, tx *datastore.Transaction, user *User) error {
	key := datastore.NameKey("User", string(user.UserID()), nil)
	row := NewUserRow(user)
	_, err := tx.Put(key, &row)
	return err
}

func (d *store) find(ctx context.Context, tx *datastore.Transaction, userID UserID) (*User, error) {
	return nil, nil
}

func (d *store) createMembers(ctx context.Context, tx *datastore.Transaction, user *User) error {
	return nil
}

func (d *store) findMembers(ctx context.Context, tx *datastore.Transaction, userID UserID) ([]*Member, error) {
	return nil, nil
}
