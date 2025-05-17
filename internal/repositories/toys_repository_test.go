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

	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"
	"github.com/DKhorkov/libs/pointers"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	mockclients "github.com/DKhorkov/hmtm-bff/mocks/clients"
)

func TestToysRepository_AddToy(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysClient := mockclients.NewMockToysClient(ctrl)
	repo := NewToysRepository(toysClient)

	testCases := []struct {
		name          string
		toyData       entities.AddToyDTO
		setupMocks    func(toysClient *mockclients.MockToysClient)
		expectedToyID uint64
		errorExpected bool
	}{
		{
			name: "success",
			toyData: entities.AddToyDTO{
				UserID:      1,
				CategoryID:  2,
				Name:        "Test Toy",
				Description: "Test Description",
				Price:       100,
				Quantity:    5,
				TagIDs:      []uint32{1, 2},
				Attachments: []string{"attachment1"},
			},
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					AddToy(
						gomock.Any(),
						&toys.AddToyIn{
							UserID:      1,
							CategoryID:  2,
							Name:        "Test Toy",
							Description: "Test Description",
							Price:       100,
							Quantity:    5,
							TagIDs:      []uint32{1, 2},
							Attachments: []string{"attachment1"},
						},
					).
					Return(&toys.AddToyOut{ToyID: 1}, nil).
					Times(1)
			},
			expectedToyID: 1,
			errorExpected: false,
		},
		{
			name: "error",
			toyData: entities.AddToyDTO{
				UserID: 1,
				Name:   "Test Toy",
			},
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					AddToy(
						gomock.Any(),
						&toys.AddToyIn{
							UserID: 1,
							Name:   "Test Toy",
						},
					).
					Return(nil, errors.New("add failed")).
					Times(1)
			},
			expectedToyID: 0,
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysClient)
			}

			toyID, err := repo.AddToy(context.Background(), tc.toyData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedToyID, toyID)
		})
	}
}

func TestToysRepository_GetToys(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	testCases := []struct {
		name          string
		pagination    *entities.Pagination
		filters       *entities.ToysFilters
		setupMocks    func(toysClient *mockclients.MockToysClient)
		expectedToys  []entities.Toy
		errorExpected bool
	}{
		{
			name: "success",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			filters: &entities.ToysFilters{
				Search:              pointers.New("Toy"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryID:          pointers.New[uint32](1),
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetToys(
						gomock.Any(),
						&toys.GetToysIn{
							Pagination: &toys.Pagination{
								Limit:  pointers.New[uint64](1),
								Offset: pointers.New[uint64](1),
							},
							Filters: &toys.ToysFilters{
								Search:              pointers.New("Toy"),
								PriceCeil:           pointers.New[float32](1000),
								PriceFloor:          pointers.New[float32](10),
								QuantityFloor:       pointers.New[uint32](1),
								CategoryID:          pointers.New[uint32](1),
								TagIDs:              []uint32{1},
								CreatedAtOrderByAsc: pointers.New(true),
							},
						},
					).
					Return(&toys.GetToysOut{
						Toys: []*toys.GetToyOut{
							{
								ID:          1,
								MasterID:    1,
								CategoryID:  2,
								Name:        "Toy1",
								Description: "Desc1",
								Price:       100,
								Quantity:    5,
								CreatedAt:   timestamppb.New(now),
								UpdatedAt:   timestamppb.New(now),
							},
						},
					}, nil).
					Times(1)
			},
			expectedToys: []entities.Toy{
				{
					ID:          1,
					MasterID:    1,
					CategoryID:  2,
					Name:        "Toy1",
					Description: "Desc1",
					Price:       100,
					Quantity:    5,
					CreatedAt:   now,
					UpdatedAt:   now,
					Tags:        make([]entities.Tag, 0),
					Attachments: make([]entities.ToyAttachment, 0),
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
			filters: &entities.ToysFilters{
				Search:              pointers.New("Toy"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryID:          pointers.New[uint32](1),
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetToys(
						gomock.Any(),
						&toys.GetToysIn{
							Pagination: &toys.Pagination{
								Limit:  pointers.New[uint64](1),
								Offset: pointers.New[uint64](1),
							},
							Filters: &toys.ToysFilters{
								Search:              pointers.New("Toy"),
								PriceCeil:           pointers.New[float32](1000),
								PriceFloor:          pointers.New[float32](10),
								QuantityFloor:       pointers.New[uint32](1),
								CategoryID:          pointers.New[uint32](1),
								TagIDs:              []uint32{1},
								CreatedAtOrderByAsc: pointers.New(true),
							},
						},
					).
					Return(nil, errors.New("fetch failed")).
					Times(1)
			},
			expectedToys:  nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	toysClient := mockclients.NewMockToysClient(ctrl)
	repo := NewToysRepository(toysClient)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysClient)
			}

			toysList, err := repo.GetToys(context.Background(), tc.pagination, tc.filters)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, toysList)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedToys, toysList)
			}
		})
	}
}

