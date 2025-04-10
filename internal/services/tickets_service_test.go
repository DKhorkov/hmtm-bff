package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	mockrepositories "github.com/DKhorkov/hmtm-bff/mocks/repositories"
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
)

func TestTicketsService_CreateTicket(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsRepository := mockrepositories.NewMockTicketsRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewTicketsService(ticketsRepository, logger)

	price := float32(100.0)
	testCases := []struct {
		name          string
		ticketData    entities.CreateTicketDTO
		setupMocks    func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger)
		expectedID    uint64
		errorExpected bool
	}{
		{
			name: "success",
			ticketData: entities.CreateTicketDTO{
				UserID:      1,
				CategoryID:  1,
				Name:        "Test Ticket",
				Description: "Test Description",
				Price:       &price,
				Quantity:    1,
				TagIDs:      []uint32{1, 2},
				Attachments: []string{"link1"},
			},
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					CreateTicket(gomock.Any(), entities.CreateTicketDTO{
						UserID:      1,
						CategoryID:  1,
						Name:        "Test Ticket",
						Description: "Test Description",
						Price:       &price,
						Quantity:    1,
						TagIDs:      []uint32{1, 2},
						Attachments: []string{"link1"},
					}).
					Return(uint64(1), nil).
					Times(1)
			},
			expectedID:    1,
			errorExpected: false,
		},
		{
			name: "error",
			ticketData: entities.CreateTicketDTO{
				UserID:      1,
				CategoryID:  1,
				Name:        "Test Ticket",
				Description: "Test Description",
				Price:       &price,
				Quantity:    1,
				TagIDs:      []uint32{1, 2},
				Attachments: []string{"link1"},
			},
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					CreateTicket(gomock.Any(), entities.CreateTicketDTO{
						UserID:      1,
						CategoryID:  1,
						Name:        "Test Ticket",
						Description: "Test Description",
						Price:       &price,
						Quantity:    1,
						TagIDs:      []uint32{1, 2},
						Attachments: []string{"link1"},
					}).
					Return(uint64(0), errors.New("create failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedID:    0,
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsRepository, logger)
			}

			ticketID, err := service.CreateTicket(context.Background(), tc.ticketData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedID, ticketID)
		})
	}
}

func TestTicketsService_GetTicketByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsRepository := mockrepositories.NewMockTicketsRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewTicketsService(ticketsRepository, logger)

	price := float32(100.0)
	now := time.Now()
	testCases := []struct {
		name           string
		id             uint64
		setupMocks     func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger)
		expectedTicket *entities.RawTicket
		errorExpected  bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(1)).
					Return(&entities.RawTicket{
						ID:          1,
						UserID:      1,
						CategoryID:  1,
						Name:        "Test Ticket",
						Description: "Test Description",
						Price:       &price,
						Quantity:    1,
						CreatedAt:   now,
						UpdatedAt:   now,
						TagIDs:      []uint32{1, 2},
						Attachments: []entities.TicketAttachment{{ID: 1, TicketID: 1, Link: "link1"}},
					}, nil).
					Times(1)
			},
			expectedTicket: &entities.RawTicket{
				ID:          1,
				UserID:      1,
				CategoryID:  1,
				Name:        "Test Ticket",
				Description: "Test Description",
				Price:       &price,
				Quantity:    1,
				CreatedAt:   now,
				UpdatedAt:   now,
				TagIDs:      []uint32{1, 2},
				Attachments: []entities.TicketAttachment{{ID: 1, TicketID: 1, Link: "link1"}},
			},
			errorExpected: false,
		},
		{
			name: "error",
			id:   1,
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(1)).
					Return(nil, errors.New("not found")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedTicket: nil,
			errorExpected:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsRepository, logger)
			}

			ticket, err := service.GetTicketByID(context.Background(), tc.id)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, ticket)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedTicket, ticket)
			}
		})
	}
}

