package usecases

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/DKhorkov/libs/pointers"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	mockusecases "github.com/DKhorkov/hmtm-bff/mocks/usecases"
	mockcache "github.com/DKhorkov/libs/cache/mocks"
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
)

func TestCacheDecorator_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	userID := uint64(1)
	user := &entities.User{ID: userID, DisplayName: "Test User"}
	cacheKey := fmt.Sprintf("%s:%d", getUserByIDPrefix, userID)

	testCases := []struct {
		name          string
		userID        uint64
		cacheValue    string
		expectedUser  *entities.User
		expectedError error
		setupMocks    func()
	}{
		{
			name:         "success from cache",
			userID:       userID,
			cacheValue:   `{"id":1,"DisplayName":"Cached User"}`,
			expectedUser: &entities.User{ID: 1, DisplayName: "Cached User"},
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return(`{"id":1,"DisplayName":"Cached User"}`, nil).
					Times(1)
			},
		},
		{
			name:         "success from db",
			userID:       userID,
			expectedUser: user,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetUserByID(gomock.Any(), userID).
					Return(user, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), cacheKey, gomock.Any(), getUserByIDTTL).
					Return(nil).
					Times(1)
			},
		},
		{
			name:         "cache ping error, success from db",
			userID:       userID,
			expectedUser: user,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetUserByID(gomock.Any(), userID).
					Return(user, nil).
					Times(1)
			},
		},
		{
			name:          "db error",
			userID:        userID,
			expectedError: errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetUserByID(gomock.Any(), userID).
					Return(nil, errors.New("db error")).
					Times(1)
			},
		},
		{
			name:         "set cache error",
			userID:       userID,
			expectedUser: user,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2) // One for cache miss, one for set cache error

				useCasesMock.
					EXPECT().
					GetUserByID(gomock.Any(), userID).
					Return(user, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), cacheKey, gomock.Any(), getUserByIDTTL).
					Return(errors.New("set cache error")).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			result, err := decorator.GetUserByID(context.Background(), tc.userID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, result)
			}
		})
	}
}

func TestCacheDecorator_GetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	email := "test@example.com"
	user := &entities.User{ID: 1, Email: email, DisplayName: "Test User"}
	cacheKey := fmt.Sprintf("%s:%s", getUserByEmailPrefix, email)

	testCases := []struct {
		name          string
		email         string
		cacheValue    string
		expectedUser  *entities.User
		expectedError error
		setupMocks    func()
	}{
		{
			name:         "success from cache",
			email:        email,
			cacheValue:   `{"id":1,"email":"test@example.com","displayName":"Cached User"}`,
			expectedUser: &entities.User{ID: 1, Email: email, DisplayName: "Cached User"},
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return(`{"id":1,"email":"test@example.com","displayName":"Cached User"}`, nil).
					Times(1)
			},
		},
		{
			name:         "success from db",
			email:        email,
			expectedUser: user,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)
				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)
				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
				useCasesMock.
					EXPECT().
					GetUserByEmail(gomock.Any(), email).
					Return(user, nil).
					Times(1)
				cacheMock.
					EXPECT().
					Set(gomock.Any(), cacheKey, gomock.Any(), getUserByEmailTTL).
					Return(nil).
					Times(1)
			},
		},
		{
			name:         "cache ping error, success from db",
			email:        email,
			expectedUser: user,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)
				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
				useCasesMock.
					EXPECT().
					GetUserByEmail(gomock.Any(), email).
					Return(user, nil).
					Times(1)
			},
		},
		{
			name:          "db error",
			email:         email,
			expectedError: errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)
				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)
				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
				useCasesMock.
					EXPECT().
					GetUserByEmail(gomock.Any(), email).
					Return(nil, errors.New("db error")).
					Times(1)
			},
		},
		{
			name:         "set cache error",
			email:        email,
			expectedUser: user,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)
				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)
				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2)
				useCasesMock.
					EXPECT().
					GetUserByEmail(gomock.Any(), email).
					Return(user, nil).
					Times(1)
				cacheMock.
					EXPECT().
					Set(gomock.Any(), cacheKey, gomock.Any(), getUserByEmailTTL).
					Return(errors.New("set cache error")).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}
			result, err := decorator.GetUserByEmail(context.Background(), tc.email)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, result)
			}
		})
	}
}

func TestCacheDecorator_GetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	pagination := &entities.Pagination{Offset: pointers.New[uint64](1), Limit: pointers.New[uint64](10)}
	users := []entities.User{{ID: 1, DisplayName: "Test User"}}

	testCases := []struct {
		name          string
		pagination    *entities.Pagination
		cacheValue    string
		expectedUsers []entities.User
		expectedError error
		setupMocks    func()
	}{
		{
			name:          "success from cache",
			pagination:    pagination,
			cacheValue:    `[{"id":1,"displayName":"Cached User"}]`,
			expectedUsers: []entities.User{{ID: 1, DisplayName: "Cached User"}},
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)
				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return `[{"id":1,"displayName":"Cached User"}]`, nil
						},
					),
				)
			},
		},
		{
			name:          "success from db",
			pagination:    pagination,
			expectedUsers: users,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)
				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetUsers(gomock.Any(), pagination).
					Return(users, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), getUsersTTL).
					Return(nil).
					Times(1)
			},
		},
		{
			name:          "cache ping error, success from db",
			pagination:    pagination,
			expectedUsers: users,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)
				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
				useCasesMock.
					EXPECT().
					GetUsers(gomock.Any(), pagination).
					Return(users, nil).
					Times(1)
			},
		},
		{
			name:          "db error",
			pagination:    pagination,
			expectedError: errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetUsers(gomock.Any(), pagination).
					Return(nil, errors.New("db error")).
					Times(1)
			},
		},
		{
			name:          "set cache error",
			pagination:    pagination,
			expectedUsers: users,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2)

				useCasesMock.
					EXPECT().
					GetUsers(gomock.Any(), pagination).
					Return(users, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), getUsersTTL).
					Return(errors.New("set cache error")).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			result, err := decorator.GetUsers(context.Background(), tc.pagination)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUsers, result)
			}
		})
	}
}

