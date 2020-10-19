package usr

import (
	"context"
	"net/http"

	"firebase.google.com/go/auth"
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

func (m *MiddleWare) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("auth-cookie")
		if err != nil || c == nil {
			next.ServeHTTP(w, r)
			return
		}

		uid, err := m.getUserIDFromCookie(r.Context(), c)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user, err := m.uc.Find(r.Context(), uid)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := SetUserIDToContext(r.Context(), string(user.UserID()))
		if user.OnCheckIn() {
			ctx = SetTeamIDToContext(ctx, string(user.OnCheckInTeamID()))
		}
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func DammyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := SetUserIDToContext(r.Context(), "User#0")
		ctx = SetTeamIDToContext(ctx, "Team#0")
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *MiddleWare) getUserIDFromCookie(ctx context.Context, ck *http.Cookie) (string, error) {
	tkn, err := m.fbsauth.VerifySessionCookie(ctx, ck.Value)
	if err != nil {
		return "", err
	}
	return tkn.UID, nil
}
