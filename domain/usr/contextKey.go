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

func SetUserIDToContext(
	ctx context.Context,
	userID string,
) context.Context {
	return context.WithValue(ctx, contextKeyUserID, UserID(userID))
}

func GetUserIDToContext(
	ctx context.Context,
) (UserID, error) {
	uid, ok := ctx.Value(contextKeyUserID).(UserID)
	if !ok {
		return "", errors.New("Unexpected UserID")
	}
	return uid, nil
}

func SetTeamIDToContext(
	ctx context.Context,
	teamID string,
) context.Context {
	return context.WithValue(ctx, contextKeyTeamID, TeamID(teamID))
}

func GetTeamIDToContext(
	ctx context.Context,
) (TeamID, error) {
	tid, ok := ctx.Value(contextKeyTeamID).(TeamID)
	if !ok {
		return "", errors.New("Unexpected TeamID")
	}
	return tid, nil
}
