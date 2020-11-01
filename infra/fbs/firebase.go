package fbs

import (
	"context"
	"encoding/base64"
	"errors"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var (
	ErrNotFoundCredentialsJsonBase64 = errors.New("not found env value FCM_CREDENTIALS_JSON_BASE64")
)

func NewAuthClient(ctx context.Context) (*auth.Client, error) {
	crd := os.Getenv("FCM_CREDENTIALS_JSON_BASE64")
	if crd == "" {
		return nil, ErrNotFoundCredentialsJsonBase64
	}
	bytes, err := base64.StdEncoding.DecodeString(crd)
	if err != nil {
		return nil, err
	}
	opt := option.WithCredentialsJSON(bytes)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}
	return app.Auth(ctx)
}
