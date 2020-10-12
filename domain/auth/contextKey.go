package auth

import "context"

type (
	ContextKey string
)

const (
	ContextKeyUserID ContextKey = "user-id"
	ContextKeyTeamID ContextKey = "team-id"
)

func AddValueToContext(
	ctx context.Context,
	key ContextKey,
	value interface{},
) context.Context {
	return context.WithValue(ctx, key, value)
}