func TestToysRepository_CountToys(t *testing.T) {
	testCases := []struct {
		name          string
		filters       *entities.ToysFilters
		setupMocks    func(toysClient *mockclients.MockToysClient)
		expected      uint64
		errorExpected bool
	}{
		{
			name: "success",
			filters: &entities.ToysFilters{
				Search:              pointers.New("Toy"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryID:          pointers.New[uint32](1),
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					CountToys(
						gomock.Any(),
						&toys.CountToysIn{
							Filters: &toys.ToysFilters{
								Search:              pointers.New("Toy"),
								PriceCeil:           pointers.New[float32](1000),
								PriceFloor:          pointers.New[float32](10),
								QuantityFloor:       pointers.New[uint32](1),
								CategoryID:          pointers.New[uint32](1),
								TagIDs:              []uint32{1},
								CreatedAtOrderByAsc: pointers.New(true),
							},
						},
					).
					Return(
						&toys.CountOut{Count: 1},
						nil,
					).
					Times(1)
			},
			expected: 1,
		},
		{
			name: "error",
			filters: &entities.ToysFilters{
				Search:              pointers.New("Toy"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryID:          pointers.New[uint32](1),
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					CountToys(
						gomock.Any(),
						&toys.CountToysIn{
							Filters: &toys.ToysFilters{
								Search:              pointers.New("Toy"),
								PriceCeil:           pointers.New[float32](1000),
								PriceFloor:          pointers.New[float32](10),
								QuantityFloor:       pointers.New[uint32](1),
								CategoryID:          pointers.New[uint32](1),
								TagIDs:              []uint32{1},
								CreatedAtOrderByAsc: pointers.New(true),
							},
						},
					).
					Return(
						&toys.CountOut{Count: 0},
						errors.New("error"),
					).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	toysClient := mockclients.NewMockToysClient(ctrl)
	repo := NewToysRepository(toysClient)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysClient)
			}

			actual, err := repo.CountToys(context.Background(), tc.filters)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestToysRepository_GetMasterToys(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysClient := mockclients.NewMockToysClient(ctrl)
	repo := NewToysRepository(toysClient)

	now := time.Now().UTC().Truncate(time.Second)

	testCases := []struct {
		name          string
		pagination    *entities.Pagination
		filters       *entities.ToysFilters
		masterID      uint64
		setupMocks    func(toysClient *mockclients.MockToysClient)
		expectedToys  []entities.Toy
		errorExpected bool
	}{
		{
			name: "success",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			filters: &entities.ToysFilters{
				Search:              pointers.New("Toy"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryID:          pointers.New[uint32](1),
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			masterID: 1,
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetMasterToys(
						gomock.Any(),
						&toys.GetMasterToysIn{
							MasterID: 1,
							Pagination: &toys.Pagination{
								Limit:  pointers.New[uint64](1),
								Offset: pointers.New[uint64](1),
							},
							Filters: &toys.ToysFilters{
								Search:              pointers.New("Toy"),
								PriceCeil:           pointers.New[float32](1000),
								PriceFloor:          pointers.New[float32](10),
								QuantityFloor:       pointers.New[uint32](1),
								CategoryID:          pointers.New[uint32](1),
								TagIDs:              []uint32{1},
								CreatedAtOrderByAsc: pointers.New(true),
							},
						},
					).
					Return(&toys.GetToysOut{
						Toys: []*toys.GetToyOut{
							{
								ID:          1,
								MasterID:    1,
								CategoryID:  2,
								Name:        "Toy1",
								Description: "Desc1",
								Price:       100,
								Quantity:    5,
								CreatedAt:   timestamppb.New(now),
								UpdatedAt:   timestamppb.New(now),
							},
						},
					}, nil).
					Times(1)
			},
			expectedToys: []entities.Toy{
				{
					ID:          1,
					MasterID:    1,
					CategoryID:  2,
					Name:        "Toy1",
					Description: "Desc1",
					Price:       100,
					Quantity:    5,
					CreatedAt:   now,
					UpdatedAt:   now,
					Tags:        make([]entities.Tag, 0),
					Attachments: make([]entities.ToyAttachment, 0),
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
			filters: &entities.ToysFilters{
				Search:              pointers.New("Toy"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryID:          pointers.New[uint32](1),
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			masterID: 1,
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetMasterToys(
						gomock.Any(),
						&toys.GetMasterToysIn{
							MasterID: 1,
							Pagination: &toys.Pagination{
								Limit:  pointers.New[uint64](1),
								Offset: pointers.New[uint64](1),
							},
							Filters: &toys.ToysFilters{
								Search:              pointers.New("Toy"),
								PriceCeil:           pointers.New[float32](1000),
								PriceFloor:          pointers.New[float32](10),
								QuantityFloor:       pointers.New[uint32](1),
								CategoryID:          pointers.New[uint32](1),
								TagIDs:              []uint32{1},
								CreatedAtOrderByAsc: pointers.New(true),
							},
						},
					).
					Return(nil, errors.New("fetch failed")).
					Times(1)
			},
			expectedToys:  nil,
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysClient)
			}

			toysList, err := repo.GetMasterToys(context.Background(), tc.masterID, tc.pagination, tc.filters)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, toysList)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedToys, toysList)
			}
		})
	}
}

func TestToysRepository_GetUserToys(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysClient := mockclients.NewMockToysClient(ctrl)
	repo := NewToysRepository(toysClient)

	now := time.Now().UTC().Truncate(time.Second)

	testCases := []struct {
		name          string
		pagination    *entities.Pagination
		filters       *entities.ToysFilters
		userID        uint64
		setupMocks    func(toysClient *mockclients.MockToysClient)
		expectedToys  []entities.Toy
		errorExpected bool
	}{
		{
			name: "success",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			filters: &entities.ToysFilters{
				Search:              pointers.New("Toy"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryID:          pointers.New[uint32](1),
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			userID: 1,
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetUserToys(
						gomock.Any(),
						&toys.GetUserToysIn{
							UserID: 1,
							Pagination: &toys.Pagination{
								Limit:  pointers.New[uint64](1),
								Offset: pointers.New[uint64](1),
							},
							Filters: &toys.ToysFilters{
								Search:              pointers.New("Toy"),
								PriceCeil:           pointers.New[float32](1000),
								PriceFloor:          pointers.New[float32](10),
								QuantityFloor:       pointers.New[uint32](1),
								CategoryID:          pointers.New[uint32](1),
								TagIDs:              []uint32{1},
								CreatedAtOrderByAsc: pointers.New(true),
							},
						},
					).
					Return(&toys.GetToysOut{
						Toys: []*toys.GetToyOut{
							{
								ID:          1,
								MasterID:    1,
								CategoryID:  2,
								Name:        "Toy1",
								Description: "Desc1",
								Price:       100,
								Quantity:    5,
								CreatedAt:   timestamppb.New(now),
								UpdatedAt:   timestamppb.New(now),
							},
						},
					}, nil).
					Times(1)
			},
			expectedToys: []entities.Toy{
				{
					ID:          1,
					MasterID:    1,
					CategoryID:  2,
					Name:        "Toy1",
					Description: "Desc1",
					Price:       100,
					Quantity:    5,
					CreatedAt:   now,
					UpdatedAt:   now,
					Tags:        make([]entities.Tag, 0),
					Attachments: make([]entities.ToyAttachment, 0),
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
			filters: &entities.ToysFilters{
				Search:              pointers.New("Toy"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryID:          pointers.New[uint32](1),
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			userID: 1,
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetUserToys(
						gomock.Any(),
						&toys.GetUserToysIn{
							UserID: 1,
							Pagination: &toys.Pagination{
								Limit:  pointers.New[uint64](1),
								Offset: pointers.New[uint64](1),
							},
							Filters: &toys.ToysFilters{
								Search:              pointers.New("Toy"),
								PriceCeil:           pointers.New[float32](1000),
								PriceFloor:          pointers.New[float32](10),
								QuantityFloor:       pointers.New[uint32](1),
								CategoryID:          pointers.New[uint32](1),
								TagIDs:              []uint32{1},
								CreatedAtOrderByAsc: pointers.New(true),
							},
						},
					).
					Return(nil, errors.New("fetch failed")).
					Times(1)
			},
			expectedToys:  nil,
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysClient)
			}

			toysList, err := repo.GetUserToys(context.Background(), tc.userID, tc.pagination, tc.filters)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, toysList)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedToys, toysList)
			}
		})
	}
}