func TestCacheDecorator_GetToys(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	pagination := &entities.Pagination{Offset: pointers.New[uint64](1), Limit: pointers.New[uint64](10)}
	filters := &entities.ToysFilters{Search: pointers.New("toy")}
	toys := []entities.Toy{{ID: 1}}

	testCases := []struct {
		name          string
		pagination    *entities.Pagination
		filters       *entities.ToysFilters
		cacheValue    string
		expectedToys  []entities.Toy
		expectedError error
		setupMocks    func()
	}{
		{
			name:         "success from cache",
			pagination:   pagination,
			filters:      filters,
			cacheValue:   `[{"id":1,"name":"Cached Toy"}]`,
			expectedToys: []entities.Toy{{ID: 1, Name: "Cached Toy"}},
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return `[{"id":1,"name":"Cached Toy"}]`, nil
						},
					),
				)
			},
		},
		{
			name:         "success from db",
			pagination:   pagination,
			filters:      filters,
			expectedToys: toys,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetToys(gomock.Any(), pagination, filters).
					Return(toys, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), getToysTTL).
					Return(nil).
					Times(1)
			},
		},
		{
			name:         "cache ping error, success from db",
			pagination:   pagination,
			filters:      filters,
			expectedToys: toys,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetToys(gomock.Any(), pagination, filters).
					Return(toys, nil).
					Times(1)
			},
		},
		{
			name:          "db error",
			pagination:    pagination,
			filters:       filters,
			expectedError: errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetToys(gomock.Any(), pagination, filters).
					Return(nil, errors.New("db error")).
					Times(1)
			},
		},
		{
			name:         "set cache error",
			pagination:   pagination,
			filters:      filters,
			expectedToys: toys,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2)

				useCasesMock.
					EXPECT().
					GetToys(gomock.Any(), pagination, filters).
					Return(toys, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), getToysTTL).
					Return(errors.New("set cache error")).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			result, err := decorator.GetToys(context.Background(), tc.pagination, tc.filters)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedToys, result)
			}
		})
	}
}

func TestCacheDecorator_CountToys(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	filters := &entities.ToysFilters{Search: pointers.New("toy")}
	count := uint64(100)

	testCases := []struct {
		name          string
		filters       *entities.ToysFilters
		cacheValue    string
		expectedCount uint64
		expectedError error
		setupMocks    func()
	}{
		{
			name:          "success from cache",
			filters:       filters,
			cacheValue:    "100",
			expectedCount: count,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "100", nil
						},
					),
				)
			},
		},
		{
			name:          "success from db",
			filters:       filters,
			expectedCount: count,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					CountToys(gomock.Any(), filters).
					Return(count, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), countToysTTL).
					Return(nil).
					Times(1)
			},
		},
		{
			name:          "cache ping error, success from db",
			filters:       filters,
			expectedCount: count,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					CountToys(gomock.Any(), filters).
					Return(count, nil).
					Times(1)
			},
		},
		{
			name:          "db error",
			filters:       filters,
			expectedError: errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					CountToys(gomock.Any(), filters).
					Return(uint64(0), errors.New("db error")).
					Times(1)
			},
		},
		{
			name:          "set cache error",
			filters:       filters,
			expectedCount: count,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2)

				useCasesMock.
					EXPECT().
					CountToys(gomock.Any(), filters).
					Return(count, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), countToysTTL).
					Return(errors.New("set cache error")).
					Times(1)
			},
		},
		{
			name:          "parse error, success from db",
			filters:       filters,
			cacheValue:    "invalid",
			expectedCount: count,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "invalid", nil
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					CountToys(gomock.Any(), filters).
					Return(count, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), countToysTTL).
					Return(nil).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			result, err := decorator.CountToys(context.Background(), tc.filters)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Equal(t, uint64(0), result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCount, result)
			}
		})
	}
}

func TestCacheDecorator_CountMasterToys(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	masterID := uint64(1)
	filters := &entities.ToysFilters{Search: pointers.New("toy")}
	count := uint64(50)

	testCases := []struct {
		name          string
		masterID      uint64
		filters       *entities.ToysFilters
		cacheValue    string
		expectedCount uint64
		expectedError error
		setupMocks    func()
	}{
		{
			name:          "success from cache",
			masterID:      masterID,
			filters:       filters,
			cacheValue:    "50",
			expectedCount: count,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "50", nil
						},
					),
				)
			},
		},
		{
			name:          "success from db",
			masterID:      masterID,
			filters:       filters,
			expectedCount: count,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					CountMasterToys(gomock.Any(), masterID, filters).
					Return(count, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), countMasterToysTTL).
					Return(nil).
					Times(1)
			},
		},
		{
			name:          "cache ping error, success from db",
			masterID:      masterID,
			filters:       filters,
			expectedCount: count,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					CountMasterToys(gomock.Any(), masterID, filters).
					Return(count, nil).
					Times(1)
			},
		},
		{
			name:          "db error",
			masterID:      masterID,
			filters:       filters,
			expectedError: errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					CountMasterToys(gomock.Any(), masterID, filters).
					Return(uint64(0), errors.New("db error")).
					Times(1)
			},
		},
		{
			name:          "set cache error",
			masterID:      masterID,
			filters:       filters,
			expectedCount: count,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2)

				useCasesMock.
					EXPECT().
					CountMasterToys(gomock.Any(), masterID, filters).
					Return(count, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), countMasterToysTTL).
					Return(errors.New("set cache error")).
					Times(1)
			},
		},
		{
			name:          "parse error, success from db",
			masterID:      masterID,
			filters:       filters,
			cacheValue:    "invalid",
			expectedCount: count,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "invalid", nil
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					CountMasterToys(gomock.Any(), masterID, filters).
					Return(count, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), countMasterToysTTL).
					Return(nil).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			result, err := decorator.CountMasterToys(context.Background(), tc.masterID, tc.filters)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Equal(t, uint64(0), result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCount, result)
			}
		})
	}
}

func TestCacheDecorator_GetMasterToys(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	masterID := uint64(1)
	pagination := &entities.Pagination{Offset: pointers.New[uint64](1), Limit: pointers.New[uint64](10)}
	filters := &entities.ToysFilters{Search: pointers.New("toy")}
	toys := []entities.Toy{{ID: 1}}

	testCases := []struct {
		name          string
		masterID      uint64
		pagination    *entities.Pagination
		filters       *entities.ToysFilters
		cacheValue    string
		expectedToys  []entities.Toy
		expectedError error
		setupMocks    func()
	}{
		{
			name:         "success from cache",
			masterID:     masterID,
			pagination:   pagination,
			filters:      filters,
			cacheValue:   `[{"id":1}]`,
			expectedToys: []entities.Toy{{ID: 1}},
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return `[{"id":1}]`, nil
						},
					),
				)
			},
		},
		{
			name:         "success from db",
			masterID:     masterID,
			pagination:   pagination,
			filters:      filters,
			expectedToys: toys,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetMasterToys(gomock.Any(), masterID, pagination, filters).
					Return(toys, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), getMasterToysTTL).
					Return(nil).
					Times(1)
			},
		},
		{
			name:         "cache ping error, success from db",
			masterID:     masterID,
			pagination:   pagination,
			filters:      filters,
			expectedToys: toys,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetMasterToys(gomock.Any(), masterID, pagination, filters).
					Return(toys, nil).
					Times(1)
			},
		},
		{
			name:          "db error",
			masterID:      masterID,
			pagination:    pagination,
			filters:       filters,
			expectedError: errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetMasterToys(gomock.Any(), masterID, pagination, filters).
					Return(nil, errors.New("db error")).
					Times(1)
			},
		},
		{
			name:         "set cache error",
			masterID:     masterID,
			pagination:   pagination,
			filters:      filters,
			expectedToys: toys,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2)

				useCasesMock.
					EXPECT().
					GetMasterToys(gomock.Any(), masterID, pagination, filters).
					Return(toys, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), getMasterToysTTL).
					Return(errors.New("set cache error")).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			result, err := decorator.GetMasterToys(context.Background(), tc.masterID, tc.pagination, tc.filters)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedToys, result)
			}
		})
	}
}

