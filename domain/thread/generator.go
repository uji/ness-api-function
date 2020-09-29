package thread

import (
	"time"

	"github.com/google/uuid"
)

type (
	Generator func(attr ThreadAttribute) (Thread, error)

	ThreadAttribute struct {
		Title string
	}
)

var _ Generator = DefaultGenerator

func DefaultGenerator(attr ThreadAttribute) (Thread, error) {
	id := "Thread#" + uuid.New().String()
	return &thread{
		id:        id,
		title:     attr.Title,
		closed:    false,
		createdAt: time.Now(),
	}, nil
}
