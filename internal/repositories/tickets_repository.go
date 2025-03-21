package repositories

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/DKhorkov/hmtm-tickets/api/protobuf/generated/go/tickets"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

func NewTicketsRepository(client interfaces.TicketsClient) *TicketsRepository {
	return &TicketsRepository{client: client}
}

type TicketsRepository struct {
	client interfaces.TicketsClient
}

func (repo *TicketsRepository) CreateTicket(
	ctx context.Context,
	ticketData entities.CreateTicketDTO,
) (uint64, error) {
	response, err := repo.client.CreateTicket(
		ctx,
		&tickets.CreateTicketIn{
			UserID:      ticketData.UserID,
			CategoryID:  ticketData.CategoryID,
			Name:        ticketData.Name,
			Description: ticketData.Description,
			Price:       ticketData.Price,
			Quantity:    ticketData.Quantity,
			TagIDs:      ticketData.TagIDs,
			Attachments: ticketData.Attachments,
		},
	)

	if err != nil {
		return 0, err
	}

	return response.GetTicketID(), nil
}

func (repo *TicketsRepository) GetTicketByID(ctx context.Context, id uint64) (*entities.RawTicket, error) {
	response, err := repo.client.GetTicket(
		ctx,
		&tickets.GetTicketIn{
			ID: id,
		},
	)

	if err != nil {
		return nil, err
	}

	return repo.processTicketResponse(response), nil
}

func (repo *TicketsRepository) GetAllTickets(ctx context.Context) ([]entities.RawTicket, error) {
	response, err := repo.client.GetTickets(
		ctx,
		&emptypb.Empty{},
	)

	if err != nil {
		return nil, err
	}

	allTickets := make([]entities.RawTicket, len(response.GetTickets()))
	for i, ticketResponse := range response.GetTickets() {
		allTickets[i] = *repo.processTicketResponse(ticketResponse)
	}

	return allTickets, nil
}

func (repo *TicketsRepository) GetUserTickets(ctx context.Context, userID uint64) ([]entities.RawTicket, error) {
	response, err := repo.client.GetUserTickets(
		ctx,
		&tickets.GetUserTicketsIn{
			UserID: userID,
		},
	)

	if err != nil {
		return nil, err
	}

	userTickets := make([]entities.RawTicket, len(response.GetTickets()))
	for i, ticketResponse := range response.GetTickets() {
		userTickets[i] = *repo.processTicketResponse(ticketResponse)
	}

	return userTickets, nil
}

func (repo *TicketsRepository) RespondToTicket(
	ctx context.Context,
	respondData entities.RespondToTicketDTO,
) (uint64, error) {
	response, err := repo.client.RespondToTicket(
		ctx,
		&tickets.RespondToTicketIn{
			UserID:   respondData.UserID,
			TicketID: respondData.TicketID,
			Price:    respondData.Price,
			Comment:  respondData.Comment,
		},
	)

	if err != nil {
		return 0, err
	}

	return response.GetRespondID(), nil
}

func (repo *TicketsRepository) GetRespondByID(ctx context.Context, id uint64) (*entities.Respond, error) {
	response, err := repo.client.GetRespond(
		ctx,
		&tickets.GetRespondIn{
			ID: id,
		},
	)

	if err != nil {
		return nil, err
	}

	return repo.processRespondResponse(response), nil
}

func (repo *TicketsRepository) GetTicketResponds(ctx context.Context, ticketID uint64) ([]entities.Respond, error) {
	response, err := repo.client.GetTicketResponds(
		ctx,
		&tickets.GetTicketRespondsIn{
			TicketID: ticketID,
		},
	)

	if err != nil {
		return nil, err
	}

	ticketResponds := make([]entities.Respond, len(response.GetResponds()))
	for i, respondResponse := range response.GetResponds() {
		ticketResponds[i] = *repo.processRespondResponse(respondResponse)
	}

	return ticketResponds, nil
}

