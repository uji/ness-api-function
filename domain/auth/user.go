package auth

type (
	UserID string
	TeamID string
)

type User struct {
	userID UserID
	teamID TeamID
}

func (u *User) UserID() UserID {
	return u.userID
}

func (u *User) TeamID() TeamID {
	return u.teamID
}
