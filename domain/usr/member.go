package usr

import "time"

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
