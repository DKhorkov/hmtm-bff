package main

import (
	"github.com/DKhorkov/hmtm-bff/internal/app"
	"github.com/DKhorkov/hmtm-bff/internal/config"
	graphqlcontroller "github.com/DKhorkov/hmtm-bff/internal/controllers/graph"
	"github.com/DKhorkov/hmtm-bff/internal/entities"
	mocks "github.com/DKhorkov/hmtm-bff/internal/mocks/repositories"
	"github.com/DKhorkov/hmtm-bff/internal/services"
	"github.com/DKhorkov/hmtm-bff/internal/usecases"
)

func main() {
	settings := config.GetConfig()
	ssoRepository := &mocks.MockedSsoRepository{UsersStorage: map[int]*entities.User{}}
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
