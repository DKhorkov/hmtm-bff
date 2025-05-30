package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/DKhorkov/hmtm-tickets/api/protobuf/generated/go/tickets"
	"github.com/DKhorkov/libs/pointers"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	mockclients "github.com/DKhorkov/hmtm-bff/mocks/clients"
)

func TestTicketsRepository_CreateTicket(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsClient := mockclients.NewMockTicketsClient(ctrl)
	repo := NewTicketsRepository(ticketsClient)

	testCases := []struct {
		name             string
		ticketData       entities.CreateTicketDTO
		setupMocks       func(ticketsClient *mockclients.MockTicketsClient)
		expectedTicketID uint64
		errorExpected    bool
	}{
		{
			name: "success",
			ticketData: entities.CreateTicketDTO{
				UserID:      1,
				CategoryID:  2,
				Name:        "Test Ticket",
				Description: "Test Description",
				Price:       pointers.New[float32](100),
				Quantity:    5,
				TagIDs:      []uint32{1, 2},
				Attachments: []string{"attachment1"},
			},
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					CreateTicket(
						gomock.Any(),
						&tickets.CreateTicketIn{
							UserID:      1,
							CategoryID:  2,
							Name:        "Test Ticket",
							Description: "Test Description",
							Price:       pointers.New[float32](100),
							Quantity:    5,
							TagIDs:      []uint32{1, 2},
							Attachments: []string{"attachment1"},
						},
					).
					Return(&tickets.CreateTicketOut{TicketID: 1}, nil).
					Times(1)
			},
			expectedTicketID: 1,
			errorExpected:    false,
		},
		{
			name: "error",
			ticketData: entities.CreateTicketDTO{
				UserID: 1,
				Name:   "Test Ticket",
			},
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					CreateTicket(
						gomock.Any(),
						&tickets.CreateTicketIn{
							UserID: 1,
							Name:   "Test Ticket",
						},
					).
					Return(nil, errors.New("creation failed")).
					Times(1)
			},
			expectedTicketID: 0,
			errorExpected:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsClient)
			}

			ticketID, err := repo.CreateTicket(context.Background(), tc.ticketData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedTicketID, ticketID)
		})
	}
}

func TestTicketsRepository_GetTicketByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsClient := mockclients.NewMockTicketsClient(ctrl)
	repo := NewTicketsRepository(ticketsClient)

	now := time.Now().UTC().Truncate(time.Second)

	testCases := []struct {
		name           string
		id             uint64
		setupMocks     func(ticketsClient *mockclients.MockTicketsClient)
		expectedTicket *entities.RawTicket
		errorExpected  bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					GetTicket(
						gomock.Any(),
						&tickets.GetTicketIn{ID: 1},
					).
					Return(&tickets.GetTicketOut{
						ID:          1,
						UserID:      1,
						CategoryID:  2,
						Name:        "Test Ticket",
						Description: "Description",
						Price:       pointers.New[float32](100),
						Quantity:    5,
						CreatedAt:   timestamppb.New(now),
						UpdatedAt:   timestamppb.New(now),
						TagIDs:      []uint32{1, 2},
						Attachments: []*tickets.Attachment{
							{
								ID:        1,
								TicketID:  1,
								Link:      "link1",
								CreatedAt: timestamppb.New(now),
								UpdatedAt: timestamppb.New(now),
							},
						},
					}, nil).
					Times(1)
			},
			expectedTicket: &entities.RawTicket{
				ID:          1,
				UserID:      1,
				CategoryID:  2,
				Name:        "Test Ticket",
				Description: "Description",
				Price:       pointers.New[float32](100),
				Quantity:    5,
				CreatedAt:   now,
				UpdatedAt:   now,
				TagIDs:      []uint32{1, 2},
				Attachments: []entities.TicketAttachment{
					{ID: 1, TicketID: 1, Link: "link1", CreatedAt: now, UpdatedAt: now},
				},
			},
			errorExpected: false,
		},
		{
			name: "error",
			id:   1,
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					GetTicket(
						gomock.Any(),
						&tickets.GetTicketIn{ID: 1},
					).
					Return(nil, errors.New("ticket not found")).
					Times(1)
			},
			expectedTicket: nil,
			errorExpected:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsClient)
			}

			ticket, err := repo.GetTicketByID(context.Background(), tc.id)
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

