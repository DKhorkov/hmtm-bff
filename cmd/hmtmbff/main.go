package main

import (
	"github.com/DKhorkov/hmtm-bff/internal/app"
	"github.com/DKhorkov/hmtm-bff/internal/config"
	graphqlcontroller "github.com/DKhorkov/hmtm-bff/internal/controllers/graph"
	mocks "github.com/DKhorkov/hmtm-bff/internal/mocks/repositories"
	"github.com/DKhorkov/hmtm-bff/internal/services"
	"github.com/DKhorkov/hmtm-bff/internal/usecases"
	ssoentities "github.com/DKhorkov/hmtm-sso/entities"
)

func main() {
	settings := config.New()

	ssoRepository := &mocks.MockedSsoRepository{UsersStorage: map[int]*ssoentities.User{}}
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
