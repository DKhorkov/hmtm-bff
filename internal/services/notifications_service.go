package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

func NewNotificationsService(
	notificationsRepository interfaces.NotificationsRepository,
	logger *slog.Logger,
) *NotificationsService {
	return &NotificationsService{
		notificationsRepository: notificationsRepository,
		logger:                  logger,
	}
}

type NotificationsService struct {
	notificationsRepository interfaces.NotificationsRepository
	logger                  *slog.Logger
}

func (service *NotificationsService) GetUserEmailCommunications(
	ctx context.Context,
	userID uint64,
) ([]entities.Email, error) {
	emailCommunications, err := service.notificationsRepository.GetUserEmailCommunications(ctx, userID)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get Email Communications for User with ID=%d", userID),
			err,
		)
	}

	return emailCommunications, err
}
