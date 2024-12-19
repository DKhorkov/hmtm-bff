package graphqlcontroller

import (
	"log/slog"

	"github.com/DKhorkov/hmtm-bff/internal/config"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

const (
	accessTokenCookieName  = "accessToken"
	refreshTokenCookieName = "refreshToken"
)

/*
Resolver

This file will not be regenerated automatically.
It serves as dependency injection for your app, add any dependencies you require here.

https://stackoverflow.com/questions/62348857/unit-testing-graphql-in-golang
*/
type Resolver struct {
	useCases      interfaces.UseCases
	logger        *slog.Logger
	cookiesConfig config.CookiesConfig
}

func NewResolver(
	useCases interfaces.UseCases,
	logger *slog.Logger,
	cookiesConfig config.CookiesConfig,
) *Resolver {
	return &Resolver{
		useCases:      useCases,
		logger:        logger,
		cookiesConfig: cookiesConfig,
	}
}
