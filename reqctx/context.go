package reqctx

import (
	"context"

	"github.com/uji/ness-api-function/nesserr"
)

type contextKey int8

const (
	contextKeyAuthenticationInfo contextKey = iota + 1
)

type AuthenticationInfo interface {
	TeamID() string
	UserID() string
}

func NewRequestContext(
	parentCtx context.Context,
	ainfo AuthenticationInfo,
) context.Context {
	return context.WithValue(parentCtx, contextKeyAuthenticationInfo, ainfo)
}

func GetAuthenticationInfo(ctx context.Context) (AuthenticationInfo, error) {
	ainfo, ok := ctx.Value(contextKeyAuthenticationInfo).(AuthenticationInfo)
	if !ok {
		return nil, nesserr.ErrUnauthorized
	}
	return ainfo, nil
}
