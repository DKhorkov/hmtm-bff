package main

import (
	"github.com/DKhorkov/hmtm-bff/internal/app"
	ssogrpcclient "github.com/DKhorkov/hmtm-bff/internal/clients/sso/grpc"
	"github.com/DKhorkov/hmtm-bff/internal/config"
	graphqlcontroller "github.com/DKhorkov/hmtm-bff/internal/controllers/graph"
	"github.com/DKhorkov/hmtm-bff/internal/repositories"
	"github.com/DKhorkov/hmtm-bff/internal/services"
	"github.com/DKhorkov/hmtm-bff/internal/usecases"
)

func main() {
	settings := config.New()

	grpcClient, err := ssogrpcclient.New(
		settings.Clients.SSO.Host,
		settings.Clients.SSO.Port,
		settings.Clients.SSO.RetriesCount,
		settings.Clients.SSO.RetryTimeout,
	)

	if err != nil {
		panic(err)
	}

	ssoRepository := &repositories.GrpcSsoRepository{Client: grpcClient}
	ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
	useCases := &usecases.CommonUseCases{SsoService: ssoService}
	controller := graphqlcontroller.New(
		settings.HTTP.Host,
		settings.HTTP.Port,
		settings.HTTP.ReadHeaderTimeout,
		useCases,
	)

	application := app.New(controller)
	application.Run()
}
