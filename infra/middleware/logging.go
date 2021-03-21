package middleware

import (
	"log"
	"net/http"

	"github.com/uji/ness-api-function/reqctx"
)

type Logging struct{}

func NewLogging() Logging {
	return Logging{}
}

func (l *Logging) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ainfo, err := reqctx.GetAuthenticationInfo(ctx)
		if err != nil {
			log.Println("request start, AuthenticationInfo is null")
		} else {
			log.Println("request start, AuthenticationInfo: ", ainfo)
		}
		next.ServeHTTP(w, r)
		log.Println("request end")
	})
}