func (repo *TicketsRepository) GetUserResponds(ctx context.Context, userID uint64) ([]entities.Respond, error) {
	response, err := repo.client.GetUserResponds(
		ctx,
		&tickets.GetUserRespondsIn{
			UserID: userID,
		},
	)

	if err != nil {
		return nil, err
	}

	userResponds := make([]entities.Respond, len(response.GetResponds()))
	for i, respondResponse := range response.GetResponds() {
		userResponds[i] = *repo.processRespondResponse(respondResponse)
	}

	return userResponds, nil
}

func (repo *TicketsRepository) processRespondResponse(respondResponse *tickets.GetRespondOut) *entities.Respond {
	return &entities.Respond{
		ID:        respondResponse.GetID(),
		MasterID:  respondResponse.GetMasterID(),
		TicketID:  respondResponse.GetTicketID(),
		Price:     respondResponse.GetPrice(),
		Comment:   respondResponse.Comment,
		CreatedAt: respondResponse.GetCreatedAt().AsTime(),
		UpdatedAt: respondResponse.GetUpdatedAt().AsTime(),
	}
}

func (repo *TicketsRepository) processTicketResponse(ticketResponse *tickets.GetTicketOut) *entities.RawTicket {
	attachments := make([]entities.TicketAttachment, len(ticketResponse.GetAttachments()))
	for i, attachment := range ticketResponse.GetAttachments() {
		attachments[i] = entities.TicketAttachment{
			ID:        attachment.GetID(),
			TicketID:  attachment.GetTicketID(),
			Link:      attachment.GetLink(),
			CreatedAt: attachment.GetCreatedAt().AsTime(),
			UpdatedAt: attachment.GetUpdatedAt().AsTime(),
		}
	}

	return &entities.RawTicket{
		ID:          ticketResponse.GetID(),
		UserID:      ticketResponse.GetUserID(),
		CategoryID:  ticketResponse.GetCategoryID(),
		Name:        ticketResponse.GetName(),
		Description: ticketResponse.GetDescription(),
		Price:       ticketResponse.Price,
		Quantity:    ticketResponse.GetQuantity(),
		CreatedAt:   ticketResponse.GetCreatedAt().AsTime(),
		UpdatedAt:   ticketResponse.GetUpdatedAt().AsTime(),
		TagIDs:      ticketResponse.GetTagIDs(),
		Attachments: attachments,
	}
}

func (repo *TicketsRepository) UpdateRespond(ctx context.Context, respondData entities.UpdateRespondDTO) error {
	_, err := repo.client.UpdateRespond(
		ctx,
		&tickets.UpdateRespondIn{
			ID:      respondData.ID,
			Price:   respondData.Price,
			Comment: respondData.Comment,
		},
	)

	return err
}

func (repo *TicketsRepository) DeleteRespond(ctx context.Context, id uint64) error {
	_, err := repo.client.DeleteRespond(
		ctx,
		&tickets.DeleteRespondIn{
			ID: id,
		},
	)

	return err
}

func (repo *TicketsRepository) UpdateTicket(ctx context.Context, ticketData entities.UpdateTicketDTO) error {
	_, err := repo.client.UpdateTicket(
		ctx,
		&tickets.UpdateTicketIn{
			ID:          ticketData.ID,
			Name:        ticketData.Name,
			Description: ticketData.Description,
			CategoryID:  ticketData.CategoryID,
			Price:       ticketData.Price,
			Quantity:    ticketData.Quantity,
			TagIDs:      ticketData.TagIDs,
			Attachments: ticketData.Attachments,
		},
	)

	return err
}

func (repo *TicketsRepository) DeleteTicket(ctx context.Context, id uint64) error {
	_, err := repo.client.DeleteTicket(
		ctx,
		&tickets.DeleteTicketIn{
			ID: id,
		},
	)

	return err
}
