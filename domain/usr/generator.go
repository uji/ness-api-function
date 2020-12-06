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
	// TODO: remove onCheckInTeamID and default members
	return &User{
		userID:          UserID(attr.userID),
		onCheckInTeamID: "Team-0",
		name:            attr.name,
		createdAt:       time.Now(),
		updatedAt:       time.Now(),
		members: []*Member{
			{
				teamID:    "Team-0",
				name:      "test team",
				createdAt: time.Now(),
				updatedAt: time.Now(),
			},
		},
	}, nil
}
