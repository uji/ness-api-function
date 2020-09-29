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
	}

	Thread interface {
		ID() string
		Title() string
		Closed() bool
		CreatedAt() time.Time
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
