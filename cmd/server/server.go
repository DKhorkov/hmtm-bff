package main

import (
	"context"
	"fmt"

	"github.com/DKhorkov/libs/cache"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/tracing"

	"github.com/DKhorkov/hmtm-bff/internal/app"
	notificationsgrpcclient "github.com/DKhorkov/hmtm-bff/internal/clients/notifications/grpc"
	ssogrpcclient "github.com/DKhorkov/hmtm-bff/internal/clients/sso/grpc"
	ticketsgrpcclient "github.com/DKhorkov/hmtm-bff/internal/clients/tickets/grpc"
	toysgrpcclient "github.com/DKhorkov/hmtm-bff/internal/clients/toys/grpc"
	"github.com/DKhorkov/hmtm-bff/internal/config"
	graphqlcontroller "github.com/DKhorkov/hmtm-bff/internal/controllers/graphql"
	"github.com/DKhorkov/hmtm-bff/internal/repositories"
	"github.com/DKhorkov/hmtm-bff/internal/services"
	"github.com/DKhorkov/hmtm-bff/internal/usecases"
)

func main() {
	settings := config.New()
	logger := logging.New(
		settings.Logging.Level,
		settings.Logging.LogFilePath,
	)

	// App configs info for frontend purposes:
	logging.LogInfo(logger, fmt.Sprintf("Application settings: %+v", settings))

	traceProvider, err := tracing.New(settings.Tracing.Server)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = traceProvider.Shutdown(context.Background()); err != nil {
			logging.LogError(logger, "Error shutting down traceProvider", err)
		}
	}()

	ssoClient, err := ssogrpcclient.New(
		settings.Clients.SSO.Host,
		settings.Clients.SSO.Port,
		settings.Clients.SSO.RetriesCount,
		settings.Clients.SSO.RetryTimeout,
		logger,
		traceProvider,
		settings.Tracing.Spans.Clients.SSO,
	)
	if err != nil {
		panic(err)
	}

	toysClient, err := toysgrpcclient.New(
		settings.Clients.Toys.Host,
		settings.Clients.Toys.Port,
		settings.Clients.Toys.RetriesCount,
		settings.Clients.Toys.RetryTimeout,
		logger,
		traceProvider,
		settings.Tracing.Spans.Clients.Toys,
	)
	if err != nil {
		panic(err)
	}

	ticketsClient, err := ticketsgrpcclient.New(
		settings.Clients.Tickets.Host,
		settings.Clients.Tickets.Port,
		settings.Clients.Tickets.RetriesCount,
		settings.Clients.Tickets.RetryTimeout,
		logger,
		traceProvider,
		settings.Tracing.Spans.Clients.Tickets,
	)
	if err != nil {
		panic(err)
	}

	notificationsClient, err := notificationsgrpcclient.New(
		settings.Clients.Notifications.Host,
		settings.Clients.Notifications.Port,
		settings.Clients.Notifications.RetriesCount,
		settings.Clients.Notifications.RetryTimeout,
		logger,
		traceProvider,
		settings.Tracing.Spans.Clients.Notifications,
	)
	if err != nil {
		panic(err)
	}

	ssoRepository := repositories.NewSsoRepository(ssoClient)
	ssoService := services.NewSsoService(ssoRepository, logger)

	toysRepository := repositories.NewToysRepository(toysClient)
	toysService := services.NewToysService(toysRepository, logger)

	fileStorageRepository, err := repositories.NewS3FileStorageRepository(settings.S3, logger)
	if err != nil {
		panic(err)
	}

	fileStorageService := services.NewFileStorageService(fileStorageRepository, logger)

	ticketsRepository := repositories.NewTicketsRepository(ticketsClient)
	ticketsService := services.NewTicketsService(ticketsRepository, logger)

	notificationsRepository := repositories.NewNotificationsRepository(notificationsClient)
	notificationsService := services.NewNotificationsService(notificationsRepository, logger)

	cacheProvider, err := cache.New(
		cache.WithHost(settings.Cache.Host),
		cache.WithPort(settings.Cache.Port),
		cache.WithPassword(settings.Cache.Password),
	)
	if err != nil {
		panic(err)
	}

	useCases := usecases.NewCacheDecorator(
		usecases.New(
			ssoService,
			toysService,
			fileStorageService,
			ticketsService,
			notificationsService,
			settings.Validation,
			logger,
			traceProvider,
		),
		cacheProvider,
		logger,
	)

	controller := graphqlcontroller.New(
		settings.HTTP,
		settings.CORS,
		settings.Cookies,
		useCases,
		logger,
		traceProvider,
		settings.Tracing,
	)

	application := app.New(controller)
	application.Run()
}
