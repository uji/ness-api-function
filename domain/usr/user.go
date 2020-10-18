package usr

import "time"

type (
	UserID string
	TeamID string

	User struct {
		userID          UserID
		onCheckInTeamID TeamID
		name            string
		createdAt       time.Time
		updatedAt       time.Time
		members         []*Member
	}
)

func (u *User) UserID() UserID {
	return u.userID
}
func (u *User) OnCheckInTeamID() TeamID {
	return u.onCheckInTeamID
}
func (u *User) Name() string {
	return u.name
}
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}
func (u *User) Members() []*Member {
	return u.members
}

func (u *User) OnCheckIn() bool {
	return u.onCheckInTeamID != ""
}
