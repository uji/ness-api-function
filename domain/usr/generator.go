package usr

import "time"

type (
	UserAttribute struct {
		userID string
		name   string
	}

	Generator func(attr UserAttribute) (*User, error)
)

var _ Generator = DefaultGenerator

func DefaultGenerator(attr UserAttribute) (*User, error) {
	return &User{
		userID:          UserID(attr.userID),
		onCheckInTeamID: "",
		name:            attr.name,
		createdAt:       time.Now(),
		updatedAt:       time.Now(),
		members:         []*Member{},
	}, nil
}
