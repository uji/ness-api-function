package fbs

import (
	"context"
	"errors"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var (
	ErrInvalidCredentialsJsonBase64 = errors.New("invalid fcm credentials json base64")
)

func NewAuthClient(ctx context.Context) (*auth.Client, error) {
	crd := os.Getenv("FCM_CREDENTIALS_JSON_BASE64")
	if crd == "" {
		return nil, ErrInvalidCredentialsJsonBase64
	}
	opt := option.WithCredentialsJSON([]byte(crd))
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}
	return app.Auth(ctx)
}
