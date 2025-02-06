package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

func NewCommonNotificationsService(
	notificationsRepository interfaces.NotificationsRepository,
	logger *slog.Logger,
) *CommonNotificationsService {
	return &CommonNotificationsService{
		notificationsRepository: notificationsRepository,
		logger:                  logger,
	}
}

type CommonNotificationsService struct {
	notificationsRepository interfaces.NotificationsRepository
	logger                  *slog.Logger
}

func (service *CommonNotificationsService) GetUserEmailCommunications(
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
