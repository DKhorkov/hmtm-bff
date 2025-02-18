package repositories

import (
	"context"

	"github.com/DKhorkov/hmtm-notifications/api/protobuf/generated/go/notifications"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

func NewNotificationsRepository(client interfaces.NotificationsClient) *NotificationsRepository {
	return &NotificationsRepository{client: client}
}

type NotificationsRepository struct {
	client interfaces.NotificationsClient
}

func (repo *NotificationsRepository) GetUserEmailCommunications(
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

func (repo *NotificationsRepository) processEmailCommunicationResponse(
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