func TestCacheDecorator_GetToyByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	toyID := uint64(1)
	toy := &entities.Toy{ID: toyID}
	cacheKey := fmt.Sprintf("%s:%d", getToyByIDPrefix, toyID)

	testCases := []struct {
		name          string
		toyID         uint64
		cacheValue    string
		expectedToy   *entities.Toy
		expectedError error
		setupMocks    func()
	}{
		{
			name:        "success from cache",
			toyID:       toyID,
			cacheValue:  `{"id":1}`,
			expectedToy: &entities.Toy{ID: 1},
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return(`{"id":1}`, nil).
					Times(1)
			},
		},
		{
			name:        "success from db",
			toyID:       toyID,
			expectedToy: toy,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetToyByID(gomock.Any(), toyID).
					Return(toy, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), cacheKey, gomock.Any(), getToyByIDTTL).
					Return(nil).
					Times(1)
			},
		},
		{
			name:        "cache ping error, success from db",
			toyID:       toyID,
			expectedToy: toy,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetToyByID(gomock.Any(), toyID).
					Return(toy, nil).
					Times(1)
			},
		},
		{
			name:          "db error",
			toyID:         toyID,
			expectedError: errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetToyByID(gomock.Any(), toyID).
					Return(nil, errors.New("db error")).
					Times(1)
			},
		},
		{
			name:        "set cache error",
			toyID:       toyID,
			expectedToy: toy,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2)

				useCasesMock.
					EXPECT().
					GetToyByID(gomock.Any(), toyID).
					Return(toy, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), cacheKey, gomock.Any(), getToyByIDTTL).
					Return(errors.New("set cache error")).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			result, err := decorator.GetToyByID(context.Background(), tc.toyID)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedToy, result)
			}
		})
	}
}

func TestCacheDecorator_GetMasters(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	pagination := &entities.Pagination{Offset: pointers.New[uint64](1), Limit: pointers.New[uint64](10)}
	masters := []entities.Master{{ID: 1}}

	testCases := []struct {
		name            string
		pagination      *entities.Pagination
		cacheValue      string
		expectedMasters []entities.Master
		expectedError   error
		setupMocks      func()
	}{
		{
			name:            "success from cache",
			pagination:      pagination,
			cacheValue:      `[{"id":1,"userId":1}]`,
			expectedMasters: masters,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return `[{"id":1}]`, nil
						},
					),
				)
			},
		},
		{
			name:            "success from db",
			pagination:      pagination,
			expectedMasters: masters,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetMasters(gomock.Any(), pagination).
					Return(masters, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), getMastersTTL).
					Return(nil).
					Times(1)
			},
		},
		{
			name:            "cache ping error, success from db",
			pagination:      pagination,
			expectedMasters: masters,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetMasters(gomock.Any(), pagination).
					Return(masters, nil).
					Times(1)
			},
		},
		{
			name:          "db error",
			pagination:    pagination,
			expectedError: errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetMasters(gomock.Any(), pagination).
					Return(nil, errors.New("db error")).
					Times(1)
			},
		},
		{
			name:            "set cache error",
			pagination:      pagination,
			expectedMasters: masters,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2)

				useCasesMock.
					EXPECT().
					GetMasters(gomock.Any(), pagination).
					Return(masters, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), getMastersTTL).
					Return(errors.New("set cache error")).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			result, err := decorator.GetMasters(context.Background(), tc.pagination)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedMasters, result)
			}
		})
	}
}

func TestCacheDecorator_GetMasterByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	masterID := uint64(1)
	master := &entities.Master{ID: masterID}
	cacheKey := fmt.Sprintf("%s:%d", getMasterByIDPrefix, masterID)

	testCases := []struct {
		name           string
		masterID       uint64
		cacheValue     string
		expectedMaster *entities.Master
		expectedError  error
		setupMocks     func()
	}{
		{
			name:           "success from cache",
			masterID:       masterID,
			cacheValue:     `{"id":1,"userId":1}`,
			expectedMaster: &entities.Master{ID: 1},
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return(`{"id":1,"name":"Cached Master"}`, nil).
					Times(1)
			},
		},
		{
			name:           "success from db",
			masterID:       masterID,
			expectedMaster: master,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetMasterByID(gomock.Any(), masterID).
					Return(master, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), cacheKey, gomock.Any(), getMasterByIDTTL).
					Return(nil).
					Times(1)
			},
		},
		{
			name:           "cache ping error, success from db",
			masterID:       masterID,
			expectedMaster: master,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetMasterByID(gomock.Any(), masterID).
					Return(master, nil).
					Times(1)
			},
		},
		{
			name:          "db error",
			masterID:      masterID,
			expectedError: errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetMasterByID(gomock.Any(), masterID).
					Return(nil, errors.New("db error")).
					Times(1)
			},
		},
		{
			name:           "set cache error",
			masterID:       masterID,
			expectedMaster: master,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2)

				useCasesMock.
					EXPECT().
					GetMasterByID(gomock.Any(), masterID).
					Return(master, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), cacheKey, gomock.Any(), getMasterByIDTTL).
					Return(errors.New("set cache error")).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			result, err := decorator.GetMasterByID(context.Background(), tc.masterID)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedMaster, result)
			}
		})
	}
}

func TestCacheDecorator_GetMasterByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	userID := uint64(1)
	master := &entities.Master{ID: 1, UserID: userID}
	cacheKey := fmt.Sprintf("%s:%d", getMasterByUserIDPrefix, userID)

	testCases := []struct {
		name           string
		userID         uint64
		cacheValue     string
		expectedMaster *entities.Master
		expectedError  error
		setupMocks     func()
	}{
		{
			name:           "success from cache",
			userID:         userID,
			cacheValue:     `{"id":1,"userId":1}`,
			expectedMaster: &entities.Master{ID: 1, UserID: userID},
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return(`{"id":1,"userId":1}`, nil).
					Times(1)
			},
		},
		{
			name:           "success from db",
			userID:         userID,
			expectedMaster: master,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetMasterByUserID(gomock.Any(), userID).
					Return(master, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), cacheKey, gomock.Any(), getMasterByUserIDTTL).
					Return(nil).
					Times(1)
			},
		},
		{
			name:           "cache ping error, success from db",
			userID:         userID,
			expectedMaster: master,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetMasterByUserID(gomock.Any(), userID).
					Return(master, nil).
					Times(1)
			},
		},
		{
			name:          "db error",
			userID:        userID,
			expectedError: errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetMasterByUserID(gomock.Any(), userID).
					Return(nil, errors.New("db error")).
					Times(1)
			},
		},
		{
			name:           "set cache error",
			userID:         userID,
			expectedMaster: master,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2)

				useCasesMock.
					EXPECT().
					GetMasterByUserID(gomock.Any(), userID).
					Return(master, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), cacheKey, gomock.Any(), getMasterByUserIDTTL).
					Return(errors.New("set cache error")).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			result, err := decorator.GetMasterByUserID(context.Background(), tc.userID)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedMaster, result)
			}
		})
	}
}

