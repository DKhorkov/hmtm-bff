package services

import (
	"context"
	"errors"
	"github.com/DKhorkov/libs/pointers"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	mockrepositories "github.com/DKhorkov/hmtm-bff/mocks/repositories"
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
)

func TestToysService_AddToy(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewToysService(toysRepository, logger)

	testCases := []struct {
		name          string
		toyData       entities.AddToyDTO
		setupMocks    func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger)
		expectedID    uint64
		errorExpected bool
	}{
		{
			name: "success",
			toyData: entities.AddToyDTO{
				UserID:      1,
				CategoryID:  1,
				Name:        "Test Toy",
				Description: "Test Description",
				Price:       100.0,
				Quantity:    10,
				TagIDs:      []uint32{1, 2},
				Attachments: []string{"link1"},
			},
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					AddToy(gomock.Any(), entities.AddToyDTO{
						UserID:      1,
						CategoryID:  1,
						Name:        "Test Toy",
						Description: "Test Description",
						Price:       100.0,
						Quantity:    10,
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
			toyData: entities.AddToyDTO{
				UserID:      1,
				CategoryID:  1,
				Name:        "Test Toy",
				Description: "Test Description",
				Price:       100.0,
				Quantity:    10,
				TagIDs:      []uint32{1, 2},
				Attachments: []string{"link1"},
			},
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					AddToy(gomock.Any(), entities.AddToyDTO{
						UserID:      1,
						CategoryID:  1,
						Name:        "Test Toy",
						Description: "Test Description",
						Price:       100.0,
						Quantity:    10,
						TagIDs:      []uint32{1, 2},
						Attachments: []string{"link1"},
					}).
					Return(uint64(0), errors.New("add failed")).
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
				tc.setupMocks(toysRepository, logger)
			}

			toyID, err := service.AddToy(context.Background(), tc.toyData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedID, toyID)
		})
	}
}

func TestToysService_GetToys(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewToysService(toysRepository, logger)

	testCases := []struct {
		name          string
		pagination    *entities.Pagination
		filters       *entities.ToysFilters
		setupMocks    func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger)
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
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetToys(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
						&entities.ToysFilters{
							Search:              pointers.New("Toy"),
							PriceCeil:           pointers.New[float32](1000),
							PriceFloor:          pointers.New[float32](10),
							QuantityFloor:       pointers.New[uint32](1),
							CategoryID:          pointers.New[uint32](1),
							TagIDs:              []uint32{1},
							CreatedAtOrderByAsc: pointers.New(true),
						},
					).
					Return(
						[]entities.Toy{
							{
								ID:          1,
								MasterID:    1,
								CategoryID:  1,
								Name:        "Test Toy",
								Description: "Test Description",
								Price:       100.0,
								Quantity:    10,
								CreatedAt:   now,
								UpdatedAt:   now,
								Tags:        []entities.Tag{{ID: 1, Name: "tag1"}},
								Attachments: []entities.ToyAttachment{{ID: 1, ToyID: 1, Link: "link1"}},
							},
						},
						nil).
					Times(1)
			},
			expectedToys: []entities.Toy{
				{
					ID:          1,
					MasterID:    1,
					CategoryID:  1,
					Name:        "Test Toy",
					Description: "Test Description",
					Price:       100.0,
					Quantity:    10,
					CreatedAt:   now,
					UpdatedAt:   now,
					Tags:        []entities.Tag{{ID: 1, Name: "tag1"}},
					Attachments: []entities.ToyAttachment{{ID: 1, ToyID: 1, Link: "link1"}},
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
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetToys(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
						&entities.ToysFilters{
							Search:              pointers.New("Toy"),
							PriceCeil:           pointers.New[float32](1000),
							PriceFloor:          pointers.New[float32](10),
							QuantityFloor:       pointers.New[uint32](1),
							CategoryID:          pointers.New[uint32](1),
							TagIDs:              []uint32{1},
							CreatedAtOrderByAsc: pointers.New(true),
						},
					).
					Return(nil, errors.New("fetch failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedToys:  nil,
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysRepository, logger)
			}

			toys, err := service.GetToys(context.Background(), tc.pagination, tc.filters)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, toys)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedToys, toys)
			}
		})
	}
}

