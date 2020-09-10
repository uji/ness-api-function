package main

import (
	"log"
	"net/http"
	"os"

	"example.com/ness-api-function/registory"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
)

const defaultPort = "3000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := registory.NewRegisterdServer()

	c := cors.AllowAll()

	http.Handle("/", c.Handler(playground.Handler("GraphQL playground", "/query")))
	http.Handle("/query", c.Handler(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
