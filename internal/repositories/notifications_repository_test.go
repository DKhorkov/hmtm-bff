package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/DKhorkov/hmtm-notifications/api/protobuf/generated/go/notifications"
	"github.com/DKhorkov/libs/pointers"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	mockclients "github.com/DKhorkov/hmtm-bff/mocks/clients"
)

func TestNotificationsRepository_GetUserEmailCommunications(t *testing.T) {
	sentAt := time.Now().UTC()

	ctrl := gomock.NewController(t)
	notificationsClient := mockclients.NewMockNotificationsClient(ctrl)
	repo := NewNotificationsRepository(notificationsClient)

	testCases := []struct {
		name          string
		pagination    *entities.Pagination
		userID        uint64
		setupMocks    func(notificationsClient *mockclients.MockNotificationsClient)
		expected      []entities.Email
		errorExpected bool
	}{
		{
			name:   "success with emails",
			userID: 1,
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			setupMocks: func(notificationsClient *mockclients.MockNotificationsClient) {
				notificationsClient.EXPECT().
					GetUserEmailCommunications(
						gomock.Any(),
						&notifications.GetUserEmailCommunicationsIn{
							UserID: 1,
							Pagination: &notifications.Pagination{
								Limit:  pointers.New[uint64](1),
								Offset: pointers.New[uint64](1),
							},
						},
					).
					Return(
						&notifications.GetUserEmailCommunicationsOut{
							Emails: []*notifications.Email{
								{
									ID:      1,
									UserID:  1,
									Email:   "test1@example.com",
									Content: "Content 1",
									SentAt:  timestamppb.New(sentAt),
								},
								{
									ID:      2,
									UserID:  1,
									Email:   "test2@example.com",
									Content: "Content 2",
									SentAt:  timestamppb.New(sentAt),
								},
							},
						},
						nil,
					).
					Times(1)
			},
			expected: []entities.Email{
				{
					ID:      1,
					UserID:  1,
					Email:   "test1@example.com",
					Content: "Content 1",
					SentAt:  sentAt.Truncate(time.Second),
				},
				{
					ID:      2,
					UserID:  1,
					Email:   "test2@example.com",
					Content: "Content 2",
					SentAt:  sentAt.Truncate(time.Second),
				},
			},
			errorExpected: false,
		},
		{
			name:   "success with empty list",
			userID: 1,
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			setupMocks: func(notificationsClient *mockclients.MockNotificationsClient) {
				notificationsClient.EXPECT().
					GetUserEmailCommunications(
						gomock.Any(),
						&notifications.GetUserEmailCommunicationsIn{
							UserID: 1,
							Pagination: &notifications.Pagination{
								Limit:  pointers.New[uint64](1),
								Offset: pointers.New[uint64](1),
							},
						},
					).
					Return(
						&notifications.GetUserEmailCommunicationsOut{
							Emails: []*notifications.Email{},
						},
						nil,
					).
					Times(1)
			},
			expected:      []entities.Email{},
			errorExpected: false,
		},
		{
			name:   "client error",
			userID: 1,
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			setupMocks: func(notificationsClient *mockclients.MockNotificationsClient) {
				notificationsClient.EXPECT().
					GetUserEmailCommunications(
						gomock.Any(),
						&notifications.GetUserEmailCommunicationsIn{
							UserID: 1,
							Pagination: &notifications.Pagination{
								Limit:  pointers.New[uint64](1),
								Offset: pointers.New[uint64](1),
							},
						},
					).
					Return(nil, errors.New("client error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(notificationsClient)
			}

			emails, err := repo.GetUserEmailCommunications(context.Background(), tc.userID, tc.pagination)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, emails)
			} else {
				require.NoError(t, err)
				// Усекаем время в полученных результатах для точного сравнения
				for i := range emails {
					emails[i].SentAt = emails[i].SentAt.Truncate(time.Second)
				}

				require.Equal(t, tc.expected, emails)
			}
		})
	}
}

func TestNotificationsRepository_CountUserEmailCommunications(t *testing.T) {
	ctrl := gomock.NewController(t)
	notificationsClient := mockclients.NewMockNotificationsClient(ctrl)
	repo := NewNotificationsRepository(notificationsClient)

	testCases := []struct {
		name          string
		userID        uint64
		setupMocks    func(notificationsClient *mockclients.MockNotificationsClient)
		expected      uint64
		errorExpected bool
	}{
		{
			name:   "success",
			userID: 1,
			setupMocks: func(notificationsClient *mockclients.MockNotificationsClient) {
				notificationsClient.EXPECT().
					CountUserEmailCommunications(
						gomock.Any(),
						&notifications.CountUserEmailCommunicationsIn{
							UserID: 1,
						},
					).
					Return(&notifications.CountOut{Count: 1}, nil).
					Times(1)
			},
			expected:      1,
			errorExpected: false,
		},
		{
			name:   "client error",
			userID: 1,
			setupMocks: func(notificationsClient *mockclients.MockNotificationsClient) {
				notificationsClient.EXPECT().
					CountUserEmailCommunications(
						gomock.Any(),
						&notifications.CountUserEmailCommunicationsIn{
							UserID: 1,
						},
					).
					Return(nil, errors.New("client error")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(notificationsClient)
			}

			actual, err := repo.CountUserEmailCommunications(context.Background(), tc.userID)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}
