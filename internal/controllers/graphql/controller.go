package graphqlcontroller

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	graphqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"

	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/middlewares"
	"github.com/DKhorkov/libs/tracing"

	graphqlapi "github.com/DKhorkov/hmtm-bff/api/graphql"
	"github.com/DKhorkov/hmtm-bff/internal/config"
	customerrors "github.com/DKhorkov/hmtm-bff/internal/errors"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

func New(
	httpConfig config.HTTPConfig,
	corsConfig config.CORSConfig,
	cookiesConfig config.CookiesConfig,
	useCases interfaces.UseCases,
	logger logging.Logger,
	traceProvider tracing.Provider,
	tracingConfig config.TracingConfig,
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

	// Configures tracing:
	httpHandler = middlewares.TracingMiddleware(httpHandler, traceProvider, tracingConfig.Spans.Root)

	// Read cookies for auth purposes:
	httpHandler = middlewares.CookiesMiddleware(httpHandler, []string{"accessToken", "refreshToken"})

	// Create request ID for request for later logging:
	httpHandler = middlewares.RequestIDMiddleware(httpHandler)

	// Configuring logging. Should be used after middlewares.RequestIDMiddleware:
	httpHandler = middlewares.GraphQLLoggingMiddleware(httpHandler, logger)

	// Protecting server from too long requests:
	httpHandler = http.TimeoutHandler(
		httpHandler,
		httpConfig.TimeoutHandlerTimeout,
		customerrors.HTTPHandlerTimeoutError{}.Error(),
	)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", httpConfig.Host, httpConfig.Port),
		ReadHeaderTimeout: httpConfig.ReadHeaderTimeout,
		ReadTimeout:       httpConfig.ReadTimeout,
		Handler:           httpHandler,
	}

	return &Controller{
		httpServer: httpServer,
		host:       httpConfig.Host,
		port:       httpConfig.Port,
		logger:     logger,
	}
}

type Controller struct {
	httpServer *http.Server
	host       string
	port       int
	logger     logging.Logger
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
