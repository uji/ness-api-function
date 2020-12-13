package main

import (
	"context"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/uji/ness-api-function/registory"
)

var muxAdpt *gorillamux.GorillaMuxAdapter

func init() {
	r := mux.NewRouter()
	srv := registory.NewRegisterdServer()

	opt := cors.Options{
		AllowedOrigins: []string{"http://localhost:8080"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}
	c := cors.New(opt)

	r.Handle("/", c.Handler(playground.Handler("GraphQL playground", "/query")))
	r.Handle("/query", c.Handler(srv))

	muxAdpt = gorillamux.New(r)
}

func handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("request info: ", req.Body, req.Headers, req.MultiValueHeaders)
	return muxAdpt.Proxy(req)
}

func main() {
	lambda.Start(handle)
}
