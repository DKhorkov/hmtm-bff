package graphqlcore

import (
	"log/slog"

	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

/*
Resolver

This file will not be regenerated automatically.
It serves as dependency injection for your app, add any dependencies you require here.

https://stackoverflow.com/questions/62348857/unit-testing-graphql-in-golang
*/
type Resolver struct {
	UseCases interfaces.UseCases
	Logger   *slog.Logger
}