func TestTicketsService_GetAllTickets(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsRepository := mockrepositories.NewMockTicketsRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewTicketsService(ticketsRepository, logger)

	price := float32(100.0)
	now := time.Now()
	testCases := []struct {
		name            string
		setupMocks      func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger)
		expectedTickets []entities.RawTicket
		errorExpected   bool
	}{
		{
			name: "success",
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					GetAllTickets(gomock.Any()).
					Return([]entities.RawTicket{
						{
							ID:          1,
							UserID:      1,
							CategoryID:  1,
							Name:        "Test Ticket",
							Description: "Test Description",
							Price:       &price,
							Quantity:    1,
							CreatedAt:   now,
							UpdatedAt:   now,
							TagIDs:      []uint32{1, 2},
							Attachments: []entities.TicketAttachment{{ID: 1, TicketID: 1, Link: "link1"}},
						},
					}, nil).
					Times(1)
			},
			expectedTickets: []entities.RawTicket{
				{
					ID:          1,
					UserID:      1,
					CategoryID:  1,
					Name:        "Test Ticket",
					Description: "Test Description",
					Price:       &price,
					Quantity:    1,
					CreatedAt:   now,
					UpdatedAt:   now,
					TagIDs:      []uint32{1, 2},
					Attachments: []entities.TicketAttachment{{ID: 1, TicketID: 1, Link: "link1"}},
				},
			},
			errorExpected: false,
		},
		{
			name: "error",
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					GetAllTickets(gomock.Any()).
					Return(nil, errors.New("fetch failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedTickets: nil,
			errorExpected:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsRepository, logger)
			}

			tickets, err := service.GetAllTickets(context.Background())
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, tickets)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedTickets, tickets)
			}
		})
	}
}

func TestTicketsService_GetUserTickets(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsRepository := mockrepositories.NewMockTicketsRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewTicketsService(ticketsRepository, logger)

	price := float32(100.0)
	now := time.Now()
	testCases := []struct {
		name            string
		userID          uint64
		setupMocks      func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger)
		expectedTickets []entities.RawTicket
		errorExpected   bool
	}{
		{
			name:   "success",
			userID: 1,
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					GetUserTickets(gomock.Any(), uint64(1)).
					Return([]entities.RawTicket{
						{
							ID:          1,
							UserID:      1,
							CategoryID:  1,
							Name:        "Test Ticket",
							Description: "Test Description",
							Price:       &price,
							Quantity:    1,
							CreatedAt:   now,
							UpdatedAt:   now,
							TagIDs:      []uint32{1, 2},
							Attachments: []entities.TicketAttachment{{ID: 1, TicketID: 1, Link: "link1"}},
						},
					}, nil).
					Times(1)
			},
			expectedTickets: []entities.RawTicket{
				{
					ID:          1,
					UserID:      1,
					CategoryID:  1,
					Name:        "Test Ticket",
					Description: "Test Description",
					Price:       &price,
					Quantity:    1,
					CreatedAt:   now,
					UpdatedAt:   now,
					TagIDs:      []uint32{1, 2},
					Attachments: []entities.TicketAttachment{{ID: 1, TicketID: 1, Link: "link1"}},
				},
			},
			errorExpected: false,
		},
		{
			name:   "error",
			userID: 1,
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					GetUserTickets(gomock.Any(), uint64(1)).
					Return(nil, errors.New("fetch failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedTickets: nil,
			errorExpected:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsRepository, logger)
			}

			tickets, err := service.GetUserTickets(context.Background(), tc.userID)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, tickets)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedTickets, tickets)
			}
		})
	}
}

func TestTicketsService_RespondToTicket(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsRepository := mockrepositories.NewMockTicketsRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewTicketsService(ticketsRepository, logger)

	price := float32(50.0)
	comment := "Test Comment"
	testCases := []struct {
		name          string
		respondData   entities.RespondToTicketDTO
		setupMocks    func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger)
		expectedID    uint64
		errorExpected bool
	}{
		{
			name: "success",
			respondData: entities.RespondToTicketDTO{
				TicketID: 1,
				UserID:   2,
				Price:    price,
				Comment:  &comment,
			},
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					RespondToTicket(gomock.Any(), entities.RespondToTicketDTO{
						TicketID: 1,
						UserID:   2,
						Price:    price,
						Comment:  &comment,
					}).
					Return(uint64(1), nil).
					Times(1)
			},
			expectedID:    1,
			errorExpected: false,
		},
		{
			name: "error",
			respondData: entities.RespondToTicketDTO{
				TicketID: 1,
				UserID:   2,
				Price:    price,
				Comment:  &comment,
			},
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					RespondToTicket(gomock.Any(), entities.RespondToTicketDTO{
						TicketID: 1,
						UserID:   2,
						Price:    price,
						Comment:  &comment,
					}).
					Return(uint64(0), errors.New("respond failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedID:    0,
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsRepository, logger)
			}

			respondID, err := service.RespondToTicket(context.Background(), tc.respondData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedID, respondID)
		})
	}
}

