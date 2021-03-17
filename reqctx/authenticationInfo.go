package reqctx

type authenticationInfo struct {
	teamID string
	userID string
}

var _ AuthenticationInfo = authenticationInfo{}

func NewAuthenticationInfo(
	teamID string,
	userID string,
) authenticationInfo {
	return authenticationInfo{
		teamID: teamID,
		userID: userID,
	}
}

func (a authenticationInfo) TeamID() string {
	return a.teamID
}

func (a authenticationInfo) UserID() string {
	return a.userID
}
