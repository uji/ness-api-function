package thread

import "github.com/google/uuid"

type (
	generator struct {
		generate func(title string) (*Thread, error)
	}
)

func NewGeneratorConfigured() *generator {
	return &generator{
		func(title string) (*Thread, error) {
			id := "Thread#" + uuid.New().String()
			return &Thread{
				id:     id,
				title:  title,
				closed: false,
			}, nil
		},
	}
}

func NewGenerator(
	genFunc func(title string) (*Thread, error),
) *generator {
	return &generator{genFunc}
}

func (f *generator) Generate(title string) (*Thread, error) {
	return f.generate(title)
}
