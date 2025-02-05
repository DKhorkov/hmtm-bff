package repositories

import (
	"context"

	"github.com/DKhorkov/hmtm-notifications/api/protobuf/generated/go/notifications"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

func NewGrpcNotificationsRepository(client interfaces.NotificationsGrpcClient) *GrpcNotificationsRepository {
	return &GrpcNotificationsRepository{client: client}
}

type GrpcNotificationsRepository struct {
	client interfaces.NotificationsGrpcClient
}

func (repo *GrpcNotificationsRepository) GetUserEmailCommunications(
	ctx context.Context,
	userID uint64,
) ([]entities.Email, error) {
	response, err := repo.client.GetUserEmailCommunications(
		ctx,
		&notifications.GetUserEmailCommunicationsIn{
			UserID: userID,
		},
	)

	if err != nil {
		return nil, err
	}

	emailCommunications := make([]entities.Email, len(response.GetEmails()))
	for i, communicationResponse := range response.GetEmails() {
		emailCommunications[i] = *repo.processEmailCommunicationResponse(communicationResponse)
	}

	return emailCommunications, nil
}

func (repo *GrpcNotificationsRepository) processEmailCommunicationResponse(
	emailCommunicationResponse *notifications.Email,
) *entities.Email {
	return &entities.Email{
		ID:      emailCommunicationResponse.GetID(),
		UserID:  emailCommunicationResponse.GetUserID(),
		Email:   emailCommunicationResponse.GetEmail(),
		Content: emailCommunicationResponse.GetContent(),
		SentAt:  emailCommunicationResponse.GetSentAt().AsTime(),
	}
}
