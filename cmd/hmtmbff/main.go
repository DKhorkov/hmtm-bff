package main

import (
	"fmt"
	"hmtmbff/configs"
	"hmtmbff/graph"
	"hmtmbff/internal/mocks"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	var config = configs.GetConfig()

	graphqlServer := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					UsersService: &mocks.MockUsersService{},
				},
			},
		),
	)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", graphqlServer)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", config.Graphql.Port),
		ReadHeaderTimeout: time.Duration(config.HTTP.ReadHeaderTimeout) * time.Second,
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", strconv.Itoa(config.Graphql.Port))
	log.Fatal(httpServer.ListenAndServe())
}
