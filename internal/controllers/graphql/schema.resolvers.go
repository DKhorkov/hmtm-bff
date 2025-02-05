package graphqlcontroller

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.55

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	graphqlapi "github.com/DKhorkov/hmtm-bff/api/graphql"
	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/libs/contextlib"
	"github.com/DKhorkov/libs/cookies"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/middlewares"
)

// User is the resolver for the user field.
func (r *emailResolver) User(ctx context.Context, obj *entities.Email) (*entities.User, error) {
	user, err := r.useCases.GetUserByID(ctx, obj.UserID)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			r.logger,
			fmt.Sprintf("Failed to get User for Email Communication with ID=%d", obj.ID),
			err,
		)
	}

	return user, err
}

// User is the resolver for the user field.
func (r *masterResolver) User(ctx context.Context, obj *entities.Master) (*entities.User, error) {
	user, err := r.useCases.GetUserByID(ctx, obj.UserID)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			r.logger,
			fmt.Sprintf("Failed to get User for Master with ID=%d", obj.ID),
			err,
		)
	}

	return user, err
}

// RegisterUser is the resolver for the registerUser field.
func (r *mutationResolver) RegisterUser(ctx context.Context, input graphqlapi.RegisterUserInput) (string, error) {
	logging.LogRequest(ctx, r.logger, input)

	userData := entities.RegisterUserDTO{
		DisplayName: input.DisplayName,
		Email:       input.Email,
		Password:    input.Password,
	}

	userID, err := r.useCases.RegisterUser(ctx, userData)
	return strconv.FormatUint(userID, 10), err
}

// LoginUser is the resolver for the loginUser field.
func (r *mutationResolver) LoginUser(ctx context.Context, input graphqlapi.LoginUserInput) (bool, error) {
	logging.LogRequest(ctx, r.logger, input)

	userData := entities.LoginUserDTO{
		Email:    input.Email,
		Password: input.Password,
	}

	tokens, err := r.useCases.LoginUser(ctx, userData)
	if err != nil {
		return false, err
	}

	writer, err := contextlib.GetValue[http.ResponseWriter](ctx, middlewares.CookiesWriterName)
	if err != nil {
		logging.LogErrorContext(ctx, r.logger, "Failed to get cookies writer", err)
		return false, contextlib.ValueNotFoundError{Message: middlewares.CookiesWriterName}
	}

	cookies.Set(writer, accessTokenCookieName, tokens.AccessToken, r.cookiesConfig.AccessToken)
	cookies.Set(writer, refreshTokenCookieName, tokens.RefreshToken, r.cookiesConfig.RefreshToken)
	return true, nil
}

// LogoutUser is the resolver for the logoutUser field.
func (r *mutationResolver) LogoutUser(ctx context.Context) (bool, error) {
	logging.LogRequest(ctx, r.logger, nil)

	accessToken, err := contextlib.GetValue[*http.Cookie](ctx, accessTokenCookieName)
	if err != nil {
		return false, cookies.NotFoundError{Message: accessTokenCookieName}
	}

	err = r.useCases.LogoutUser(ctx, accessToken.Value)
	if err != nil {
		return false, err
	}

	writer, err := contextlib.GetValue[http.ResponseWriter](ctx, middlewares.CookiesWriterName)
	if err != nil {
		logging.LogErrorContext(ctx, r.logger, "Failed to get cookies writer", err)
		return false, contextlib.ValueNotFoundError{Message: middlewares.CookiesWriterName}
	}

	// Deleting cookies:
	cookies.Set(writer, accessTokenCookieName, "", cookies.Config{MaxAge: -1})
	cookies.Set(writer, refreshTokenCookieName, "", cookies.Config{MaxAge: -1})
	return true, nil
}

// RefreshTokens is the resolver for the refreshTokens field.
func (r *mutationResolver) RefreshTokens(ctx context.Context, input any) (bool, error) {
	logging.LogRequest(ctx, r.logger, input)

	refreshToken, err := contextlib.GetValue[*http.Cookie](ctx, refreshTokenCookieName)
	if err != nil {
		return false, cookies.NotFoundError{Message: refreshTokenCookieName}
	}

	tokens, err := r.useCases.RefreshTokens(ctx, refreshToken.Value)
	if err != nil {
		return false, err
	}

	writer, err := contextlib.GetValue[http.ResponseWriter](ctx, middlewares.CookiesWriterName)
	if err != nil {
		logging.LogErrorContext(ctx, r.logger, "Failed to get cookies writer", err)
		return false, contextlib.ValueNotFoundError{Message: middlewares.CookiesWriterName}
	}

	cookies.Set(writer, accessTokenCookieName, tokens.AccessToken, r.cookiesConfig.AccessToken)
	cookies.Set(writer, refreshTokenCookieName, tokens.RefreshToken, r.cookiesConfig.RefreshToken)
	return true, nil
}

