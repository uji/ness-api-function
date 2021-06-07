package usr

//go:generate moq -out repository_moq.go . Store

import (
	"context"

	"cloud.google.com/go/datastore"
)

type repository struct {
	datastore *datastore.Client
	cmd       Store
}

type Store interface {
	create(ctx context.Context, tx *datastore.Transaction, user *User) error
	find(ctx context.Context, tx *datastore.Transaction, userID UserID) (*User, error)
	createMembers(ctx context.Context, tx *datastore.Transaction, user *User) error
	findMembers(ctx context.Context, tx *datastore.Transaction, userID UserID) ([]*Member, error)
}

var _ Repository = &repository{}

func NewRepository(datastore *datastore.Client, cmd Store) *repository {
	return &repository{datastore, cmd}
}

func (r *repository) create(ctx context.Context, user *User) error {
	_, err := r.datastore.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		if err := r.cmd.create(ctx, tx, user); err != nil {
			return err
		}
		return r.cmd.createMembers(ctx, tx, user)
	})
	return err
}

func (r *repository) find(ctx context.Context, userID UserID) (*User, error) {
	var res *User
	_, err := r.datastore.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		user, err := r.cmd.find(ctx, tx, userID)
		if err != nil {
			return err
		}
		members, err := r.cmd.findMembers(ctx, tx, userID)
		if err != nil {
			return err
		}
		res = user
		res.members = members
		return err
	})
	return res, err
}