func TestTicketsService_GetRespondByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsRepository := mockrepositories.NewMockTicketsRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewTicketsService(ticketsRepository, logger)

	price := float32(50.0)
	comment := "Test Comment"
	now := time.Now()
	testCases := []struct {
		name            string
		id              uint64
		setupMocks      func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger)
		expectedRespond *entities.Respond
		errorExpected   bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					GetRespondByID(gomock.Any(), uint64(1)).
					Return(&entities.Respond{
						ID:        1,
						TicketID:  1,
						MasterID:  2,
						Price:     price,
						Comment:   &comment,
						CreatedAt: now,
						UpdatedAt: now,
					}, nil).
					Times(1)
			},
			expectedRespond: &entities.Respond{
				ID:        1,
				TicketID:  1,
				MasterID:  2,
				Price:     price,
				Comment:   &comment,
				CreatedAt: now,
				UpdatedAt: now,
			},
			errorExpected: false,
		},
		{
			name: "error",
			id:   1,
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					GetRespondByID(gomock.Any(), uint64(1)).
					Return(nil, errors.New("not found")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedRespond: nil,
			errorExpected:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsRepository, logger)
			}

			respond, err := service.GetRespondByID(context.Background(), tc.id)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, respond)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedRespond, respond)
			}
		})
	}
}

func TestTicketsService_GetTicketResponds(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsRepository := mockrepositories.NewMockTicketsRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewTicketsService(ticketsRepository, logger)

	price := float32(50.0)
	comment := "Test Comment"
	now := time.Now()
	testCases := []struct {
		name             string
		ticketID         uint64
		setupMocks       func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger)
		expectedResponds []entities.Respond
		errorExpected    bool
	}{
		{
			name:     "success",
			ticketID: 1,
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					GetTicketResponds(gomock.Any(), uint64(1)).
					Return([]entities.Respond{
						{
							ID:        1,
							TicketID:  1,
							MasterID:  2,
							Price:     price,
							Comment:   &comment,
							CreatedAt: now,
							UpdatedAt: now,
						},
					}, nil).
					Times(1)
			},
			expectedResponds: []entities.Respond{
				{
					ID:        1,
					TicketID:  1,
					MasterID:  2,
					Price:     price,
					Comment:   &comment,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			errorExpected: false,
		},
		{
			name:     "error",
			ticketID: 1,
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					GetTicketResponds(gomock.Any(), uint64(1)).
					Return(nil, errors.New("fetch failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedResponds: nil,
			errorExpected:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsRepository, logger)
			}

			responds, err := service.GetTicketResponds(context.Background(), tc.ticketID)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, responds)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResponds, responds)
			}
		})
	}
}

func TestTicketsService_GetUserResponds(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsRepository := mockrepositories.NewMockTicketsRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewTicketsService(ticketsRepository, logger)

	price := float32(50.0)
	comment := "Test Comment"
	now := time.Now()
	testCases := []struct {
		name             string
		userID           uint64
		setupMocks       func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger)
		expectedResponds []entities.Respond
		errorExpected    bool
	}{
		{
			name:   "success",
			userID: 2,
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					GetUserResponds(gomock.Any(), uint64(2)).
					Return([]entities.Respond{
						{
							ID:        1,
							TicketID:  1,
							MasterID:  2,
							Price:     price,
							Comment:   &comment,
							CreatedAt: now,
							UpdatedAt: now,
						},
					}, nil).
					Times(1)
			},
			expectedResponds: []entities.Respond{
				{
					ID:        1,
					TicketID:  1,
					MasterID:  2,
					Price:     price,
					Comment:   &comment,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			errorExpected: false,
		},
		{
			name:   "error",
			userID: 2,
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					GetUserResponds(gomock.Any(), uint64(2)).
					Return(nil, errors.New("fetch failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedResponds: nil,
			errorExpected:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsRepository, logger)
			}

			responds, err := service.GetUserResponds(context.Background(), tc.userID)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, responds)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResponds, responds)
			}
		})
	}
}

