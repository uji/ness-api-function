package thread

import (
	"time"
)

type (
	thread struct {
		id        string
		title     string
		closed    bool
		createdAt time.Time
		updatedAt time.Time
	}

	Thread interface {
		ID() string
		Title() string
		Closed() bool
		CreatedAt() time.Time
		UpdatedAt() time.Time
		Close()
	}
)

var _ Thread = &thread{}

func (t *thread) ID() string {
	return t.id
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

func (t *thread) Close() {
	t.closed = true
	t.updatedAt = time.Now()
}
