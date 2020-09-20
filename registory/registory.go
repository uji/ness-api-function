package registory

import (
	"example.com/ness-api-function/domain/thread"
	"example.com/ness-api-function/graph"
	"example.com/ness-api-function/graph/generated"
	"example.com/ness-api-function/infra/db"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/guregu/dynamo"
)

func newThreadUsecase(db *dynamo.DB) *thread.Usecase {
	gen := thread.NewGeneratorConfigured()
	rp := thread.NewDynamoRepository(db, "Thread", gen)
	return thread.NewUsecase(gen, rp)
}

func NewRegisterdServer() *handler.Server {
	db := db.NewDynamoDB()

	thrd := newThreadUsecase(db)

	rslv := graph.NewResolver(thrd)
	schm := generated.NewExecutableSchema(generated.Config{Resolvers: rslv})
	return handler.NewDefaultServer(schm)
}