func TestTicketsRepository_GetTickets(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsClient := mockclients.NewMockTicketsClient(ctrl)
	repo := NewTicketsRepository(ticketsClient)

	now := time.Now().UTC().Truncate(time.Second)

	testCases := []struct {
		name            string
		pagination      *entities.Pagination
		filters         *entities.TicketsFilters
		setupMocks      func(ticketsClient *mockclients.MockTicketsClient)
		expectedTickets []entities.RawTicket
		errorExpected   bool
	}{
		{
			name: "success",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			filters: &entities.TicketsFilters{
				Search:              pointers.New("Ticket"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryIDs:         []uint32{1},
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					GetTickets(
						gomock.Any(),
						&tickets.GetTicketsIn{
							Pagination: &tickets.Pagination{
								Limit:  pointers.New[uint64](1),
								Offset: pointers.New[uint64](1),
							},
							Filters: &tickets.TicketsFilters{
								Search:              pointers.New("Ticket"),
								PriceCeil:           pointers.New[float32](1000),
								PriceFloor:          pointers.New[float32](10),
								QuantityFloor:       pointers.New[uint32](1),
								CategoryIDs:         []uint32{1},
								TagIDs:              []uint32{1},
								CreatedAtOrderByAsc: pointers.New(true),
							},
						},
					).
					Return(&tickets.GetTicketsOut{
						Tickets: []*tickets.GetTicketOut{
							{
								ID:          1,
								UserID:      1,
								CategoryID:  2,
								Name:        "Ticket1",
								Description: "Desc1",
								Price:       pointers.New[float32](100),
								Quantity:    5,
								CreatedAt:   timestamppb.New(now),
								UpdatedAt:   timestamppb.New(now),
							},
							{
								ID:          2,
								UserID:      2,
								CategoryID:  3,
								Name:        "Ticket2",
								Description: "Desc2",
								Price:       pointers.New[float32](100),
								Quantity:    10,
								CreatedAt:   timestamppb.New(now),
								UpdatedAt:   timestamppb.New(now),
							},
						},
					}, nil).
					Times(1)
			},
			expectedTickets: []entities.RawTicket{
				{
					ID:          1,
					UserID:      1,
					CategoryID:  2,
					Name:        "Ticket1",
					Description: "Desc1",
					Price:       pointers.New[float32](100),
					Quantity:    5,
					CreatedAt:   now,
					UpdatedAt:   now,
					Attachments: make([]entities.TicketAttachment, 0),
				},
				{
					ID:          2,
					UserID:      2,
					CategoryID:  3,
					Name:        "Ticket2",
					Description: "Desc2",
					Price:       pointers.New[float32](100),
					Quantity:    10,
					CreatedAt:   now,
					UpdatedAt:   now,
					Attachments: make([]entities.TicketAttachment, 0),
				},
			},
			errorExpected: false,
		},
		{
			name: "error",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			filters: &entities.TicketsFilters{
				Search:              pointers.New("Ticket"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryIDs:         []uint32{1},
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					GetTickets(
						gomock.Any(),
						&tickets.GetTicketsIn{
							Pagination: &tickets.Pagination{
								Limit:  pointers.New[uint64](1),
								Offset: pointers.New[uint64](1),
							},
							Filters: &tickets.TicketsFilters{
								Search:              pointers.New("Ticket"),
								PriceCeil:           pointers.New[float32](1000),
								PriceFloor:          pointers.New[float32](10),
								QuantityFloor:       pointers.New[uint32](1),
								CategoryIDs:         []uint32{1},
								TagIDs:              []uint32{1},
								CreatedAtOrderByAsc: pointers.New(true),
							},
						},
					).
					Return(nil, errors.New("fetch failed")).
					Times(1)
			},
			expectedTickets: nil,
			errorExpected:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsClient)
			}

			ticketsList, err := repo.GetTickets(context.Background(), tc.pagination, tc.filters)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, ticketsList)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedTickets, ticketsList)
			}
		})
	}
}