func TestToysService_CountToys(t *testing.T) {
	testCases := []struct {
		name          string
		filters       *entities.ToysFilters
		setupMocks    func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger)
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
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					CountToys(
						gomock.Any(),
						&entities.ToysFilters{
							Search:              pointers.New("Toy"),
							PriceCeil:           pointers.New[float32](1000),
							PriceFloor:          pointers.New[float32](10),
							QuantityFloor:       pointers.New[uint32](1),
							CategoryID:          pointers.New[uint32](1),
							TagIDs:              []uint32{1},
							CreatedAtOrderByAsc: pointers.New(true),
						},
					).
					Return(uint64(1), nil).
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
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					CountToys(
						gomock.Any(),
						&entities.ToysFilters{
							Search:              pointers.New("Toy"),
							PriceCeil:           pointers.New[float32](1000),
							PriceFloor:          pointers.New[float32](10),
							QuantityFloor:       pointers.New[uint32](1),
							CategoryID:          pointers.New[uint32](1),
							TagIDs:              []uint32{1},
							CreatedAtOrderByAsc: pointers.New(true),
						},
					).
					Return(uint64(0), errors.New("error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewToysService(toysRepository, logger)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysRepository, logger)
			}

			actual, err := service.CountToys(context.Background(), tc.filters)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestToysService_GetMasterToys(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewToysService(toysRepository, logger)

	testCases := []struct {
		name          string
		pagination    *entities.Pagination
		filters       *entities.ToysFilters
		masterID      uint64
		setupMocks    func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger)
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
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetMasterToys(
						gomock.Any(),
						uint64(1),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
						&entities.ToysFilters{
							Search:              pointers.New("Toy"),
							PriceCeil:           pointers.New[float32](1000),
							PriceFloor:          pointers.New[float32](10),
							QuantityFloor:       pointers.New[uint32](1),
							CategoryID:          pointers.New[uint32](1),
							TagIDs:              []uint32{1},
							CreatedAtOrderByAsc: pointers.New(true),
						},
					).
					Return(
						[]entities.Toy{
							{
								ID:          1,
								MasterID:    1,
								CategoryID:  1,
								Name:        "Test Toy",
								Description: "Test Description",
								Price:       100.0,
								Quantity:    10,
								CreatedAt:   now,
								UpdatedAt:   now,
								Tags:        []entities.Tag{{ID: 1, Name: "tag1"}},
							},
						},
						nil,
					).
					Times(1)
			},
			expectedToys: []entities.Toy{
				{
					ID:          1,
					MasterID:    1,
					CategoryID:  1,
					Name:        "Test Toy",
					Description: "Test Description",
					Price:       100.0,
					Quantity:    10,
					CreatedAt:   now,
					UpdatedAt:   now,
					Tags:        []entities.Tag{{ID: 1, Name: "tag1"}},
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
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetMasterToys(
						gomock.Any(),
						uint64(1),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
						&entities.ToysFilters{
							Search:              pointers.New("Toy"),
							PriceCeil:           pointers.New[float32](1000),
							PriceFloor:          pointers.New[float32](10),
							QuantityFloor:       pointers.New[uint32](1),
							CategoryID:          pointers.New[uint32](1),
							TagIDs:              []uint32{1},
							CreatedAtOrderByAsc: pointers.New(true),
						},
					).
					Return(nil, errors.New("fetch failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedToys:  nil,
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysRepository, logger)
			}

			toys, err := service.GetMasterToys(context.Background(), tc.masterID, tc.pagination, tc.filters)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, toys)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedToys, toys)
			}
		})
	}
}

