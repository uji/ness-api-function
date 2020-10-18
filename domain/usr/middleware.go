package usr

import (
	"context"
	"net/http"

	firebase "firebase.google.com/go"
)

type MiddleWare struct {
	uc *Usecase
}

func NewMiddleWare(uc *Usecase) *MiddleWare {
	return &MiddleWare{uc}
}

func (m *MiddleWare) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("auth-cookie")
		if err != nil || c == nil {
			next.ServeHTTP(w, r)
			return
		}

		uid, err := getUserIDFromCookie(r.Context(), c)
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

func getUserIDFromCookie(ctx context.Context, ck *http.Cookie) (string, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return "", err
	}
	c, err := app.Auth(ctx)
	if err != nil {
		return "", err
	}
	tkn, err := c.VerifySessionCookie(ctx, ck.Value)
	if err != nil {
		return "", err
	}
	return tkn.UID, nil
}
