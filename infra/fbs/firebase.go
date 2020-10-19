package fbs

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

func NewAuthClient(ctx context.Context) (*auth.Client, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, err
	}
	return app.Auth(ctx)
}
