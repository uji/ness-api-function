package usr

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/uji/ness-api-function/reqctx"
)

type MiddleWare struct {
	fbsauth *auth.Client
	uc      *Usecase
}

func NewMiddleWare(
	fbc *auth.Client,
	uc *Usecase,
) *MiddleWare {
	return &MiddleWare{fbc, uc}
}

const (
	bearerPrefix = "Bearer "
)

func (m *MiddleWare) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if h == "" {
			next.ServeHTTP(w, r)
			return
		}

		if !strings.HasPrefix(h, bearerPrefix) {
			next.ServeHTTP(w, r)
			return
		}
		tkn := strings.Replace(h, bearerPrefix, "", 1)

		uid, err := m.getUserIDFromCookie(r.Context(), tkn)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := m.uc.Find(r.Context(), uid)
		if errors.Is(err, ErrNotFoundUser) {
			ctx := reqctx.NewRequestContext(r.Context(), reqctx.NewAuthenticationInfo("Team#0", uid))
			user, err = m.uc.Create(ctx, CreateRequest{
				Name: "test",
			})
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
		} else if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ainfo := reqctx.NewAuthenticationInfo(string(user.OnCheckInTeamID()), string(user.UserID()))
		ctx := reqctx.NewRequestContext(r.Context(), ainfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func DammyMiddleware(userID, teamID string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ainfo := reqctx.NewAuthenticationInfo(teamID, userID)
		ctx := reqctx.NewRequestContext(context.Background(), ainfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *MiddleWare) getUserIDFromCookie(ctx context.Context, idToken string) (string, error) {
	tkn, err := m.fbsauth.VerifyIDTokenAndCheckRevoked(ctx, idToken)
	if err != nil {
		return "", err
	}
	return tkn.UID, nil
}
