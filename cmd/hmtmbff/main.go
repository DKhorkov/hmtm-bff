package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/DKhorkov/hmtm-bff/internal/config"

	"github.com/DKhorkov/hmtm-bff/graph"
	"github.com/DKhorkov/hmtm-bff/internal/mocks"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	settings := config.GetConfig()

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
		Addr:              fmt.Sprintf(":%d", settings.Graphql.Port),
		ReadHeaderTimeout: time.Duration(settings.HTTP.ReadHeaderTimeout) * time.Second,
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", strconv.Itoa(settings.Graphql.Port))
	log.Fatal(httpServer.ListenAndServe())
}
