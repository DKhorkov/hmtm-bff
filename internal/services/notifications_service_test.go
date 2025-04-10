package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	mocklogging "github.com/DKhorkov/libs/logging/mocks"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	mockrepositories "github.com/DKhorkov/hmtm-bff/mocks/repositories"
)

var (
	now = time.Now()
)

func TestNotificationsService_GetUserEmailCommunications(t *testing.T) {
	ctrl := gomock.NewController(t)
	notificationsRepository := mockrepositories.NewMockNotificationsRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewNotificationsService(notificationsRepository, logger)

	testCases := []struct {
		name           string
		userID         uint64
		setupMocks     func(notificationsRepository *mockrepositories.MockNotificationsRepository, logger *mocklogging.MockLogger)
		expectedEmails []entities.Email
		errorExpected  bool
	}{
		{
			name:   "success",
			userID: 1,
			setupMocks: func(notificationsRepository *mockrepositories.MockNotificationsRepository, logger *mocklogging.MockLogger) {
				notificationsRepository.
					EXPECT().
					GetUserEmailCommunications(gomock.Any(), uint64(1)).
					Return([]entities.Email{
						{
							ID:      1,
							UserID:  1,
							Content: "Test Email",
							SentAt:  now,
						},
					}, nil).
					Times(1)
			},
			expectedEmails: []entities.Email{
				{
					ID:      1,
					UserID:  1,
					Content: "Test Email",
					SentAt:  now,
				},
			},
			errorExpected: false,
		},
		{
			name:   "error",
			userID: 1,
			setupMocks: func(notificationsRepository *mockrepositories.MockNotificationsRepository, logger *mocklogging.MockLogger) {
				notificationsRepository.
					EXPECT().
					GetUserEmailCommunications(gomock.Any(), uint64(1)).
					Return(nil, errors.New("fetch failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedEmails: nil,
			errorExpected:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(notificationsRepository, logger)
			}

			emails, err := service.GetUserEmailCommunications(context.Background(), tc.userID)
			if tc.errorExpected {
				require.Error(t, err)
				require.Equal(t, "fetch failed", err.Error())
				require.Nil(t, emails)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedEmails, emails)
			}
		})
	}
}