func TestCacheDecorator_GetAllCategories(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	categories := []entities.Category{{ID: 1, Name: "Test Category"}}
	cacheKey := getCategoriesPrefix

	testCases := []struct {
		name               string
		cacheValue         string
		expectedCategories []entities.Category
		expectedError      error
		setupMocks         func()
	}{
		{
			name:               "success from cache",
			cacheValue:         `[{"id":1,"name":"Cached Category"}]`,
			expectedCategories: []entities.Category{{ID: 1, Name: "Cached Category"}},
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return(`[{"id":1,"name":"Cached Category"}]`, nil).
					Times(1)
			},
		},
		{
			name:               "success from db",
			expectedCategories: categories,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetAllCategories(gomock.Any()).
					Return(categories, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), cacheKey, gomock.Any(), getCategoriesTTL).
					Return(nil).
					Times(1)
			},
		},
		{
			name:               "cache ping error, success from db",
			expectedCategories: categories,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetAllCategories(gomock.Any()).
					Return(categories, nil).
					Times(1)
			},
		},
		{
			name:          "db error",
			expectedError: errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetAllCategories(gomock.Any()).
					Return(nil, errors.New("db error")).
					Times(1)
			},
		},
		{
			name:               "set cache error",
			expectedCategories: categories,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2)

				useCasesMock.
					EXPECT().
					GetAllCategories(gomock.Any()).
					Return(categories, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), cacheKey, gomock.Any(), getCategoriesTTL).
					Return(errors.New("set cache error")).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			result, err := decorator.GetAllCategories(context.Background())
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCategories, result)
			}
		})
	}
}

func TestCacheDecorator_GetCategoryByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	categoryID := uint32(1)
	category := &entities.Category{ID: categoryID, Name: "Test Category"}
	cacheKey := fmt.Sprintf("%s:%d", getCategoryByIDPrefix, categoryID)

	testCases := []struct {
		name             string
		categoryID       uint32
		cacheValue       string
		expectedCategory *entities.Category
		expectedError    error
		setupMocks       func()
	}{
		{
			name:             "success from cache",
			categoryID:       categoryID,
			cacheValue:       `{"id":1,"name":"Cached Category"}`,
			expectedCategory: &entities.Category{ID: 1, Name: "Cached Category"},
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return(`{"id":1,"name":"Cached Category"}`, nil).
					Times(1)
			},
		},
		{
			name:             "success from db",
			categoryID:       categoryID,
			expectedCategory: category,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetCategoryByID(gomock.Any(), categoryID).
					Return(category, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), cacheKey, gomock.Any(), getCategoryByIDTTL).
					Return(nil).
					Times(1)
			},
		},
		{
			name:             "cache ping error, success from db",
			categoryID:       categoryID,
			expectedCategory: category,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetCategoryByID(gomock.Any(), categoryID).
					Return(category, nil).
					Times(1)
			},
		},
		{
			name:          "db error",
			categoryID:    categoryID,
			expectedError: errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetCategoryByID(gomock.Any(), categoryID).
					Return(nil, errors.New("db error")).
					Times(1)
			},
		},
		{
			name:             "set cache error",
			categoryID:       categoryID,
			expectedCategory: category,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2)

				useCasesMock.
					EXPECT().
					GetCategoryByID(gomock.Any(), categoryID).
					Return(category, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), cacheKey, gomock.Any(), getCategoryByIDTTL).
					Return(errors.New("set cache error")).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			result, err := decorator.GetCategoryByID(context.Background(), tc.categoryID)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCategory, result)
			}
		})
	}
}

func TestCacheDecorator_GetAllTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	tags := []entities.Tag{{ID: 1, Name: "Test Tag"}}
	cacheKey := getTagsPrefix

	testCases := []struct {
		name          string
		cacheValue    string
		expectedTags  []entities.Tag
		expectedError error
		setupMocks    func()
	}{
		{
			name:         "success from cache",
			cacheValue:   `[{"id":1,"name":"Cached Tag"}]`,
			expectedTags: []entities.Tag{{ID: 1, Name: "Cached Tag"}},
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return(`[{"id":1,"name":"Cached Tag"}]`, nil).
					Times(1)
			},
		},
		{
			name:         "success from db",
			expectedTags: tags,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(tags, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), cacheKey, gomock.Any(), getTagsTTL).
					Return(nil).
					Times(1)
			},
		},
		{
			name:         "cache ping error, success from db",
			expectedTags: tags,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(tags, nil).
					Times(1)
			},
		},
		{
			name:          "db error",
			expectedError: errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(nil, errors.New("db error")).
					Times(1)
			},
		},
		{
			name:         "set cache error",
			expectedTags: tags,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2)

				useCasesMock.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(tags, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), cacheKey, gomock.Any(), getTagsTTL).
					Return(errors.New("set cache error")).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			result, err := decorator.GetAllTags(context.Background())
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTags, result)
			}
		})
	}
}

func TestCacheDecorator_GetTagByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	tagID := uint32(1)
	tag := &entities.Tag{ID: tagID, Name: "Test Tag"}
	cacheKey := fmt.Sprintf("%s:%d", getTagByIDPrefix, tagID)

	testCases := []struct {
		name          string
		tagID         uint32
		cacheValue    string
		expectedTag   *entities.Tag
		expectedError error
		setupMocks    func()
	}{
		{
			name:        "success from cache",
			tagID:       tagID,
			cacheValue:  `{"id":1,"name":"Cached Tag"}`,
			expectedTag: &entities.Tag{ID: 1, Name: "Cached Tag"},
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return(`{"id":1,"name":"Cached Tag"}`, nil).
					Times(1)
			},
		},
		{
			name:        "success from db",
			tagID:       tagID,
			expectedTag: tag,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetTagByID(gomock.Any(), tagID).
					Return(tag, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), cacheKey, gomock.Any(), getTagByIDTTL).
					Return(nil).
					Times(1)
			},
		},
		{
			name:        "cache ping error, success from db",
			tagID:       tagID,
			expectedTag: tag,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetTagByID(gomock.Any(), tagID).
					Return(tag, nil).
					Times(1)
			},
		},
		{
			name:          "db error",
			tagID:         tagID,
			expectedError: errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetTagByID(gomock.Any(), tagID).
					Return(nil, errors.New("db error")).
					Times(1)
			},
		},
		{
			name:        "set cache error",
			tagID:       tagID,
			expectedTag: tag,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2)

				useCasesMock.
					EXPECT().
					GetTagByID(gomock.Any(), tagID).
					Return(tag, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), cacheKey, gomock.Any(), getTagByIDTTL).
					Return(errors.New("set cache error")).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			result, err := decorator.GetTagByID(context.Background(), tc.tagID)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTag, result)
			}
		})
	}
}

