package thread

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
