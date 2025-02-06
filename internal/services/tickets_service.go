package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

func NewCommonTicketsService(
	ticketsRepository interfaces.TicketsRepository,
	logger *slog.Logger,
) *CommonTicketsService {
	return &CommonTicketsService{
		ticketsRepository: ticketsRepository,
		logger:            logger,
	}
}

type CommonTicketsService struct {
	ticketsRepository interfaces.TicketsRepository
	logger            *slog.Logger
}

func (service *CommonTicketsService) CreateTicket(
	ctx context.Context,
	ticketData entities.CreateTicketDTO,
) (uint64, error) {
	ticketID, err := service.ticketsRepository.CreateTicket(ctx, ticketData)
	if err != nil {
		logging.LogErrorContext(ctx, service.logger, "Error occurred while trying to create new Ticket", err)
	}

	return ticketID, err
}

func (service *CommonTicketsService) GetTicketByID(ctx context.Context, id uint64) (*entities.RawTicket, error) {
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

func (service *CommonTicketsService) GetAllTickets(ctx context.Context) ([]entities.RawTicket, error) {
	tickets, err := service.ticketsRepository.GetAllTickets(ctx)
	if err != nil {
		logging.LogErrorContext(ctx, service.logger, "Error occurred while trying to get all Tickets", err)
	}

	return tickets, err
}

func (service *CommonTicketsService) GetUserTickets(ctx context.Context, userID uint64) ([]entities.RawTicket, error) {
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

func (service *CommonTicketsService) RespondToTicket(
	ctx context.Context,
	respondData entities.RespondToTicketDTO,
) (uint64, error) {
	respondID, err := service.ticketsRepository.RespondToTicket(ctx, respondData)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to respond to Ticket with ID=%d", respondData.TicketID),
			err,
		)
	}

	return respondID, err
}

func (service *CommonTicketsService) GetRespondByID(ctx context.Context, id uint64) (*entities.Respond, error) {
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

func (service *CommonTicketsService) GetTicketResponds(
	ctx context.Context,
	ticketID uint64,
) ([]entities.Respond, error) {
	responds, err := service.ticketsRepository.GetTicketResponds(ctx, ticketID)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get Responds for Ticket with ID=%d", ticketID),
			err,
		)
	}

	return responds, err
}

func (service *CommonTicketsService) GetUserResponds(ctx context.Context, userID uint64) ([]entities.Respond, error) {
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
