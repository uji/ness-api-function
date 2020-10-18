package registory

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/guregu/dynamo"
	"github.com/uji/ness-api-function/domain/thread"
	"github.com/uji/ness-api-function/domain/usr"
	"github.com/uji/ness-api-function/graph"
	"github.com/uji/ness-api-function/graph/generated"
	"github.com/uji/ness-api-function/infra/db"
)

func newThreadUsecase(dnmdb *dynamo.DB) *thread.Usecase {
	rp := thread.NewDynamoRepository(dnmdb, db.ThreadTableName)
	return thread.NewUsecase(thread.DefaultGenerator, rp)
}

func newUserUsecase(dnmdb *dynamo.DB) *usr.Usecase {
	rp := usr.NewDynamoRepository(dnmdb, db.UserTableName)
	return usr.NewUsecase(rp)
}

func NewRegisterdServer() http.Handler {
	db := db.NewDynamoDB()

	uc := newUserUsecase(db)
	thrd := newThreadUsecase(db)

	rslv := graph.NewResolver(thrd)
	schm := generated.NewExecutableSchema(generated.Config{Resolvers: rslv})

	usrMiddleWare := usr.NewMiddleWare(uc)
	return usrMiddleWare.Handle(handler.NewDefaultServer(schm))
}

func NewRegisterdServerWithDammyAuth() http.Handler {
	db := db.NewDynamoDB()

	_ = newUserUsecase(db)
	thrd := newThreadUsecase(db)

	rslv := graph.NewResolver(thrd)
	schm := generated.NewExecutableSchema(generated.Config{Resolvers: rslv})

	return usr.DammyMiddleware(handler.NewDefaultServer(schm))
}
