package graphqlcontroller

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.55

import (
	"context"
	"github.com/DKhorkov/hmtm-bff/internal/middlewares"
	"strconv"

	graphqlapi "github.com/DKhorkov/hmtm-bff/api/graphql"
	customerrors "github.com/DKhorkov/hmtm-bff/internal/errors"
	ssoentities "github.com/DKhorkov/hmtm-sso/pkg/entities"
	toysentities "github.com/DKhorkov/hmtm-toys/pkg/entities"
	"github.com/DKhorkov/libs/logging"
)

// RegisterUser is the resolver for the registerUser field.
func (r *mutationResolver) RegisterUser(ctx context.Context, input graphqlapi.RegisterUserInput) (int, error) {
	r.logger.Info(
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
			Email:    input.Email,
			Password: input.Password,
		},
	}

	userID, err := r.useCases.RegisterUser(userData)
	return int(userID), err
}

// LoginUser is the resolver for the loginUser field.
func (r *mutationResolver) LoginUser(ctx context.Context, input graphqlapi.LoginUserInput) (*ssoentities.TokensDTO, error) {
	r.logger.Info(
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

	writer, ok := getHTTPWriterFromContext(ctx)
	if !ok {
		r.logger.ErrorContext(
			ctx,
			"Failed to get cookies writer",
			"Context",
			ctx,
			"Traceback",
			logging.GetLogTraceback(),
		)

		return nil, customerrors.ContextValueNotFoundError{Message: middlewares.CookiesWriterName}
	}

	tokens, err := r.useCases.LoginUser(userData)
	if err != nil {
		return nil, err
	}

	setCookie(writer, accessTokenCookieName, tokens.AccessToken, r.cookiesConfig.AccessToken)
	setCookie(writer, refreshTokenCookieName, tokens.RefreshToken, r.cookiesConfig.RefreshToken)
	return tokens, nil
}

// RefreshTokens is the resolver for the refreshTokens field.
func (r *mutationResolver) RefreshTokens(ctx context.Context, input graphqlapi.RefreshTokensInput) (*ssoentities.TokensDTO, error) {
	r.logger.Info(
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

	return r.useCases.RefreshTokens(refreshTokensData)
}

// RegisterMaster is the resolver for the registerMaster field.
func (r *mutationResolver) RegisterMaster(ctx context.Context, input graphqlapi.RegisterMasterInput) (int, error) {
	r.logger.Info(
		"Received new request",
		"Request",
		input,
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	masterData := toysentities.RawRegisterMasterDTO{
		AccessToken: input.AccessToken,
		Info:        input.Info,
	}

	masterID, err := r.useCases.RegisterMaster(masterData)
	return int(masterID), err
}

// AddToy is the resolver for the addToy field.
func (r *mutationResolver) AddToy(ctx context.Context, input graphqlapi.AddToyInput) (int, error) {
	r.logger.Info(
		"Received new request",
		"Request",
		input,
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	tagsIDs := make([]uint32, len(input.TagsIDs))
	for i, id := range input.TagsIDs {
		tagsIDs[i] = uint32(id)
	}

	toyData := toysentities.RawAddToyDTO{
		AccessToken: input.AccessToken,
		CategoryID:  uint32(input.CategoryID),
		Name:        input.Name,
		Description: input.Description,
		Price:       float32(input.Price),
		Quantity:    uint32(input.Quantity),
		TagsIDs:     tagsIDs,
	}

	toyID, err := r.useCases.AddToy(toyData)
	return int(toyID), err
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*ssoentities.User, error) {
	r.logger.Info(
		"Received new request",
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	return r.useCases.GetAllUsers()
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*ssoentities.User, error) {
	r.logger.Info(
		"Received new request",
		"Request",
		id,
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	userID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	return r.useCases.GetUserByID(uint64(userID))
}

// Me is the resolver for me field.
func (r *queryResolver) Me(ctx context.Context, accessToken string) (*ssoentities.User, error) {
	r.logger.Info(
		"Received new request",
		"Request",
		accessToken,
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	cookiesAccessToken, ok := getCookieFromContext(ctx, accessTokenCookieName)
	if !ok {
		return nil, customerrors.CookieNotFoundError{Message: accessTokenCookieName}
	}

	return r.useCases.GetMe(cookiesAccessToken.Value)
}

// Master is the resolver for the master field.
func (r *queryResolver) Master(ctx context.Context, id string) (*toysentities.Master, error) {
	r.logger.Info(
		"Received new request",
		"Request",
		id,
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	masterID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	return r.useCases.GetMasterByID(uint64(masterID))
}

// Masters is the resolver for the masters field.
func (r *queryResolver) Masters(ctx context.Context) ([]*toysentities.Master, error) {
	r.logger.Info(
		"Received new request",
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	return r.useCases.GetAllMasters()
}

// Toy is the resolver for the toy field.
func (r *queryResolver) Toy(ctx context.Context, id string) (*toysentities.Toy, error) {
	r.logger.Info(
		"Received new request",
		"Request",
		id,
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	toyID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	return r.useCases.GetToyByID(uint64(toyID))
}

// Toys is the resolver for the toys field.
func (r *queryResolver) Toys(ctx context.Context) ([]*toysentities.Toy, error) {
	r.logger.Info(
		"Received new request",
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	return r.useCases.GetAllToys()
}

// Tag is the resolver for the tag field.
func (r *queryResolver) Tag(ctx context.Context, id string) (*toysentities.Tag, error) {
	r.logger.Info(
		"Received new request",
		"Request",
		id,
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	tagID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	return r.useCases.GetTagByID(uint32(tagID))
}

// Tags is the resolver for the tags field.
func (r *queryResolver) Tags(ctx context.Context) ([]*toysentities.Tag, error) {
	r.logger.Info(
		"Received new request",
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	return r.useCases.GetAllTags()
}

// Category is the resolver for the category field.
func (r *queryResolver) Category(ctx context.Context, id string) (*toysentities.Category, error) {
	r.logger.Info(
		"Received new request",
		"Request",
		id,
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	categoryID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	return r.useCases.GetCategoryByID(uint32(categoryID))
}

// Categories is the resolver for the categories field.
func (r *queryResolver) Categories(ctx context.Context) ([]*toysentities.Category, error) {
	r.logger.Info(
		"Received new request",
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	return r.useCases.GetAllCategories()
}

// Price is the resolver for the price field.
func (r *toyResolver) Price(ctx context.Context, obj *toysentities.Toy) (float64, error) {
	return float64(obj.Price), nil
}

// Quantity is the resolver for the quantity field.
func (r *toyResolver) Quantity(ctx context.Context, obj *toysentities.Toy) (int, error) {
	return int(obj.Quantity), nil
}

// Mutation returns graphqlapi.MutationResolver implementation.
func (r *Resolver) Mutation() graphqlapi.MutationResolver { return &mutationResolver{r} }

// Query returns graphqlapi.QueryResolver implementation.
func (r *Resolver) Query() graphqlapi.QueryResolver { return &queryResolver{r} }

// Toy returns graphqlapi.ToyResolver implementation.
func (r *Resolver) Toy() graphqlapi.ToyResolver { return &toyResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type toyResolver struct{ *Resolver }