// RegisterMaster is the resolver for the registerMaster field.
func (r *mutationResolver) RegisterMaster(ctx context.Context, input graphqlapi.RegisterMasterInput) (string, error) {
	logging.LogRequest(ctx, r.logger, input)

	accessToken, err := contextlib.GetValue[*http.Cookie](ctx, accessTokenCookieName)
	if err != nil {
		return "", cookies.NotFoundError{Message: accessTokenCookieName}
	}

	masterData := entities.RawRegisterMasterDTO{
		AccessToken: accessToken.Value,
		Info:        input.Info,
	}

	masterID, err := r.useCases.RegisterMaster(ctx, masterData)
	return strconv.FormatUint(masterID, 10), err
}

// AddToy is the resolver for the addToy field.
func (r *mutationResolver) AddToy(ctx context.Context, input graphqlapi.AddToyInput) (string, error) {
	logging.LogRequest(ctx, r.logger, input)

	accessToken, err := contextlib.GetValue[*http.Cookie](ctx, accessTokenCookieName)
	if err != nil {
		return "", cookies.NotFoundError{Message: accessTokenCookieName}
	}

	tagIDs := make([]uint32, len(input.TagIds))
	for i, id := range input.TagIds {
		tagID, err := strconv.Atoi(id)
		if err != nil {
			return "", err
		}

		tagIDs[i] = uint32(tagID)
	}

	categoryID, err := strconv.Atoi(input.CategoryID)
	if err != nil {
		return "", err
	}

	toyData := entities.RawAddToyDTO{
		AccessToken: accessToken.Value,
		CategoryID:  uint32(categoryID),
		Name:        input.Name,
		Description: input.Description,
		Price:       float32(input.Price),
		Quantity:    uint32(input.Quantity),
		TagIDs:      tagIDs,
		Attachments: input.Attachments,
	}

	toyID, err := r.useCases.AddToy(ctx, toyData)
	return strconv.FormatUint(toyID, 10), err
}

// CreateTicket is the resolver for the createTicket field.
func (r *mutationResolver) CreateTicket(ctx context.Context, input graphqlapi.CreateTicketInput) (string, error) {
	logging.LogRequest(ctx, r.logger, input)

	accessToken, err := contextlib.GetValue[*http.Cookie](ctx, accessTokenCookieName)
	if err != nil {
		return "", cookies.NotFoundError{Message: accessTokenCookieName}
	}

	tagIDs := make([]uint32, len(input.TagIds))
	for i, id := range input.TagIds {
		tagID, err := strconv.Atoi(id)
		if err != nil {
			return "", err
		}

		tagIDs[i] = uint32(tagID)
	}

	categoryID, err := strconv.Atoi(input.CategoryID)
	if err != nil {
		return "", err
	}

	ticketData := entities.RawCreateTicketDTO{
		AccessToken: accessToken.Value,
		CategoryID:  uint32(categoryID),
		Name:        input.Name,
		Description: input.Description,
		Price:       float32(input.Price),
		Quantity:    uint32(input.Quantity),
		TagIDs:      tagIDs,
		Attachments: input.Attachments,
	}

	ticketID, err := r.useCases.CreateTicket(ctx, ticketData)
	return strconv.FormatUint(ticketID, 10), err
}

// RespondToTicket is the resolver for the respondToTicket field.
func (r *mutationResolver) RespondToTicket(ctx context.Context, input graphqlapi.RespondToTicketInput) (string, error) {
	logging.LogRequest(ctx, r.logger, input)

	accessToken, err := contextlib.GetValue[*http.Cookie](ctx, accessTokenCookieName)
	if err != nil {
		return "", cookies.NotFoundError{Message: accessTokenCookieName}
	}

	ticketID, err := strconv.Atoi(input.TicketID)
	if err != nil {
		return "", err
	}

	respondData := entities.RawRespondToTicketDTO{
		AccessToken: accessToken.Value,
		TicketID:    uint64(ticketID),
	}

	respondID, err := r.useCases.RespondToTicket(ctx, respondData)
	return strconv.FormatUint(respondID, 10), err
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*entities.User, error) {
	logging.LogRequest(ctx, r.logger, nil)

	users, err := r.useCases.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]*entities.User, len(users))
	for i, user := range users {
		response[i] = &user
	}

	return response, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*entities.User, error) {
	logging.LogRequest(ctx, r.logger, id)

	userID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	return r.useCases.GetUserByID(ctx, uint64(userID))
}

