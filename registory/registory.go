package registory

import (
	"example.com/ness-api-function/domain/thread"
	"example.com/ness-api-function/graph"
	"example.com/ness-api-function/graph/generated"
	"example.com/ness-api-function/infra/db"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/jmoiron/sqlx"
)

func newThreadUsecase(db *sqlx.DB) *thread.Usecase {
	rp := thread.NewPsqlRepository(db)
	return thread.NewUsecase(rp)
}

func NewRegisterdServer() *handler.Server {
	db := db.NewDB()

	thrd := newThreadUsecase(db)

	rslv := graph.NewResolver(thrd)
	schm := generated.NewExecutableSchema(generated.Config{Resolvers: rslv})
	return handler.NewDefaultServer(schm)
}
