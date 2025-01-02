package main

import (
	"fmt"

	"github.com/DKhorkov/hmtm-bff/internal/app"
	ssogrpcclient "github.com/DKhorkov/hmtm-bff/internal/clients/sso/grpc"
	ticketsgrpcclient "github.com/DKhorkov/hmtm-bff/internal/clients/tickets/grpc"
	toysgrpcclient "github.com/DKhorkov/hmtm-bff/internal/clients/toys/grpc"
	"github.com/DKhorkov/hmtm-bff/internal/config"
	graphqlcontroller "github.com/DKhorkov/hmtm-bff/internal/controllers/graphql"
	"github.com/DKhorkov/hmtm-bff/internal/repositories"
	"github.com/DKhorkov/hmtm-bff/internal/services"
	"github.com/DKhorkov/hmtm-bff/internal/usecases"
	"github.com/DKhorkov/libs/logging"
)

func main() {
	settings := config.New()
	logger := logging.GetInstance(
		settings.Logging.Level,
		settings.Logging.LogFilePath,
	)

	// App configs info for frontend purposes:
	logging.LogInfo(logger, fmt.Sprintf("Application settings: %+v", settings))

	ssoClient, err := ssogrpcclient.New(
		settings.Clients.SSO.Host,
		settings.Clients.SSO.Port,
		settings.Clients.SSO.RetriesCount,
		settings.Clients.SSO.RetryTimeout,
		logger,
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
	)

	if err != nil {
		panic(err)
	}

	ssoRepository := repositories.NewGrpcSsoRepository(ssoClient)
	ssoService := services.NewCommonSsoService(ssoRepository, logger)

	toysRepository := repositories.NewGrpcToysRepository(toysClient)
	toysService := services.NewCommonToysService(toysRepository, logger)

	fileStorageRepository := repositories.NewS3FileStorageRepository(settings.S3, logger)
	fileStorageService := services.NewCommonFileStorageService(fileStorageRepository, logger)

	ticketsRepository := repositories.NewGrpcTicketsRepository(ticketsClient)
	ticketsService := services.NewCommonTicketsService(ticketsRepository, logger)

	useCases := usecases.NewCommonUseCases(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
	)

	controller := graphqlcontroller.New(
		settings.HTTP,
		settings.CORS,
		settings.Cookies,
		useCases,
		logger,
	)

	application := app.New(controller)
	application.Run()
}
