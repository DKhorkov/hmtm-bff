package graph

import "hmtmbff/internal/services"

/*
Resolver

This file will not be regenerated automatically.
It serves as dependency injection for your app, add any dependencies you require here.

https://stackoverflow.com/questions/62348857/unit-testing-graphql-in-golang
*/
type Resolver struct {
	UsersService services.UsersService
}
