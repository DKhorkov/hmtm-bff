package services

import (
	"context"
	"fmt"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

type NotificationsService struct {
	notificationsRepository interfaces.NotificationsRepository
	logger                  logging.Logger
}

func NewNotificationsService(
	notificationsRepository interfaces.NotificationsRepository,
	logger logging.Logger,
) *NotificationsService {
	return &NotificationsService{
		notificationsRepository: notificationsRepository,
		logger:                  logger,
	}
}

func (service *NotificationsService) GetUserEmailCommunications(
	ctx context.Context,
	userID uint64,
) ([]entities.Email, error) {
	emailCommunications, err := service.notificationsRepository.GetUserEmailCommunications(
		ctx,
		userID,
	)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf(
				"Error occurred while trying to get Email Communications for User with ID=%d",
				userID,
			),
			err,
		)
	}

	return emailCommunications, err
}