func TestToysService_GetUserToys(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewToysService(toysRepository, logger)

	testCases := []struct {
		name          string
		pagination    *entities.Pagination
		filters       *entities.ToysFilters
		userID        uint64
		setupMocks    func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger)
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
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetUserToys(
						gomock.Any(),
						uint64(1),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
						&entities.ToysFilters{
							Search:              pointers.New("Toy"),
							PriceCeil:           pointers.New[float32](1000),
							PriceFloor:          pointers.New[float32](10),
							QuantityFloor:       pointers.New[uint32](1),
							CategoryID:          pointers.New[uint32](1),
							TagIDs:              []uint32{1},
							CreatedAtOrderByAsc: pointers.New(true),
						},
					).
					Return([]entities.Toy{
						{
							ID:          1,
							MasterID:    1,
							CategoryID:  1,
							Name:        "Test Toy",
							Description: "Test Description",
							Price:       100.0,
							Quantity:    10,
							CreatedAt:   now,
							UpdatedAt:   now,
						},
					}, nil).
					Times(1)
			},
			expectedToys: []entities.Toy{
				{
					ID:          1,
					MasterID:    1,
					CategoryID:  1,
					Name:        "Test Toy",
					Description: "Test Description",
					Price:       100.0,
					Quantity:    10,
					CreatedAt:   now,
					UpdatedAt:   now,
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
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetUserToys(
						gomock.Any(),
						uint64(1),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
						&entities.ToysFilters{
							Search:              pointers.New("Toy"),
							PriceCeil:           pointers.New[float32](1000),
							PriceFloor:          pointers.New[float32](10),
							QuantityFloor:       pointers.New[uint32](1),
							CategoryID:          pointers.New[uint32](1),
							TagIDs:              []uint32{1},
							CreatedAtOrderByAsc: pointers.New(true),
						},
					).
					Return(nil, errors.New("fetch failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedToys:  nil,
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysRepository, logger)
			}

			toys, err := service.GetUserToys(context.Background(), tc.userID, tc.pagination, tc.filters)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, toys)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedToys, toys)
			}
		})
	}
}

func TestToysService_GetToyByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewToysService(toysRepository, logger)

	testCases := []struct {
		name          string
		id            uint64
		setupMocks    func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger)
		expectedToy   *entities.Toy
		errorExpected bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetToyByID(gomock.Any(), uint64(1)).
					Return(&entities.Toy{
						ID:          1,
						MasterID:    1,
						CategoryID:  1,
						Name:        "Test Toy",
						Description: "Test Description",
						Price:       100.0,
						Quantity:    10,
						CreatedAt:   now,
						UpdatedAt:   now,
						Tags:        []entities.Tag{{ID: 1, Name: "tag1"}},
					}, nil).
					Times(1)
			},
			expectedToy: &entities.Toy{
				ID:          1,
				MasterID:    1,
				CategoryID:  1,
				Name:        "Test Toy",
				Description: "Test Description",
				Price:       100.0,
				Quantity:    10,
				CreatedAt:   now,
				UpdatedAt:   now,
				Tags:        []entities.Tag{{ID: 1, Name: "tag1"}},
			},
			errorExpected: false,
		},
		{
			name: "error",
			id:   1,
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetToyByID(gomock.Any(), uint64(1)).
					Return(nil, errors.New("not found")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedToy:   nil,
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysRepository, logger)
			}

			toy, err := service.GetToyByID(context.Background(), tc.id)
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

func TestToysService_GetMasters(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewToysService(toysRepository, logger)

	info := "Master Info"
	testCases := []struct {
		name            string
		pagination      *entities.Pagination
		setupMocks      func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger)
		expectedMasters []entities.Master
		errorExpected   bool
	}{
		{
			name: "success",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetMasters(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return([]entities.Master{
						{
							ID:        1,
							UserID:    1,
							Info:      &info,
							CreatedAt: now,
							UpdatedAt: now,
						},
					}, nil).
					Times(1)
			},
			expectedMasters: []entities.Master{
				{
					ID:        1,
					UserID:    1,
					Info:      &info,
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
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetMasters(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return(nil, errors.New("fetch failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedMasters: nil,
			errorExpected:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysRepository, logger)
			}

			masters, err := service.GetMasters(context.Background(), tc.pagination)
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

func TestToysService_GetMasterByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewToysService(toysRepository, logger)

	info := "Master Info"
	testCases := []struct {
		name           string
		id             uint64
		setupMocks     func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger)
		expectedMaster *entities.Master
		errorExpected  bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetMasterByID(gomock.Any(), uint64(1)).
					Return(&entities.Master{
						ID:        1,
						UserID:    1,
						Info:      &info,
						CreatedAt: now,
						UpdatedAt: now,
					}, nil).
					Times(1)
			},
			expectedMaster: &entities.Master{
				ID:        1,
				UserID:    1,
				Info:      &info,
				CreatedAt: now,
				UpdatedAt: now,
			},
			errorExpected: false,
		},
		{
			name: "error",
			id:   1,
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetMasterByID(gomock.Any(), uint64(1)).
					Return(nil, errors.New("not found")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedMaster: nil,
			errorExpected:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysRepository, logger)
			}

			master, err := service.GetMasterByID(context.Background(), tc.id)
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

func TestToysService_RegisterMaster(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewToysService(toysRepository, logger)

	info := "Master Info"
	testCases := []struct {
		name          string
		masterData    entities.RegisterMasterDTO
		setupMocks    func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger)
		expectedID    uint64
		errorExpected bool
	}{
		{
			name: "success",
			masterData: entities.RegisterMasterDTO{
				UserID: 1,
				Info:   &info,
			},
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					RegisterMaster(gomock.Any(), entities.RegisterMasterDTO{
						UserID: 1,
						Info:   &info,
					}).
					Return(uint64(1), nil).
					Times(1)
			},
			expectedID:    1,
			errorExpected: false,
		},
		{
			name: "error",
			masterData: entities.RegisterMasterDTO{
				UserID: 1,
				Info:   &info,
			},
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					RegisterMaster(gomock.Any(), entities.RegisterMasterDTO{
						UserID: 1,
						Info:   &info,
					}).
					Return(uint64(0), errors.New("register failed")).
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
				tc.setupMocks(toysRepository, logger)
			}

			masterID, err := service.RegisterMaster(context.Background(), tc.masterData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedID, masterID)
		})
	}
}

