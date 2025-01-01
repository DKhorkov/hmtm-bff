package repositories

import (
	"context"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
	"github.com/DKhorkov/hmtm-tickets/api/protobuf/generated/go/tickets"
	"github.com/DKhorkov/libs/contextlib"
	"github.com/DKhorkov/libs/requestid"
)

func NewGrpcTicketsRepository(client interfaces.TicketsGrpcClient) *GrpcTicketsRepository {
	return &GrpcTicketsRepository{client: client}
}

type GrpcTicketsRepository struct {
	client interfaces.TicketsGrpcClient
}

func (repo *GrpcTicketsRepository) CreateTicket(
	ctx context.Context,
	ticketData entities.CreateTicketDTO,
) (uint64, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.CreateTicket(
		ctx,
		&tickets.CreateTicketIn{
			RequestID:   requestID,
			UserID:      ticketData.UserID,
			CategoryID:  ticketData.CategoryID,
			Name:        ticketData.Name,
			Description: ticketData.Description,
			Price:       ticketData.Price,
			Quantity:    ticketData.Quantity,
			TagIDs:      ticketData.TagIDs,
		},
	)

	if err != nil {
		return 0, err
	}

	return response.GetTicketID(), nil
}

func (repo *GrpcTicketsRepository) GetTicketByID(ctx context.Context, id uint64) (*entities.RawTicket, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.GetTicket(
		ctx,
		&tickets.GetTicketIn{
			RequestID: requestID,
			ID:        id,
		},
	)

	if err != nil {
		return nil, err
	}

	return repo.processTicketResponse(response), nil
}

func (repo *GrpcTicketsRepository) GetAllTickets(ctx context.Context) ([]entities.RawTicket, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.GetTickets(
		ctx,
		&tickets.GetTicketsIn{RequestID: requestID},
	)

	if err != nil {
		return nil, err
	}

	allTickets := make([]entities.RawTicket, len(response.GetTickets()))
	for index, ticketResponse := range response.GetTickets() {
		allTickets[index] = *repo.processTicketResponse(ticketResponse)
	}

	return allTickets, nil
}

func (repo *GrpcTicketsRepository) GetUserTickets(ctx context.Context, userID uint64) ([]entities.RawTicket, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.GetUserTickets(
		ctx,
		&tickets.GetUserTicketsIn{
			RequestID: requestID,
			UserID:    userID,
		},
	)

	if err != nil {
		return nil, err
	}

	userTickets := make([]entities.RawTicket, len(response.GetTickets()))
	for index, ticketResponse := range response.GetTickets() {
		userTickets[index] = *repo.processTicketResponse(ticketResponse)
	}

	return userTickets, nil
}

func (repo *GrpcTicketsRepository) RespondToTicket(
	ctx context.Context,
	respondData entities.RespondToTicketDTO,
) (uint64, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.RespondToTicket(
		ctx,
		&tickets.RespondToTicketIn{
			RequestID: requestID,
			UserID:    respondData.UserID,
			TicketID:  respondData.TicketID,
		},
	)

	if err != nil {
		return 0, err
	}

	return response.GetRespondID(), nil
}

func (repo *GrpcTicketsRepository) GetRespondByID(ctx context.Context, id uint64) (*entities.Respond, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.GetRespond(
		ctx,
		&tickets.GetRespondIn{
			RequestID: requestID,
			ID:        id,
		},
	)

	if err != nil {
		return nil, err
	}

	return repo.processRespondResponse(response), nil
}

func (repo *GrpcTicketsRepository) GetTicketResponds(ctx context.Context, ticketID uint64) ([]entities.Respond, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.GetTicketResponds(
		ctx,
		&tickets.GetTicketRespondsIn{
			RequestID: requestID,
			TicketID:  ticketID,
		},
	)

	if err != nil {
		return nil, err
	}

	ticketResponds := make([]entities.Respond, len(response.GetResponds()))
	for index, respondResponse := range response.GetResponds() {
		ticketResponds[index] = *repo.processRespondResponse(respondResponse)
	}

	return ticketResponds, nil
}

func (repo *GrpcTicketsRepository) GetUserResponds(ctx context.Context, userID uint64) ([]entities.Respond, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.GetUserResponds(
		ctx,
		&tickets.GetUserRespondsIn{
			RequestID: requestID,
			UserID:    userID,
		},
	)

	if err != nil {
		return nil, err
	}

	userResponds := make([]entities.Respond, len(response.GetResponds()))
	for index, respondResponse := range response.GetResponds() {
		userResponds[index] = *repo.processRespondResponse(respondResponse)
	}

	return userResponds, nil
}

func (repo *GrpcTicketsRepository) processRespondResponse(respondResponse *tickets.GetRespondOut) *entities.Respond {
	return &entities.Respond{
		ID:        respondResponse.GetID(),
		MasterID:  respondResponse.GetMasterID(),
		TicketID:  respondResponse.GetTicketID(),
		CreatedAt: respondResponse.GetCreatedAt().AsTime(),
		UpdatedAt: respondResponse.GetUpdatedAt().AsTime(),
	}
}

func (repo *GrpcTicketsRepository) processTicketResponse(ticketResponse *tickets.GetTicketOut) *entities.RawTicket {
	return &entities.RawTicket{
		ID:          ticketResponse.GetID(),
		UserID:      ticketResponse.GetUserID(),
		CategoryID:  ticketResponse.GetCategoryID(),
		Name:        ticketResponse.GetName(),
		Description: ticketResponse.GetDescription(),
		Price:       ticketResponse.GetPrice(),
		Quantity:    ticketResponse.GetQuantity(),
		CreatedAt:   ticketResponse.GetCreatedAt().AsTime(),
		UpdatedAt:   ticketResponse.GetUpdatedAt().AsTime(),
		TagIDs:      ticketResponse.GetTagIDs(),
	}
}
