//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package thread

import (
	"time"
)

type (
	UserID string
	TeamID string

	thread struct {
		id        string
		teamID    TeamID
		createrID UserID
		title     string
		closed    bool
		createdAt time.Time
		updatedAt time.Time
	}

	Thread interface {
		ID() string
		TeamID() TeamID
		CreatorID() UserID
		Title() string
		Closed() bool
		CreatedAt() time.Time
		UpdatedAt() time.Time
		Open()
		Close()
	}
)

var _ Thread = &thread{}

func (t *thread) ID() string {
	return t.id
}
func (t *thread) TeamID() TeamID {
	return t.teamID
}
func (t *thread) CreatorID() UserID {
	return t.createrID
}
func (t *thread) Title() string {
	return t.title
}
func (t *thread) Closed() bool {
	return t.closed
}
func (t *thread) CreatedAt() time.Time {
	return t.createdAt
}
func (t *thread) UpdatedAt() time.Time {
	return t.updatedAt
}

func (t *thread) Open() {
	t.closed = false
	t.updatedAt = time.Now()
}
func (t *thread) Close() {
	t.closed = true
	t.updatedAt = time.Now()
}
