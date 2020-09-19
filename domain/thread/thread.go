package thread

import "github.com/google/uuid"

type Thread struct {
	id     string
	title  string
	closed bool
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

var NewThread func(title string) (*Thread, error) = func(title string) (*Thread, error) {
	id := "Thread#" + uuid.New().String()
	return &Thread{
		id:     id,
		title:  title,
		closed: false,
	}, nil
}
