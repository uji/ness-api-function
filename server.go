package main

import (
	"context"

	"example.com/ness-api-function/graph"
	"example.com/ness-api-function/graph/generated"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
)

var muxAdpt *gorillamux.GorillaMuxAdapter

func init() {
  r := mux.NewRouter()
  schm := generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}})
	srv := handler.NewDefaultServer(schm)

	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

  muxAdpt = gorillamux.New(r)
}

func handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
  return muxAdpt.Proxy(req)
}

func main() {
	lambda.Start(handle)
}