func TestToysRepository_GetToyByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysClient := mockclients.NewMockToysClient(ctrl)
	repo := NewToysRepository(toysClient)

	now := time.Now().UTC().Truncate(time.Second)

	testCases := []struct {
		name          string
		id            uint64
		setupMocks    func(toysClient *mockclients.MockToysClient)
		expectedToy   *entities.Toy
		errorExpected bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetToy(
						gomock.Any(),
						&toys.GetToyIn{ID: 1},
					).
					Return(&toys.GetToyOut{
						ID:          1,
						MasterID:    1,
						CategoryID:  2,
						Name:        "Test Toy",
						Description: "Desc",
						Price:       100,
						Quantity:    5,
						CreatedAt:   timestamppb.New(now),
						UpdatedAt:   timestamppb.New(now),
						Tags: []*toys.GetTagOut{
							{
								ID:        1,
								Name:      "Tag1",
								CreatedAt: timestamppb.New(now),
								UpdatedAt: timestamppb.New(now),
							},
						},
						Attachments: []*toys.Attachment{
							{
								ID:        1,
								ToyID:     1,
								Link:      "test",
								CreatedAt: timestamppb.New(now),
								UpdatedAt: timestamppb.New(now),
							},
						},
					}, nil).
					Times(1)
			},
			expectedToy: &entities.Toy{
				ID:          1,
				MasterID:    1,
				CategoryID:  2,
				Name:        "Test Toy",
				Description: "Desc",
				Price:       100,
				Quantity:    5,
				CreatedAt:   now,
				UpdatedAt:   now,
				Tags: []entities.Tag{
					{
						ID:   1,
						Name: "Tag1",
					},
				},
				Attachments: []entities.ToyAttachment{
					{
						ID:        1,
						ToyID:     1,
						Link:      "test",
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
			errorExpected: false,
		},
		{
			name: "error",
			id:   1,
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetToy(
						gomock.Any(),
						&toys.GetToyIn{ID: 1},
					).
					Return(nil, errors.New("not found")).
					Times(1)
			},
			expectedToy:   nil,
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysClient)
			}

			toy, err := repo.GetToyByID(context.Background(), tc.id)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, toy)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedToy, toy)
			}
		})
	}
}

