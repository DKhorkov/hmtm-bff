package main

import (
	"fmt"

	"github.com/DKhorkov/hmtm-bff/internal/app"
	ssogrpcclient "github.com/DKhorkov/hmtm-bff/internal/clients/sso/grpc"
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
	logger.Info(fmt.Sprintf("Application settings: %+v", settings))

	grpcClient, err := ssogrpcclient.New(
		settings.Clients.SSO.Host,
		settings.Clients.SSO.Port,
		settings.Clients.SSO.RetriesCount,
		settings.Clients.SSO.RetryTimeout,
		logger,
	)

	if err != nil {
		panic(err)
	}

	ssoRepository := repositories.NewGrpcSsoRepository(grpcClient)
	ssoService := services.NewCommonSsoService(ssoRepository)
	useCases := usecases.NewCommonUseCases(ssoService)
	controller := graphqlcontroller.New(
		settings.HTTP,
		settings.CORS,
		useCases,
		logger,
	)

	application := app.New(controller)
	application.Run()
}
