package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
	"github.com/uji/ness-api-function/domain/usr"
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

	srv := registory.NewRegisterdServer()

	c := cors.AllowAll()

	http.Handle("/", c.Handler(playground.Handler("GraphQL playground", "/query")))
	if *TeamID != "" && *UserID != "" {
		http.Handle("/query", c.Handler(usr.DammyMiddleware(srv)))
	} else {
		http.Handle("/query", c.Handler(srv))
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
