package graphqlcontroller

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/DKhorkov/hmtm-bff/internal/interfaces"

	"github.com/99designs/gqlgen/graphql/handler"
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
	host string,
	port int,
	readHeaderTimeout time.Duration,
	useCases interfaces.UseCases,
	logger *slog.Logger,
) *Controller {
	graphqlServer := handler.NewDefaultServer(
		graphqlcore.NewExecutableSchema(
			graphqlcore.Config{
				Resolvers: &graphqlcore.Resolver{
					UseCases: useCases,
				},
			},
		),
	)

	http.Handle("/", playground.Handler("GraphQL playground", "/query")) // TODO should be deleted on prod
	http.Handle("/query", graphqlServer)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", host, port),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	return &Controller{
		httpServer: httpServer,
		port:       port,
		host:       host,
		logger:     logger,
	}
}