func TestTicketsRepository_GetUserTickets(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsClient := mockclients.NewMockTicketsClient(ctrl)
	repo := NewTicketsRepository(ticketsClient)

	now := time.Now().UTC().Truncate(time.Second)

	testCases := []struct {
		name            string
		pagination      *entities.Pagination
		filters         *entities.TicketsFilters
		userID          uint64
		setupMocks      func(ticketsClient *mockclients.MockTicketsClient)
		expectedTickets []entities.RawTicket
		errorExpected   bool
	}{
		{
			name: "success",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			filters: &entities.TicketsFilters{
				Search:              pointers.New("Ticket"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryIDs:         []uint32{1},
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			userID: 1,
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					GetUserTickets(
						gomock.Any(),
						&tickets.GetUserTicketsIn{
							UserID: 1,
							Pagination: &tickets.Pagination{
								Limit:  pointers.New[uint64](1),
								Offset: pointers.New[uint64](1),
							},
							Filters: &tickets.TicketsFilters{
								Search:              pointers.New("Ticket"),
								PriceCeil:           pointers.New[float32](1000),
								PriceFloor:          pointers.New[float32](10),
								QuantityFloor:       pointers.New[uint32](1),
								CategoryIDs:         []uint32{1},
								TagIDs:              []uint32{1},
								CreatedAtOrderByAsc: pointers.New(true),
							},
						},
					).
					Return(
						&tickets.GetTicketsOut{
							Tickets: []*tickets.GetTicketOut{
								{
									ID:          1,
									UserID:      1,
									CategoryID:  2,
									Name:        "Ticket1",
									Description: "Desc1",
									Price:       pointers.New[float32](100),
									Quantity:    5,
									CreatedAt:   timestamppb.New(now),
									UpdatedAt:   timestamppb.New(now),
								},
							},
						},
						nil,
					).
					Times(1)
			},
			expectedTickets: []entities.RawTicket{
				{
					ID:          1,
					UserID:      1,
					CategoryID:  2,
					Name:        "Ticket1",
					Description: "Desc1",
					Price:       pointers.New[float32](100),
					Quantity:    5,
					CreatedAt:   now,
					UpdatedAt:   now,
					Attachments: make([]entities.TicketAttachment, 0),
				},
			},
			errorExpected: false,
		},
		{
			name: "error",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			filters: &entities.TicketsFilters{
				Search:              pointers.New("Ticket"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryIDs:         []uint32{1},
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			userID: 1,
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					GetUserTickets(
						gomock.Any(),
						&tickets.GetUserTicketsIn{
							UserID: 1,
							Pagination: &tickets.Pagination{
								Limit:  pointers.New[uint64](1),
								Offset: pointers.New[uint64](1),
							},
							Filters: &tickets.TicketsFilters{
								Search:              pointers.New("Ticket"),
								PriceCeil:           pointers.New[float32](1000),
								PriceFloor:          pointers.New[float32](10),
								QuantityFloor:       pointers.New[uint32](1),
								CategoryIDs:         []uint32{1},
								TagIDs:              []uint32{1},
								CreatedAtOrderByAsc: pointers.New(true),
							},
						},
					).
					Return(nil, errors.New("fetch failed")).
					Times(1)
			},
			expectedTickets: nil,
			errorExpected:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsClient)
			}

			ticketsList, err := repo.GetUserTickets(context.Background(), tc.userID, tc.pagination, tc.filters)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, ticketsList)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedTickets, ticketsList)
			}
		})
	}
}

func TestTicketsRepository_RespondToTicket(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsClient := mockclients.NewMockTicketsClient(ctrl)
	repo := NewTicketsRepository(ticketsClient)

	testCases := []struct {
		name              string
		respondData       entities.RespondToTicketDTO
		setupMocks        func(ticketsClient *mockclients.MockTicketsClient)
		expectedRespondID uint64
		errorExpected     bool
	}{
		{
			name: "success",
			respondData: entities.RespondToTicketDTO{
				UserID:   1,
				TicketID: 1,
				Price:    200,
				Comment:  pointers.New("Test Comment"),
			},
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					RespondToTicket(
						gomock.Any(),
						&tickets.RespondToTicketIn{
							UserID:   1,
							TicketID: 1,
							Price:    200,
							Comment:  pointers.New("Test Comment"),
						},
					).
					Return(&tickets.RespondToTicketOut{RespondID: 1}, nil).
					Times(1)
			},
			expectedRespondID: 1,
			errorExpected:     false,
		},
		{
			name: "error",
			respondData: entities.RespondToTicketDTO{
				UserID:   1,
				TicketID: 1,
			},
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					RespondToTicket(
						gomock.Any(),
						&tickets.RespondToTicketIn{
							UserID:   1,
							TicketID: 1,
						},
					).
					Return(nil, errors.New("respond failed")).
					Times(1)
			},
			expectedRespondID: 0,
			errorExpected:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsClient)
			}

			respondID, err := repo.RespondToTicket(context.Background(), tc.respondData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedRespondID, respondID)
		})
	}
}

