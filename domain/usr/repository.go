package usr

import (
	"context"
	"time"

	"github.com/guregu/dynamo"
)

type repository struct {
	db  *dynamo.DB
	tbl *dynamo.Table
}

var _ Repository = &repository{}

func NewDynamoRepository(
	db *dynamo.DB,
	tableName string,
) *repository {
	tbl := db.Table(tableName)
	return &repository{db, &tbl}
}

func (d *repository) create(ctx context.Context, user *User) error {
	uitms := newUserItems(user)
	condition := "attribute_not_exists(PK) AND attribute_not_exists(SK)"

	wtx := d.db.WriteTx()
	for _, uitm := range uitms {
		put := d.tbl.Put(uitm).If(condition)
		wtx = wtx.Put(put)
	}

	return wtx.Run()
}

func (d *repository) find(ctx context.Context, userID UserID) (*User, error) {
	var uitms []*userItem
	if err := d.tbl.Get("PK", string(userID)).All(&uitms); err != nil {
		return nil, err
	}
	return toUser(uitms), nil
}

type (
	userItem struct {
		PK              string    // Hash key
		SK              string    // Range key
		Name            string    `dynamo:"Name"`
		OnCheckInTeamID string    `dynamo:"OnCheckInTeamID"`
		CreatedAt       time.Time `dynamo:"CreatedAt"`
		UpdatedAt       time.Time `dynamo:"UpdatedAt"`
	}
)

func newUserItems(user *User) []*userItem {
	uitms := make([]*userItem, len(user.Members())+1)

	uitms[0] = &userItem{
		PK:              string(user.UserID()),
		SK:              "0",
		Name:            user.Name(),
		OnCheckInTeamID: string(user.OnCheckInTeamID()),
		CreatedAt:       user.CreatedAt(),
		UpdatedAt:       user.UpdatedAt(),
	}

	for i, m := range user.Members() {
		uitms[i+1] = &userItem{
			PK:        string(user.UserID()),
			SK:        string(m.TeamID()),
			Name:      m.Name(),
			CreatedAt: m.CreatedAt(),
			UpdatedAt: m.UpdatedAt(),
		}
	}
	return uitms
}

func toUser(userItems []*userItem) *User {
	var user User
	mmbs := make([]*Member, 0, len(userItems)-1)
	for _, u := range userItems {
		if u.SK == "0" {
			user = User{
				userID:          UserID(u.PK),
				onCheckInTeamID: TeamID(u.OnCheckInTeamID),
				name:            u.Name,
				createdAt:       u.CreatedAt,
				updatedAt:       u.UpdatedAt,
			}
			continue
		}
		mmbs = append(mmbs, &Member{
			teamID:    TeamID(u.SK),
			name:      u.Name,
			createdAt: u.CreatedAt,
			updatedAt: u.UpdatedAt,
		})
	}
	user.members = mmbs
	return &user
}
