package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
	"github.com/uji/ness-api-function/registory"
)

const defaultPort = "3000"

func usage() {
	fmt.Fprintf(os.Stderr, "%s is runner of graphql server for debug\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\t%s [flag] # run graphql server\n", os.Args[0])
	flag.PrintDefaults()
}

var (
	TeamID = flag.String("teamID", "", "dammy value of authentication TeamID")
	UserID = flag.String("userID", "", "dammy value of authentication UserID")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

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

	http.Handle("/", c.Handler(playground.Handler("GraphQL playground", "/query")))

	var srv http.Handler
	if UserID != nil && *UserID != "" &&
		TeamID != nil && *TeamID != "" {
		srv = registory.NewRegisterdServerWithDammyAuth(*TeamID, *UserID)
		log.Printf("use dammy authentication middleware TeamID=%s UserID=%s", *TeamID, *UserID)
	} else {
		srv = registory.NewRegisterdServer()
	}
	http.Handle("/query", c.Handler(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
