package thread

import (
	"time"

	"github.com/google/uuid"
)

type (
	Generator func(attr threadAttribute) (Thread, error)

	threadAttribute struct {
		Title     string
		TeamID    TeamID
		CreatorID UserID
	}
)

var _ Generator = DefaultGenerator

func DefaultGenerator(attr threadAttribute) (Thread, error) {
	id := "Thread#" + uuid.New().String()
	return &thread{
		id:        id,
		teamID:    attr.TeamID,
		createrID: attr.CreatorID,
		title:     attr.Title,
		closed:    false,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}