func TestCacheDecorator_GetTicketByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	ticketID := uint64(1)
	ticket := &entities.Ticket{ID: ticketID, Name: "Test Ticket"}
	cacheKey := fmt.Sprintf("%s:%d", getTicketByIDPrefix, ticketID)

	testCases := []struct {
		name           string
		ticketID       uint64
		cacheValue     string
		expectedTicket *entities.Ticket
		expectedError  error
		setupMocks     func()
	}{
		{
			name:           "success from cache",
			ticketID:       ticketID,
			cacheValue:     `{"id":1,"name":"Cached Ticket"}`,
			expectedTicket: &entities.Ticket{ID: 1, Name: "Cached Ticket"},
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return(`{"id":1,"name":"Cached Ticket"}`, nil).
					Times(1)
			},
		},
		{
			name:           "success from db",
			ticketID:       ticketID,
			expectedTicket: ticket,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetTicketByID(gomock.Any(), ticketID).
					Return(ticket, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), cacheKey, gomock.Any(), getTicketByIDTTL).
					Return(nil).
					Times(1)
			},
		},
		{
			name:           "cache ping error, success from db",
			ticketID:       ticketID,
			expectedTicket: ticket,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetTicketByID(gomock.Any(), ticketID).
					Return(ticket, nil).
					Times(1)
			},
		},
		{
			name:          "db error",
			ticketID:      ticketID,
			expectedError: errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetTicketByID(gomock.Any(), ticketID).
					Return(nil, errors.New("db error")).
					Times(1)
			},
		},
		{
			name:           "set cache error",
			ticketID:       ticketID,
			expectedTicket: ticket,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				cacheMock.
					EXPECT().
					Get(gomock.Any(), cacheKey).
					Return("", errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2)

				useCasesMock.
					EXPECT().
					GetTicketByID(gomock.Any(), ticketID).
					Return(ticket, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), cacheKey, gomock.Any(), getTicketByIDTTL).
					Return(errors.New("set cache error")).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			result, err := decorator.GetTicketByID(context.Background(), tc.ticketID)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTicket, result)
			}
		})
	}
}

func TestCacheDecorator_GetTickets(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	pagination := &entities.Pagination{Offset: pointers.New[uint64](1), Limit: pointers.New[uint64](10)}
	filters := &entities.TicketsFilters{Search: pointers.New("ticket")}
	tickets := []entities.Ticket{{ID: 1, Name: "Test Ticket"}}

	testCases := []struct {
		name            string
		pagination      *entities.Pagination
		filters         *entities.TicketsFilters
		cacheValue      string
		expectedTickets []entities.Ticket
		expectedError   error
		setupMocks      func()
	}{
		{
			name:            "success from cache",
			pagination:      pagination,
			filters:         filters,
			cacheValue:      `[{"id":1,"name":"Cached Ticket"}]`,
			expectedTickets: []entities.Ticket{{ID: 1, Name: "Cached Ticket"}},
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return `[{"id":1,"name":"Cached Ticket"}]`, nil
						},
					),
				)
			},
		},
		{
			name:            "success from db",
			pagination:      pagination,
			filters:         filters,
			expectedTickets: tickets,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetTickets(gomock.Any(), pagination, filters).
					Return(tickets, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), getTicketsTTL).
					Return(nil).
					Times(1)
			},
		},
		{
			name:            "cache ping error, success from db",
			pagination:      pagination,
			filters:         filters,
			expectedTickets: tickets,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetTickets(gomock.Any(), pagination, filters).
					Return(tickets, nil).
					Times(1)
			},
		},
		{
			name:          "db error",
			pagination:    pagination,
			filters:       filters,
			expectedError: errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetTickets(gomock.Any(), pagination, filters).
					Return(nil, errors.New("db error")).
					Times(1)
			},
		},
		{
			name:            "set cache error",
			pagination:      pagination,
			filters:         filters,
			expectedTickets: tickets,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2)

				useCasesMock.
					EXPECT().
					GetTickets(gomock.Any(), pagination, filters).
					Return(tickets, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), getTicketsTTL).
					Return(errors.New("set cache error")).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			result, err := decorator.GetTickets(context.Background(), tc.pagination, tc.filters)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTickets, result)
			}
		})
	}
}

func TestCacheDecorator_GetUserTickets(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	userID := uint64(1)
	pagination := &entities.Pagination{Offset: pointers.New[uint64](1), Limit: pointers.New[uint64](10)}
	filters := &entities.TicketsFilters{Search: pointers.New("ticket")}
	tickets := []entities.Ticket{{ID: 1, Name: "Test Ticket", UserID: userID}}

	testCases := []struct {
		name            string
		userID          uint64
		pagination      *entities.Pagination
		filters         *entities.TicketsFilters
		cacheValue      string
		expectedTickets []entities.Ticket
		expectedError   error
		setupMocks      func()
	}{
		{
			name:            "success from cache",
			userID:          userID,
			pagination:      pagination,
			filters:         filters,
			cacheValue:      `[{"id":1,"name":"Cached Ticket","userId":1}]`,
			expectedTickets: []entities.Ticket{{ID: 1, Name: "Cached Ticket", UserID: userID}},
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return `[{"id":1,"name":"Cached Ticket","userId":1}]`, nil
						},
					),
				)
			},
		},
		{
			name:            "success from db",
			userID:          userID,
			pagination:      pagination,
			filters:         filters,
			expectedTickets: tickets,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetUserTickets(gomock.Any(), userID, pagination, filters).
					Return(tickets, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), getUserTicketsTTL).
					Return(nil).
					Times(1)
			},
		},
		{
			name:            "cache ping error, success from db",
			userID:          userID,
			pagination:      pagination,
			filters:         filters,
			expectedTickets: tickets,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetUserTickets(gomock.Any(), userID, pagination, filters).
					Return(tickets, nil).
					Times(1)
			},
		},
		{
			name:          "db error",
			userID:        userID,
			pagination:    pagination,
			filters:       filters,
			expectedError: errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)
				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					GetUserTickets(gomock.Any(), userID, pagination, filters).
					Return(nil, errors.New("db error")).
					Times(1)
			},
		},
		{
			name:            "set cache error",
			userID:          userID,
			pagination:      pagination,
			filters:         filters,
			expectedTickets: tickets,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2)

				useCasesMock.
					EXPECT().
					GetUserTickets(gomock.Any(), userID, pagination, filters).
					Return(tickets, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), getUserTicketsTTL).
					Return(errors.New("set cache error")).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			result, err := decorator.GetUserTickets(context.Background(), tc.userID, tc.pagination, tc.filters)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTickets, result)
			}
		})
	}
}