func TestToysRepository_GetMasters(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysClient := mockclients.NewMockToysClient(ctrl)
	repo := NewToysRepository(toysClient)

	now := time.Now().UTC().Truncate(time.Second)

	testCases := []struct {
		name            string
		pagination      *entities.Pagination
		setupMocks      func(toysClient *mockclients.MockToysClient)
		expectedMasters []entities.Master
		errorExpected   bool
	}{
		{
			name: "success",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetMasters(
						gomock.Any(),
						&toys.GetMastersIn{
							Pagination: &toys.Pagination{
								Limit:  pointers.New[uint64](1),
								Offset: pointers.New[uint64](1),
							},
						},
					).
					Return(&toys.GetMastersOut{
						Masters: []*toys.GetMasterOut{
							{
								ID:        1,
								UserID:    1,
								Info:      pointers.New("Master Info"),
								CreatedAt: timestamppb.New(now),
								UpdatedAt: timestamppb.New(now),
							},
						},
					},
						nil,
					).
					Times(1)
			},
			expectedMasters: []entities.Master{
				{
					ID:        1,
					UserID:    1,
					Info:      pointers.New("Master Info"),
					CreatedAt: now,
					UpdatedAt: now,
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
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetMasters(
						gomock.Any(),
						&toys.GetMastersIn{
							Pagination: &toys.Pagination{
								Limit:  pointers.New[uint64](1),
								Offset: pointers.New[uint64](1),
							},
						},
					).
					Return(nil, errors.New("fetch failed")).
					Times(1)
			},
			expectedMasters: nil,
			errorExpected:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysClient)
			}

			masters, err := repo.GetMasters(context.Background(), tc.pagination)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, masters)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedMasters, masters)
			}
		})
	}
}

