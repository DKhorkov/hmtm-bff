package graphqlcontroller

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.55

import (
	"context"
	"strconv"

	graphqlapi "github.com/DKhorkov/hmtm-bff/api/graphql"
	customerrors "github.com/DKhorkov/hmtm-bff/internal/errors"
	"github.com/DKhorkov/hmtm-bff/internal/middlewares"
	"github.com/DKhorkov/hmtm-bff/internal/models"
	"github.com/DKhorkov/libs/logging"
)

// User is the resolver for the user field.
func (r *masterResolver) User(ctx context.Context, obj *models.Master) (*models.User, error) {
	return r.useCases.GetUserByID(obj.UserID)
}

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

	userData := models.RegisterUserDTO{
		Email:    input.Email,
		Password: input.Password,
	}

	userID, err := r.useCases.RegisterUser(userData)
	return int(userID), err
}

// LoginUser is the resolver for the loginUser field.
func (r *mutationResolver) LoginUser(ctx context.Context, input graphqlapi.LoginUserInput) (bool, error) {
	r.logger.Info(
		"Received new request",
		"Request",
		input,
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	userData := models.LoginUserDTO{
		Email:    input.Email,
		Password: input.Password,
	}

	tokens, err := r.useCases.LoginUser(userData)
	if err != nil {
		return false, err
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

		return false, customerrors.ContextValueNotFoundError{Message: middlewares.CookiesWriterName}
	}

	setCookie(writer, accessTokenCookieName, tokens.AccessToken, r.cookiesConfig.AccessToken)
	setCookie(writer, refreshTokenCookieName, tokens.RefreshToken, r.cookiesConfig.RefreshToken)
	return true, nil
}

// RefreshTokens is the resolver for the refreshTokens field.
func (r *mutationResolver) RefreshTokens(ctx context.Context, input any) (bool, error) {
	r.logger.Info(
		"Received new request",
		"Request",
		input,
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	refreshToken, ok := getCookieFromContext(ctx, refreshTokenCookieName)
	if !ok {
		return false, customerrors.CookieNotFoundError{Message: refreshTokenCookieName}
	}

	tokens, err := r.useCases.RefreshTokens(refreshToken.Value)
	if err != nil {
		return false, err
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

		return false, customerrors.ContextValueNotFoundError{Message: middlewares.CookiesWriterName}
	}

	setCookie(writer, accessTokenCookieName, tokens.AccessToken, r.cookiesConfig.AccessToken)
	setCookie(writer, refreshTokenCookieName, tokens.RefreshToken, r.cookiesConfig.RefreshToken)
	return true, nil
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

	accessToken, ok := getCookieFromContext(ctx, accessTokenCookieName)
	if !ok {
		return 0, customerrors.CookieNotFoundError{Message: accessTokenCookieName}
	}

	masterData := models.RegisterMasterDTO{
		AccessToken: accessToken.Value,
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

	accessToken, ok := getCookieFromContext(ctx, accessTokenCookieName)
	if !ok {
		return 0, customerrors.CookieNotFoundError{Message: accessTokenCookieName}
	}

	tagsIDs := make([]uint32, len(input.TagsIDs))
	for i, id := range input.TagsIDs {
		tagsIDs[i] = uint32(*id)
	}

	toyData := models.AddToyDTO{
		AccessToken: accessToken.Value,
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
func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	r.logger.Info(
		"Received new request",
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	users, err := r.useCases.GetAllUsers()
	if err != nil {
		return nil, err
	}

	response := make([]*models.User, len(users))
	for index, user := range users {
		response[index] = &user
	}

	return response, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
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
func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
	r.logger.Info(
		"Received new request",
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	accessToken, ok := getCookieFromContext(ctx, accessTokenCookieName)
	if !ok {
		return nil, customerrors.CookieNotFoundError{Message: accessTokenCookieName}
	}

	return r.useCases.GetMe(accessToken.Value)
}

// Master is the resolver for the master field.
func (r *queryResolver) Master(ctx context.Context, id string) (*models.Master, error) {
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
func (r *queryResolver) Masters(ctx context.Context) ([]*models.Master, error) {
	r.logger.Info(
		"Received new request",
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	masters, err := r.useCases.GetAllMasters()
	if err != nil {
		return nil, err
	}

	response := make([]*models.Master, len(masters))
	for index, master := range masters {
		response[index] = &master
	}

	return response, nil
}

// MasterToys is the resolver for the masterToys field.
func (r *queryResolver) MasterToys(ctx context.Context, masterID string) ([]*models.Toy, error) {
	r.logger.Info(
		"Received new request",
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	processedMasterID, err := strconv.Atoi(masterID)
	if err != nil {
		return nil, err
	}

	toys, err := r.useCases.GetMasterToys(uint64(processedMasterID))
	if err != nil {
		return nil, err
	}

	response := make([]*models.Toy, len(toys))
	for index, toy := range toys {
		response[index] = &toy
	}

	return response, nil
}

// Toy is the resolver for the toy field.
func (r *queryResolver) Toy(ctx context.Context, id string) (*models.Toy, error) {
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
func (r *queryResolver) Toys(ctx context.Context) ([]*models.Toy, error) {
	r.logger.Info(
		"Received new request",
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	toys, err := r.useCases.GetAllToys()
	if err != nil {
		return nil, err
	}

	response := make([]*models.Toy, len(toys))
	for index, toy := range toys {
		response[index] = &toy
	}

	return response, nil
}

// Tag is the resolver for the tag field.
func (r *queryResolver) Tag(ctx context.Context, id string) (*models.Tag, error) {
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
func (r *queryResolver) Tags(ctx context.Context) ([]*models.Tag, error) {
	r.logger.Info(
		"Received new request",
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	tags, err := r.useCases.GetAllTags()
	if err != nil {
		return nil, err
	}

	response := make([]*models.Tag, len(tags))
	for index, tag := range tags {
		response[index] = &tag
	}

	return response, nil
}

// Category is the resolver for the category field.
func (r *queryResolver) Category(ctx context.Context, id string) (*models.Category, error) {
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
func (r *queryResolver) Categories(ctx context.Context) ([]*models.Category, error) {
	r.logger.Info(
		"Received new request",
		"Context",
		ctx,
		"Traceback",
		logging.GetLogTraceback(),
	)

	categories, err := r.useCases.GetAllCategories()
	if err != nil {
		return nil, err
	}

	response := make([]*models.Category, len(categories))
	for index, category := range categories {
		response[index] = &category
	}

	return response, nil
}

// Master is the resolver for the master field.
func (r *toyResolver) Master(ctx context.Context, obj *models.Toy) (*models.Master, error) {
	return r.useCases.GetMasterByID(obj.MasterID)
}

// Category is the resolver for the category field.
func (r *toyResolver) Category(ctx context.Context, obj *models.Toy) (*models.Category, error) {
	return r.useCases.GetCategoryByID(obj.CategoryID)
}

// Price is the resolver for the price field.
func (r *toyResolver) Price(ctx context.Context, obj *models.Toy) (float64, error) {
	return float64(obj.Price), nil
}

// Quantity is the resolver for the quantity field.
func (r *toyResolver) Quantity(ctx context.Context, obj *models.Toy) (int, error) {
	return int(obj.Quantity), nil
}

// Master returns graphqlapi.MasterResolver implementation.
func (r *Resolver) Master() graphqlapi.MasterResolver { return &masterResolver{r} }

// Mutation returns graphqlapi.MutationResolver implementation.
func (r *Resolver) Mutation() graphqlapi.MutationResolver { return &mutationResolver{r} }

// Query returns graphqlapi.QueryResolver implementation.
func (r *Resolver) Query() graphqlapi.QueryResolver { return &queryResolver{r} }

// Toy returns graphqlapi.ToyResolver implementation.
func (r *Resolver) Toy() graphqlapi.ToyResolver { return &toyResolver{r} }

type masterResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type toyResolver struct{ *Resolver }
