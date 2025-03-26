package services

import (
	"context"
	"fmt"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

func NewTicketsService(
	ticketsRepository interfaces.TicketsRepository,
	logger logging.Logger,
) *TicketsService {
	return &TicketsService{
		ticketsRepository: ticketsRepository,
		logger:            logger,
	}
}

type TicketsService struct {
	ticketsRepository interfaces.TicketsRepository
	logger            logging.Logger
}

func (service *TicketsService) CreateTicket(
	ctx context.Context,
	ticketData entities.CreateTicketDTO,
) (uint64, error) {
	ticketID, err := service.ticketsRepository.CreateTicket(ctx, ticketData)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to create new Ticket",
			err,
		)
	}

	return ticketID, err
}

func (service *TicketsService) GetTicketByID(
	ctx context.Context,
	id uint64,
) (*entities.RawTicket, error) {
	ticket, err := service.ticketsRepository.GetTicketByID(ctx, id)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get Ticket with ID=%d", id),
			err,
		)
	}

	return ticket, err
}

func (service *TicketsService) GetAllTickets(ctx context.Context) ([]entities.RawTicket, error) {
	tickets, err := service.ticketsRepository.GetAllTickets(ctx)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to get all Tickets",
			err,
		)
	}

	return tickets, err
}

func (service *TicketsService) GetUserTickets(
	ctx context.Context,
	userID uint64,
) ([]entities.RawTicket, error) {
	tickets, err := service.ticketsRepository.GetUserTickets(ctx, userID)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get Tickets for User with ID=%d", userID),
			err,
		)
	}

	return tickets, err
}

func (service *TicketsService) RespondToTicket(
	ctx context.Context,
	respondData entities.RespondToTicketDTO,
) (uint64, error) {
	respondID, err := service.ticketsRepository.RespondToTicket(ctx, respondData)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf(
				"Error occurred while trying to respond to Ticket with ID=%d",
				respondData.TicketID,
			),
			err,
		)
	}

	return respondID, err
}

func (service *TicketsService) GetRespondByID(
	ctx context.Context,
	id uint64,
) (*entities.Respond, error) {
	respond, err := service.ticketsRepository.GetRespondByID(ctx, id)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get Respond with ID=%d", id),
			err,
		)
	}

	return respond, err
}

func (service *TicketsService) GetTicketResponds(
	ctx context.Context,
	ticketID uint64,
) ([]entities.Respond, error) {
	responds, err := service.ticketsRepository.GetTicketResponds(ctx, ticketID)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf(
				"Error occurred while trying to get Responds for Ticket with ID=%d",
				ticketID,
			),
			err,
		)
	}

	return responds, err
}

func (service *TicketsService) GetUserResponds(
	ctx context.Context,
	userID uint64,
) ([]entities.Respond, error) {
	responds, err := service.ticketsRepository.GetUserResponds(ctx, userID)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get Responds for User with ID=%d", userID),
			err,
		)
	}

	return responds, err
}

func (service *TicketsService) UpdateRespond(
	ctx context.Context,
	respondData entities.UpdateRespondDTO,
) error {
	err := service.ticketsRepository.UpdateRespond(ctx, respondData)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to update Respond with ID=%d", respondData.ID),
			err,
		)
	}

	return err
}

func (service *TicketsService) DeleteRespond(ctx context.Context, id uint64) error {
	err := service.ticketsRepository.DeleteRespond(ctx, id)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to delete Respond with ID=%d", id),
			err,
		)
	}

	return err
}

func (service *TicketsService) UpdateTicket(
	ctx context.Context,
	ticketData entities.UpdateTicketDTO,
) error {
	err := service.ticketsRepository.UpdateTicket(ctx, ticketData)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to update Ticket with ID=%d", ticketData.ID),
			err,
		)
	}

	return err
}

func (service *TicketsService) DeleteTicket(ctx context.Context, id uint64) error {
	err := service.ticketsRepository.DeleteTicket(ctx, id)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to delete Ticket with ID=%d", id),
			err,
		)
	}

	return err
}