func TestToysRepository_GetMasterByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysClient := mockclients.NewMockToysClient(ctrl)
	repo := NewToysRepository(toysClient)

	now := time.Now().UTC().Truncate(time.Second)

	testCases := []struct {
		name           string
		id             uint64
		setupMocks     func(toysClient *mockclients.MockToysClient)
		expectedMaster *entities.Master
		errorExpected  bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetMaster(
						gomock.Any(),
						&toys.GetMasterIn{ID: 1},
					).
					Return(
						&toys.GetMasterOut{
							ID:        1,
							UserID:    1,
							Info:      pointers.New("Master Info"),
							CreatedAt: timestamppb.New(now),
							UpdatedAt: timestamppb.New(now),
						},
						nil,
					).
					Times(1)
			},
			expectedMaster: &entities.Master{
				ID:        1,
				UserID:    1,
				Info:      pointers.New("Master Info"),
				CreatedAt: now,
				UpdatedAt: now,
			},
			errorExpected: false,
		},
		{
			name: "error",
			id:   1,
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetMaster(
						gomock.Any(),
						&toys.GetMasterIn{ID: 1},
					).
					Return(nil, errors.New("not found")).
					Times(1)
			},
			expectedMaster: nil,
			errorExpected:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysClient)
			}

			master, err := repo.GetMasterByID(context.Background(), tc.id)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, master)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedMaster, master)
			}
		})
	}
}

