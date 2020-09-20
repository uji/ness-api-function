package thread

import "github.com/google/uuid"

type (
	Generator struct {
		generate func(title string) (*Thread, error)
	}
)

func NewGeneratorConfigured() *Generator {
	return &Generator{
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
) *Generator {
	return &Generator{genFunc}
}

func (f *Generator) Generate(title string) (*Thread, error) {
	return f.generate(title)
}