func TestTicketsService_UpdateRespond(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsRepository := mockrepositories.NewMockTicketsRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewTicketsService(ticketsRepository, logger)

	price := float32(75.0)
	comment := "Updated Comment"
	testCases := []struct {
		name          string
		respondData   entities.UpdateRespondDTO
		setupMocks    func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger)
		errorExpected bool
	}{
		{
			name: "success",
			respondData: entities.UpdateRespondDTO{
				ID:      1,
				Price:   &price,
				Comment: &comment,
			},
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					UpdateRespond(gomock.Any(), entities.UpdateRespondDTO{
						ID:      1,
						Price:   &price,
						Comment: &comment,
					}).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "error",
			respondData: entities.UpdateRespondDTO{
				ID:      1,
				Price:   &price,
				Comment: &comment,
			},
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					UpdateRespond(gomock.Any(), entities.UpdateRespondDTO{
						ID:      1,
						Price:   &price,
						Comment: &comment,
					}).
					Return(errors.New("update failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsRepository, logger)
			}

			err := service.UpdateRespond(context.Background(), tc.respondData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTicketsService_DeleteRespond(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsRepository := mockrepositories.NewMockTicketsRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewTicketsService(ticketsRepository, logger)

	testCases := []struct {
		name          string
		id            uint64
		setupMocks    func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger)
		errorExpected bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					DeleteRespond(gomock.Any(), uint64(1)).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "error",
			id:   1,
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					DeleteRespond(gomock.Any(), uint64(1)).
					Return(errors.New("delete failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsRepository, logger)
			}

			err := service.DeleteRespond(context.Background(), tc.id)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTicketsService_UpdateTicket(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsRepository := mockrepositories.NewMockTicketsRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewTicketsService(ticketsRepository, logger)

	price := float32(150.0)
	name := "Updated Ticket"
	quantity := uint32(2)
	testCases := []struct {
		name          string
		ticketData    entities.UpdateTicketDTO
		setupMocks    func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger)
		errorExpected bool
	}{
		{
			name: "success",
			ticketData: entities.UpdateTicketDTO{
				ID:       1,
				Name:     &name,
				Price:    &price,
				Quantity: &quantity,
				TagIDs:   []uint32{1, 3},
			},
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					UpdateTicket(gomock.Any(), entities.UpdateTicketDTO{
						ID:       1,
						Name:     &name,
						Price:    &price,
						Quantity: &quantity,
						TagIDs:   []uint32{1, 3},
					}).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "error",
			ticketData: entities.UpdateTicketDTO{
				ID:       1,
				Name:     &name,
				Price:    &price,
				Quantity: &quantity,
				TagIDs:   []uint32{1, 3},
			},
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					UpdateTicket(gomock.Any(), entities.UpdateTicketDTO{
						ID:       1,
						Name:     &name,
						Price:    &price,
						Quantity: &quantity,
						TagIDs:   []uint32{1, 3},
					}).
					Return(errors.New("update failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsRepository, logger)
			}

			err := service.UpdateTicket(context.Background(), tc.ticketData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTicketsService_DeleteTicket(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsRepository := mockrepositories.NewMockTicketsRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewTicketsService(ticketsRepository, logger)

	testCases := []struct {
		name          string
		id            uint64
		setupMocks    func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger)
		errorExpected bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					DeleteTicket(gomock.Any(), uint64(1)).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "error",
			id:   1,
			setupMocks: func(ticketsRepository *mockrepositories.MockTicketsRepository, logger *mocklogging.MockLogger) {
				ticketsRepository.
					EXPECT().
					DeleteTicket(gomock.Any(), uint64(1)).
					Return(errors.New("delete failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsRepository, logger)
			}

			err := service.DeleteTicket(context.Background(), tc.id)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
