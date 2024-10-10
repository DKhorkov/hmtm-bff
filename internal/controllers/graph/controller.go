package graphqlcontroller

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	graphqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/DKhorkov/hmtm-bff/internal/config"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
	"github.com/rs/cors"

	"github.com/99designs/gqlgen/graphql/playground"
	graphqlcore "github.com/DKhorkov/hmtm-bff/internal/controllers/graph/core"
	"github.com/DKhorkov/hmtm-sso/pkg/logging"
)

type Controller struct {
	httpServer *http.Server
	host       string
	port       int
	logger     *slog.Logger
}

// Run gRPC server.
func (controller *Controller) Run() {
	controller.logger.Info(
		fmt.Sprintf("Starting GraphQL Server at http://%s:%d", controller.host, controller.port),
		"Traceback",
		logging.GetLogTraceback(),
	)

	if err := controller.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		controller.logger.Error(
			"HTTP server error",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)
	}

	controller.logger.Info("Stopped serving new connections.")
}

// Stop http server gracefully (graceful shutdown).
func (controller *Controller) Stop() {
	// Stops accepting new requests and processes already received requests:
	err := controller.httpServer.Shutdown(context.Background())
	if err != nil {
		controller.logger.Error(
			"HTTP shutdown error",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)
	}

	controller.logger.Info("Graceful shutdown completed.")
}

func New(
	httpConfig config.HTTPConfig,
	corsConfig config.CORSConfig,
	useCases interfaces.UseCases,
	logger *slog.Logger,
) *Controller {
	graphqlServer := graphqlhandler.NewDefaultServer(
		graphqlcore.NewExecutableSchema(
			graphqlcore.Config{
				Resolvers: &graphqlcore.Resolver{
					UseCases: useCases,
					Logger:   logger,
				},
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