func TestToysRepository_RegisterMaster(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysClient := mockclients.NewMockToysClient(ctrl)
	repo := NewToysRepository(toysClient)

	testCases := []struct {
		name             string
		masterData       entities.RegisterMasterDTO
		setupMocks       func(toysClient *mockclients.MockToysClient)
		expectedMasterID uint64
		errorExpected    bool
	}{
		{
			name: "success",
			masterData: entities.RegisterMasterDTO{
				UserID: 1,
				Info:   pointers.New("Master Info"),
			},
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					RegisterMaster(
						gomock.Any(),
						&toys.RegisterMasterIn{
							UserID: 1,
							Info:   pointers.New("Master Info"),
						},
					).
					Return(&toys.RegisterMasterOut{MasterID: 1}, nil).
					Times(1)
			},
			expectedMasterID: 1,
			errorExpected:    false,
		},
		{
			name: "error",
			masterData: entities.RegisterMasterDTO{
				UserID: 1,
			},
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					RegisterMaster(
						gomock.Any(),
						&toys.RegisterMasterIn{
							UserID: 1,
						},
					).
					Return(nil, errors.New("registration failed")).
					Times(1)
			},
			expectedMasterID: 0,
			errorExpected:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysClient)
			}

			masterID, err := repo.RegisterMaster(context.Background(), tc.masterData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedMasterID, masterID)
		})
	}
}

func TestToysRepository_GetAllCategories(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysClient := mockclients.NewMockToysClient(ctrl)
	repo := NewToysRepository(toysClient)

	testCases := []struct {
		name               string
		setupMocks         func(toysClient *mockclients.MockToysClient)
		expectedCategories []entities.Category
		errorExpected      bool
	}{
		{
			name: "success",
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetCategories(
						gomock.Any(),
						&emptypb.Empty{},
					).
					Return(&toys.GetCategoriesOut{
						Categories: []*toys.GetCategoryOut{
							{ID: 1, Name: "Category1"},
						},
					}, nil).
					Times(1)
			},
			expectedCategories: []entities.Category{
				{ID: 1, Name: "Category1"},
			},
			errorExpected: false,
		},
		{
			name: "error",
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetCategories(
						gomock.Any(),
						&emptypb.Empty{},
					).
					Return(nil, errors.New("fetch failed")).
					Times(1)
			},
			expectedCategories: nil,
			errorExpected:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysClient)
			}

			categories, err := repo.GetAllCategories(context.Background())
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, categories)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedCategories, categories)
			}
		})
	}
}

func TestToysRepository_GetCategoryByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysClient := mockclients.NewMockToysClient(ctrl)
	repo := NewToysRepository(toysClient)

	testCases := []struct {
		name             string
		id               uint32
		setupMocks       func(toysClient *mockclients.MockToysClient)
		expectedCategory *entities.Category
		errorExpected    bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetCategory(
						gomock.Any(),
						&toys.GetCategoryIn{ID: 1},
					).
					Return(&toys.GetCategoryOut{
						ID:   1,
						Name: "Test Category",
					}, nil).
					Times(1)
			},
			expectedCategory: &entities.Category{
				ID:   1,
				Name: "Test Category",
			},
			errorExpected: false,
		},
		{
			name: "error",
			id:   1,
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetCategory(
						gomock.Any(),
						&toys.GetCategoryIn{ID: 1},
					).
					Return(nil, errors.New("not found")).
					Times(1)
			},
			expectedCategory: nil,
			errorExpected:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysClient)
			}

			category, err := repo.GetCategoryByID(context.Background(), tc.id)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, category)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedCategory, category)
			}
		})
	}
}

func TestToysRepository_GetAllTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysClient := mockclients.NewMockToysClient(ctrl)
	repo := NewToysRepository(toysClient)

	testCases := []struct {
		name          string
		setupMocks    func(toysClient *mockclients.MockToysClient)
		expectedTags  []entities.Tag
		errorExpected bool
	}{
		{
			name: "success",
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetTags(
						gomock.Any(),
						&emptypb.Empty{},
					).
					Return(&toys.GetTagsOut{
						Tags: []*toys.GetTagOut{
							{ID: 1, Name: "Tag1"},
						},
					}, nil).
					Times(1)
			},
			expectedTags: []entities.Tag{
				{ID: 1, Name: "Tag1"},
			},
			errorExpected: false,
		},
		{
			name: "error",
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetTags(
						gomock.Any(),
						&emptypb.Empty{},
					).
					Return(nil, errors.New("fetch failed")).
					Times(1)
			},
			expectedTags:  nil,
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysClient)
			}

			tags, err := repo.GetAllTags(context.Background())
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, tags)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedTags, tags)
			}
		})
	}
}

