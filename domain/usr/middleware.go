package usr

import (
	"context"
	"log"
	"net/http"
	"strings"

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

const (
	bearerPrefix = "Bearer "
)

func (m *MiddleWare) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if h == "" {
			log.Println("not found Bearer token")
			next.ServeHTTP(w, r)
			return
		}

		if !strings.HasPrefix(h, bearerPrefix) {
			log.Println("not found Bearer prefix: ", h)
			next.ServeHTTP(w, r)
			return
		}
		tkn := strings.Replace(h, bearerPrefix, "", 1)

		uid, err := m.getUserIDFromCookie(r.Context(), tkn)
		if err != nil {
			log.Println("can not get userID", err)
			next.ServeHTTP(w, r)
			return
		}
		user, err := m.uc.Find(r.Context(), uid)
		if err != nil {
			log.Println("failed find user", err)
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

func DammyMiddleware(userID, teamID string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := SetUserIDToContext(r.Context(), userID)
		ctx = SetTeamIDToContext(ctx, teamID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *MiddleWare) getUserIDFromCookie(ctx context.Context, idToken string) (string, error) {
	tkn, err := m.fbsauth.VerifyIDTokenAndCheckRevoked(ctx, idToken)
	if err != nil {
		return "", err
	}
	return tkn.UID, nil
}