func TestCacheDecorator_CountTickets(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	filters := &entities.TicketsFilters{Search: pointers.New("ticket")}
	count := uint64(100)

	testCases := []struct {
		name          string
		filters       *entities.TicketsFilters
		cacheValue    string
		expectedCount uint64
		expectedError error
		setupMocks    func()
	}{
		{
			name:          "success from cache",
			filters:       filters,
			cacheValue:    "100",
			expectedCount: count,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "100", nil
						},
					),
				)
			},
		},
		{
			name:          "success from db",
			filters:       filters,
			expectedCount: count,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					CountTickets(gomock.Any(), filters).
					Return(count, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), countTicketsTTL).
					Return(nil).
					Times(1)
			},
		},
		{
			name:          "cache ping error, success from db",
			filters:       filters,
			expectedCount: count,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					CountTickets(gomock.Any(), filters).
					Return(count, nil).
					Times(1)
			},
		},
		{
			name:          "db error",
			filters:       filters,
			expectedError: errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					CountTickets(gomock.Any(), filters).
					Return(uint64(0), errors.New("db error")).
					Times(1)
			},
		},
		{
			name:          "set cache error",
			filters:       filters,
			expectedCount: count,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2)

				useCasesMock.
					EXPECT().
					CountTickets(gomock.Any(), filters).
					Return(count, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), countTicketsTTL).
					Return(errors.New("set cache error")).
					Times(1)
			},
		},
		{
			name:          "parse error, success from db",
			filters:       filters,
			cacheValue:    "invalid",
			expectedCount: count,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "invalid", nil
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					CountTickets(gomock.Any(), filters).
					Return(count, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), countTicketsTTL).
					Return(nil).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			result, err := decorator.CountTickets(context.Background(), tc.filters)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Equal(t, uint64(0), result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCount, result)
			}
		})
	}
}

func TestCacheDecorator_CountUserTickets(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	userID := uint64(1)
	filters := &entities.TicketsFilters{Search: pointers.New("ticket")}
	count := uint64(50)

	testCases := []struct {
		name          string
		userID        uint64
		filters       *entities.TicketsFilters
		cacheValue    string
		expectedCount uint64
		expectedError error
		setupMocks    func()
	}{
		{
			name:          "success from cache",
			userID:        userID,
			filters:       filters,
			cacheValue:    "50",
			expectedCount: count,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "50", nil
						},
					),
				)
			},
		},
		{
			name:          "success from db",
			userID:        userID,
			filters:       filters,
			expectedCount: count,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					CountUserTickets(gomock.Any(), userID, filters).
					Return(count, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), countUserTicketsTTL).
					Return(nil).
					Times(1)
			},
		},
		{
			name:          "cache ping error, success from db",
			userID:        userID,
			filters:       filters,
			expectedCount: count,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					CountUserTickets(gomock.Any(), userID, filters).
					Return(count, nil).
					Times(1)
			},
		},
		{
			name:          "db error",
			userID:        userID,
			filters:       filters,
			expectedError: errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					CountUserTickets(gomock.Any(), userID, filters).
					Return(uint64(0), errors.New("db error")).
					Times(1)
			},
		},
		{
			name:          "set cache error",
			userID:        userID,
			filters:       filters,
			expectedCount: count,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "", errors.New("not found")
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(2)

				useCasesMock.
					EXPECT().
					CountUserTickets(gomock.Any(), userID, filters).
					Return(count, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), countUserTicketsTTL).
					Return(errors.New("set cache error")).
					Times(1)
			},
		},
		{
			name:          "parse error, success from db",
			userID:        userID,
			filters:       filters,
			cacheValue:    "invalid",
			expectedCount: count,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				gomock.InOrder(
					cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(
						func(_ context.Context, key string) (string, error) {
							return "invalid", nil
						},
					),
				)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					CountUserTickets(gomock.Any(), userID, filters).
					Return(count, nil).
					Times(1)

				cacheMock.
					EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), countUserTicketsTTL).
					Return(nil).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			result, err := decorator.CountUserTickets(context.Background(), tc.userID, tc.filters)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Equal(t, uint64(0), result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCount, result)
			}
		})
	}
}

