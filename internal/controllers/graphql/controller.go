package graphqlcontroller

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	graphqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	graphqlapi "github.com/DKhorkov/hmtm-bff/api/graphql"
	"github.com/DKhorkov/hmtm-bff/internal/config"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/middlewares"
	"github.com/rs/cors"
)

type Controller struct {
	httpServer *http.Server
	host       string
	port       int
	logger     *slog.Logger
}

// Run gRPC server.
func (controller *Controller) Run() {
	logging.LogInfo(
		controller.logger,
		fmt.Sprintf("Starting GraphQL Server at http://%s:%d", controller.host, controller.port),
	)

	if err := controller.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logging.LogError(controller.logger, "HTTP server error", err)
	}

	logging.LogInfo(controller.logger, "Stopped serving new connections.")
}

// Stop http server gracefully (graceful shutdown).
func (controller *Controller) Stop() {
	// Stops accepting new requests and processes already received requests:
	err := controller.httpServer.Shutdown(context.Background())
	if err != nil {
		logging.LogError(controller.logger, "HTTP shutdown error", err)
	}

	logging.LogInfo(controller.logger, "Graceful shutdown completed.")
}

func New(
	httpConfig config.HTTPConfig,
	corsConfig config.CORSConfig,
	cookiesConfig config.CookiesConfig,
	useCases interfaces.UseCases,
	logger *slog.Logger,
) *Controller {
	graphqlServer := graphqlhandler.NewDefaultServer(
		graphqlapi.NewExecutableSchema(
			graphqlapi.Config{
				Resolvers: NewResolver(
					useCases,
					logger,
					cookiesConfig,
				),
			},
		),
	)

	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GraphQL playground", "/query")) // TODO should be deleted on prod
	mux.Handle("/query", graphqlServer)

	httpHandler := cors.New(
		cors.Options{
			AllowedOrigins:   corsConfig.AllowedOrigins,
			AllowedMethods:   corsConfig.AllowedMethods,
			AllowedHeaders:   corsConfig.AllowedHeaders,
			MaxAge:           corsConfig.MaxAge,
			AllowCredentials: corsConfig.AllowCredentials,
		},
	).Handler(mux)

	// Read cookies for auth purposes:
	httpHandler = middlewares.CookiesMiddleware(httpHandler, []string{"accessToken", "refreshToken"})

	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", httpConfig.Host, httpConfig.Port),
		ReadHeaderTimeout: httpConfig.ReadHeaderTimeout,
		Handler:           httpHandler,
	}

	return &Controller{
		httpServer: httpServer,
		host:       httpConfig.Host,
		port:       httpConfig.Port,
		logger:     logger,
	}
}