// Me is the resolver for me field.
func (r *queryResolver) Me(ctx context.Context) (*entities.User, error) {
	logging.LogRequest(ctx, r.logger, nil)

	accessToken, err := contextlib.GetValue[*http.Cookie](ctx, accessTokenCookieName)
	if err != nil {
		return nil, cookies.NotFoundError{Message: accessTokenCookieName}
	}

	return r.useCases.GetMe(ctx, accessToken.Value)
}

// Master is the resolver for the master field.
func (r *queryResolver) Master(ctx context.Context, id string) (*entities.Master, error) {
	logging.LogRequest(ctx, r.logger, id)

	masterID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	return r.useCases.GetMasterByID(ctx, uint64(masterID))
}

// Masters is the resolver for the masters field.
func (r *queryResolver) Masters(ctx context.Context) ([]*entities.Master, error) {
	logging.LogRequest(ctx, r.logger, nil)

	masters, err := r.useCases.GetAllMasters(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]*entities.Master, len(masters))
	for i, master := range masters {
		response[i] = &master
	}

	return response, nil
}

// MasterToys is the resolver for the masterToys field.
func (r *queryResolver) MasterToys(ctx context.Context, masterID string) ([]*entities.Toy, error) {
	logging.LogRequest(ctx, r.logger, masterID)

	processedMasterID, err := strconv.Atoi(masterID)
	if err != nil {
		return nil, err
	}

	toys, err := r.useCases.GetMasterToys(ctx, uint64(processedMasterID))
	if err != nil {
		return nil, err
	}

	response := make([]*entities.Toy, len(toys))
	for i, toy := range toys {
		response[i] = &toy
	}

	return response, nil
}

// Toy is the resolver for the toy field.
func (r *queryResolver) Toy(ctx context.Context, id string) (*entities.Toy, error) {
	logging.LogRequest(ctx, r.logger, id)

	toyID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	return r.useCases.GetToyByID(ctx, uint64(toyID))
}

// Toys is the resolver for the toys field.
func (r *queryResolver) Toys(ctx context.Context) ([]*entities.Toy, error) {
	logging.LogRequest(ctx, r.logger, nil)

	toys, err := r.useCases.GetAllToys(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]*entities.Toy, len(toys))
	for i, toy := range toys {
		response[i] = &toy
	}

	return response, nil
}

// MyToys is the resolver for the myToys field.
func (r *queryResolver) MyToys(ctx context.Context) ([]*entities.Toy, error) {
	logging.LogRequest(ctx, r.logger, nil)

	accessToken, err := contextlib.GetValue[*http.Cookie](ctx, accessTokenCookieName)
	if err != nil {
		return nil, cookies.NotFoundError{Message: accessTokenCookieName}
	}

	toys, err := r.useCases.GetMyToys(ctx, accessToken.Value)
	if err != nil {
		return nil, err
	}

	response := make([]*entities.Toy, len(toys))
	for i, toy := range toys {
		response[i] = &toy
	}

	return response, nil
}

// Tag is the resolver for the tag field.
func (r *queryResolver) Tag(ctx context.Context, id string) (*entities.Tag, error) {
	logging.LogRequest(ctx, r.logger, id)

	tagID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	return r.useCases.GetTagByID(ctx, uint32(tagID))
}

// Tags is the resolver for the tags field.
func (r *queryResolver) Tags(ctx context.Context) ([]*entities.Tag, error) {
	logging.LogRequest(ctx, r.logger, nil)

	tags, err := r.useCases.GetAllTags(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]*entities.Tag, len(tags))
	for i, tag := range tags {
		response[i] = &tag
	}

	return response, nil
}

// Category is the resolver for the category field.
func (r *queryResolver) Category(ctx context.Context, id string) (*entities.Category, error) {
	logging.LogRequest(ctx, r.logger, id)

	categoryID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	return r.useCases.GetCategoryByID(ctx, uint32(categoryID))
}

// Categories is the resolver for the categories field.
func (r *queryResolver) Categories(ctx context.Context) ([]*entities.Category, error) {
	logging.LogRequest(ctx, r.logger, nil)

	categories, err := r.useCases.GetAllCategories(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]*entities.Category, len(categories))
	for i, category := range categories {
		response[i] = &category
	}

	return response, nil
}

