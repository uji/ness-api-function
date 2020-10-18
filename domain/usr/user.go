package usr

import "time"

type UserID string
type TeamID string

type Member struct {
	teamID    TeamID
	name      string
	createdAt time.Time
	updatedAt time.Time
}

func (m *Member) TeamID() TeamID {
	return m.teamID
}
func (m *Member) Name() string {
	return m.name
}
func (m *Member) CreatedAt() time.Time {
	return m.createdAt
}
func (m *Member) UpdatedAt() time.Time {
	return m.updatedAt
}

type User struct {
	userID          UserID
	onCheckInTeamID TeamID
	name            string
	createdAt       time.Time
	updatedAt       time.Time
	members         []*Member
}

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
