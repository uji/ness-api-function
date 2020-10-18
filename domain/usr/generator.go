package usr

import "time"

type (
	userAttribute struct {
		userID string
		name   string
	}

	Generator func(attr userAttribute) (*User, error)
)

var _ Generator = DefaultGenerator

func DefaultGenerator(attr userAttribute) (*User, error) {
	return &User{
		userID:          UserID(attr.userID),
		onCheckInTeamID: "",
		name:            attr.name,
		createdAt:       time.Now(),
		updatedAt:       time.Now(),
		members:         []*Member{},
	}, nil
}