// Ticket is the resolver for the ticket field.
func (r *queryResolver) Ticket(ctx context.Context, id string) (*entities.Ticket, error) {
	logging.LogRequest(ctx, r.logger, id)

	ticketID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	return r.useCases.GetTicketByID(ctx, uint64(ticketID))
}

// Tickets is the resolver for the tickets field.
func (r *queryResolver) Tickets(ctx context.Context) ([]*entities.Ticket, error) {
	logging.LogRequest(ctx, r.logger, nil)

	tickets, err := r.useCases.GetAllTickets(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]*entities.Ticket, len(tickets))
	for i, ticket := range tickets {
		response[i] = &ticket
	}

	return response, nil
}

// UserTickets is the resolver for the userTickets field.
func (r *queryResolver) UserTickets(ctx context.Context, userID string) ([]*entities.Ticket, error) {
	logging.LogRequest(ctx, r.logger, nil)

	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}

	tickets, err := r.useCases.GetUserTickets(ctx, uint64(intUserID))
	if err != nil {
		return nil, err
	}

	response := make([]*entities.Ticket, len(tickets))
	for i, ticket := range tickets {
		response[i] = &ticket
	}

	return response, nil
}

// MyTickets is the resolver for the myTickets field.
func (r *queryResolver) MyTickets(ctx context.Context) ([]*entities.Ticket, error) {
	logging.LogRequest(ctx, r.logger, nil)

	accessToken, err := contextlib.GetValue[*http.Cookie](ctx, accessTokenCookieName)
	if err != nil {
		return nil, cookies.NotFoundError{Message: accessTokenCookieName}
	}

	tickets, err := r.useCases.GetMyTickets(ctx, accessToken.Value)
	if err != nil {
		return nil, err
	}

	response := make([]*entities.Ticket, len(tickets))
	for i, ticket := range tickets {
		response[i] = &ticket
	}

	return response, nil
}

// Respond is the resolver for the respond field.
func (r *queryResolver) Respond(ctx context.Context, id string) (*entities.Respond, error) {
	logging.LogRequest(ctx, r.logger, id)

	accessToken, err := contextlib.GetValue[*http.Cookie](ctx, accessTokenCookieName)
	if err != nil {
		return nil, cookies.NotFoundError{Message: accessTokenCookieName}
	}

	respondID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	return r.useCases.GetRespondByID(ctx, uint64(respondID), accessToken.Value)
}

// TicketResponds is the resolver for the ticketResponds field.
func (r *queryResolver) TicketResponds(ctx context.Context, ticketID string) ([]*entities.Respond, error) {
	logging.LogRequest(ctx, r.logger, nil)

	accessToken, err := contextlib.GetValue[*http.Cookie](ctx, accessTokenCookieName)
	if err != nil {
		return nil, cookies.NotFoundError{Message: accessTokenCookieName}
	}

	intTicketID, err := strconv.Atoi(ticketID)
	if err != nil {
		return nil, err
	}

	responds, err := r.useCases.GetTicketResponds(ctx, uint64(intTicketID), accessToken.Value)
	if err != nil {
		return nil, err
	}

	response := make([]*entities.Respond, len(responds))
	for i, respond := range responds {
		response[i] = &respond
	}

	return response, nil
}

// MyResponds is the resolver for the myResponds field.
func (r *queryResolver) MyResponds(ctx context.Context) ([]*entities.Respond, error) {
	logging.LogRequest(ctx, r.logger, nil)

	accessToken, err := contextlib.GetValue[*http.Cookie](ctx, accessTokenCookieName)
	if err != nil {
		return nil, cookies.NotFoundError{Message: accessTokenCookieName}
	}

	responds, err := r.useCases.GetMyResponds(ctx, accessToken.Value)
	if err != nil {
		return nil, err
	}

	response := make([]*entities.Respond, len(responds))
	for i, respond := range responds {
		response[i] = &respond
	}

	return response, nil
}

// MyEmailCommunications is the resolver for the myEmailCommunications field.
func (r *queryResolver) MyEmailCommunications(ctx context.Context) ([]*entities.Email, error) {
	logging.LogRequest(ctx, r.logger, nil)

	accessToken, err := contextlib.GetValue[*http.Cookie](ctx, accessTokenCookieName)
	if err != nil {
		return nil, cookies.NotFoundError{Message: accessTokenCookieName}
	}

	emailCommunications, err := r.useCases.GetMyEmailCommunications(ctx, accessToken.Value)
	if err != nil {
		return nil, err
	}

	response := make([]*entities.Email, len(emailCommunications))
	for i, communication := range emailCommunications {
		response[i] = &communication
	}

	return response, nil
}

