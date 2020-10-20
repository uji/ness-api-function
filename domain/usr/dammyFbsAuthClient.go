package usr

import (
	"context"

	"firebase.google.com/go/auth"
)

type DammyFireBaseAuthClient struct{}

var _ FireBaseAuthClient = &DammyFireBaseAuthClient{}

func (d *DammyFireBaseAuthClient) GetUser(
	ctx context.Context,
	uid string,
) (*auth.UserRecord, error) {
	return &auth.UserRecord{
		UserInfo: &auth.UserInfo{
			DisplayName: "display name",
			Email:       "email",
			PhoneNumber: "phone number",
			PhotoURL:    "photo url",
			ProviderID:  "provider id",
			UID:         uid,
		},
		CustomClaims:           map[string]interface{}{},
		Disabled:               false,
		EmailVerified:          false,
		ProviderUserInfo:       []*auth.UserInfo{},
		TokensValidAfterMillis: 0,
		UserMetadata:           &auth.UserMetadata{},
		TenantID:               "",
	}, nil
}
