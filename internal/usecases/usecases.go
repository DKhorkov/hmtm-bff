package usecases

import (
	"context"
	"fmt"
	"path"

	customerrors "github.com/DKhorkov/hmtm-bff/internal/errors"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
	"github.com/DKhorkov/libs/security"
)

func NewCommonUseCases(
	ssoService interfaces.SsoService,
	toysService interfaces.ToysService,
	fileStorageService interfaces.FileStorageService,
	ticketsService interfaces.TicketsService,
) *CommonUseCases {
	return &CommonUseCases{
		ssoService:         ssoService,
		toysService:        toysService,
		fileStorageService: fileStorageService,
		ticketsService:     ticketsService,
	}
}

type CommonUseCases struct {
	ssoService         interfaces.SsoService
	toysService        interfaces.ToysService
	fileStorageService interfaces.FileStorageService
	ticketsService     interfaces.TicketsService
}

func (useCases *CommonUseCases) RegisterUser(ctx context.Context, userData entities.RegisterUserDTO) (uint64, error) {
	return useCases.ssoService.RegisterUser(ctx, userData)
}

func (useCases *CommonUseCases) LoginUser(
	ctx context.Context,
	userData entities.LoginUserDTO,
) (*entities.TokensDTO, error) {
	return useCases.ssoService.LoginUser(ctx, userData)
}

func (useCases *CommonUseCases) GetMe(ctx context.Context, accessToken string) (*entities.User, error) {
	return useCases.ssoService.GetMe(ctx, accessToken)
}

func (useCases *CommonUseCases) RefreshTokens(ctx context.Context, refreshToken string) (*entities.TokensDTO, error) {
	return useCases.ssoService.RefreshTokens(ctx, refreshToken)
}

func (useCases *CommonUseCases) GetUserByID(ctx context.Context, id uint64) (*entities.User, error) {
	return useCases.ssoService.GetUserByID(ctx, id)
}

func (useCases *CommonUseCases) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	return useCases.ssoService.GetAllUsers(ctx)
}

func (useCases *CommonUseCases) AddToy(ctx context.Context, rawToyData entities.RawAddToyDTO) (uint64, error) {
	user, err := useCases.GetMe(ctx, rawToyData.AccessToken)
	if err != nil {
		return 0, err
	}

	toyData := entities.AddToyDTO{
		UserID:      user.ID,
		CategoryID:  rawToyData.CategoryID,
		Name:        rawToyData.Name,
		Description: rawToyData.Description,
		Quantity:    rawToyData.Quantity,
		Price:       rawToyData.Price,
		TagIDs:      rawToyData.TagIDs,
	}

	return useCases.toysService.AddToy(ctx, toyData)
}

func (useCases *CommonUseCases) GetAllToys(ctx context.Context) ([]entities.Toy, error) {
	return useCases.toysService.GetAllToys(ctx)
}

func (useCases *CommonUseCases) GetMasterToys(ctx context.Context, masterID uint64) ([]entities.Toy, error) {
	return useCases.toysService.GetMasterToys(ctx, masterID)
}

func (useCases *CommonUseCases) GetToyByID(ctx context.Context, id uint64) (*entities.Toy, error) {
	return useCases.toysService.GetToyByID(ctx, id)
}

func (useCases *CommonUseCases) GetAllMasters(ctx context.Context) ([]entities.Master, error) {
	return useCases.toysService.GetAllMasters(ctx)
}

func (useCases *CommonUseCases) GetMasterByID(ctx context.Context, id uint64) (*entities.Master, error) {
	return useCases.toysService.GetMasterByID(ctx, id)
}

func (useCases *CommonUseCases) RegisterMaster(
	ctx context.Context,
	rawMasterData entities.RawRegisterMasterDTO,
) (uint64, error) {
	user, err := useCases.GetMe(ctx, rawMasterData.AccessToken)
	if err != nil {
		return 0, err
	}

	masterData := entities.RegisterMasterDTO{
		UserID: user.ID,
		Info:   rawMasterData.Info,
	}

	return useCases.toysService.RegisterMaster(ctx, masterData)
}

func (useCases *CommonUseCases) GetAllCategories(ctx context.Context) ([]entities.Category, error) {
	return useCases.toysService.GetAllCategories(ctx)
}

func (useCases *CommonUseCases) GetCategoryByID(ctx context.Context, id uint32) (*entities.Category, error) {
	return useCases.toysService.GetCategoryByID(ctx, id)
}

func (useCases *CommonUseCases) GetAllTags(ctx context.Context) ([]entities.Tag, error) {
	return useCases.toysService.GetAllTags(ctx)
}

func (useCases *CommonUseCases) GetTagByID(ctx context.Context, id uint32) (*entities.Tag, error) {
	return useCases.toysService.GetTagByID(ctx, id)
}

func (useCases *CommonUseCases) UploadFile(ctx context.Context, filename string, file []byte) (string, error) {
	key := security.RawEncode([]byte(filename)) + path.Ext(filename)
	return useCases.fileStorageService.Upload(ctx, key, file)
}