func TestToysService_GetAllCategories(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewToysService(toysRepository, logger)

	testCases := []struct {
		name               string
		setupMocks         func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger)
		expectedCategories []entities.Category
		errorExpected      bool
	}{
		{
			name: "success",
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetAllCategories(gomock.Any()).
					Return([]entities.Category{
						{ID: 1, Name: "Category 1"},
					}, nil).
					Times(1)
			},
			expectedCategories: []entities.Category{
				{ID: 1, Name: "Category 1"},
			},
			errorExpected: false,
		},
		{
			name: "error",
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetAllCategories(gomock.Any()).
					Return(nil, errors.New("fetch failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedCategories: nil,
			errorExpected:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysRepository, logger)
			}

			categories, err := service.GetAllCategories(context.Background())
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

func TestToysService_GetCategoryByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewToysService(toysRepository, logger)

	testCases := []struct {
		name             string
		id               uint32
		setupMocks       func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger)
		expectedCategory *entities.Category
		errorExpected    bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetCategoryByID(gomock.Any(), uint32(1)).
					Return(&entities.Category{ID: 1, Name: "Category 1"}, nil).
					Times(1)
			},
			expectedCategory: &entities.Category{ID: 1, Name: "Category 1"},
			errorExpected:    false,
		},
		{
			name: "error",
			id:   1,
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetCategoryByID(gomock.Any(), uint32(1)).
					Return(nil, errors.New("not found")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedCategory: nil,
			errorExpected:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysRepository, logger)
			}

			category, err := service.GetCategoryByID(context.Background(), tc.id)
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

func TestToysService_GetAllTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewToysService(toysRepository, logger)

	testCases := []struct {
		name          string
		setupMocks    func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger)
		expectedTags  []entities.Tag
		errorExpected bool
	}{
		{
			name: "success",
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{
						{ID: 1, Name: "tag1"},
					}, nil).
					Times(1)
			},
			expectedTags: []entities.Tag{
				{ID: 1, Name: "tag1"},
			},
			errorExpected: false,
		},
		{
			name: "error",
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(nil, errors.New("fetch failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedTags:  nil,
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysRepository, logger)
			}

			tags, err := service.GetAllTags(context.Background())
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

func TestToysService_GetTagByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewToysService(toysRepository, logger)

	testCases := []struct {
		name          string
		id            uint32
		setupMocks    func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger)
		expectedTag   *entities.Tag
		errorExpected bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetTagByID(gomock.Any(), uint32(1)).
					Return(&entities.Tag{ID: 1, Name: "tag1"}, nil).
					Times(1)
			},
			expectedTag:   &entities.Tag{ID: 1, Name: "tag1"},
			errorExpected: false,
		},
		{
			name: "error",
			id:   1,
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetTagByID(gomock.Any(), uint32(1)).
					Return(nil, errors.New("not found")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedTag:   nil,
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysRepository, logger)
			}

			tag, err := service.GetTagByID(context.Background(), tc.id)
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

