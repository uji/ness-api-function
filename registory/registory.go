package registory

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/uji/ness-api-function/domain/thread"
	"github.com/uji/ness-api-function/domain/usr"
	"github.com/uji/ness-api-function/graph"
	"github.com/uji/ness-api-function/graph/generated"
	"github.com/uji/ness-api-function/infra/elsch"
	"github.com/uji/ness-api-function/infra/fbs"
	"github.com/uji/ness-api-function/infra/gcp/dtstr"
	"github.com/uji/ness-api-function/infra/middleware"
)

func NewRegisterdServer() http.Handler {
	fbsauth, err := fbs.NewAuthClient(context.Background())
	if err != nil {
		panic(err)
	}
	es, err := elsch.NewClient(elsch.ThreadIndexName)
	if err != nil {
		panic(err)
	}
	datastore, err := dtstr.NewClient()
	if err != nil {
		panic(err)
	}
	store := usr.NewStore()
	usrRp := usr.NewRepository(datastore, store)
	user := usr.NewUsecase(fbsauth, usr.DefaultGenerator, usrRp)
	thrdRp := thread.NewRepository(es)
	thrd := thread.NewUsecase(thread.DefaultGenerator, thrdRp)

	rslv := graph.NewResolver(user, thrd)
	schm := generated.NewExecutableSchema(generated.Config{Resolvers: rslv})
	usrMdl := usr.NewMiddleWare(fbsauth, user)
	logMdl := middleware.NewLogging()
	return usrMdl.Handle(logMdl.Handle(handler.NewDefaultServer(schm)))
}

func NewRegisterdServerWithDammyAuth(teamID, userID string) http.Handler {
	fbsauth := &usr.DammyFireBaseAuthClient{}
	es, err := elsch.NewClient(elsch.ThreadIndexName)
	if err != nil {
		panic(err)
	}
	datastore, err := dtstr.NewClient()
	if err != nil {
		panic(err)
	}
	store := usr.NewStore()
	usrRp := usr.NewRepository(datastore, store)
	user := usr.NewUsecase(fbsauth, usr.DefaultGenerator, usrRp)
	thrdRp := thread.NewRepository(es)
	thrd := thread.NewUsecase(thread.DefaultGenerator, thrdRp)

	rslv := graph.NewResolver(user, thrd)
	schm := generated.NewExecutableSchema(generated.Config{Resolvers: rslv})
	logMdl := middleware.NewLogging()
	return usr.DammyMiddleware(userID, teamID, logMdl.Handle(handler.NewDefaultServer(schm)))
}
