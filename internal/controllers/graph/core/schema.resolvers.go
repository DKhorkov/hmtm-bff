package graphqlcore

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.55

import (
	"context"
	"github.com/DKhorkov/libs/logging"
	"strconv"

	"github.com/DKhorkov/hmtm-bff/internal/controllers/graph/schemas"
	ssoentities "github.com/DKhorkov/hmtm-sso/pkg/entities"
)

// RegisterUser is the resolver for the registerUser field.
func (r *mutationResolver) RegisterUser(ctx context.Context, input schemas.RegisterUserInput) (int, error) {
	r.Logger.Info(
		"Received new request",
		"Request",
		input,
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	userData := ssoentities.RegisterUserDTO{
		Credentials: ssoentities.LoginUserDTO{
			Email:    input.Credentials.Email,
			Password: input.Credentials.Password,
		},
	}

	return r.UseCases.RegisterUser(userData)
}

// LoginUser is the resolver for the loginUser field.
func (r *mutationResolver) LoginUser(ctx context.Context, input schemas.LoginUserInput) (*ssoentities.TokensDTO, error) {
	r.Logger.Info(
		"Received new request",
		"Request",
		input,
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	userData := ssoentities.LoginUserDTO{
		Email:    input.Email,
		Password: input.Password,
	}

	return r.UseCases.LoginUser(userData)
}

// RefreshTokens is the resolver for the refreshTokens field.
func (r *mutationResolver) RefreshTokens(ctx context.Context, input schemas.RefreshTokensInput) (*ssoentities.TokensDTO, error) {
	r.Logger.Info(
		"Received new request",
		"Request",
		input,
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	refreshTokensData := ssoentities.TokensDTO{
		AccessToken:  input.AccessToken,
		RefreshToken: input.RefreshToken,
	}

	return r.UseCases.RefreshTokens(refreshTokensData)
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*ssoentities.User, error) {
	r.Logger.Info(
		"Received new request",
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	return r.UseCases.GetAllUsers()
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*ssoentities.User, error) {
	r.Logger.Info(
		"Received new request",
		"Request",
		id,
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	userId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	return r.UseCases.GetUserByID(userId)
}

// Me is the resolver for me field.
func (r *queryResolver) Me(ctx context.Context, accessToken string) (*ssoentities.User, error) {
	r.Logger.Info(
		"Received new request",
		"Request",
		accessToken,
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	return r.UseCases.GetMe(accessToken)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