func TestToysService_CreateTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewToysService(toysRepository, logger)

	testCases := []struct {
		name          string
		tagsData      []entities.CreateTagDTO
		setupMocks    func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger)
		expectedIDs   []uint32
		errorExpected bool
	}{
		{
			name: "success",
			tagsData: []entities.CreateTagDTO{
				{Name: "tag1"},
				{Name: "tag2"},
			},
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					CreateTags(gomock.Any(), []entities.CreateTagDTO{
						{Name: "tag1"},
						{Name: "tag2"},
					}).
					Return([]uint32{1, 2}, nil).
					Times(1)
			},
			expectedIDs:   []uint32{1, 2},
			errorExpected: false,
		},
		{
			name: "error",
			tagsData: []entities.CreateTagDTO{
				{Name: "tag1"},
				{Name: "tag2"},
			},
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					CreateTags(gomock.Any(), []entities.CreateTagDTO{
						{Name: "tag1"},
						{Name: "tag2"},
					}).
					Return(nil, errors.New("create failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedIDs:   nil,
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysRepository, logger)
			}

			tagIDs, err := service.CreateTags(context.Background(), tc.tagsData)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, tagIDs)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedIDs, tagIDs)
			}
		})
	}
}

func TestToysService_UpdateToy(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewToysService(toysRepository, logger)

	price := float32(150.0)
	name := "Updated Toy"
	testCases := []struct {
		name          string
		toyData       entities.UpdateToyDTO
		setupMocks    func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger)
		errorExpected bool
	}{
		{
			name: "success",
			toyData: entities.UpdateToyDTO{
				ID:    1,
				Name:  &name,
				Price: &price,
			},
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					UpdateToy(gomock.Any(), entities.UpdateToyDTO{
						ID:    1,
						Name:  &name,
						Price: &price,
					}).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "error",
			toyData: entities.UpdateToyDTO{
				ID:    1,
				Name:  &name,
				Price: &price,
			},
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					UpdateToy(gomock.Any(), entities.UpdateToyDTO{
						ID:    1,
						Name:  &name,
						Price: &price,
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
				tc.setupMocks(toysRepository, logger)
			}

			err := service.UpdateToy(context.Background(), tc.toyData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestToysService_DeleteToy(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewToysService(toysRepository, logger)

	testCases := []struct {
		name          string
		id            uint64
		setupMocks    func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger)
		errorExpected bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					DeleteToy(gomock.Any(), uint64(1)).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "error",
			id:   1,
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					DeleteToy(gomock.Any(), uint64(1)).
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
				tc.setupMocks(toysRepository, logger)
			}

			err := service.DeleteToy(context.Background(), tc.id)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestToysService_GetMasterByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewToysService(toysRepository, logger)

	info := "Master Info"
	testCases := []struct {
		name           string
		userID         uint64
		setupMocks     func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger)
		expectedMaster *entities.Master
		errorExpected  bool
	}{
		{
			name:   "success",
			userID: 1,
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetMasterByUserID(gomock.Any(), uint64(1)).
					Return(&entities.Master{
						ID:        1,
						UserID:    1,
						Info:      &info,
						CreatedAt: now,
						UpdatedAt: now,
					}, nil).
					Times(1)
			},
			expectedMaster: &entities.Master{
				ID:        1,
				UserID:    1,
				Info:      &info,
				CreatedAt: now,
				UpdatedAt: now,
			},
			errorExpected: false,
		},
		{
			name:   "error",
			userID: 1,
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					GetMasterByUserID(gomock.Any(), uint64(1)).
					Return(nil, errors.New("not found")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedMaster: nil,
			errorExpected:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysRepository, logger)
			}

			master, err := service.GetMasterByUserID(context.Background(), tc.userID)
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

func TestToysService_UpdateMaster(t *testing.T) {
	ctrl := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewToysService(toysRepository, logger)

	info := "Updated Master Info"
	testCases := []struct {
		name          string
		masterData    entities.UpdateMasterDTO
		setupMocks    func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger)
		errorExpected bool
	}{
		{
			name: "success",
			masterData: entities.UpdateMasterDTO{
				ID:   1,
				Info: &info,
			},
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					UpdateMaster(gomock.Any(), entities.UpdateMasterDTO{
						ID:   1,
						Info: &info,
					}).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "error",
			masterData: entities.UpdateMasterDTO{
				ID:   1,
				Info: &info,
			},
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *mocklogging.MockLogger) {
				toysRepository.
					EXPECT().
					UpdateMaster(gomock.Any(), entities.UpdateMasterDTO{
						ID:   1,
						Info: &info,
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
				tc.setupMocks(toysRepository, logger)
			}

			err := service.UpdateMaster(context.Background(), tc.masterData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
