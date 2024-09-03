package main

import (
	"fmt"
	"hmtmbff/configs"
	"hmtmbff/graph"
	"log"
	"net/http"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	var config = configs.GetConfig()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", strconv.Itoa(config.Graphql.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Graphql.Port), nil))
}
