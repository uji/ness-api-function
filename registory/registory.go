package registory

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/guregu/dynamo"
	"github.com/uji/ness-api-function/domain/thread"
	"github.com/uji/ness-api-function/domain/usr"
	"github.com/uji/ness-api-function/graph"
	"github.com/uji/ness-api-function/graph/generated"
	"github.com/uji/ness-api-function/infra/db"
	"github.com/uji/ness-api-function/infra/fbs"
)

func newThreadUsecase(dnmdb *dynamo.DB) *thread.Usecase {
	rp := thread.NewDynamoRepository(dnmdb, db.ThreadTableName)
	return thread.NewUsecase(thread.DefaultGenerator, rp)
}

func newUserUsecase(dnmdb *dynamo.DB) *usr.Usecase {
	rp := usr.NewDynamoRepository(dnmdb, db.UserTableName)
	return usr.NewUsecase(usr.DefaultGenerator, rp)
}

func NewRegisterdServer() http.Handler {
	db := db.NewDynamoDB()
	fbsauth, err := fbs.NewAuthClient(context.Background())
	if err != nil {
		panic(err)
	}

	user := newUserUsecase(db)
	thrd := newThreadUsecase(db)

	rslv := graph.NewResolver(user, thrd)
	schm := generated.NewExecutableSchema(generated.Config{Resolvers: rslv})

	usrMiddleWare := usr.NewMiddleWare(fbsauth, user)
	return usrMiddleWare.Handle(handler.NewDefaultServer(schm))
}

func NewRegisterdServerWithDammyAuth() http.Handler {
	db := db.NewDynamoDB()

	user := newUserUsecase(db)
	thrd := newThreadUsecase(db)

	rslv := graph.NewResolver(user, thrd)
	schm := generated.NewExecutableSchema(generated.Config{Resolvers: rslv})

	return usr.DammyMiddleware(handler.NewDefaultServer(schm))
}