// Ticket is the resolver for the ticket field.
func (r *respondResolver) Ticket(ctx context.Context, obj *entities.Respond) (*entities.Ticket, error) {
	ticket, err := r.useCases.GetTicketByID(ctx, obj.TicketID)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			r.logger,
			fmt.Sprintf("Failed to get Ticket for Respond with ID=%d", obj.ID),
			err,
		)
	}

	return ticket, err
}

// Master is the resolver for the master field.
func (r *respondResolver) Master(ctx context.Context, obj *entities.Respond) (*entities.Master, error) {
	master, err := r.useCases.GetMasterByID(ctx, obj.MasterID)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			r.logger,
			fmt.Sprintf("Failed to get Master for Respond with ID=%d", obj.ID),
			err,
		)
	}

	return master, err
}

// User is the resolver for the user field.
func (r *ticketResolver) User(ctx context.Context, obj *entities.Ticket) (*entities.User, error) {
	user, err := r.useCases.GetUserByID(ctx, obj.UserID)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			r.logger,
			fmt.Sprintf("Failed to get User for Ticket with ID=%d", obj.ID),
			err,
		)
	}

	return user, err
}

// Category is the resolver for the category field.
func (r *ticketResolver) Category(ctx context.Context, obj *entities.Ticket) (*entities.Category, error) {
	category, err := r.useCases.GetCategoryByID(ctx, obj.CategoryID)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			r.logger,
			fmt.Sprintf("Failed to get Category for Ticket with ID=%d", obj.ID),
			err,
		)
	}

	return category, err
}

// Price is the resolver for the price field.
func (r *ticketResolver) Price(ctx context.Context, obj *entities.Ticket) (float64, error) {
	return float64(obj.Price), nil
}

// Quantity is the resolver for the quantity field.
func (r *ticketResolver) Quantity(ctx context.Context, obj *entities.Ticket) (int, error) {
	return int(obj.Quantity), nil
}

// Master is the resolver for the master field.
func (r *toyResolver) Master(ctx context.Context, obj *entities.Toy) (*entities.Master, error) {
	master, err := r.useCases.GetMasterByID(ctx, obj.MasterID)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			r.logger,
			fmt.Sprintf("Failed to get Master for Toy with ID=%d", obj.MasterID),
			err,
		)
	}

	return master, err
}

// Category is the resolver for the category field.
func (r *toyResolver) Category(ctx context.Context, obj *entities.Toy) (*entities.Category, error) {
	category, err := r.useCases.GetCategoryByID(ctx, obj.CategoryID)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			r.logger,
			fmt.Sprintf("Failed to get Category for Toy with ID=%d", obj.ID),
			err,
		)
	}

	return category, err
}

// Price is the resolver for the price field.
func (r *toyResolver) Price(ctx context.Context, obj *entities.Toy) (float64, error) {
	return float64(obj.Price), nil
}

// Quantity is the resolver for the quantity field.
func (r *toyResolver) Quantity(ctx context.Context, obj *entities.Toy) (int, error) {
	return int(obj.Quantity), nil
}

// Email returns graphqlapi.EmailResolver implementation.
func (r *Resolver) Email() graphqlapi.EmailResolver { return &emailResolver{r} }

// Master returns graphqlapi.MasterResolver implementation.
func (r *Resolver) Master() graphqlapi.MasterResolver { return &masterResolver{r} }

// Mutation returns graphqlapi.MutationResolver implementation.
func (r *Resolver) Mutation() graphqlapi.MutationResolver { return &mutationResolver{r} }

// Query returns graphqlapi.QueryResolver implementation.
func (r *Resolver) Query() graphqlapi.QueryResolver { return &queryResolver{r} }

// Respond returns graphqlapi.RespondResolver implementation.
func (r *Resolver) Respond() graphqlapi.RespondResolver { return &respondResolver{r} }

// Ticket returns graphqlapi.TicketResolver implementation.
func (r *Resolver) Ticket() graphqlapi.TicketResolver { return &ticketResolver{r} }

// Toy returns graphqlapi.ToyResolver implementation.
func (r *Resolver) Toy() graphqlapi.ToyResolver { return &toyResolver{r} }

type emailResolver struct{ *Resolver }
type masterResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type respondResolver struct{ *Resolver }
type ticketResolver struct{ *Resolver }
type toyResolver struct{ *Resolver }