func TestTicketsRepository_GetRespondByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsClient := mockclients.NewMockTicketsClient(ctrl)
	repo := NewTicketsRepository(ticketsClient)

	now := time.Now().UTC().Truncate(time.Second)

	testCases := []struct {
		name            string
		id              uint64
		setupMocks      func(ticketsClient *mockclients.MockTicketsClient)
		expectedRespond *entities.Respond
		errorExpected   bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					GetRespond(
						gomock.Any(),
						&tickets.GetRespondIn{ID: 1},
					).
					Return(&tickets.GetRespondOut{
						ID:        1,
						MasterID:  2,
						TicketID:  1,
						Price:     200,
						Comment:   pointers.New("Test Comment"),
						CreatedAt: timestamppb.New(now),
						UpdatedAt: timestamppb.New(now),
					}, nil).
					Times(1)
			},
			expectedRespond: &entities.Respond{
				ID:        1,
				MasterID:  2,
				TicketID:  1,
				Price:     200,
				Comment:   pointers.New("Test Comment"),
				CreatedAt: now,
				UpdatedAt: now,
			},
			errorExpected: false,
		},
		{
			name: "error",
			id:   1,
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					GetRespond(
						gomock.Any(),
						&tickets.GetRespondIn{ID: 1},
					).
					Return(nil, errors.New("respond not found")).
					Times(1)
			},
			expectedRespond: nil,
			errorExpected:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsClient)
			}

			respond, err := repo.GetRespondByID(context.Background(), tc.id)
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

