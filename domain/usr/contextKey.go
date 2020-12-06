package usr

import (
	"context"
	"errors"
)

type (
	contextKey string
)

const (
	contextKeyUserID contextKey = "user-id"
	contextKeyTeamID contextKey = "team-id"
)

var (
	ErrUnexpectedUserID = errors.New("Unexpected UserID")
	ErrUnexpectedTeamID = errors.New("Unexpected TeamID")
)

func SetUserIDToContext(
	ctx context.Context,
	userID string,
) context.Context {
	return context.WithValue(ctx, contextKeyUserID, UserID(userID))
}

func GetUserIDToContext(
	ctx context.Context,
) (string, error) {
	uid, ok := ctx.Value(contextKeyUserID).(UserID)
	if !ok {
		return "", ErrUnexpectedUserID
	}
	return string(uid), nil
}

func SetTeamIDToContext(
	ctx context.Context,
	teamID string,
) context.Context {
	return context.WithValue(ctx, contextKeyTeamID, TeamID(teamID))
}

func GetTeamIDToContext(
	ctx context.Context,
) (string, error) {
	tid, ok := ctx.Value(contextKeyTeamID).(TeamID)
	if !ok {
		return "", ErrUnexpectedTeamID
	}
	return string(tid), nil
}