func TestToysRepository_GetTagByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysClient := mockclients.NewMockToysClient(ctrl)
	repo := NewToysRepository(toysClient)

	testCases := []struct {
		name          string
		id            uint32
		setupMocks    func(toysClient *mockclients.MockToysClient)
		expectedTag   *entities.Tag
		errorExpected bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetTag(
						gomock.Any(),
						&toys.GetTagIn{ID: 1},
					).
					Return(&toys.GetTagOut{
						ID:   1,
						Name: "Test Tag",
					}, nil).
					Times(1)
			},
			expectedTag: &entities.Tag{
				ID:   1,
				Name: "Test Tag",
			},
			errorExpected: false,
		},
		{
			name: "error",
			id:   1,
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetTag(
						gomock.Any(),
						&toys.GetTagIn{ID: 1},
					).
					Return(nil, errors.New("not found")).
					Times(1)
			},
			expectedTag:   nil,
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysClient)
			}

			tag, err := repo.GetTagByID(context.Background(), tc.id)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, tag)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedTag, tag)
			}
		})
	}
}

func TestToysRepository_CreateTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysClient := mockclients.NewMockToysClient(ctrl)
	repo := NewToysRepository(toysClient)

	testCases := []struct {
		name           string
		tagsData       []entities.CreateTagDTO
		setupMocks     func(toysClient *mockclients.MockToysClient)
		expectedTagIDs []uint32
		errorExpected  bool
	}{
		{
			name: "success",
			tagsData: []entities.CreateTagDTO{
				{Name: "Tag1"},
				{Name: "Tag2"},
			},
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					CreateTags(
						gomock.Any(),
						&toys.CreateTagsIn{
							Tags: []*toys.CreateTagIn{
								{Name: "Tag1"},
								{Name: "Tag2"},
							},
						},
					).
					Return(&toys.CreateTagsOut{
						Tags: []*toys.CreateTagOut{
							{ID: 1},
							{ID: 2},
						},
					}, nil).
					Times(1)
			},
			expectedTagIDs: []uint32{1, 2},
			errorExpected:  false,
		},
		{
			name: "error",
			tagsData: []entities.CreateTagDTO{
				{Name: "Tag1"},
			},
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					CreateTags(
						gomock.Any(),
						&toys.CreateTagsIn{
							Tags: []*toys.CreateTagIn{
								{Name: "Tag1"},
							},
						},
					).
					Return(nil, errors.New("creation failed")).
					Times(1)
			},
			expectedTagIDs: nil,
			errorExpected:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysClient)
			}

			tagIDs, err := repo.CreateTags(context.Background(), tc.tagsData)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, tagIDs)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedTagIDs, tagIDs)
			}
		})
	}
}

