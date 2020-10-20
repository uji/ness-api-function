package registory

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/uji/ness-api-function/domain/thread"
	"github.com/uji/ness-api-function/domain/usr"
	"github.com/uji/ness-api-function/graph"
	"github.com/uji/ness-api-function/graph/generated"
	"github.com/uji/ness-api-function/infra/db"
	"github.com/uji/ness-api-function/infra/fbs"
)

func NewRegisterdServer() http.Handler {
	dnmdb := db.NewDynamoDB()
	fbsauth, err := fbs.NewAuthClient(context.Background())
	if err != nil {
		panic(err)
	}

	usrRp := usr.NewDynamoRepository(dnmdb, db.UserTableName)
	user := usr.NewUsecase(fbsauth, usr.DefaultGenerator, usrRp)
	thrdRp := thread.NewDynamoRepository(dnmdb, db.ThreadTableName)
	thrd := thread.NewUsecase(thread.DefaultGenerator, thrdRp)

	rslv := graph.NewResolver(user, thrd)
	schm := generated.NewExecutableSchema(generated.Config{Resolvers: rslv})
	usrMdl := usr.NewMiddleWare(fbsauth, user)
	return usrMdl.Handle(handler.NewDefaultServer(schm))
}

func NewRegisterdServerWithDammyAuth() http.Handler {
	dnmdb := db.NewDynamoDB()
	fbsauth := &usr.DammyFireBaseAuthClient{}

	usrRp := usr.NewDynamoRepository(dnmdb, db.UserTableName)
	user := usr.NewUsecase(fbsauth, usr.DefaultGenerator, usrRp)
	thrdRp := thread.NewDynamoRepository(dnmdb, db.ThreadTableName)
	thrd := thread.NewUsecase(thread.DefaultGenerator, thrdRp)

	rslv := graph.NewResolver(user, thrd)
	schm := generated.NewExecutableSchema(generated.Config{Resolvers: rslv})
	return usr.DammyMiddleware(handler.NewDefaultServer(schm))
}
