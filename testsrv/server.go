package main

import (
	"log"
	"net/http"
	"os"

	"example.com/ness-api-function/graph"
	"example.com/ness-api-function/graph/generated"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "3000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	schm := generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}})
	srv := handler.NewDefaultServer(schm)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
