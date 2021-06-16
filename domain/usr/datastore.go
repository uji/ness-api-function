package usr

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
)

type store struct {
	client            *datastore.Client
	userEntityKeyName string
	userInfoKeyName   string
	memberKeyName     string
}

var _ Store = &store{}

func NewStore(client *datastore.Client, userEntityKey, userInfoKey, memberKeyName string) *store {
	return &store{client, userEntityKey, userInfoKey, memberKeyName}
}

type userRow struct {
	UserID          string
	OnCheckInTeamID string
	Name            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (u userRow) toUser() *User {
	return &User{
		userID:          UserID(u.UserID),
		onCheckInTeamID: TeamID(u.OnCheckInTeamID),
		name:            u.Name,
		createdAt:       u.CreatedAt,
		updatedAt:       u.UpdatedAt,
		members:         []*Member{},
	}
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

type memberRow struct {
	TeamID    string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m memberRow) toMember() *Member {
	return &Member{
		teamID:    TeamID(m.TeamID),
		name:      m.Name,
		createdAt: m.CreatedAt,
		updatedAt: m.UpdatedAt,
	}
}

func NewMemberRow(member *Member) memberRow {
	return memberRow{
		TeamID:    string(member.teamID),
		Name:      member.name,
		CreatedAt: member.createdAt,
		UpdatedAt: member.updatedAt,
	}
}

func (s *store) create(ctx context.Context, tx *datastore.Transaction, user *User) error {
	parentKey := datastore.NameKey(s.userEntityKeyName, string(user.UserID()), nil)
	key := datastore.NameKey(s.userInfoKeyName, string(user.UserID()), parentKey)
	row := NewUserRow(user)
	_, err := tx.Put(key, &row)
	return err
}

func (s *store) find(ctx context.Context, tx *datastore.Transaction, userID UserID) (*User, error) {
	parentKey := datastore.NameKey(s.userEntityKeyName, string(userID), nil)
	key := datastore.NameKey(s.userInfoKeyName, string(userID), parentKey)
	rslt := userRow{}
	if err := tx.Get(key, &rslt); err != nil {
		return nil, err
	}
	return rslt.toUser(), nil
}

func (s *store) createMembers(ctx context.Context, tx *datastore.Transaction, user *User) error {
	parentKey := datastore.NameKey(s.userEntityKeyName, string(user.UserID()), nil)
	muts := make([]*datastore.Mutation, len(user.members))
	for i, m := range user.members {
		key := datastore.NameKey(s.memberKeyName, string(m.teamID), parentKey)
		muts[i] = datastore.NewUpsert(key, NewMemberRow(m))
	}
	_, err := tx.Mutate(muts...)
	return err
}

func (s *store) findMembers(ctx context.Context, tx *datastore.Transaction, userID UserID) ([]*Member, error) {
	parentKey := datastore.NameKey(s.userEntityKeyName, string(userID), nil)
	q := datastore.NewQuery(s.memberKeyName).Transaction(tx).Ancestor(parentKey)
	var rslt []*memberRow
	if _, err := s.client.GetAll(ctx, q, &rslt); err != nil {
		return nil, err
	}

	res := make([]*Member, len(rslt))
	for i, r := range rslt {
		res[i] = r.toMember()
	}
	return res, nil
}