func (useCases *CommonUseCases) CreateTicket(
	ctx context.Context,
	rawTicketData entities.RawCreateTicketDTO,
) (uint64, error) {
	user, err := useCases.GetMe(ctx, rawTicketData.AccessToken)
	if err != nil {
		return 0, err
	}

	ticketData := entities.CreateTicketDTO{
		UserID:      user.ID,
		CategoryID:  rawTicketData.CategoryID,
		Name:        rawTicketData.Name,
		Description: rawTicketData.Description,
		Price:       rawTicketData.Price,
		Quantity:    rawTicketData.Quantity,
		TagIDs:      rawTicketData.TagIDs,
	}

	return useCases.ticketsService.CreateTicket(ctx, ticketData)
}

func (useCases *CommonUseCases) GetTicketByID(ctx context.Context, id uint64) (*entities.Ticket, error) {
	rawTicket, err := useCases.ticketsService.GetTicketByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return useCases.processRawTicket(ctx, *rawTicket), nil
}

func (useCases *CommonUseCases) processRawTicket(ctx context.Context, ticket entities.RawTicket) *entities.Ticket {
	processedTags := make([]entities.Tag, len(ticket.TagIDs))
	for tagIndex := range ticket.TagIDs {
		processedTags[tagIndex] = entities.Tag{ID: ticket.TagIDs[tagIndex]}
	}

	tags, err := useCases.toysService.GetAllTags(ctx)
	if err == nil { // Soft processing if tags were received not to have distributed monolith antipattern.
		tagsMap := make(map[uint32]entities.Tag)
		for _, tag := range tags {
			tagsMap[tag.ID] = tag
		}

		for index, tag := range processedTags {
			if _, ok := tagsMap[tag.ID]; ok {
				processedTags[index].Name = tagsMap[tag.ID].Name
			}
		}
	}

	return &entities.Ticket{
		ID:          ticket.ID,
		UserID:      ticket.UserID,
		CategoryID:  ticket.CategoryID,
		Name:        ticket.Name,
		Description: ticket.Description,
		Price:       ticket.Price,
		Quantity:    ticket.Quantity,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
		Tags:        processedTags,
	}
}

func (useCases *CommonUseCases) GetAllTickets(ctx context.Context) ([]entities.Ticket, error) {
	rawTickets, err := useCases.ticketsService.GetAllTickets(ctx)
	if err != nil {
		return nil, err
	}

	tickets := make([]entities.Ticket, len(rawTickets))
	for index, rawTicket := range rawTickets {
		tickets[index] = *useCases.processRawTicket(ctx, rawTicket)
	}

	return tickets, err
}

func (useCases *CommonUseCases) GetUserTickets(ctx context.Context, userID uint64) ([]entities.Ticket, error) {
	rawTickets, err := useCases.ticketsService.GetUserTickets(ctx, userID)
	if err != nil {
		return nil, err
	}

	tickets := make([]entities.Ticket, len(rawTickets))
	for index, rawTicket := range rawTickets {
		tickets[index] = *useCases.processRawTicket(ctx, rawTicket)
	}

	return tickets, err
}

func (useCases *CommonUseCases) GetMyTickets(
	ctx context.Context,
	accessToken string,
) ([]entities.Ticket, error) {
	user, err := useCases.GetMe(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return useCases.GetUserTickets(ctx, user.ID)
}

func (useCases *CommonUseCases) RespondToTicket(
	ctx context.Context,
	rawRespondData entities.RawRespondToTicketDTO,
) (uint64, error) {
	user, err := useCases.GetMe(ctx, rawRespondData.AccessToken)
	if err != nil {
		return 0, err
	}

	respondData := entities.RespondToTicketDTO{
		UserID:   user.ID,
		TicketID: rawRespondData.TicketID,
	}

	return useCases.ticketsService.RespondToTicket(ctx, respondData)
}

func (useCases *CommonUseCases) GetRespondByID(
	ctx context.Context,
	id uint64,
	accessToken string,
) (*entities.Respond, error) {
	respond, err := useCases.ticketsService.GetRespondByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user, err := useCases.GetMe(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	ticket, err := useCases.ticketsService.GetTicketByID(ctx, respond.TicketID)
	if err != nil {
		return nil, err
	}

	master, err := useCases.toysService.GetMasterByID(ctx, respond.MasterID)
	if err != nil {
		return nil, err
	}

	// Check if Respond belongs to Ticket owner or to Master, which responded to Ticket.
	if ticket.UserID != user.ID && master.UserID != user.ID {
		return nil, &customerrors.PermissionDeniedError{
			Message: fmt.Sprintf(
				"User with ID=%d is not rather owner of Respond with ID=%d, or owner of Ticket with ID=%d",
				user.ID,
				id,
				ticket.ID,
			),
		}
	}

	return respond, nil
}

func (useCases *CommonUseCases) GetTicketResponds(
	ctx context.Context,
	ticketID uint64,
	accessToken string,
) ([]entities.Respond, error) {
	user, err := useCases.GetMe(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	ticket, err := useCases.ticketsService.GetTicketByID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	// Check if Ticket belongs to current User.
	if ticket.UserID != user.ID {
		return nil, &customerrors.PermissionDeniedError{
			Message: fmt.Sprintf(
				"Ticket with ID=%d does not belong to current User with ID=%d",
				ticketID,
				user.ID,
			),
		}
	}

	return useCases.ticketsService.GetTicketResponds(ctx, ticketID)
}

func (useCases *CommonUseCases) GetMyResponds(
	ctx context.Context,
	accessToken string,
) ([]entities.Respond, error) {
	user, err := useCases.GetMe(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return useCases.ticketsService.GetUserResponds(ctx, user.ID)
}