func TestCacheDecorator_UpdateUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	userProfileData := entities.RawUpdateUserProfileDTO{AccessToken: "token123"}
	user := &entities.User{ID: 1, Email: "test@example.com"}

	testCases := []struct {
		name            string
		userProfileData entities.RawUpdateUserProfileDTO
		expectedError   error
		setupMocks      func()
	}{
		{
			name:            "success with cache deletion",
			userProfileData: userProfileData,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					UpdateUserProfile(gomock.Any(), userProfileData).
					Return(nil).
					Times(1)

				useCasesMock.
					EXPECT().
					GetMe(gomock.Any(), userProfileData.AccessToken).
					Return(user, nil).
					Times(1)

				cacheMock.
					EXPECT().
					DelByPattern(gomock.Any(), fmt.Sprintf("%s:%d", getUserByIDPrefix, user.ID), nil).
					Return(nil).
					Times(1)

				cacheMock.
					EXPECT().
					DelByPattern(gomock.Any(), fmt.Sprintf("%s:%s", getUserByEmailPrefix, user.Email), nil).
					Return(nil).
					Times(1)
			},
		},
		{
			name:            "cache ping error, success from db",
			userProfileData: userProfileData,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					UpdateUserProfile(gomock.Any(), userProfileData).
					Return(nil).
					Times(1)
			},
		},
		{
			name:            "db error",
			userProfileData: userProfileData,
			expectedError:   errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					UpdateUserProfile(gomock.Any(), userProfileData).
					Return(errors.New("db error")).
					Times(1)
			},
		},
		{
			name:            "get me error, no cache deletion",
			userProfileData: userProfileData,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					UpdateUserProfile(gomock.Any(), userProfileData).
					Return(nil).
					Times(1)

				useCasesMock.
					EXPECT().
					GetMe(gomock.Any(), userProfileData.AccessToken).
					Return(nil, errors.New("get me error")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
		},
		{
			name:            "delete cache error",
			userProfileData: userProfileData,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					UpdateUserProfile(gomock.Any(), userProfileData).
					Return(nil).
					Times(1)

				useCasesMock.
					EXPECT().
					GetMe(gomock.Any(), userProfileData.AccessToken).
					Return(user, nil).
					Times(1)

				cacheMock.
					EXPECT().
					DelByPattern(gomock.Any(), fmt.Sprintf("%s:%d", getUserByIDPrefix, user.ID), nil).
					Return(errors.New("delete cache error")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				cacheMock.
					EXPECT().
					DelByPattern(gomock.Any(), fmt.Sprintf("%s:%s", getUserByEmailPrefix, user.Email), nil).
					Return(nil).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			err := decorator.UpdateUserProfile(context.Background(), tc.userProfileData)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCacheDecorator_UpdateToy(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	toyData := entities.RawUpdateToyDTO{ID: 1}
	toy := &entities.Toy{ID: 1, MasterID: 2}

	testCases := []struct {
		name          string
		toyData       entities.RawUpdateToyDTO
		expectedError error
		setupMocks    func()
	}{
		{
			name:    "success with cache deletion",
			toyData: toyData,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					UpdateToy(gomock.Any(), toyData).
					Return(nil).
					Times(1)

				useCasesMock.
					EXPECT().
					GetToyByID(gomock.Any(), toyData.ID).
					Return(toy, nil).
					Times(1)

				cacheMock.
					EXPECT().
					DelByPattern(gomock.Any(), fmt.Sprintf("%s:%d", getToyByIDPrefix, toy.ID), nil).
					Return(nil).
					Times(1)

				cacheMock.
					EXPECT().
					DelByPattern(gomock.Any(), fmt.Sprintf("%s:%d*", getMasterToysPrefix, toy.MasterID), nil).
					Return(nil).
					Times(1)

				cacheMock.
					EXPECT().
					DelByPattern(gomock.Any(), fmt.Sprintf("%s:%d*", countMasterToysPrefix, toy.MasterID), nil).
					Return(nil).
					Times(1)
			},
		},
		{
			name:    "cache ping error, success from db",
			toyData: toyData,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					UpdateToy(gomock.Any(), toyData).
					Return(nil).
					Times(1)
			},
		},
		{
			name:          "db error",
			toyData:       toyData,
			expectedError: errors.New("db error"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					UpdateToy(gomock.Any(), toyData).
					Return(errors.New("db error")).
					Times(1)
			},
		},
		{
			name:    "get toy error, no cache deletion",
			toyData: toyData,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					UpdateToy(gomock.Any(), toyData).
					Return(nil).
					Times(1)

				useCasesMock.
					EXPECT().
					GetToyByID(gomock.Any(), toyData.ID).
					Return(nil, errors.New("get toy error")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
		},
		{
			name:    "delete cache error",
			toyData: toyData,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					UpdateToy(gomock.Any(), toyData).
					Return(nil).
					Times(1)

				useCasesMock.
					EXPECT().
					GetToyByID(gomock.Any(), toyData.ID).
					Return(toy, nil).
					Times(1)

				cacheMock.
					EXPECT().
					DelByPattern(gomock.Any(), fmt.Sprintf("%s:%d", getToyByIDPrefix, toy.ID), nil).
					Return(errors.New("delete cache error")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				cacheMock.
					EXPECT().
					DelByPattern(gomock.Any(), fmt.Sprintf("%s:%d*", getMasterToysPrefix, toy.MasterID), nil).
					Return(nil).
					Times(1)

				cacheMock.
					EXPECT().
					DelByPattern(gomock.Any(), fmt.Sprintf("%s:%d*", countMasterToysPrefix, toy.MasterID), nil).
					Return(nil).
					Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			err := decorator.UpdateToy(context.Background(), tc.toyData)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCacheDecorator_DeleteToy(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	accessToken := "test-token"
	toyID := uint64(1)
	toy := &entities.Toy{ID: toyID, MasterID: 100, Name: "Test Toy"}
	patterns := []string{
		fmt.Sprintf("%s:%d", getToyByIDPrefix, toy.ID),
		fmt.Sprintf("%s:%d*", getMasterToysPrefix, toy.MasterID),
		fmt.Sprintf("%s:%d*", countMasterToysPrefix, toy.MasterID),
	}

	testCases := []struct {
		name          string
		accessToken   string
		toyID         uint64
		expectedError error
		setupMocks    func()
	}{
		{
			name:        "success with cache delete",
			accessToken: accessToken,
			toyID:       toyID,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					DeleteToy(gomock.Any(), accessToken, toyID).
					Return(nil).
					Times(1)

				useCasesMock.
					EXPECT().
					GetToyByID(gomock.Any(), toyID).
					Return(toy, nil).
					Times(1)

				for _, pattern := range patterns {
					cacheMock.
						EXPECT().
						DelByPattern(gomock.Any(), pattern, nil).
						Return(nil).
						Times(1)
				}
			},
		},
		{
			name:        "cache ping error",
			accessToken: accessToken,
			toyID:       toyID,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					DeleteToy(gomock.Any(), accessToken, toyID).
					Return(nil).
					Times(1)
			},
		},
		{
			name:          "delete toy error",
			accessToken:   accessToken,
			toyID:         toyID,
			expectedError: errors.New("delete failed"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					DeleteToy(gomock.Any(), accessToken, toyID).
					Return(errors.New("delete failed")).
					Times(1)
			},
		},
		{
			name:        "get toy error after delete",
			accessToken: accessToken,
			toyID:       toyID,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					DeleteToy(gomock.Any(), accessToken, toyID).
					Return(nil).
					Times(1)

				useCasesMock.
					EXPECT().
					GetToyByID(gomock.Any(), toyID).
					Return(nil, errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
		},
		{
			name:        "cache delete error",
			accessToken: accessToken,
			toyID:       toyID,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					DeleteToy(gomock.Any(), accessToken, toyID).
					Return(nil).
					Times(1)

				useCasesMock.
					EXPECT().
					GetToyByID(gomock.Any(), toyID).
					Return(toy, nil).
					Times(1)

				for _, pattern := range patterns {
					cacheMock.
						EXPECT().
						DelByPattern(gomock.Any(), pattern, nil).
						Return(errors.New("cache delete failed")).
						Times(1)

					loggerMock.
						EXPECT().
						ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
						Times(1)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			err := decorator.DeleteToy(context.Background(), tc.accessToken, tc.toyID)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCacheDecorator_UpdateMaster(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	masterID := uint64(1)
	userID := uint64(10)
	rawMasterData := entities.RawUpdateMasterDTO{ID: masterID}
	master := &entities.Master{ID: masterID, UserID: userID}
	patterns := []string{
		fmt.Sprintf("%s:%d", getMasterByIDPrefix, master.ID),
		fmt.Sprintf("%s:%d", getMasterByUserIDPrefix, master.UserID),
		fmt.Sprintf("%s:%d*", getMasterToysPrefix, master.ID),
		fmt.Sprintf("%s:%d*", countMasterToysPrefix, master.ID),
		fmt.Sprintf("%s*", getMastersPrefix),
	}

	testCases := []struct {
		name          string
		rawMasterData entities.RawUpdateMasterDTO
		expectedError error
		setupMocks    func()
	}{
		{
			name:          "success with cache delete",
			rawMasterData: rawMasterData,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					UpdateMaster(gomock.Any(), rawMasterData).
					Return(nil).
					Times(1)

				useCasesMock.
					EXPECT().
					GetMasterByID(gomock.Any(), masterID).
					Return(master, nil).
					Times(1)

				for _, pattern := range patterns {
					cacheMock.
						EXPECT().
						DelByPattern(gomock.Any(), pattern, nil).
						Return(nil).
						Times(1)
				}
			},
		},
		{
			name:          "cache ping error",
			rawMasterData: rawMasterData,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					UpdateMaster(gomock.Any(), rawMasterData).
					Return(nil).
					Times(1)
			},
		},
		{
			name:          "update master error",
			rawMasterData: rawMasterData,
			expectedError: errors.New("update failed"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					UpdateMaster(gomock.Any(), rawMasterData).
					Return(errors.New("update failed")).
					Times(1)
			},
		},
		{
			name:          "get master error after update",
			rawMasterData: rawMasterData,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					UpdateMaster(gomock.Any(), rawMasterData).
					Return(nil).
					Times(1)

				useCasesMock.
					EXPECT().
					GetMasterByID(gomock.Any(), masterID).
					Return(nil, errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
		},
		{
			name:          "cache delete error",
			rawMasterData: rawMasterData,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					UpdateMaster(gomock.Any(), rawMasterData).
					Return(nil).
					Times(1)

				useCasesMock.
					EXPECT().
					GetMasterByID(gomock.Any(), masterID).
					Return(master, nil).
					Times(1)

				for _, pattern := range patterns {
					cacheMock.
						EXPECT().
						DelByPattern(gomock.Any(), pattern, nil).
						Return(errors.New("cache delete failed")).
						Times(1)

					loggerMock.
						EXPECT().
						ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
						Times(1)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			err := decorator.UpdateMaster(context.Background(), tc.rawMasterData)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCacheDecorator_UpdateTicket(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	ticketID := uint64(1)
	userID := uint64(10)
	rawTicketData := entities.RawUpdateTicketDTO{ID: ticketID}
	ticket := &entities.Ticket{ID: ticketID, UserID: userID}
	patterns := []string{
		fmt.Sprintf("%s:%d", getTicketByIDPrefix, ticket.ID),
		fmt.Sprintf("%s:%d", getUserTicketsPrefix, ticket.UserID),
		fmt.Sprintf("%s:%d*", countUserTicketsPrefix, ticket.UserID),
	}

	testCases := []struct {
		name          string
		rawTicketData entities.RawUpdateTicketDTO
		expectedError error
		setupMocks    func()
	}{
		{
			name:          "success with cache delete",
			rawTicketData: rawTicketData,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					UpdateTicket(gomock.Any(), rawTicketData).
					Return(nil).
					Times(1)

				useCasesMock.
					EXPECT().
					GetTicketByID(gomock.Any(), ticketID).
					Return(ticket, nil).
					Times(1)

				for _, pattern := range patterns {
					cacheMock.
						EXPECT().
						DelByPattern(gomock.Any(), pattern, nil).
						Return(nil).
						Times(1)
				}
			},
		},
		{
			name:          "cache ping error",
			rawTicketData: rawTicketData,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					UpdateTicket(gomock.Any(), rawTicketData).
					Return(nil).
					Times(1)
			},
		},
		{
			name:          "update ticket error",
			rawTicketData: rawTicketData,
			expectedError: errors.New("update failed"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					UpdateTicket(gomock.Any(), rawTicketData).
					Return(errors.New("update failed")).
					Times(1)
			},
		},
		{
			name:          "get ticket error after update",
			rawTicketData: rawTicketData,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					UpdateTicket(gomock.Any(), rawTicketData).
					Return(nil).
					Times(1)

				useCasesMock.
					EXPECT().
					GetTicketByID(gomock.Any(), ticketID).
					Return(nil, errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
		},
		{
			name:          "cache delete error",
			rawTicketData: rawTicketData,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					UpdateTicket(gomock.Any(), rawTicketData).
					Return(nil).
					Times(1)

				useCasesMock.
					EXPECT().
					GetTicketByID(gomock.Any(), ticketID).
					Return(ticket, nil).
					Times(1)

				for _, pattern := range patterns {
					cacheMock.
						EXPECT().
						DelByPattern(gomock.Any(), pattern, nil).
						Return(errors.New("cache delete failed")).
						Times(1)

					loggerMock.
						EXPECT().
						ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
						Times(1)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			err := decorator.UpdateTicket(context.Background(), tc.rawTicketData)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCacheDecorator_DeleteTicket(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheMock := mockcache.NewMockProvider(ctrl)
	loggerMock := mocklogging.NewMockLogger(ctrl)
	useCasesMock := mockusecases.NewMockUseCases(ctrl)
	decorator := NewCacheDecorator(useCasesMock, cacheMock, loggerMock)

	accessToken := "test-token"
	ticketID := uint64(1)
	ticket := &entities.Ticket{ID: ticketID, UserID: 10}
	patterns := []string{
		fmt.Sprintf("%s:%d", getTicketByIDPrefix, ticket.ID),
		fmt.Sprintf("%s:%d", getUserTicketsPrefix, ticket.UserID),
		fmt.Sprintf("%s:%d*", countUserTicketsPrefix, ticket.UserID),
	}

	testCases := []struct {
		name          string
		accessToken   string
		ticketID      uint64
		expectedError error
		setupMocks    func()
	}{
		{
			name:        "success with cache delete",
			accessToken: accessToken,
			ticketID:    ticketID,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					DeleteTicket(gomock.Any(), accessToken, ticketID).
					Return(nil).
					Times(1)

				useCasesMock.
					EXPECT().
					GetTicketByID(gomock.Any(), ticketID).
					Return(ticket, nil).
					Times(1)

				for _, pattern := range patterns {
					cacheMock.
						EXPECT().
						DelByPattern(gomock.Any(), pattern, nil).
						Return(nil).
						Times(1)
				}
			},
		},
		{
			name:        "cache ping error",
			accessToken: accessToken,
			ticketID:    ticketID,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", errors.New("cache unavailable")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				useCasesMock.
					EXPECT().
					DeleteTicket(gomock.Any(), accessToken, ticketID).
					Return(nil).
					Times(1)
			},
		},
		{
			name:          "delete ticket error",
			accessToken:   accessToken,
			ticketID:      ticketID,
			expectedError: errors.New("delete failed"),
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					DeleteTicket(gomock.Any(), accessToken, ticketID).
					Return(errors.New("delete failed")).
					Times(1)
			},
		},
		{
			name:        "get ticket error after delete",
			accessToken: accessToken,
			ticketID:    ticketID,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					DeleteTicket(gomock.Any(), accessToken, ticketID).
					Return(nil).
					Times(1)

				useCasesMock.
					EXPECT().
					GetTicketByID(gomock.Any(), ticketID).
					Return(nil, errors.New("not found")).
					Times(1)

				loggerMock.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
		},
		{
			name:        "cache delete error",
			accessToken: accessToken,
			ticketID:    ticketID,
			setupMocks: func() {
				cacheMock.
					EXPECT().
					Ping(gomock.Any()).
					Return("", nil).
					Times(1)

				useCasesMock.
					EXPECT().
					DeleteTicket(gomock.Any(), accessToken, ticketID).
					Return(nil).
					Times(1)

				useCasesMock.
					EXPECT().
					GetTicketByID(gomock.Any(), ticketID).
					Return(ticket, nil).
					Times(1)

				for _, pattern := range patterns {
					cacheMock.
						EXPECT().
						DelByPattern(gomock.Any(), pattern, nil).
						Return(errors.New("cache delete failed")).
						Times(1)

					loggerMock.
						EXPECT().
						ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
						Times(1)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			err := decorator.DeleteTicket(context.Background(), tc.accessToken, tc.ticketID)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
