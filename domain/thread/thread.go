package thread

import "time"

type Thread struct {
	id        string
	title     string
	closed    bool
	createdAt time.Time
}

func (t Thread) ID() string {
	return t.id
}
func (t Thread) Title() string {
	return t.title
}
func (t Thread) Closed() bool {
	return t.closed
}
func (t Thread) CreatedAt() time.Time {
	return t.createdAt
}