func TestTicketsRepository_GetTicketResponds(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsClient := mockclients.NewMockTicketsClient(ctrl)
	repo := NewTicketsRepository(ticketsClient)

	now := time.Now().UTC().Truncate(time.Second)

	testCases := []struct {
		name             string
		ticketID         uint64
		setupMocks       func(ticketsClient *mockclients.MockTicketsClient)
		expectedResponds []entities.Respond
		errorExpected    bool
	}{
		{
			name:     "success",
			ticketID: 1,
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					GetTicketResponds(
						gomock.Any(),
						&tickets.GetTicketRespondsIn{TicketID: 1},
					).
					Return(&tickets.GetRespondsOut{
						Responds: []*tickets.GetRespondOut{
							{
								ID:        1,
								MasterID:  2,
								TicketID:  1,
								Price:     200,
								Comment:   pointers.New("Test Comment"),
								CreatedAt: timestamppb.New(now),
								UpdatedAt: timestamppb.New(now),
							},
						},
					},
						nil,
					).
					Times(1)
			},
			expectedResponds: []entities.Respond{
				{
					ID:        1,
					MasterID:  2,
					TicketID:  1,
					Price:     200,
					Comment:   pointers.New("Test Comment"),
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			errorExpected: false,
		},
		{
			name:     "error",
			ticketID: 1,
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					GetTicketResponds(
						gomock.Any(),
						&tickets.GetTicketRespondsIn{TicketID: 1},
					).
					Return(nil, errors.New("fetch failed")).
					Times(1)
			},
			expectedResponds: nil,
			errorExpected:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsClient)
			}

			responds, err := repo.GetTicketResponds(context.Background(), tc.ticketID)
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

func TestTicketsRepository_GetUserResponds(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsClient := mockclients.NewMockTicketsClient(ctrl)
	repo := NewTicketsRepository(ticketsClient)

	now := time.Now().UTC().Truncate(time.Second)

	testCases := []struct {
		name             string
		userID           uint64
		setupMocks       func(ticketsClient *mockclients.MockTicketsClient)
		expectedResponds []entities.Respond
		errorExpected    bool
	}{
		{
			name:   "success",
			userID: 1,
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					GetUserResponds(
						gomock.Any(),
						&tickets.GetUserRespondsIn{UserID: 1},
					).
					Return(&tickets.GetRespondsOut{
						Responds: []*tickets.GetRespondOut{
							{
								ID:        1,
								MasterID:  1,
								TicketID:  2,
								Price:     200,
								Comment:   pointers.New("Test Comment"),
								CreatedAt: timestamppb.New(now),
								UpdatedAt: timestamppb.New(now),
							},
						},
					},
						nil,
					).
					Times(1)
			},
			expectedResponds: []entities.Respond{
				{
					ID:        1,
					MasterID:  1,
					TicketID:  2,
					Price:     200,
					Comment:   pointers.New("Test Comment"),
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			errorExpected: false,
		},
		{
			name:   "error",
			userID: 1,
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					GetUserResponds(
						gomock.Any(),
						&tickets.GetUserRespondsIn{UserID: 1},
					).
					Return(nil, errors.New("fetch failed")).
					Times(1)
			},
			expectedResponds: nil,
			errorExpected:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsClient)
			}

			responds, err := repo.GetUserResponds(context.Background(), tc.userID)
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

func TestTicketsRepository_UpdateRespond(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsClient := mockclients.NewMockTicketsClient(ctrl)
	repo := NewTicketsRepository(ticketsClient)

	testCases := []struct {
		name          string
		respondData   entities.UpdateRespondDTO
		setupMocks    func(ticketsClient *mockclients.MockTicketsClient)
		errorExpected bool
	}{
		{
			name: "success",
			respondData: entities.UpdateRespondDTO{
				ID:      1,
				Price:   pointers.New[float32](100),
				Comment: pointers.New("Test Comment"),
			},
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					UpdateRespond(
						gomock.Any(),
						&tickets.UpdateRespondIn{
							ID:      1,
							Price:   pointers.New[float32](100),
							Comment: pointers.New("Test Comment"),
						},
					).
					Return(&emptypb.Empty{}, nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "error",
			respondData: entities.UpdateRespondDTO{
				ID: 1,
			},
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					UpdateRespond(
						gomock.Any(),
						&tickets.UpdateRespondIn{
							ID: 1,
						},
					).
					Return(nil, errors.New("update failed")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsClient)
			}

			err := repo.UpdateRespond(context.Background(), tc.respondData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTicketsRepository_DeleteRespond(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsClient := mockclients.NewMockTicketsClient(ctrl)
	repo := NewTicketsRepository(ticketsClient)

	testCases := []struct {
		name          string
		id            uint64
		setupMocks    func(ticketsClient *mockclients.MockTicketsClient)
		errorExpected bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					DeleteRespond(
						gomock.Any(),
						&tickets.DeleteRespondIn{ID: 1},
					).
					Return(&emptypb.Empty{}, nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "error",
			id:   1,
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					DeleteRespond(
						gomock.Any(),
						&tickets.DeleteRespondIn{ID: 1},
					).
					Return(nil, errors.New("delete failed")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsClient)
			}

			err := repo.DeleteRespond(context.Background(), tc.id)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTicketsRepository_UpdateTicket(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsClient := mockclients.NewMockTicketsClient(ctrl)
	repo := NewTicketsRepository(ticketsClient)

	testCases := []struct {
		name          string
		ticketData    entities.UpdateTicketDTO
		setupMocks    func(ticketsClient *mockclients.MockTicketsClient)
		errorExpected bool
	}{
		{
			name: "success",
			ticketData: entities.UpdateTicketDTO{
				ID:          1,
				Name:        pointers.New("Updated Ticket"),
				Description: pointers.New("Updated Description"),
				CategoryID:  pointers.New[uint32](2),
				Price:       pointers.New[float32](100),
				Quantity:    pointers.New[uint32](10),
				TagIDs:      []uint32{1, 3},
				Attachments: []string{"updated_attachment"},
			},
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					UpdateTicket(
						gomock.Any(),
						&tickets.UpdateTicketIn{
							ID:          1,
							Name:        pointers.New("Updated Ticket"),
							Description: pointers.New("Updated Description"),
							CategoryID:  pointers.New[uint32](2),
							Price:       pointers.New[float32](100),
							Quantity:    pointers.New[uint32](10),
							TagIDs:      []uint32{1, 3},
							Attachments: []string{"updated_attachment"},
						},
					).
					Return(&emptypb.Empty{}, nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "error",
			ticketData: entities.UpdateTicketDTO{
				ID: 1,
			},
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					UpdateTicket(
						gomock.Any(),
						&tickets.UpdateTicketIn{
							ID: 1,
						},
					).
					Return(nil, errors.New("update failed")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsClient)
			}

			err := repo.UpdateTicket(context.Background(), tc.ticketData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTicketsRepository_DeleteTicket(t *testing.T) {
	ctrl := gomock.NewController(t)
	ticketsClient := mockclients.NewMockTicketsClient(ctrl)
	repo := NewTicketsRepository(ticketsClient)

	testCases := []struct {
		name          string
		id            uint64
		setupMocks    func(ticketsClient *mockclients.MockTicketsClient)
		errorExpected bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					DeleteTicket(
						gomock.Any(),
						&tickets.DeleteTicketIn{ID: 1},
					).
					Return(&emptypb.Empty{}, nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "error",
			id:   1,
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					DeleteTicket(
						gomock.Any(),
						&tickets.DeleteTicketIn{ID: 1},
					).
					Return(nil, errors.New("delete failed")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsClient)
			}

			err := repo.DeleteTicket(context.Background(), tc.id)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTicketsRepository_CountTickets(t *testing.T) {
	testCases := []struct {
		name          string
		filters       *entities.TicketsFilters
		setupMocks    func(ticketsClient *mockclients.MockTicketsClient)
		expected      uint64
		errorExpected bool
	}{
		{
			name: "success",
			filters: &entities.TicketsFilters{
				Search:              pointers.New("Ticket"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryIDs:         []uint32{1},
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					CountTickets(
						gomock.Any(),
						&tickets.CountTicketsIn{
							Filters: &tickets.TicketsFilters{
								Search:              pointers.New("Ticket"),
								PriceCeil:           pointers.New[float32](1000),
								PriceFloor:          pointers.New[float32](10),
								QuantityFloor:       pointers.New[uint32](1),
								CategoryIDs:         []uint32{1},
								TagIDs:              []uint32{1},
								CreatedAtOrderByAsc: pointers.New(true),
							},
						},
					).
					Return(
						&tickets.CountOut{Count: 1},
						nil,
					).
					Times(1)
			},
			expected: 1,
		},
		{
			name: "error",
			filters: &entities.TicketsFilters{
				Search:              pointers.New("Ticket"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryIDs:         []uint32{1},
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					CountTickets(
						gomock.Any(),
						&tickets.CountTicketsIn{
							Filters: &tickets.TicketsFilters{
								Search:              pointers.New("Ticket"),
								PriceCeil:           pointers.New[float32](1000),
								PriceFloor:          pointers.New[float32](10),
								QuantityFloor:       pointers.New[uint32](1),
								CategoryIDs:         []uint32{1},
								TagIDs:              []uint32{1},
								CreatedAtOrderByAsc: pointers.New(true),
							},
						},
					).
					Return(
						&tickets.CountOut{Count: 0},
						errors.New("error"),
					).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ticketsClient := mockclients.NewMockTicketsClient(ctrl)
	repo := NewTicketsRepository(ticketsClient)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsClient)
			}

			actual, err := repo.CountTickets(context.Background(), tc.filters)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestTicketsRepository_CountUserTickets(t *testing.T) {
	testCases := []struct {
		name          string
		userID        uint64
		filters       *entities.TicketsFilters
		setupMocks    func(ticketsClient *mockclients.MockTicketsClient)
		expected      uint64
		errorExpected bool
	}{
		{
			name:   "success",
			userID: 1,
			filters: &entities.TicketsFilters{
				Search:              pointers.New("Ticket"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryIDs:         []uint32{1},
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					CountUserTickets(
						gomock.Any(),
						&tickets.CountUserTicketsIn{
							UserID: 1,
							Filters: &tickets.TicketsFilters{
								Search:              pointers.New("Ticket"),
								PriceCeil:           pointers.New[float32](1000),
								PriceFloor:          pointers.New[float32](10),
								QuantityFloor:       pointers.New[uint32](1),
								CategoryIDs:         []uint32{1},
								TagIDs:              []uint32{1},
								CreatedAtOrderByAsc: pointers.New(true),
							},
						},
					).
					Return(
						&tickets.CountOut{Count: 1},
						nil,
					).
					Times(1)
			},
			expected: 1,
		},
		{
			name:   "error",
			userID: 1,
			filters: &entities.TicketsFilters{
				Search:              pointers.New("Ticket"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryIDs:         []uint32{1},
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			setupMocks: func(ticketsClient *mockclients.MockTicketsClient) {
				ticketsClient.
					EXPECT().
					CountUserTickets(
						gomock.Any(),
						&tickets.CountUserTicketsIn{
							UserID: 1,
							Filters: &tickets.TicketsFilters{
								Search:              pointers.New("Ticket"),
								PriceCeil:           pointers.New[float32](1000),
								PriceFloor:          pointers.New[float32](10),
								QuantityFloor:       pointers.New[uint32](1),
								CategoryIDs:         []uint32{1},
								TagIDs:              []uint32{1},
								CreatedAtOrderByAsc: pointers.New(true),
							},
						},
					).
					Return(
						&tickets.CountOut{Count: 0},
						errors.New("error"),
					).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ticketsClient := mockclients.NewMockTicketsClient(ctrl)
	repo := NewTicketsRepository(ticketsClient)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsClient)
			}

			actual, err := repo.CountUserTickets(context.Background(), tc.userID, tc.filters)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}
