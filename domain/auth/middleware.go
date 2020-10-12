package auth

import (
	"context"
	"net/http"

	firebase "firebase.google.com/go"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("auth-cookie")
		if err != nil || c == nil {
			next.ServeHTTP(w, r)
			return
		}

		uid, err := verifyCookie(r.Context(), c)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user := getUser(uid)

		ctx := AddValueToContext(r.Context(), ContextKeyUserID, user.UserID())
		ctx = AddValueToContext(ctx, ContextKeyTeamID, user.TeamID())
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func verifyCookie(ctx context.Context, ck *http.Cookie) (UserID, error) {
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
	return UserID(tkn.UID), nil
}

func getUser(id UserID) User {
	return User{
		userID: id,
		teamID: "Team#0",
	}
}
