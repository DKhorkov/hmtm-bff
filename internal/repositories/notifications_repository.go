package repositories

import (
	"context"

	"github.com/DKhorkov/hmtm-notifications/api/protobuf/generated/go/notifications"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

type NotificationsRepository struct {
	client interfaces.NotificationsClient
}

func NewNotificationsRepository(client interfaces.NotificationsClient) *NotificationsRepository {
	return &NotificationsRepository{client: client}
}

func (repo *NotificationsRepository) GetUserEmailCommunications(
	ctx context.Context,
	userID uint64,
	pagination *entities.Pagination,
) ([]entities.Email, error) {
	in := &notifications.GetUserEmailCommunicationsIn{UserID: userID}
	if pagination != nil {
		in.Pagination = &notifications.Pagination{
			Limit:  pagination.Limit,
			Offset: pagination.Offset,
		}
	}

	response, err := repo.client.GetUserEmailCommunications(
		ctx,
		in,
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

func (repo *NotificationsRepository) CountUserEmailCommunications(ctx context.Context, userID uint64) (uint64, error) {
	response, err := repo.client.CountUserEmailCommunications(
		ctx,
		&notifications.CountUserEmailCommunicationsIn{
			UserID: userID,
		},
	)
	if err != nil {
		return 0, err
	}

	return response.GetCount(), nil
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