func TestToysRepository_UpdateToy(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysClient := mockclients.NewMockToysClient(ctrl)
	repo := NewToysRepository(toysClient)

	testCases := []struct {
		name          string
		toyData       entities.UpdateToyDTO
		setupMocks    func(toysClient *mockclients.MockToysClient)
		errorExpected bool
	}{
		{
			name: "success",
			toyData: entities.UpdateToyDTO{
				ID:          1,
				Name:        pointers.New("Updated Toy"),
				Description: pointers.New("Updated Desc"),
				CategoryID:  pointers.New[uint32](2),
				Price:       pointers.New[float32](150),
				Quantity:    pointers.New[uint32](10),
				TagIDs:      []uint32{1, 2},
				Attachments: []string{"updated_attachment"},
			},
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					UpdateToy(
						gomock.Any(),
						&toys.UpdateToyIn{
							ID:          1,
							Name:        pointers.New("Updated Toy"),
							Description: pointers.New("Updated Desc"),
							CategoryID:  pointers.New[uint32](2),
							Price:       pointers.New[float32](150),
							Quantity:    pointers.New[uint32](10),
							TagIDs:      []uint32{1, 2},
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
			toyData: entities.UpdateToyDTO{
				ID: 1,
			},
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					UpdateToy(
						gomock.Any(),
						&toys.UpdateToyIn{
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
				tc.setupMocks(toysClient)
			}

			err := repo.UpdateToy(context.Background(), tc.toyData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestToysRepository_DeleteToy(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysClient := mockclients.NewMockToysClient(ctrl)
	repo := NewToysRepository(toysClient)

	testCases := []struct {
		name          string
		id            uint64
		setupMocks    func(toysClient *mockclients.MockToysClient)
		errorExpected bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					DeleteToy(
						gomock.Any(),
						&toys.DeleteToyIn{ID: 1},
					).
					Return(&emptypb.Empty{}, nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "error",
			id:   1,
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					DeleteToy(
						gomock.Any(),
						&toys.DeleteToyIn{ID: 1},
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
				tc.setupMocks(toysClient)
			}

			err := repo.DeleteToy(context.Background(), tc.id)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestToysRepository_GetMasterByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysClient := mockclients.NewMockToysClient(ctrl)
	repo := NewToysRepository(toysClient)

	now := time.Now().UTC().Truncate(time.Second)

	testCases := []struct {
		name           string
		userID         uint64
		setupMocks     func(toysClient *mockclients.MockToysClient)
		expectedMaster *entities.Master
		errorExpected  bool
	}{
		{
			name:   "success",
			userID: 1,
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetMasterByUser(
						gomock.Any(),
						&toys.GetMasterByUserIn{UserID: 1},
					).
					Return(&toys.GetMasterOut{
						ID:        1,
						UserID:    1,
						Info:      pointers.New("Master Info"),
						CreatedAt: timestamppb.New(now),
						UpdatedAt: timestamppb.New(now),
					}, nil).
					Times(1)
			},
			expectedMaster: &entities.Master{
				ID:        1,
				UserID:    1,
				Info:      pointers.New("Master Info"),
				CreatedAt: now,
				UpdatedAt: now,
			},
			errorExpected: false,
		},
		{
			name:   "error",
			userID: 1,
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					GetMasterByUser(
						gomock.Any(),
						&toys.GetMasterByUserIn{UserID: 1},
					).
					Return(nil, errors.New("not found")).
					Times(1)
			},
			expectedMaster: nil,
			errorExpected:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysClient)
			}

			master, err := repo.GetMasterByUserID(context.Background(), tc.userID)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, master)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedMaster, master)
			}
		})
	}
}

func TestToysRepository_UpdateMaster(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysClient := mockclients.NewMockToysClient(ctrl)
	repo := NewToysRepository(toysClient)

	testCases := []struct {
		name          string
		masterData    entities.UpdateMasterDTO
		setupMocks    func(toysClient *mockclients.MockToysClient)
		errorExpected bool
	}{
		{
			name: "success",
			masterData: entities.UpdateMasterDTO{
				ID:   1,
				Info: pointers.New("Master Info"),
			},
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					UpdateMaster(
						gomock.Any(),
						&toys.UpdateMasterIn{
							ID:   1,
							Info: pointers.New("Master Info"),
						},
					).
					Return(&emptypb.Empty{}, nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "error",
			masterData: entities.UpdateMasterDTO{
				ID: 1,
			},
			setupMocks: func(toysClient *mockclients.MockToysClient) {
				toysClient.
					EXPECT().
					UpdateMaster(
						gomock.Any(),
						&toys.UpdateMasterIn{
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
				tc.setupMocks(toysClient)
			}

			err := repo.UpdateMaster(context.Background(), tc.masterData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
