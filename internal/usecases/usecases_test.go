package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	customerrors "github.com/DKhorkov/hmtm-bff/internal/errors"
	"github.com/DKhorkov/libs/pointers"
	"github.com/DKhorkov/libs/security"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	mockservices "github.com/DKhorkov/hmtm-bff/mocks/services"
	mocklogger "github.com/DKhorkov/libs/logging/mocks"
	tracingmock "github.com/DKhorkov/libs/tracing/mocks"

	"github.com/DKhorkov/hmtm-bff/internal/config"
	"github.com/DKhorkov/hmtm-bff/internal/entities"
)

var (
	ctx              = context.Background()
	validationConfig = config.ValidationConfig{
		FileMaxSize: int64(5 * 1024 * 1024), // 5 Mb
		FileAllowedExtensions: []string{
			".png",
			".svg",
			".gif",
			".jpg",
			".jpeg",
			".jfif",
			".pjpeg",
			".pjp",
		},
	}
)

func TestUseCases_RegisterUser(t *testing.T) {
	testCases := []struct {
		name       string
		userData   entities.RegisterUserDTO
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		expected      uint64
		errorExpected bool
	}{
		{
			name: "success",
			userData: entities.RegisterUserDTO{
				DisplayName: "test",
				Email:       "test@test.com",
				Password:    "test",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					RegisterUser(
						gomock.Any(),
						entities.RegisterUserDTO{
							DisplayName: "test",
							Email:       "test@test.com",
							Password:    "test",
						},
					).
					Return(uint64(1), nil).
					Times(1)
			},
			expected: uint64(1),
		},
		{
			name: "error",
			userData: entities.RegisterUserDTO{
				DisplayName: "test",
				Email:       "test@test.com",
				Password:    "test",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					RegisterUser(
						gomock.Any(),
						entities.RegisterUserDTO{
							DisplayName: "test",
							Email:       "test@test.com",
							Password:    "test",
						},
					).
					Return(uint64(0), errors.New("test")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			actual, err := useCases.RegisterUser(ctx, tc.userData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_LoginUser(t *testing.T) {
	testCases := []struct {
		name       string
		userData   entities.LoginUserDTO
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		expected      *entities.TokensDTO
		errorExpected bool
	}{
		{
			name: "success",
			userData: entities.LoginUserDTO{
				Email:    "test@test.com",
				Password: "test",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				tokens := &entities.TokensDTO{
					AccessToken:  "access_token",
					RefreshToken: "refresh_token",
				}
				ssoService.
					EXPECT().
					LoginUser(
						gomock.Any(),
						entities.LoginUserDTO{
							Email:    "test@test.com",
							Password: "test",
						},
					).
					Return(tokens, nil).
					Times(1)
			},
			expected: &entities.TokensDTO{
				AccessToken:  "access_token",
				RefreshToken: "refresh_token",
			},
		},
		{
			name: "invalid credentials",
			userData: entities.LoginUserDTO{
				Email:    "wrong@test.com",
				Password: "wrong",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					LoginUser(
						gomock.Any(),
						entities.LoginUserDTO{
							Email:    "wrong@test.com",
							Password: "wrong",
						},
					).
					Return(nil, errors.New("invalid credentials")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			actual, err := useCases.LoginUser(ctx, tc.userData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_LogoutUser(t *testing.T) {
	testCases := []struct {
		name        string
		accessToken string
		setupMocks  func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		errorExpected bool
	}{
		{
			name:        "success",
			accessToken: "valid_access_token",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					LogoutUser(
						gomock.Any(),
						"valid_access_token",
					).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:        "invalid token",
			accessToken: "invalid_access_token",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					LogoutUser(
						gomock.Any(),
						"invalid_access_token",
					).
					Return(errors.New("invalid token")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			err := useCases.LogoutUser(ctx, tc.accessToken)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUseCases_VerifyUserEmail(t *testing.T) {
	testCases := []struct {
		name             string
		verifyEmailToken string
		setupMocks       func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		errorExpected bool
	}{
		{
			name:             "success",
			verifyEmailToken: "valid_verify_token",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					VerifyUserEmail(
						gomock.Any(),
						"valid_verify_token",
					).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:             "invalid token",
			verifyEmailToken: "invalid_verify_token",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					VerifyUserEmail(
						gomock.Any(),
						"invalid_verify_token",
					).
					Return(errors.New("invalid or expired token")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			err := useCases.VerifyUserEmail(ctx, tc.verifyEmailToken)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUseCases_SendVerifyEmailMessage(t *testing.T) {
	testCases := []struct {
		name       string
		email      string
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		errorExpected bool
	}{
		{
			name:  "success",
			email: "test@example.com",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					SendVerifyEmailMessage(
						gomock.Any(),
						"test@example.com",
					).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:  "invalid email",
			email: "invalid-email",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					SendVerifyEmailMessage(
						gomock.Any(),
						"invalid-email",
					).
					Return(errors.New("invalid email format")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			err := useCases.SendVerifyEmailMessage(ctx, tc.email)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUseCases_SendForgetPasswordMessage(t *testing.T) {
	testCases := []struct {
		name       string
		email      string
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		errorExpected bool
	}{
		{
			name:  "success",
			email: "test@example.com",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					SendForgetPasswordMessage(
						gomock.Any(),
						"test@example.com",
					).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:  "user not found",
			email: "nonexistent@example.com",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					SendForgetPasswordMessage(
						gomock.Any(),
						"nonexistent@example.com",
					).
					Return(errors.New("user not found")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			err := useCases.SendForgetPasswordMessage(ctx, tc.email)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUseCases_ChangePassword(t *testing.T) {
	testCases := []struct {
		name        string
		accessToken string
		oldPassword string
		newPassword string
		setupMocks  func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		errorExpected bool
	}{
		{
			name:        "success",
			accessToken: "valid_access_token",
			oldPassword: "old_pass_123",
			newPassword: "new_pass_456",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					ChangePassword(
						gomock.Any(),
						"valid_access_token",
						"old_pass_123",
						"new_pass_456",
					).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:        "invalid token",
			accessToken: "invalid_access_token",
			oldPassword: "old_pass_123",
			newPassword: "new_pass_456",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					ChangePassword(
						gomock.Any(),
						"invalid_access_token",
						"old_pass_123",
						"new_pass_456",
					).
					Return(errors.New("invalid or expired token")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name:        "wrong old password",
			accessToken: "valid_access_token",
			oldPassword: "wrong_old_pass",
			newPassword: "new_pass_456",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					ChangePassword(
						gomock.Any(),
						"valid_access_token",
						"wrong_old_pass",
						"new_pass_456",
					).
					Return(errors.New("incorrect old password")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			err := useCases.ChangePassword(ctx, tc.accessToken, tc.oldPassword, tc.newPassword)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUseCases_ForgetPassword(t *testing.T) {
	testCases := []struct {
		name                string
		forgetPasswordToken string
		newPassword         string
		setupMocks          func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		errorExpected bool
	}{
		{
			name:                "success",
			forgetPasswordToken: "valid_reset_token",
			newPassword:         "new_pass_123",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					ForgetPassword(
						gomock.Any(),
						"valid_reset_token",
						"new_pass_123",
					).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:                "invalid token",
			forgetPasswordToken: "invalid_reset_token",
			newPassword:         "new_pass_123",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					ForgetPassword(
						gomock.Any(),
						"invalid_reset_token",
						"new_pass_123",
					).
					Return(errors.New("invalid or expired token")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			err := useCases.ForgetPassword(ctx, tc.forgetPasswordToken, tc.newPassword)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUseCases_GetMe(t *testing.T) {
	testCases := []struct {
		name        string
		accessToken string
		setupMocks  func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		expected      *entities.User
		errorExpected bool
	}{
		{
			name:        "success",
			accessToken: "valid_access_token",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{
					ID:          1,
					DisplayName: "Test User",
					Email:       "test@example.com",
				}
				ssoService.
					EXPECT().
					GetMe(
						gomock.Any(),
						"valid_access_token",
					).
					Return(user, nil).
					Times(1)
			},
			expected: &entities.User{
				ID:          1,
				DisplayName: "Test User",
				Email:       "test@example.com",
			},
			errorExpected: false,
		},
		{
			name:        "invalid token",
			accessToken: "invalid_access_token",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					GetMe(
						gomock.Any(),
						"invalid_access_token",
					).
					Return(nil, errors.New("invalid or expired token")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			actual, err := useCases.GetMe(ctx, tc.accessToken)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_RefreshTokens(t *testing.T) {
	testCases := []struct {
		name         string
		refreshToken string
		setupMocks   func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		expected      *entities.TokensDTO
		errorExpected bool
	}{
		{
			name:         "success",
			refreshToken: "valid_refresh_token",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				tokens := &entities.TokensDTO{
					AccessToken:  "new_access_token",
					RefreshToken: "new_refresh_token",
				}
				ssoService.
					EXPECT().
					RefreshTokens(
						gomock.Any(),
						"valid_refresh_token",
					).
					Return(tokens, nil).
					Times(1)
			},
			expected: &entities.TokensDTO{
				AccessToken:  "new_access_token",
				RefreshToken: "new_refresh_token",
			},
			errorExpected: false,
		},
		{
			name:         "invalid refresh token",
			refreshToken: "invalid_refresh_token",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					RefreshTokens(
						gomock.Any(),
						"invalid_refresh_token",
					).
					Return(nil, errors.New("invalid or expired refresh token")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			actual, err := useCases.RefreshTokens(ctx, tc.refreshToken)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetUserByID(t *testing.T) {
	testCases := []struct {
		name       string
		id         uint64
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		expected      *entities.User
		errorExpected bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{
					ID:          1,
					DisplayName: "Test User",
					Email:       "test@example.com",
				}
				ssoService.
					EXPECT().
					GetUserByID(
						gomock.Any(),
						uint64(1),
					).
					Return(user, nil).
					Times(1)
			},
			expected: &entities.User{
				ID:          1,
				DisplayName: "Test User",
				Email:       "test@example.com",
			},
			errorExpected: false,
		},
		{
			name: "user not found",
			id:   999,
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					GetUserByID(
						gomock.Any(),
						uint64(999),
					).
					Return(nil, errors.New("user not found")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			actual, err := useCases.GetUserByID(ctx, tc.id)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetUserByEmail(t *testing.T) {
	testCases := []struct {
		name       string
		email      string
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		expected      *entities.User
		errorExpected bool
	}{
		{
			name:  "success",
			email: "test@example.com",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{
					ID:          1,
					DisplayName: "Test User",
					Email:       "test@example.com",
				}
				ssoService.
					EXPECT().
					GetUserByEmail(
						gomock.Any(),
						"test@example.com",
					).
					Return(user, nil).
					Times(1)
			},
			expected: &entities.User{
				ID:          1,
				DisplayName: "Test User",
				Email:       "test@example.com",
			},
			errorExpected: false,
		},
		{
			name:  "user not found",
			email: "nonexistent@example.com",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					GetUserByEmail(
						gomock.Any(),
						"nonexistent@example.com",
					).
					Return(nil, errors.New("user not found")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			actual, err := useCases.GetUserByEmail(ctx, tc.email)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetAllUsers(t *testing.T) {
	testCases := []struct {
		name       string
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		expected      []entities.User
		errorExpected bool
	}{
		{
			name: "success with users",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				users := []entities.User{
					{
						ID:          1,
						DisplayName: "User One",
						Email:       "user1@example.com",
					},
					{
						ID:          2,
						DisplayName: "User Two",
						Email:       "user2@example.com",
					},
				}
				ssoService.
					EXPECT().
					GetAllUsers(
						gomock.Any(),
					).
					Return(users, nil).
					Times(1)
			},
			expected: []entities.User{
				{
					ID:          1,
					DisplayName: "User One",
					Email:       "user1@example.com",
				},
				{
					ID:          2,
					DisplayName: "User Two",
					Email:       "user2@example.com",
				},
			},
			errorExpected: false,
		},
		{
			name: "success with empty list",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					GetAllUsers(
						gomock.Any(),
					).
					Return([]entities.User{}, nil).
					Times(1)
			},
			expected:      []entities.User{},
			errorExpected: false,
		},
		{
			name: "error from service",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					GetAllUsers(gomock.Any()).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			actual, err := useCases.GetAllUsers(ctx)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_AddToy(t *testing.T) {
	testCases := []struct {
		name       string
		rawToyData entities.RawAddToyDTO
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		expected      uint64
		errorExpected bool
	}{
		{
			name: "success",
			rawToyData: entities.RawAddToyDTO{
				AccessToken: "valid_access_token",
				CategoryID:  1,
				Name:        "Test Toy",
				Description: "Test Description",
				Price:       100,
				Quantity:    10,
				Tags:        []string{"tag1", "tag2"},
				Attachments: []*graphql.Upload{
					{File: strings.NewReader("test content"), Filename: "file1.jpg"},
				},
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}
				tagIDs := []uint32{1, 2}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("uploaded/file1.jpg", nil).
					Times(1)

				toysService.
					EXPECT().
					CreateTags(gomock.Any(), []entities.CreateTagDTO{
						{Name: "tag1"},
						{Name: "tag2"},
					}).
					Return(tagIDs, nil).
					Times(1)

				toysService.
					EXPECT().
					AddToy(
						gomock.Any(),
						entities.AddToyDTO{
							UserID:      1,
							CategoryID:  1,
							Name:        "Test Toy",
							Description: "Test Description",
							Price:       100,
							Quantity:    10,
							TagIDs:      tagIDs,
							Attachments: []string{"uploaded/file1.jpg"},
						},
					).
					Return(uint64(123), nil).
					Times(1)
			},
			expected:      123,
			errorExpected: false,
		},
		{
			name: "invalid access token",
			rawToyData: entities.RawAddToyDTO{
				AccessToken: "invalid_access_token",
				CategoryID:  1,
				Name:        "Test Toy",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "invalid_access_token").
					Return(nil, errors.New("invalid token")).
					Times(1)
			},
			expected:      0,
			errorExpected: true,
		},
		{
			name: "file upload error",
			rawToyData: entities.RawAddToyDTO{
				AccessToken: "valid_access_token",
				CategoryID:  1,
				Name:        "Test Toy",
				Attachments: []*graphql.Upload{
					{File: strings.NewReader("test content"), Filename: "file1.jpg"},
				},
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				logger *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", errors.New("test")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expected:      0,
			errorExpected: true,
		},
		{
			name: "create tags error",
			rawToyData: entities.RawAddToyDTO{
				AccessToken: "valid_access_token",
				CategoryID:  1,
				Name:        "Test Toy",
				Description: "Test Description",
				Price:       100,
				Quantity:    10,
				Tags:        []string{"tag1", "tag2"},
				Attachments: []*graphql.Upload{
					{File: strings.NewReader("test content"), Filename: "file1.jpg"},
				},
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("uploaded/file1.jpg", nil).
					Times(1)

				toysService.
					EXPECT().
					CreateTags(gomock.Any(), []entities.CreateTagDTO{
						{Name: "tag1"},
						{Name: "tag2"},
					}).
					Return(nil, errors.New("test")).
					Times(1)

			},
			expected:      0,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			actual, err := useCases.AddToy(ctx, tc.rawToyData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_UploadFiles(t *testing.T) {
	testCases := []struct {
		name       string
		userID     uint64
		files      []*graphql.Upload
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		expected      []string
		errorExpected bool
	}{
		{
			name:   "success all files uploaded",
			userID: 1,
			files: []*graphql.Upload{
				{File: strings.NewReader("test content"), Filename: "file1.jpg"},
				{File: strings.NewReader("test content"), Filename: "file2.jpg"},
			},
			setupMocks: func(
				_ *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				logger *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("uploaded/file1.jpg", nil).
					Times(1)

				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("uploaded/file2.jpg", nil).
					Times(1)
			},
			expected:      []string{"uploaded/file1.jpg", "uploaded/file2.jpg"},
			errorExpected: false,
		},
		{
			name:   "partial success",
			userID: 1,
			files: []*graphql.Upload{
				{File: strings.NewReader("test content"), Filename: "file1.jpg"},
				{File: strings.NewReader("test content"), Filename: "file2.jpg"},
			},
			setupMocks: func(
				_ *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				logger *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("uploaded/file1.jpg", nil).
					Times(1)

				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", errors.New("upload failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expected:      []string{"uploaded/file1.jpg"},
			errorExpected: false,
		},
		{
			name:   "all files failed",
			userID: 1,
			files: []*graphql.Upload{
				{File: strings.NewReader("test content"), Filename: "file1.jpg"},
				{File: strings.NewReader("test content"), Filename: "file2.jpg"},
			},
			setupMocks: func(
				_ *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				logger *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", errors.New("upload failed")).
					AnyTimes()

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			actual, err := useCases.UploadFiles(ctx, tc.userID, tc.files)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_UploadFile(t *testing.T) {
	testCases := []struct {
		name       string
		userID     uint64
		file       *graphql.Upload
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		expected      string
		errorExpected bool
	}{
		{
			name:   "success",
			userID: 1,
			file: &graphql.Upload{
				File:     strings.NewReader("test content"),
				Filename: "test.jpg",
			},
			setupMocks: func(
				_ *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("uploaded/test.jpg", nil).
					Times(1)
			},
			expected:      "uploaded/test.jpg",
			errorExpected: false,
		},
		{
			name:   "create filename error",
			userID: 1,
			file: &graphql.Upload{
				File:     strings.NewReader("test content"),
				Filename: "invalid/file/name",
			},
			expected:      "",
			errorExpected: true,
		},
		{
			name:   "upload error",
			userID: 1,
			file: &graphql.Upload{
				File:     strings.NewReader("test content"),
				Filename: "test.jpg",
			},
			setupMocks: func(
				_ *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", errors.New("upload failed")).
					Times(1)
			},
			expected:      "",
			errorExpected: true,
		},
		{
			name:   "read file error",
			userID: 1,
			file: &graphql.Upload{
				File:     &errorReader{err: errors.New("read error")},
				Filename: "test.jpg",
				Size:     1024,
			},
			expected:      "",
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			actual, err := useCases.UploadFile(ctx, tc.userID, tc.file)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

// errorReader - вспомогательная структура для имитации ошибки чтения
type errorReader struct {
	err error
}

func (r *errorReader) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}

func TestUseCases_createFilename(t *testing.T) {
	testCases := []struct {
		name             string
		userID           uint64
		file             *graphql.Upload
		validationConfig config.ValidationConfig
		expected         string
		errorExpected    bool
		expectedError    error
	}{
		{
			name:   "success",
			userID: 1,
			file: &graphql.Upload{
				Filename: "test.jpg",
				Size:     1024,
			},
			validationConfig: validationConfig,
			expected:         security.RawEncode([]byte("1:test.jpg")) + ".jpg",
			errorExpected:    false,
		},
		{
			name:   "invalid extension",
			userID: 1,
			file: &graphql.Upload{
				Filename: "test.exe",
				Size:     1024,
			},
			validationConfig: validationConfig,
			expected:         "",
			errorExpected:    true,
			expectedError:    &customerrors.InvalidFileExtensionError{},
		},
		{
			name:   "file too large",
			userID: 1,
			file: &graphql.Upload{
				Filename: "test.jpg",
				Size:     8 * 1042 * 1024,
			},
			validationConfig: validationConfig,
			expected:         "",
			errorExpected:    true,
			expectedError:    &customerrors.InvalidFileSizeError{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ssoService := mockservices.NewMockSsoService(ctrl)
			toysService := mockservices.NewMockToysService(ctrl)
			ticketsService := mockservices.NewMockTicketsService(ctrl)
			notificationsService := mockservices.NewMockNotificationsService(ctrl)
			fileStorageService := mockservices.NewMockFileStorageService(ctrl)
			logger := mocklogger.NewMockLogger(ctrl)
			traceProvider := tracingmock.NewMockProvider(ctrl)

			useCases := New(
				ssoService,
				toysService,
				fileStorageService,
				ticketsService,
				notificationsService,
				tc.validationConfig,
				logger,
				traceProvider,
			)

			actual, err := useCases.createFilename(tc.userID, tc.file)
			if tc.errorExpected {
				require.Error(t, err)
				require.IsType(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetAllToys(t *testing.T) {
	testCases := []struct {
		name       string
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		expected      []entities.Toy
		errorExpected bool
	}{
		{
			name: "success with toys",
			setupMocks: func(
				_ *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				toys := []entities.Toy{
					{
						ID:          1,
						CategoryID:  1,
						Name:        "Toy One",
						Description: "Description One",
						Price:       100,
						Quantity:    10,
					},
					{
						ID:          2,
						CategoryID:  2,
						Name:        "Toy Two",
						Description: "Description Two",
						Price:       200,
						Quantity:    5,
					},
				}

				toysService.
					EXPECT().
					GetAllToys(gomock.Any()).
					Return(toys, nil).
					Times(1)
			},
			expected: []entities.Toy{
				{
					ID:          1,
					CategoryID:  1,
					Name:        "Toy One",
					Description: "Description One",
					Price:       100,
					Quantity:    10,
				},
				{
					ID:          2,
					CategoryID:  2,
					Name:        "Toy Two",
					Description: "Description Two",
					Price:       200,
					Quantity:    5,
				},
			},
			errorExpected: false,
		},
		{
			name: "success with empty list",
			setupMocks: func(
				_ *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				toysService.
					EXPECT().
					GetAllToys(gomock.Any()).
					Return([]entities.Toy{}, nil).
					Times(1)
			},
			expected:      []entities.Toy{},
			errorExpected: false,
		},
		{
			name: "error from service",
			setupMocks: func(
				_ *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				toysService.
					EXPECT().
					GetAllToys(gomock.Any()).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			actual, err := useCases.GetAllToys(ctx)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetMasterToys(t *testing.T) {
	testCases := []struct {
		name       string
		masterID   uint64
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		expected      []entities.Toy
		errorExpected bool
	}{
		{
			name:     "success with toys",
			masterID: 1,
			setupMocks: func(
				_ *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				toys := []entities.Toy{
					{
						ID:          1,
						CategoryID:  1,
						Name:        "Toy One",
						Description: "Description One",
						Price:       100,
						Quantity:    10,
					},
					{
						ID:          2,
						CategoryID:  2,
						Name:        "Toy Two",
						Description: "Description Two",
						Price:       200,
						Quantity:    5,
					},
				}

				toysService.
					EXPECT().
					GetMasterToys(gomock.Any(), uint64(1)).
					Return(toys, nil).
					Times(1)
			},
			expected: []entities.Toy{
				{
					ID:          1,
					CategoryID:  1,
					Name:        "Toy One",
					Description: "Description One",
					Price:       100,
					Quantity:    10,
				},
				{
					ID:          2,
					CategoryID:  2,
					Name:        "Toy Two",
					Description: "Description Two",
					Price:       200,
					Quantity:    5,
				},
			},
			errorExpected: false,
		},
		{
			name:     "success with empty list",
			masterID: 2,
			setupMocks: func(
				_ *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				toysService.
					EXPECT().
					GetMasterToys(gomock.Any(), uint64(2)).
					Return([]entities.Toy{}, nil).
					Times(1)
			},
			expected:      []entities.Toy{},
			errorExpected: false,
		},
		{
			name:     "error from service",
			masterID: 3,
			setupMocks: func(
				_ *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				toysService.
					EXPECT().
					GetMasterToys(gomock.Any(), uint64(3)).
					Return(nil, errors.New("toys not found")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			actual, err := useCases.GetMasterToys(ctx, tc.masterID)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetMyToys(t *testing.T) {
	testCases := []struct {
		name        string
		accessToken string
		setupMocks  func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		expected      []entities.Toy
		errorExpected bool
	}{
		{
			name:        "success with toys",
			accessToken: "valid_access_token",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}
				toys := []entities.Toy{
					{
						ID:          1,
						CategoryID:  1,
						Name:        "My Toy One",
						Description: "Description One",
						Price:       100,
						Quantity:    10,
					},
					{
						ID:          2,
						CategoryID:  2,
						Name:        "My Toy Two",
						Description: "Description Two",
						Price:       200,
						Quantity:    5,
					},
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetUserToys(gomock.Any(), uint64(1)).
					Return(toys, nil).
					Times(1)
			},
			expected: []entities.Toy{
				{
					ID:          1,
					CategoryID:  1,
					Name:        "My Toy One",
					Description: "Description One",
					Price:       100,
					Quantity:    10,
				},
				{
					ID:          2,
					CategoryID:  2,
					Name:        "My Toy Two",
					Description: "Description Two",
					Price:       200,
					Quantity:    5,
				},
			},
			errorExpected: false,
		},
		{
			name:        "success with empty list",
			accessToken: "valid_access_token",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 2}
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetUserToys(gomock.Any(), uint64(2)).
					Return([]entities.Toy{}, nil).
					Times(1)
			},
			expected:      []entities.Toy{},
			errorExpected: false,
		},
		{
			name:        "invalid access token",
			accessToken: "invalid_access_token",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "invalid_access_token").
					Return(nil, errors.New("invalid token")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
		{
			name:        "error from toys service",
			accessToken: "valid_access_token",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 3}
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetUserToys(gomock.Any(), uint64(3)).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			actual, err := useCases.GetMyToys(ctx, tc.accessToken)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetToyByID(t *testing.T) {
	testCases := []struct {
		name       string
		id         uint64
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		expected      *entities.Toy
		errorExpected bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(
				_ *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				toy := &entities.Toy{
					ID:          1,
					CategoryID:  1,
					Name:        "Test Toy",
					Description: "Test Description",
					Price:       100,
					Quantity:    10,
				}

				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), uint64(1)).
					Return(toy, nil).
					Times(1)
			},
			expected: &entities.Toy{
				ID:          1,
				CategoryID:  1,
				Name:        "Test Toy",
				Description: "Test Description",
				Price:       100,
				Quantity:    10,
			},
			errorExpected: false,
		},
		{
			name: "toy not found",
			id:   999,
			setupMocks: func(
				_ *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), uint64(999)).
					Return(nil, errors.New("toy not found")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			actual, err := useCases.GetToyByID(ctx, tc.id)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_DeleteTicket(t *testing.T) {
	type args struct {
		ctx         context.Context
		accessToken string
		id          uint64
	}

	// Test data
	testUser := &entities.User{
		ID:          1,
		DisplayName: "Test User",
		Email:       "test@example.com",
	}

	testRawTicket := &entities.RawTicket{
		ID:          1,
		UserID:      1,
		CategoryID:  1,
		Name:        "Test Ticket",
		Description: "Test Description",
		Quantity:    1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		TagIDs:      []uint32{1, 2, 3},
		Attachments: []entities.TicketAttachment{
			{
				Link: "test/file1.jpg",
			},
			{
				Link: "test/file2.jpg",
			},
		},
	}

	otherTestRawTicket := &entities.RawTicket{
		ID:          2,
		UserID:      2,
		CategoryID:  2,
		Name:        "Test Ticket",
		Description: "Test Description",
		Quantity:    1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		TagIDs:      []uint32{1, 2, 3},
	}

	otherUserTicket := &entities.Ticket{
		ID:     2,
		UserID: 2,
	}

	testCases := []struct {
		name       string
		args       args
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			ticketsService *mockservices.MockTicketsService,
			fileStorageService *mockservices.MockFileStorageService,
			toysService *mockservices.MockToysService,
			logger *mocklogger.MockLogger,
		)
		expectedError error
	}{
		{
			name: "successful deletion with attachments",
			args: args{
				ctx:         context.Background(),
				accessToken: "valid_token",
				id:          1,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				fileStorageService *mockservices.MockFileStorageService,
				toysService *mockservices.MockToysService,
				_ *mocklogger.MockLogger,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
						{ID: 3, Name: "Tag3"},
					}, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(1)).
					Return(testRawTicket, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					DeleteMany(gomock.Any(), []string{"file1.jpg", "file2.jpg"}).
					Return(nil).
					Times(1)

				ticketsService.
					EXPECT().
					DeleteTicket(gomock.Any(), uint64(1)).
					Return(nil).
					Times(1)
			},
			expectedError: nil,
		},
		{
			name: "invalid access token",
			args: args{
				ctx:         context.Background(),
				accessToken: "invalid_token",
				id:          1,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockToysService,
				_ *mocklogger.MockLogger,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "invalid_token").
					Return(nil, errors.New("invalid token")).
					Times(1)
			},
			expectedError: errors.New("invalid token"),
		},
		{
			name: "ticket not found",
			args: args{
				ctx:         context.Background(),
				accessToken: "valid_token",
				id:          999,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockToysService,
				_ *mocklogger.MockLogger,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(999)).
					Return(nil, errors.New("ticket not found")).
					Times(1)
			},
			expectedError: errors.New("ticket not found"),
		},
		{
			name: "permission denied - not owner",
			args: args{
				ctx:         context.Background(),
				accessToken: "valid_token",
				id:          2,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				_ *mockservices.MockFileStorageService,
				toysService *mockservices.MockToysService,
				_ *mocklogger.MockLogger,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(2)).
					Return(otherTestRawTicket, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
						{ID: 3, Name: "Tag3"},
					}, nil).
					Times(1)
			},
			expectedError: &customerrors.PermissionDeniedError{
				Message: fmt.Sprintf(
					"User with ID=%d is not owner of Ticket with ID=%d",
					testUser.ID,
					otherUserTicket.ID,
				),
			},
		},
		{
			name: "file deletion errors but ticket still deleted",
			args: args{
				ctx:         context.Background(),
				accessToken: "valid_token",
				id:          1,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				fileStorageService *mockservices.MockFileStorageService,
				toysService *mockservices.MockToysService,
				logger *mocklogger.MockLogger,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(1)).
					Return(testRawTicket, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
						{ID: 3, Name: "Tag3"},
					}, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					DeleteMany(gomock.Any(), []string{"file1.jpg", "file2.jpg"}).
					Return([]error{errors.New("file1 deletion failed"), nil}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				ticketsService.
					EXPECT().
					DeleteTicket(gomock.Any(), uint64(1)).
					Return(nil).
					Times(1)
			},
			expectedError: nil,
		},
		{
			name: "ticket deletion fails after successful file deletion",
			args: args{
				ctx:         context.Background(),
				accessToken: "valid_token",
				id:          1,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				fileStorageService *mockservices.MockFileStorageService,
				toysService *mockservices.MockToysService,
				_ *mocklogger.MockLogger,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(1)).
					Return(testRawTicket, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
						{ID: 3, Name: "Tag3"},
					}, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					DeleteMany(gomock.Any(), []string{"file1.jpg", "file2.jpg"}).
					Return(nil).
					Times(1)

				ticketsService.
					EXPECT().
					DeleteTicket(gomock.Any(), uint64(1)).
					Return(errors.New("database error")).
					Times(1)
			},
			expectedError: errors.New("database error"),
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)

	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		nil,              // notificationsService not needed
		validationConfig, // validationConfig not needed
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					ticketsService,
					fileStorageService,
					toysService,
					logger,
				)
			}

			err := useCases.DeleteTicket(tc.args.ctx, tc.args.accessToken, tc.args.id)

			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)

				var expectedErr *customerrors.PermissionDeniedError
				switch {
				case errors.As(tc.expectedError, &expectedErr):
					var actualErr *customerrors.PermissionDeniedError
					ok := errors.As(err, &actualErr)
					require.True(t, ok, "Expected PermissionDeniedError")
					require.Equal(t, expectedErr.Message, actualErr.Message)
				default:
					require.EqualError(t, err, tc.expectedError.Error())
				}
			}
		})
	}
}

func TestUseCases_processRawTicket(t *testing.T) {
	type args struct {
		ticket entities.RawTicket
	}

	// Тестовые данные
	testRawTicket := entities.RawTicket{
		ID:          1,
		UserID:      1,
		CategoryID:  1,
		Name:        "Test Ticket",
		Description: "Test Description",
		Quantity:    1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		TagIDs:      []uint32{1, 2, 3},
	}

	testTags := []entities.Tag{
		{ID: 1, Name: "Tag1"},
		{ID: 2, Name: "Tag2"},
		{ID: 3, Name: "Tag3"},
	}

	testCases := []struct {
		name     string
		args     args
		expected *entities.Ticket
	}{
		{
			name: "successful processing with tags",
			args: args{
				ticket: testRawTicket,
			},
			expected: &entities.Ticket{
				ID:          testRawTicket.ID,
				UserID:      testRawTicket.UserID,
				CategoryID:  testRawTicket.CategoryID,
				Name:        testRawTicket.Name,
				Description: testRawTicket.Description,
				Price:       testRawTicket.Price,
				Quantity:    testRawTicket.Quantity,
				CreatedAt:   testRawTicket.CreatedAt,
				UpdatedAt:   testRawTicket.UpdatedAt,
				Tags: []entities.Tag{
					{ID: 1, Name: "Tag1"},
					{ID: 2, Name: "Tag2"},
					{ID: 3, Name: "Tag3"},
				},
			},
		},
		{
			name: "processing with empty tags",
			args: args{
				ticket: entities.RawTicket{
					ID:          1,
					UserID:      1,
					CategoryID:  1,
					Name:        "Test Ticket",
					Description: "Test Description",
					Quantity:    1,
					CreatedAt:   testRawTicket.CreatedAt,
					UpdatedAt:   testRawTicket.UpdatedAt,
					TagIDs:      []uint32{},
				},
			},
			expected: &entities.Ticket{
				ID:          testRawTicket.ID,
				UserID:      testRawTicket.UserID,
				CategoryID:  testRawTicket.CategoryID,
				Name:        testRawTicket.Name,
				Description: testRawTicket.Description,
				Price:       testRawTicket.Price,
				Quantity:    testRawTicket.Quantity,
				CreatedAt:   testRawTicket.CreatedAt,
				UpdatedAt:   testRawTicket.UpdatedAt,
				Tags:        []entities.Tag{},
			},
		},
	}

	ctrl := gomock.NewController(t)
	toysService := mockservices.NewMockToysService(ctrl)
	// Другие сервисы не нужны для этого теста
	useCases := &UseCases{
		toysService: toysService,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := useCases.processRawTicket(tc.args.ticket, testTags)

			require.Equal(t, tc.expected.ID, result.ID)
			require.Equal(t, tc.expected.UserID, result.UserID)
			require.Equal(t, tc.expected.CategoryID, result.CategoryID)
			require.Equal(t, tc.expected.Name, result.Name)
			require.Equal(t, tc.expected.Description, result.Description)
			require.Equal(t, tc.expected.Price, result.Price)
			require.Equal(t, tc.expected.Quantity, result.Quantity)
			require.Equal(t, tc.expected.CreatedAt, result.CreatedAt)
			require.Equal(t, tc.expected.UpdatedAt, result.UpdatedAt)

			// Проверяем теги
			require.Len(t, result.Tags, len(tc.expected.Tags))
			for i, expectedTag := range tc.expected.Tags {
				require.Equal(t, expectedTag.ID, result.Tags[i].ID)
				require.Equal(t, expectedTag.Name, result.Tags[i].Name)
			}
		})
	}
}

func TestUseCases_GetTicketByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uint64
	}

	// Тестовые данные
	testTime := time.Now()
	testRawTicket := &entities.RawTicket{
		ID:          1,
		UserID:      1,
		CategoryID:  1,
		Name:        "Test Ticket",
		Description: "Test Description",
		Quantity:    1,
		CreatedAt:   testTime,
		UpdatedAt:   testTime,
		TagIDs:      []uint32{1, 2, 3},
	}

	testProcessedTicket := &entities.Ticket{
		ID:          1,
		UserID:      1,
		CategoryID:  1,
		Name:        "Test Ticket",
		Description: "Test Description",
		Quantity:    1,
		CreatedAt:   testTime,
		UpdatedAt:   testTime,
		Tags: []entities.Tag{
			{ID: 1, Name: "Tag1"},
			{ID: 2, Name: "Tag2"},
			{ID: 3, Name: "Tag3"},
		},
	}

	testCases := []struct {
		name       string
		args       args
		setupMocks func(
			ticketsService *mockservices.MockTicketsService,
			toysService *mockservices.MockToysService,
		)
		expected      *entities.Ticket
		errorExpected bool
		expectedError error
	}{
		{
			name: "successful get ticket",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			setupMocks: func(
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
			) {
				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(1)).
					Return(testRawTicket, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
						{ID: 3, Name: "Tag3"},
					}, nil).
					Times(1)
			},
			expected:      testProcessedTicket,
			errorExpected: false,
		},
		{
			name: "ticket not found",
			args: args{
				ctx: context.Background(),
				id:  999,
			},
			setupMocks: func(
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
			) {
				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(999)).
					Return(nil, errors.New("ticket not found")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
			expectedError: errors.New("ticket not found"),
		},
		{
			name: "tags service error (soft processing)",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			setupMocks: func(
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
			) {
				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(1)).
					Return(testRawTicket, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(nil, errors.New("service unavailable")).
					Times(1)
			},
			expected: &entities.Ticket{
				ID:          testRawTicket.ID,
				UserID:      testRawTicket.UserID,
				CategoryID:  testRawTicket.CategoryID,
				Name:        testRawTicket.Name,
				Description: testRawTicket.Description,
				Price:       testRawTicket.Price,
				Quantity:    testRawTicket.Quantity,
				CreatedAt:   testRawTicket.CreatedAt,
				UpdatedAt:   testRawTicket.UpdatedAt,
				Tags: []entities.Tag{
					{ID: 1},
					{ID: 2},
					{ID: 3},
				},
			},
			errorExpected: false,
		},
	}

	ctrl := gomock.NewController(t)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	useCases := &UseCases{
		ticketsService: ticketsService,
		toysService:    toysService,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsService, toysService)
			}

			result, err := useCases.GetTicketByID(tc.args.ctx, tc.args.id)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.EqualError(t, err, tc.expectedError.Error())
				}
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestUseCases_UpdateTicket(t *testing.T) {
	type args struct {
		ctx           context.Context
		rawTicketData entities.RawUpdateTicketDTO
	}

	// Test data
	testUser := &entities.User{
		ID:          1,
		DisplayName: "Test User",
		Email:       "test@example.com",
	}

	// Тестовые данные
	testRawTicket := &entities.RawTicket{
		ID:          1,
		UserID:      1,
		CategoryID:  1,
		Name:        "Test Ticket",
		Description: "Test Description",
		Quantity:    1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		TagIDs:      []uint32{1, 2, 3},
		Attachments: []entities.TicketAttachment{
			{Link: "http://storage.com/old_file1.jpg"},
			{Link: "http://storage.com/old_file2.jpg"},
			{Link: "http://storage.com/MTpzdGlsbF91c2VkLmpwZw.jpg"},
		},
	}

	testRawUpdateDTO := entities.RawUpdateTicketDTO{
		AccessToken: "valid_token",
		ID:          1,
		CategoryID:  pointers.New[uint32](1),
		Name:        pointers.New("Updated Title"),
		Description: pointers.New("Updated Description"),
		Price:       pointers.New[float32](200),
		Quantity:    pointers.New[uint32](2),
		Tags:        []string{"tag1", "tag2"},
		Attachments: []*graphql.Upload{
			{
				File:     strings.NewReader("test content"),
				Size:     1024,
				Filename: "new_file1.jpg",
			},
			{
				File:     strings.NewReader("test content"),
				Size:     1024,
				Filename: "still_used.jpg",
			},
		},
	}

	testTagIDs := []uint32{1, 2}

	testCases := []struct {
		name       string
		args       args
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			ticketsService *mockservices.MockTicketsService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			logger *mocklogger.MockLogger,
		)
		errorExpected bool
		expectedError error
	}{
		{
			name: "successful update with new attachments",
			args: args{
				ctx:           context.Background(),
				rawTicketData: testRawUpdateDTO,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mocklogger.MockLogger,
			) {
				// Authentication
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				// Get existing ticket
				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(1)).
					Return(testRawTicket, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
						{ID: 3, Name: "Tag3"},
					}, nil).
					Times(1)

				// Create tags
				toysService.
					EXPECT().
					CreateTags(gomock.Any(), []entities.CreateTagDTO{
						{Name: "tag1"},
						{Name: "tag2"},
					}).
					Return(testTagIDs, nil).
					Times(1)

				// Delete old attachments
				fileStorageService.
					EXPECT().
					DeleteMany(gomock.Any(), []string{"old_file1.jpg", "old_file2.jpg"}).
					Return(nil).
					Times(1)

				// Upload new files
				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("http://storage.com/new_file1.jpg", nil).
					Times(1)

				// Update ticket
				ticketsService.
					EXPECT().
					UpdateTicket(gomock.Any(), entities.UpdateTicketDTO{
						ID:          1,
						CategoryID:  pointers.New[uint32](1),
						Name:        pointers.New("Updated Title"),
						Description: pointers.New("Updated Description"),
						Price:       pointers.New[float32](200),
						Quantity:    pointers.New[uint32](2),
						TagIDs:      testTagIDs,
						Attachments: []string{
							"http://storage.com/MTpzdGlsbF91c2VkLmpwZw.jpg",
							"http://storage.com/new_file1.jpg",
						},
					}).
					Return(nil).
					Times(1)
			},
		},
		{
			name: "invalid access token",
			args: args{
				ctx: context.Background(),
				rawTicketData: entities.RawUpdateTicketDTO{
					AccessToken: "invalid_token",
					ID:          1,
				},
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mocklogger.MockLogger,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "invalid_token").
					Return(nil, errors.New("invalid token")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "ticket not found",
			args: args{
				ctx:           context.Background(),
				rawTicketData: testRawUpdateDTO,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mocklogger.MockLogger,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(1)).
					Return(nil, errors.New("ticket not found")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "permission denied - not owner",
			args: args{
				ctx:           context.Background(),
				rawTicketData: testRawUpdateDTO,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mocklogger.MockLogger,
			) {
				otherUserTicket := &entities.RawTicket{
					ID:     1,
					UserID: 2, // Different user
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(1)).
					Return(otherUserTicket, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
						{ID: 3, Name: "Tag3"},
					}, nil).
					Times(1)
			},
			errorExpected: true,
			expectedError: &customerrors.PermissionDeniedError{
				Message: fmt.Sprintf(
					"User with ID=%d is not owner of Ticket with ID=%d",
					testUser.ID,
					1,
				),
			},
		},
		{
			name: "tag creation fails",
			args: args{
				ctx:           context.Background(),
				rawTicketData: testRawUpdateDTO,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mocklogger.MockLogger,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(1)).
					Return(testRawTicket, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
						{ID: 3, Name: "Tag3"},
					}, nil).
					Times(1)

				toysService.
					EXPECT().
					CreateTags(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("tag service error")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "file deletion fails but update continues",
			args: args{
				ctx:           context.Background(),
				rawTicketData: testRawUpdateDTO,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				logger *mocklogger.MockLogger,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(1)).
					Return(testRawTicket, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
						{ID: 3, Name: "Tag3"},
					}, nil).
					Times(1)

				toysService.
					EXPECT().
					CreateTags(gomock.Any(), gomock.Any()).
					Return(testTagIDs, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					DeleteMany(gomock.Any(), gomock.Any()).
					Return([]error{errors.New("file deletion error")}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				// Upload new files
				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("http://storage.com/new_file1.jpg", nil).
					Times(1)

				ticketsService.
					EXPECT().
					UpdateTicket(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name: "file upload fails",
			args: args{
				ctx:           context.Background(),
				rawTicketData: testRawUpdateDTO,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				logger *mocklogger.MockLogger,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(1)).
					Return(testRawTicket, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
						{ID: 3, Name: "Tag3"},
					}, nil).
					Times(1)

				toysService.
					EXPECT().
					CreateTags(gomock.Any(), gomock.Any()).
					Return(testTagIDs, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					DeleteMany(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				// Upload new files
				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", errors.New("upload error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "extension error",
			args: args{
				ctx: context.Background(),
				rawTicketData: entities.RawUpdateTicketDTO{
					ID:          testRawTicket.ID,
					AccessToken: "valid_token",
					Attachments: []*graphql.Upload{
						{
							Filename: "test.exe",
						},
					},
				},
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mocklogger.MockLogger,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(1)).
					Return(testRawTicket, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
						{ID: 3, Name: "Tag3"},
					}, nil).
					Times(1)

				toysService.
					EXPECT().
					CreateTags(gomock.Any(), gomock.Any()).
					Return(testTagIDs, nil).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "ticket update fails",
			args: args{
				ctx:           context.Background(),
				rawTicketData: testRawUpdateDTO,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				logger *mocklogger.MockLogger,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(1)).
					Return(testRawTicket, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
						{ID: 3, Name: "Tag3"},
					}, nil).
					Times(1)

				toysService.
					EXPECT().
					CreateTags(gomock.Any(), gomock.Any()).
					Return(testTagIDs, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					DeleteMany(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				// Upload new files
				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("http://storage.com/new_file1.jpg", nil).
					Times(1)

				ticketsService.
					EXPECT().
					UpdateTicket(gomock.Any(), gomock.Any()).
					Return(errors.New("database error")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)

	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		nil, // notificationsService not needed
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					ticketsService,
					toysService,
					fileStorageService,
					logger,
				)
			}

			err := useCases.UpdateTicket(tc.args.ctx, tc.args.rawTicketData)

			if tc.errorExpected {
				require.Error(t, err)

				var expectedErr *customerrors.PermissionDeniedError
				switch {
				case errors.As(tc.expectedError, &expectedErr):
					var actualErr *customerrors.PermissionDeniedError
					ok := errors.As(err, &actualErr)
					require.True(t, ok, "Expected PermissionDeniedError")
					require.Equal(t, expectedErr.Message, actualErr.Message)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUseCases_GetMasterByID(t *testing.T) {
	testCases := []struct {
		name       string
		id         uint64
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		expected      *entities.Master
		errorExpected bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(
				_ *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				master := &entities.Master{
					ID:     1,
					UserID: 1,
					Info:   pointers.New("Test Master"),
				}

				toysService.
					EXPECT().
					GetMasterByID(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)
			},
			expected: &entities.Master{
				ID:     1,
				UserID: 1,
				Info:   pointers.New("Test Master"),
			},
			errorExpected: false,
		},
		{
			name: "master not found",
			id:   999,
			setupMocks: func(
				_ *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				toysService.
					EXPECT().
					GetMasterByID(gomock.Any(), uint64(999)).
					Return(nil, errors.New("master not found")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			actual, err := useCases.GetMasterByID(ctx, tc.id)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetAllMasters(t *testing.T) {
	testCases := []struct {
		name       string
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		expected      []entities.Master
		errorExpected bool
	}{
		{
			name: "success with masters",
			setupMocks: func(
				_ *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				masters := []entities.Master{
					{
						ID:     1,
						UserID: 1,
						Info:   pointers.New("Test Master"),
					},
					{
						ID:     2,
						UserID: 2,
						Info:   pointers.New("Test Master"),
					},
				}

				toysService.
					EXPECT().
					GetAllMasters(gomock.Any()).
					Return(masters, nil).
					Times(1)
			},
			expected: []entities.Master{
				{
					ID:     1,
					UserID: 1,
					Info:   pointers.New("Test Master"),
				},
				{
					ID:     2,
					UserID: 2,
					Info:   pointers.New("Test Master"),
				},
			},
			errorExpected: false,
		},
		{
			name: "success with empty list",
			setupMocks: func(
				_ *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				toysService.
					EXPECT().
					GetAllMasters(gomock.Any()).
					Return([]entities.Master{}, nil).
					Times(1)
			},
			expected:      []entities.Master{},
			errorExpected: false,
		},
		{
			name: "error from service",
			setupMocks: func(
				_ *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				toysService.
					EXPECT().
					GetAllMasters(gomock.Any()).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			actual, err := useCases.GetAllMasters(ctx)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_RegisterMaster(t *testing.T) {
	testCases := []struct {
		name          string
		rawMasterData entities.RawRegisterMasterDTO
		setupMocks    func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		expected      uint64
		errorExpected bool
	}{
		{
			name: "success",
			rawMasterData: entities.RawRegisterMasterDTO{
				AccessToken: "valid_access_token",
				Info:        pointers.New("Test Master"),
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{
					ID: 1,
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					RegisterMaster(gomock.Any(), entities.RegisterMasterDTO{
						UserID: 1,
						Info:   pointers.New("Test Master"),
					}).
					Return(uint64(123), nil).
					Times(1)
			},
			expected:      123,
			errorExpected: false,
		},
		{
			name: "invalid access token",
			rawMasterData: entities.RawRegisterMasterDTO{
				AccessToken: "invalid_access_token",
				Info:        pointers.New("Test Master"),
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "invalid_access_token").
					Return(nil, errors.New("invalid token")).
					Times(1)
			},
			expected:      0,
			errorExpected: true,
		},
		{
			name: "error from toys service",
			rawMasterData: entities.RawRegisterMasterDTO{
				AccessToken: "valid_access_token",
				Info:        pointers.New("Test Master"),
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{
					ID: 1,
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					RegisterMaster(gomock.Any(), entities.RegisterMasterDTO{
						UserID: 1,
						Info:   pointers.New("Test Master"),
					}).
					Return(uint64(0), errors.New("registration failed")).
					Times(1)
			},
			expected:      0,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			actual, err := useCases.RegisterMaster(ctx, tc.rawMasterData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_CreateTicket(t *testing.T) {
	testCases := []struct {
		name          string
		rawTicketData entities.RawCreateTicketDTO
		setupMocks    func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		expected      uint64
		errorExpected bool
	}{
		{
			name: "success",
			rawTicketData: entities.RawCreateTicketDTO{
				AccessToken: "valid_access_token",
				CategoryID:  1,
				Name:        "Test Ticket",
				Description: "Test Description",
				Price:       pointers.New[float32](100),
				Quantity:    10,
				Tags:        []string{"tag1", "tag2"},
				Attachments: []*graphql.Upload{
					{Filename: "test.jpg", Size: 1024, File: strings.NewReader("test content")},
				},
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				ticketsService *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("uploaded/test.jpg", nil).
					Times(1)

				toysService.
					EXPECT().
					CreateTags(gomock.Any(), []entities.CreateTagDTO{
						{Name: "tag1"},
						{Name: "tag2"},
					}).
					Return([]uint32{1, 2}, nil).
					Times(1)

				ticketsService.
					EXPECT().
					CreateTicket(
						gomock.Any(),
						entities.CreateTicketDTO{
							UserID:      1,
							CategoryID:  1,
							Name:        "Test Ticket",
							Description: "Test Description",
							Price:       pointers.New[float32](100),
							Quantity:    10,
							TagIDs:      []uint32{1, 2},
							Attachments: []string{"uploaded/test.jpg"},
						},
					).
					Return(uint64(123), nil).
					Times(1)
			},
			expected:      123,
			errorExpected: false,
		},
		{
			name: "invalid access token",
			rawTicketData: entities.RawCreateTicketDTO{
				AccessToken: "invalid_access_token",
				CategoryID:  1,
				Name:        "Test Ticket",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "invalid_access_token").
					Return(nil, errors.New("invalid token")).
					Times(1)
			},
			expected:      0,
			errorExpected: true,
		},
		{
			name: "upload files error",
			rawTicketData: entities.RawCreateTicketDTO{
				AccessToken: "valid_access_token",
				CategoryID:  1,
				Name:        "Test Ticket",
				Attachments: []*graphql.Upload{
					{Filename: "test.exe", Size: 1024, File: strings.NewReader("test content")},
				},
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				logger *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expected:      0,
			errorExpected: true,
		},
		{
			name: "create tags error",
			rawTicketData: entities.RawCreateTicketDTO{
				AccessToken: "valid_access_token",
				CategoryID:  1,
				Name:        "Test Ticket",
				Tags:        []string{"tag1"},
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0) // Нет вложений

				toysService.
					EXPECT().
					CreateTags(gomock.Any(), []entities.CreateTagDTO{{Name: "tag1"}}).
					Return(nil, errors.New("tags creation failed")).
					Times(1)
			},
			expected:      0,
			errorExpected: true,
		},
		{
			name: "create ticket error",
			rawTicketData: entities.RawCreateTicketDTO{
				AccessToken: "valid_access_token",
				CategoryID:  1,
				Name:        "Test Ticket",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				ticketsService *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0) // Нет вложений

				toysService.
					EXPECT().
					CreateTags(gomock.Any(), []entities.CreateTagDTO{}).
					Return(nil, nil).
					Times(1)

				ticketsService.
					EXPECT().
					CreateTicket(gomock.Any(), entities.CreateTicketDTO{
						UserID:      1,
						CategoryID:  1,
						Name:        "Test Ticket",
						Description: "",
						Quantity:    0,
						Attachments: []string{},
					}).
					Return(uint64(0), errors.New("ticket creation failed")).
					Times(1)
			},
			expected:      0,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			actual, err := useCases.CreateTicket(ctx, tc.rawTicketData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetRespondByID(t *testing.T) {
	testCases := []struct {
		name        string
		id          uint64
		accessToken string
		setupMocks  func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		expected      *entities.Respond
		errorExpected bool
	}{
		{
			name:        "success as ticket owner",
			id:          1,
			accessToken: "valid_access_token",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				ticketsService *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				respond := &entities.Respond{
					ID:       1,
					TicketID: 1,
					MasterID: 2,
				}

				user := &entities.User{ID: 1}

				ticket := &entities.RawTicket{
					ID:     1,
					UserID: 1,
				}

				master := &entities.Master{
					ID:     2,
					UserID: 3,
				}

				ticketsService.
					EXPECT().
					GetRespondByID(gomock.Any(), uint64(1)).
					Return(respond, nil).
					Times(1)

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(1)).
					Return(ticket, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
						{ID: 3, Name: "Tag3"},
					}, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByID(gomock.Any(), uint64(2)).
					Return(master, nil).
					Times(1)
			},
			expected: &entities.Respond{
				ID:       1,
				TicketID: 1,
				MasterID: 2,
			},
			errorExpected: false,
		},
		{
			name:        "permission denied",
			id:          1,
			accessToken: "valid_access_token",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				ticketsService *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				respond := &entities.Respond{
					ID:       1,
					TicketID: 1,
					MasterID: 2,
				}

				user := &entities.User{ID: 4}

				ticket := &entities.RawTicket{
					ID:     1,
					UserID: 1,
				}

				master := &entities.Master{
					ID:     2,
					UserID: 3,
				}

				ticketsService.
					EXPECT().
					GetRespondByID(gomock.Any(), uint64(1)).
					Return(respond, nil).
					Times(1)

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(1)).
					Return(ticket, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
						{ID: 3, Name: "Tag3"},
					}, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
						{ID: 3, Name: "Tag3"},
					}, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByID(gomock.Any(), uint64(2)).
					Return(master, nil).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
		{
			name:        "respond not found",
			id:          999,
			accessToken: "valid_access_token",
			setupMocks: func(
				_ *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				ticketsService *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ticketsService.
					EXPECT().
					GetRespondByID(gomock.Any(), uint64(999)).
					Return(nil, errors.New("respond not found")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
		{
			name:        "invalid access token",
			id:          1,
			accessToken: "invalid_access_token",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				ticketsService *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				respond := &entities.Respond{
					ID:       1,
					TicketID: 1,
					MasterID: 2,
				}

				ticketsService.
					EXPECT().
					GetRespondByID(gomock.Any(), uint64(1)).
					Return(respond, nil).
					Times(1)

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "invalid_access_token").
					Return(nil, errors.New("invalid token")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
		{
			name:        "ticket not found",
			id:          1,
			accessToken: "valid_access_token",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				ticketsService *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				respond := &entities.Respond{
					ID:       1,
					TicketID: 1,
					MasterID: 2,
				}

				user := &entities.User{ID: 1}

				ticketsService.
					EXPECT().
					GetRespondByID(gomock.Any(), uint64(1)).
					Return(respond, nil).
					Times(1)

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(1)).
					Return(nil, errors.New("ticket not found")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
		{
			name:        "master not found",
			id:          1,
			accessToken: "valid_access_token",
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				ticketsService *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				respond := &entities.Respond{
					ID:       1,
					TicketID: 1,
					MasterID: 2,
				}

				user := &entities.User{ID: 1}

				ticket := &entities.RawTicket{
					ID:     1,
					UserID: 1,
				}

				ticketsService.
					EXPECT().
					GetRespondByID(gomock.Any(), uint64(1)).
					Return(respond, nil).
					Times(1)

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(1)).
					Return(ticket, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByID(gomock.Any(), uint64(2)).
					Return(nil, errors.New("master not found")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			actual, err := useCases.GetRespondByID(ctx, tc.id, tc.accessToken)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_UpdateUserProfile(t *testing.T) {
	testCases := []struct {
		name               string
		rawUserProfileData entities.RawUpdateUserProfileDTO
		setupMocks         func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		errorExpected bool
	}{
		{
			name: "success with new avatar",
			rawUserProfileData: entities.RawUpdateUserProfileDTO{
				AccessToken: "valid_access_token",
				DisplayName: pointers.New("New Name"),
				Phone:       pointers.New("1234567890"),
				Telegram:    pointers.New("@newtelegram"),
				Avatar:      &graphql.Upload{Filename: "new_avatar.jpg", Size: 1024, File: strings.NewReader("new content")},
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{
					ID:     1,
					Avatar: pointers.New("old_avatar.jpg"),
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					Delete(gomock.Any(), "old_avatar.jpg").
					Return(nil).
					Times(1)

				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("uploaded/new_avatar.jpg", nil).
					Times(1)

				ssoService.
					EXPECT().
					UpdateUserProfile(gomock.Any(), entities.UpdateUserProfileDTO{
						AccessToken: "valid_access_token",
						DisplayName: pointers.New("New Name"),
						Phone:       pointers.New("1234567890"),
						Telegram:    pointers.New("@newtelegram"),
						Avatar:      pointers.New("uploaded/new_avatar.jpg"),
					}).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "success without avatar change",
			rawUserProfileData: entities.RawUpdateUserProfileDTO{
				AccessToken: "valid_access_token",
				DisplayName: pointers.New("New Name"),
				Phone:       pointers.New("1234567890"),
				Telegram:    pointers.New("@newtelegram"),
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{
					ID:     1,
					Avatar: pointers.New("old_avatar.jpg"),
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				ssoService.
					EXPECT().
					UpdateUserProfile(
						gomock.Any(),
						entities.UpdateUserProfileDTO{
							AccessToken: "valid_access_token",
							DisplayName: pointers.New("New Name"),
							Phone:       pointers.New("1234567890"),
							Telegram:    pointers.New("@newtelegram"),
						}).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "invalid access token",
			rawUserProfileData: entities.RawUpdateUserProfileDTO{
				AccessToken: "invalid_access_token",
				DisplayName: pointers.New("New Name"),
				Phone:       pointers.New("1234567890"),
				Telegram:    pointers.New("@newtelegram"),
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "invalid_access_token").
					Return(nil, errors.New("invalid token")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "invalid file extension",
			rawUserProfileData: entities.RawUpdateUserProfileDTO{
				AccessToken: "valid_access_token",
				DisplayName: pointers.New("New Name"),
				Phone:       pointers.New("1234567890"),
				Telegram:    pointers.New("@newtelegram"),
				Avatar:      &graphql.Upload{Filename: "new_avatar.exe", Size: 1024, File: strings.NewReader("new content")},
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{
					ID:     1,
					Avatar: pointers.New("old_avatar.jpg"),
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "delete old avatar error",
			rawUserProfileData: entities.RawUpdateUserProfileDTO{
				AccessToken: "valid_access_token",
				DisplayName: pointers.New("New Name"),
				Phone:       pointers.New("1234567890"),
				Telegram:    pointers.New("@newtelegram"),
				Avatar:      &graphql.Upload{Filename: "new_avatar.jpg", Size: 1024, File: strings.NewReader("new content")},
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{
					ID:     1,
					Avatar: pointers.New("old_avatar.jpg"),
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					Delete(gomock.Any(), "old_avatar.jpg").
					Return(errors.New("delete failed")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "upload new avatar error",
			rawUserProfileData: entities.RawUpdateUserProfileDTO{
				AccessToken: "valid_access_token",
				DisplayName: pointers.New("New Name"),
				Phone:       pointers.New("1234567890"),
				Telegram:    pointers.New("@newtelegram"),
				Avatar:      &graphql.Upload{Filename: "new_avatar.jpg", Size: 1024, File: strings.NewReader("new content")},
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{
					ID:     1,
					Avatar: pointers.New("old_avatar.jpg"),
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					Delete(gomock.Any(), "old_avatar.jpg").
					Return(nil).
					Times(1)

				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", errors.New("upload failed")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "update profile error",
			rawUserProfileData: entities.RawUpdateUserProfileDTO{
				AccessToken: "valid_access_token",
				DisplayName: pointers.New("New Name"),
				Phone:       pointers.New("1234567890"),
				Telegram:    pointers.New("@newtelegram"),
				Avatar:      &graphql.Upload{Filename: "new_avatar.jpg", Size: 1024, File: strings.NewReader("new content")},
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{
					ID:     1,
					Avatar: pointers.New("old_avatar.jpg"),
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					Delete(gomock.Any(), "old_avatar.jpg").
					Return(nil).
					Times(1)

				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("uploaded/new_avatar.jpg", nil).
					Times(1)

				ssoService.
					EXPECT().
					UpdateUserProfile(gomock.Any(), entities.UpdateUserProfileDTO{
						AccessToken: "valid_access_token",
						DisplayName: pointers.New("New Name"),
						Phone:       pointers.New("1234567890"),
						Telegram:    pointers.New("@newtelegram"),
						Avatar:      pointers.New("uploaded/new_avatar.jpg"),
					}).
					Return(errors.New("update failed")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			err := useCases.UpdateUserProfile(ctx, tc.rawUserProfileData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUseCases_UpdateToy(t *testing.T) {
	testCases := []struct {
		name       string
		rawToyData entities.RawUpdateToyDTO
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		errorExpected bool
	}{
		{
			name: "success with new and old attachments",
			rawToyData: entities.RawUpdateToyDTO{
				AccessToken: "valid_access_token",
				ID:          1,
				CategoryID:  pointers.New[uint32](2),
				Name:        pointers.New("Updated Toy"),
				Description: pointers.New("Updated Description"),
				Price:       pointers.New[float32](150),
				Quantity:    pointers.New[uint32](20),
				Tags:        []string{"tag1", "tag2"},
				Attachments: []*graphql.Upload{
					{Filename: "new_attachment.jpg", Size: 1024, File: strings.NewReader("new content")},
					{Filename: "old_attachment.jpg", Size: 1024, File: strings.NewReader("old content")},
				},
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				toy := &entities.Toy{
					ID:       1,
					MasterID: 1,
					Attachments: []entities.ToyAttachment{
						{Link: "MTpuZXdfYXR0YWNobWVudC5qcGc.jpg"},
						{Link: "path/to/to_delete.jpg"},
					},
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), uint64(1)).
					Return(toy, nil).
					Times(1)

				toysService.
					EXPECT().
					CreateTags(gomock.Any(), []entities.CreateTagDTO{
						{Name: "tag1"},
						{Name: "tag2"},
					}).
					Return([]uint32{1, 2}, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					DeleteMany(gomock.Any(), []string{"to_delete.jpg"}).
					Return([]error{}).
					Times(1)

				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("uploaded/new_attachment.jpg", nil).
					Times(1)

				toysService.
					EXPECT().
					UpdateToy(gomock.Any(), entities.UpdateToyDTO{
						ID:          1,
						CategoryID:  pointers.New[uint32](2),
						Name:        pointers.New("Updated Toy"),
						Description: pointers.New("Updated Description"),
						Price:       pointers.New[float32](150),
						Quantity:    pointers.New[uint32](20),
						TagIDs:      []uint32{1, 2},
						Attachments: []string{"MTpuZXdfYXR0YWNobWVudC5qcGc.jpg", "uploaded/new_attachment.jpg"},
					}).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "invalid access token",
			rawToyData: entities.RawUpdateToyDTO{
				AccessToken: "invalid_access_token",
				ID:          1,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "invalid_access_token").
					Return(nil, errors.New("invalid token")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "master not found",
			rawToyData: entities.RawUpdateToyDTO{
				AccessToken: "valid_access_token",
				ID:          1,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(nil, errors.New("master not found")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "toy not found",
			rawToyData: entities.RawUpdateToyDTO{
				AccessToken: "valid_access_token",
				ID:          1,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), uint64(1)).
					Return(nil, errors.New("toy not found")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "permission denied",
			rawToyData: entities.RawUpdateToyDTO{
				AccessToken: "valid_access_token",
				ID:          1,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				toy := &entities.Toy{
					ID:       1,
					MasterID: 2, // Different master
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), uint64(1)).
					Return(toy, nil).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "create tags error",
			rawToyData: entities.RawUpdateToyDTO{
				AccessToken: "valid_access_token",
				ID:          1,
				Tags:        []string{"tag1"},
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				toy := &entities.Toy{
					ID:       1,
					MasterID: 1,
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), uint64(1)).
					Return(toy, nil).
					Times(1)

				toysService.
					EXPECT().
					CreateTags(gomock.Any(), []entities.CreateTagDTO{{Name: "tag1"}}).
					Return(nil, errors.New("tags creation failed")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "invalid attachment extension",
			rawToyData: entities.RawUpdateToyDTO{
				AccessToken: "valid_access_token",
				ID:          1,
				Attachments: []*graphql.Upload{
					{Filename: "new_attachment.exe", Size: 1024, File: strings.NewReader("new content")},
				},
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				toy := &entities.Toy{
					ID:       1,
					MasterID: 1,
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), uint64(1)).
					Return(toy, nil).
					Times(1)

				toysService.
					EXPECT().
					CreateTags(gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "delete attachments error",
			rawToyData: entities.RawUpdateToyDTO{
				AccessToken: "valid_access_token",
				ID:          1,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				logger *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				toy := &entities.Toy{
					ID:       1,
					MasterID: 1,
					Attachments: []entities.ToyAttachment{
						{Link: "path/to/old_attachment.jpg"},
					},
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), uint64(1)).
					Return(toy, nil).
					Times(1)

				toysService.
					EXPECT().
					CreateTags(gomock.Any(), []entities.CreateTagDTO{}).
					Return([]uint32{}, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					DeleteMany(gomock.Any(), []string{"old_attachment.jpg"}).
					Return([]error{errors.New("delete failed")}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				toysService.
					EXPECT().
					UpdateToy(gomock.Any(), entities.UpdateToyDTO{
						ID:          1,
						CategoryID:  nil,
						Name:        nil,
						Description: nil,
						Price:       nil,
						Quantity:    nil,
						TagIDs:      []uint32{},
						Attachments: []string{},
					}).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "upload files error",
			rawToyData: entities.RawUpdateToyDTO{
				AccessToken: "valid_access_token",
				ID:          1,
				Attachments: []*graphql.Upload{
					{Filename: "new_attachment.jpg", Size: 1024, File: strings.NewReader("new content")},
				},
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				logger *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				toy := &entities.Toy{
					ID:       1,
					MasterID: 1,
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), uint64(1)).
					Return(toy, nil).
					Times(1)

				toysService.
					EXPECT().
					CreateTags(gomock.Any(), []entities.CreateTagDTO{}).
					Return([]uint32{}, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("", errors.New("upload failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "update toy error",
			rawToyData: entities.RawUpdateToyDTO{
				AccessToken: "valid_access_token",
				ID:          1,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				toy := &entities.Toy{
					ID:       1,
					MasterID: 1,
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), uint64(1)).
					Return(toy, nil).
					Times(1)

				toysService.
					EXPECT().
					CreateTags(gomock.Any(), []entities.CreateTagDTO{}).
					Return([]uint32{}, nil).
					Times(1)

				toysService.
					EXPECT().
					UpdateToy(gomock.Any(), entities.UpdateToyDTO{
						ID:          1,
						CategoryID:  nil,
						Name:        nil,
						Description: nil,
						Price:       nil,
						Quantity:    nil,
						TagIDs:      []uint32{},
						Attachments: []string{},
					}).
					Return(errors.New("update failed")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			err := useCases.UpdateToy(ctx, tc.rawToyData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUseCases_DeleteToy(t *testing.T) {
	testCases := []struct {
		name        string
		accessToken string
		id          uint64
		setupMocks  func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		errorExpected bool
	}{
		{
			name:        "success with attachments",
			accessToken: "valid_access_token",
			id:          1,
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				toy := &entities.Toy{
					ID:       1,
					MasterID: 1,
					Attachments: []entities.ToyAttachment{
						{Link: "path/to/attachment1.jpg"},
						{Link: "path/to/attachment2.jpg"},
					},
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), uint64(1)).
					Return(toy, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					DeleteMany(gomock.Any(), []string{"attachment1.jpg", "attachment2.jpg"}).
					Return([]error{}).
					Times(1)

				toysService.
					EXPECT().
					DeleteToy(gomock.Any(), uint64(1)).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:        "success without attachments",
			accessToken: "valid_access_token",
			id:          1,
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				toy := &entities.Toy{
					ID:       1,
					MasterID: 1,
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), uint64(1)).
					Return(toy, nil).
					Times(1)

				toysService.
					EXPECT().
					DeleteToy(gomock.Any(), uint64(1)).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:        "invalid access token",
			accessToken: "invalid_access_token",
			id:          1,
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "invalid_access_token").
					Return(nil, errors.New("invalid token")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name:        "master not found",
			accessToken: "valid_access_token",
			id:          1,
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(nil, errors.New("master not found")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name:        "toy not found",
			accessToken: "valid_access_token",
			id:          1,
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), uint64(1)).
					Return(nil, errors.New("toy not found")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name:        "permission denied",
			accessToken: "valid_access_token",
			id:          1,
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				toy := &entities.Toy{
					ID:       1,
					MasterID: 2, // Different master
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), uint64(1)).
					Return(toy, nil).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name:        "delete attachments error",
			accessToken: "valid_access_token",
			id:          1,
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				logger *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				toy := &entities.Toy{
					ID:       1,
					MasterID: 1,
					Attachments: []entities.ToyAttachment{
						{Link: "path/to/attachment1.jpg"},
					},
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), uint64(1)).
					Return(toy, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					DeleteMany(gomock.Any(), []string{"attachment1.jpg"}).
					Return([]error{errors.New("delete failed")}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				toysService.
					EXPECT().
					DeleteToy(gomock.Any(), uint64(1)).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:        "delete toy error",
			accessToken: "valid_access_token",
			id:          1,
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				fileStorageService *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				toy := &entities.Toy{
					ID:       1,
					MasterID: 1,
					Attachments: []entities.ToyAttachment{
						{Link: "path/to/attachment1.jpg"},
					},
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), uint64(1)).
					Return(toy, nil).
					Times(1)

				fileStorageService.
					EXPECT().
					DeleteMany(gomock.Any(), []string{"attachment1.jpg"}).
					Return([]error{}).
					Times(1)

				toysService.
					EXPECT().
					DeleteToy(gomock.Any(), uint64(1)).
					Return(errors.New("delete failed")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			err := useCases.DeleteToy(ctx, tc.accessToken, tc.id)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUseCases_UpdateRespond(t *testing.T) {
	testCases := []struct {
		name           string
		rawRespondData entities.RawUpdateRespondDTO
		setupMocks     func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		errorExpected bool
	}{
		{
			name: "success",
			rawRespondData: entities.RawUpdateRespondDTO{
				AccessToken: "valid_access_token",
				ID:          1,
				Price:       pointers.New[float32](200),
				Comment:     pointers.New("Updated comment"),
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				ticketsService *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				respond := &entities.Respond{
					ID:       1,
					MasterID: 1,
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetRespondByID(gomock.Any(), uint64(1)).
					Return(respond, nil).
					Times(1)

				ticketsService.
					EXPECT().
					UpdateRespond(
						gomock.Any(),
						entities.UpdateRespondDTO{
							ID:      1,
							Price:   pointers.New[float32](200),
							Comment: pointers.New("Updated comment"),
						},
					).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "invalid access token",
			rawRespondData: entities.RawUpdateRespondDTO{
				AccessToken: "invalid_access_token",
				ID:          1,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "invalid_access_token").
					Return(nil, errors.New("invalid token")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "master not found",
			rawRespondData: entities.RawUpdateRespondDTO{
				AccessToken: "valid_access_token",
				ID:          1,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(nil, errors.New("master not found")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "respond not found",
			rawRespondData: entities.RawUpdateRespondDTO{
				AccessToken: "valid_access_token",
				ID:          1,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				ticketsService *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetRespondByID(gomock.Any(), uint64(1)).
					Return(nil, errors.New("respond not found")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "permission denied",
			rawRespondData: entities.RawUpdateRespondDTO{
				AccessToken: "valid_access_token",
				ID:          1,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				ticketsService *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				respond := &entities.Respond{
					ID:       1,
					MasterID: 2, // Different master
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetRespondByID(gomock.Any(), uint64(1)).
					Return(respond, nil).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "update respond error",
			rawRespondData: entities.RawUpdateRespondDTO{
				AccessToken: "valid_access_token",
				ID:          1,
				Price:       pointers.New[float32](200),
				Comment:     pointers.New("Updated comment"),
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				ticketsService *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				respond := &entities.Respond{
					ID:       1,
					MasterID: 1,
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetRespondByID(gomock.Any(), uint64(1)).
					Return(respond, nil).
					Times(1)

				ticketsService.
					EXPECT().
					UpdateRespond(gomock.Any(), entities.UpdateRespondDTO{
						ID:      1,
						Price:   pointers.New[float32](200),
						Comment: pointers.New("Updated comment"),
					}).
					Return(errors.New("update failed")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			err := useCases.UpdateRespond(ctx, tc.rawRespondData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUseCases_DeleteRespond(t *testing.T) {
	testCases := []struct {
		name        string
		accessToken string
		id          uint64
		setupMocks  func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		errorExpected bool
	}{
		{
			name:        "success",
			accessToken: "valid_access_token",
			id:          1,
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				ticketsService *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				respond := &entities.Respond{
					ID:       1,
					MasterID: 1,
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetRespondByID(gomock.Any(), uint64(1)).
					Return(respond, nil).
					Times(1)

				ticketsService.
					EXPECT().
					DeleteRespond(gomock.Any(), uint64(1)).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:        "invalid access token",
			accessToken: "invalid_access_token",
			id:          1,
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "invalid_access_token").
					Return(nil, errors.New("invalid token")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name:        "master not found",
			accessToken: "valid_access_token",
			id:          1,
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(nil, errors.New("master not found")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name:        "respond not found",
			accessToken: "valid_access_token",
			id:          1,
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				ticketsService *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetRespondByID(gomock.Any(), uint64(1)).
					Return(nil, errors.New("respond not found")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name:        "permission denied",
			accessToken: "valid_access_token",
			id:          1,
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				ticketsService *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				respond := &entities.Respond{
					ID:       1,
					MasterID: 2, // Different master
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetRespondByID(gomock.Any(), uint64(1)).
					Return(respond, nil).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name:        "delete respond error",
			accessToken: "valid_access_token",
			id:          1,
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				ticketsService *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{ID: 1}

				respond := &entities.Respond{
					ID:       1,
					MasterID: 1,
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByUser(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetRespondByID(gomock.Any(), uint64(1)).
					Return(respond, nil).
					Times(1)

				ticketsService.
					EXPECT().
					DeleteRespond(gomock.Any(), uint64(1)).
					Return(errors.New("delete failed")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			err := useCases.DeleteRespond(ctx, tc.accessToken, tc.id)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUseCases_UpdateMaster(t *testing.T) {
	testCases := []struct {
		name          string
		rawMasterData entities.RawUpdateMasterDTO
		setupMocks    func(
			ssoService *mockservices.MockSsoService,
			toysService *mockservices.MockToysService,
			fileStorageService *mockservices.MockFileStorageService,
			ticketsService *mockservices.MockTicketsService,
			notificationsService *mockservices.MockNotificationsService,
			logger *mocklogger.MockLogger,
			traceProvider *tracingmock.MockProvider,
		)
		errorExpected bool
	}{
		{
			name: "success",
			rawMasterData: entities.RawUpdateMasterDTO{
				AccessToken: "valid_access_token",
				ID:          1,
				Info:        pointers.New("Updated Master Info"),
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{
					ID:     1,
					UserID: 1,
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByID(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				toysService.
					EXPECT().
					UpdateMaster(
						gomock.Any(),
						entities.UpdateMasterDTO{
							ID:   1,
							Info: pointers.New("Updated Master Info"),
						},
					).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "invalid access token",
			rawMasterData: entities.RawUpdateMasterDTO{
				AccessToken: "invalid_access_token",
				ID:          1,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "invalid_access_token").
					Return(nil, errors.New("invalid token")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "master not found",
			rawMasterData: entities.RawUpdateMasterDTO{
				AccessToken: "valid_access_token",
				ID:          1,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByID(gomock.Any(), uint64(1)).
					Return(nil, errors.New("master not found")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "permission denied",
			rawMasterData: entities.RawUpdateMasterDTO{
				AccessToken: "valid_access_token",
				ID:          1,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{
					ID:     1,
					UserID: 2, // Different user
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByID(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "update master error",
			rawMasterData: entities.RawUpdateMasterDTO{
				AccessToken: "valid_access_token",
				ID:          1,
				Info:        pointers.New("Updated Master Info"),
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockFileStorageService,
				_ *mockservices.MockTicketsService,
				_ *mockservices.MockNotificationsService,
				_ *mocklogger.MockLogger,
				_ *tracingmock.MockProvider,
			) {
				user := &entities.User{ID: 1}

				master := &entities.Master{
					ID:     1,
					UserID: 1,
				}

				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_access_token").
					Return(user, nil).
					Times(1)

				toysService.
					EXPECT().
					GetMasterByID(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)

				toysService.
					EXPECT().
					UpdateMaster(
						gomock.Any(),
						entities.UpdateMasterDTO{
							ID:   1,
							Info: pointers.New("Updated Master Info"),
						},
					).
					Return(errors.New("update failed")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	fileStorageService := mockservices.NewMockFileStorageService(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	traceProvider := tracingmock.NewMockProvider(ctrl)
	useCases := New(
		ssoService,
		toysService,
		fileStorageService,
		ticketsService,
		notificationsService,
		validationConfig,
		logger,
		traceProvider,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					ssoService,
					toysService,
					fileStorageService,
					ticketsService,
					notificationsService,
					logger,
					traceProvider,
				)
			}

			err := useCases.UpdateMaster(ctx, tc.rawMasterData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUseCases_GetMyEmailCommunications(t *testing.T) {
	type args struct {
		ctx         context.Context
		accessToken string
	}

	// Test data
	testUser := &entities.User{
		ID:          1,
		DisplayName: "Test User",
		Email:       "test@example.com",
	}

	testEmails := []entities.Email{
		{
			ID:      1,
			UserID:  1,
			Content: "Test Subject 1",
			SentAt:  time.Now(),
		},
		{
			ID:      2,
			UserID:  1,
			Content: "Test Subject 2",
			SentAt:  time.Now(),
		},
	}

	testCases := []struct {
		name       string
		args       args
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			notificationsService *mockservices.MockNotificationsService,
		)
		expected      []entities.Email
		errorExpected bool
		expectedError error
	}{
		{
			name: "successful get emails",
			args: args{
				ctx:         context.Background(),
				accessToken: "valid_token",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				notificationsService *mockservices.MockNotificationsService,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				notificationsService.
					EXPECT().
					GetUserEmailCommunications(gomock.Any(), testUser.ID).
					Return(testEmails, nil).
					Times(1)
			},
			expected:      testEmails,
			errorExpected: false,
		},
		{
			name: "invalid access token",
			args: args{
				ctx:         context.Background(),
				accessToken: "invalid_token",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockNotificationsService,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "invalid_token").
					Return(nil, errors.New("invalid token")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
			expectedError: errors.New("invalid token"),
		},
		{
			name: "failed to get emails",
			args: args{
				ctx:         context.Background(),
				accessToken: "valid_token",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				notificationsService *mockservices.MockNotificationsService,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				notificationsService.
					EXPECT().
					GetUserEmailCommunications(gomock.Any(), testUser.ID).
					Return(nil, errors.New("service unavailable")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
			expectedError: errors.New("service unavailable"),
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	notificationsService := mockservices.NewMockNotificationsService(ctrl)
	// Other services not needed for this test
	useCases := &UseCases{
		ssoService:           ssoService,
		notificationsService: notificationsService,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoService, notificationsService)
			}

			actual, err := useCases.GetMyEmailCommunications(tc.args.ctx, tc.args.accessToken)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.EqualError(t, err, tc.expectedError.Error())
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetMyResponds(t *testing.T) {
	type args struct {
		ctx         context.Context
		accessToken string
	}

	// Test data
	testUser := &entities.User{
		ID:          1,
		DisplayName: "Test User",
		Email:       "test@example.com",
	}

	testResponds := []entities.Respond{
		{
			ID:        1,
			MasterID:  1,
			TicketID:  100,
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			MasterID:  1,
			TicketID:  101,
			CreatedAt: time.Now().Add(-time.Hour),
		},
	}

	testCases := []struct {
		name       string
		args       args
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			ticketsService *mockservices.MockTicketsService,
		)
		expected      []entities.Respond
		errorExpected bool
		expectedError error
	}{
		{
			name: "successful get responds",
			args: args{
				ctx:         context.Background(),
				accessToken: "valid_token",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetUserResponds(gomock.Any(), testUser.ID).
					Return(testResponds, nil).
					Times(1)
			},
			expected:      testResponds,
			errorExpected: false,
		},
		{
			name: "invalid access token",
			args: args{
				ctx:         context.Background(),
				accessToken: "invalid_token",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockTicketsService,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "invalid_token").
					Return(nil, errors.New("invalid token")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
			expectedError: errors.New("invalid token"),
		},
		{
			name: "failed to get responds",
			args: args{
				ctx:         context.Background(),
				accessToken: "valid_token",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetUserResponds(gomock.Any(), testUser.ID).
					Return(nil, errors.New("service unavailable")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
			expectedError: errors.New("service unavailable"),
		},
		{
			name: "empty responds list",
			args: args{
				ctx:         context.Background(),
				accessToken: "valid_token",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetUserResponds(gomock.Any(), testUser.ID).
					Return([]entities.Respond{}, nil).
					Times(1)
			},
			expected:      []entities.Respond{},
			errorExpected: false,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	// Other services not needed for this test
	useCases := &UseCases{
		ssoService:     ssoService,
		ticketsService: ticketsService,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoService, ticketsService)
			}

			actual, err := useCases.GetMyResponds(tc.args.ctx, tc.args.accessToken)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.EqualError(t, err, tc.expectedError.Error())
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetTicketResponds(t *testing.T) {
	type args struct {
		ctx         context.Context
		ticketID    uint64
		accessToken string
	}

	// Test data
	testUser := &entities.User{
		ID:          1,
		DisplayName: "Test User",
		Email:       "test@example.com",
	}

	testTicket := &entities.RawTicket{
		ID:     100,
		UserID: 1,
	}

	otherUserTicket := &entities.RawTicket{
		ID:     101,
		UserID: 2, // Different user
	}

	testResponds := []entities.Respond{
		{
			ID:        1,
			MasterID:  3,
			TicketID:  100,
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			MasterID:  4,
			TicketID:  100,
			CreatedAt: time.Now().Add(-time.Hour),
		},
	}

	testCases := []struct {
		name       string
		args       args
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			ticketsService *mockservices.MockTicketsService,
			toysService *mockservices.MockToysService,
		)
		expected      []entities.Respond
		errorExpected bool
		expectedError error
	}{
		{
			name: "successful get responds",
			args: args{
				ctx:         context.Background(),
				ticketID:    100,
				accessToken: "valid_token",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
			) {
				// Authentication
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				// Get ticket
				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(100)).
					Return(testTicket, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
						{ID: 3, Name: "Tag3"},
					}, nil).
					Times(1)

				// Get responds
				ticketsService.
					EXPECT().
					GetTicketResponds(gomock.Any(), uint64(100)).
					Return(testResponds, nil).
					Times(1)
			},
			expected:      testResponds,
			errorExpected: false,
		},
		{
			name: "invalid access token",
			args: args{
				ctx:         context.Background(),
				ticketID:    100,
				accessToken: "invalid_token",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "invalid_token").
					Return(nil, errors.New("invalid token")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
			expectedError: errors.New("invalid token"),
		},
		{
			name: "ticket not found",
			args: args{
				ctx:         context.Background(),
				ticketID:    999,
				accessToken: "valid_token",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(999)).
					Return(nil, errors.New("ticket not found")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
			expectedError: errors.New("ticket not found"),
		},
		{
			name: "permission denied - not owner",
			args: args{
				ctx:         context.Background(),
				ticketID:    101,
				accessToken: "valid_token",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(101)).
					Return(otherUserTicket, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
						{ID: 3, Name: "Tag3"},
					}, nil).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
			expectedError: &customerrors.PermissionDeniedError{
				Message: fmt.Sprintf(
					"Ticket with ID=%d does not belong to current User with ID=%d",
					101,
					testUser.ID,
				),
			},
		},
		{
			name: "failed to get responds",
			args: args{
				ctx:         context.Background(),
				ticketID:    100,
				accessToken: "valid_token",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(100)).
					Return(testTicket, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
						{ID: 3, Name: "Tag3"},
					}, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketResponds(gomock.Any(), uint64(100)).
					Return(nil, errors.New("service error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
			expectedError: errors.New("service error"),
		},
		{
			name: "empty responds list",
			args: args{
				ctx:         context.Background(),
				ticketID:    100,
				accessToken: "valid_token",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketByID(gomock.Any(), uint64(100)).
					Return(testTicket, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{
						{ID: 1, Name: "Tag1"},
						{ID: 2, Name: "Tag2"},
						{ID: 3, Name: "Tag3"},
					}, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetTicketResponds(gomock.Any(), uint64(100)).
					Return([]entities.Respond{}, nil).
					Times(1)
			},
			expected:      []entities.Respond{},
			errorExpected: false,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	useCases := &UseCases{
		ssoService:     ssoService,
		ticketsService: ticketsService,
		toysService:    toysService,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoService, ticketsService, toysService)
			}

			actual, err := useCases.GetTicketResponds(tc.args.ctx, tc.args.ticketID, tc.args.accessToken)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					var expectedErr *customerrors.PermissionDeniedError
					switch {
					case errors.As(tc.expectedError, &expectedErr):
						var actualErr *customerrors.PermissionDeniedError
						ok := errors.As(err, &actualErr)
						require.True(t, ok, "Expected PermissionDeniedError")
						require.Equal(t, expectedErr.Message, actualErr.Message)
					default:
						require.EqualError(t, err, tc.expectedError.Error())
					}
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_RespondToTicket(t *testing.T) {
	type args struct {
		ctx            context.Context
		rawRespondData entities.RawRespondToTicketDTO
	}

	// Test data
	testUser := &entities.User{
		ID:          1,
		DisplayName: "Test User",
		Email:       "test@example.com",
	}

	testRawRespond := entities.RawRespondToTicketDTO{
		AccessToken: "valid_token",
		TicketID:    100,
		Price:       500,
		Comment:     pointers.New("Test comment"),
	}

	testRespond := entities.RespondToTicketDTO{
		UserID:   testUser.ID,
		TicketID: 100,
		Price:    500,
		Comment:  pointers.New("Test comment"),
	}

	testRespondID := uint64(123)

	testCases := []struct {
		name       string
		args       args
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			ticketsService *mockservices.MockTicketsService,
		)
		expectedID    uint64
		errorExpected bool
		expectedError error
	}{
		{
			name: "successful respond to ticket",
			args: args{
				ctx:            context.Background(),
				rawRespondData: testRawRespond,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					RespondToTicket(gomock.Any(), testRespond).
					Return(testRespondID, nil).
					Times(1)
			},
			expectedID:    testRespondID,
			errorExpected: false,
		},
		{
			name: "invalid access token",
			args: args{
				ctx: context.Background(),
				rawRespondData: entities.RawRespondToTicketDTO{
					AccessToken: "invalid_token",
					TicketID:    100,
				},
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				_ *mockservices.MockTicketsService,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "invalid_token").
					Return(nil, errors.New("invalid token")).
					Times(1)
			},
			expectedID:    0,
			errorExpected: true,
			expectedError: errors.New("invalid token"),
		},
		{
			name: "failed to respond to ticket",
			args: args{
				ctx:            context.Background(),
				rawRespondData: testRawRespond,
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					RespondToTicket(gomock.Any(), testRespond).
					Return(uint64(0), errors.New("service error")).
					Times(1)
			},
			expectedID:    0,
			errorExpected: true,
			expectedError: errors.New("service error"),
		},
		{
			name: "empty comment",
			args: args{
				ctx: context.Background(),
				rawRespondData: entities.RawRespondToTicketDTO{
					AccessToken: "valid_token",
					TicketID:    100,
					Price:       500,
					Comment:     pointers.New(""),
				},
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					RespondToTicket(
						gomock.Any(),
						entities.RespondToTicketDTO{
							UserID:   testUser.ID,
							TicketID: 100,
							Price:    500,
							Comment:  pointers.New(""),
						},
					).
					Return(testRespondID, nil).
					Times(1)
			},
			expectedID:    testRespondID,
			errorExpected: false,
		},
		{
			name: "zero price",
			args: args{
				ctx: context.Background(),
				rawRespondData: entities.RawRespondToTicketDTO{
					AccessToken: "valid_token",
					TicketID:    100,
					Price:       0,
					Comment:     pointers.New("Test comment"),
				},
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					RespondToTicket(
						gomock.Any(),
						entities.RespondToTicketDTO{
							UserID:   testUser.ID,
							TicketID: 100,
							Price:    0,
							Comment:  pointers.New("Test comment"),
						},
					).
					Return(testRespondID, nil).
					Times(1)
			},
			expectedID:    testRespondID,
			errorExpected: false,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	useCases := &UseCases{
		ssoService:     ssoService,
		ticketsService: ticketsService,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoService, ticketsService)
			}

			actualID, err := useCases.RespondToTicket(tc.args.ctx, tc.args.rawRespondData)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.EqualError(t, err, tc.expectedError.Error())
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expectedID, actualID)
		})
	}
}

func TestUseCases_GetAllTickets(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	// Test data
	testRawTickets := []entities.RawTicket{
		{
			ID:          1,
			UserID:      1,
			CategoryID:  1,
			Name:        "Ticket 1",
			Description: "Description 1",
			Price:       pointers.New[float32](100),
			Quantity:    1,
			TagIDs:      []uint32{1, 2},
		},
		{
			ID:          2,
			UserID:      2,
			CategoryID:  2,
			Name:        "Ticket 2",
			Description: "Description 2",
			Price:       pointers.New[float32](100),
			Quantity:    2,
			TagIDs:      []uint32{3, 4},
		},
	}

	testTags := []entities.Tag{
		{ID: 1, Name: "Tag1"},
		{ID: 2, Name: "Tag2"},
		{ID: 3, Name: "Tag3"},
		{ID: 4, Name: "Tag4"},
	}

	testProcessedTickets := []entities.Ticket{
		{
			ID:          1,
			UserID:      1,
			CategoryID:  1,
			Name:        "Ticket 1",
			Description: "Description 1",
			Price:       pointers.New[float32](100),
			Quantity:    1,
			Tags: []entities.Tag{
				{ID: 1, Name: "Tag1"},
				{ID: 2, Name: "Tag2"},
			},
		},
		{
			ID:          2,
			UserID:      2,
			CategoryID:  2,
			Name:        "Ticket 2",
			Description: "Description 2",
			Price:       pointers.New[float32](100),
			Quantity:    2,
			Tags: []entities.Tag{
				{ID: 3, Name: "Tag3"},
				{ID: 4, Name: "Tag4"},
			},
		},
	}

	testCases := []struct {
		name       string
		args       args
		setupMocks func(
			ticketsService *mockservices.MockTicketsService,
			toysService *mockservices.MockToysService,
		)
		expected      []entities.Ticket
		errorExpected bool
		expectedError error
	}{
		{
			name: "successful get all tickets",
			args: args{
				ctx: context.Background(),
			},
			setupMocks: func(
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
			) {
				ticketsService.
					EXPECT().
					GetAllTickets(gomock.Any()).
					Return(testRawTickets, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(testTags, nil).
					Times(1)
			},
			expected:      testProcessedTickets,
			errorExpected: false,
		},
		{
			name: "failed to get tickets",
			args: args{
				ctx: context.Background(),
			},
			setupMocks: func(
				ticketsService *mockservices.MockTicketsService,
				_ *mockservices.MockToysService,
			) {
				ticketsService.
					EXPECT().
					GetAllTickets(gomock.Any()).
					Return(nil, errors.New("service error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
			expectedError: errors.New("service error"),
		},
		{
			name: "failed to get tags but still returns tickets",
			args: args{
				ctx: context.Background(),
			},
			setupMocks: func(
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
			) {
				ticketsService.
					EXPECT().
					GetAllTickets(gomock.Any()).
					Return(testRawTickets, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(nil, errors.New("tags service error")).
					Times(1)
			},
			expected: []entities.Ticket{
				{
					ID:          1,
					UserID:      1,
					CategoryID:  1,
					Name:        "Ticket 1",
					Description: "Description 1",
					Price:       pointers.New[float32](100),
					Quantity:    1,
					Tags: []entities.Tag{
						{ID: 1},
						{ID: 2},
					},
				},
				{
					ID:          2,
					UserID:      2,
					CategoryID:  2,
					Name:        "Ticket 2",
					Description: "Description 2",
					Price:       pointers.New[float32](100),
					Quantity:    2,
					Tags: []entities.Tag{
						{ID: 3},
						{ID: 4},
					},
				},
			},
			errorExpected: false,
		},
		{
			name: "empty tickets list",
			args: args{
				ctx: context.Background(),
			},
			setupMocks: func(
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
			) {
				ticketsService.
					EXPECT().
					GetAllTickets(gomock.Any()).
					Return([]entities.RawTicket{}, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(testTags, nil).
					Times(1)
			},
			expected:      []entities.Ticket{},
			errorExpected: false,
		},
	}

	ctrl := gomock.NewController(t)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	useCases := &UseCases{
		ticketsService: ticketsService,
		toysService:    toysService,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsService, toysService)
			}

			actual, err := useCases.GetAllTickets(tc.args.ctx)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.EqualError(t, err, tc.expectedError.Error())
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetUserTickets(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID uint64
	}

	// Test data
	testUserID := uint64(1)
	testRawTickets := []entities.RawTicket{
		{
			ID:          1,
			UserID:      testUserID,
			CategoryID:  1,
			Name:        "Ticket 1",
			Description: "Description 1",
			Price:       pointers.New[float32](100),
			Quantity:    1,
			TagIDs:      []uint32{1, 2},
		},
		{
			ID:          2,
			UserID:      testUserID,
			CategoryID:  2,
			Name:        "Ticket 2",
			Description: "Description 2",
			Price:       pointers.New[float32](100),
			Quantity:    2,
			TagIDs:      []uint32{3, 4},
		},
	}

	testTags := []entities.Tag{
		{ID: 1, Name: "Tag1"},
		{ID: 2, Name: "Tag2"},
		{ID: 3, Name: "Tag3"},
		{ID: 4, Name: "Tag4"},
	}

	testProcessedTickets := []entities.Ticket{
		{
			ID:          1,
			UserID:      testUserID,
			CategoryID:  1,
			Name:        "Ticket 1",
			Description: "Description 1",
			Price:       pointers.New[float32](100),
			Quantity:    1,
			Tags: []entities.Tag{
				{ID: 1, Name: "Tag1"},
				{ID: 2, Name: "Tag2"},
			},
		},
		{
			ID:          2,
			UserID:      testUserID,
			CategoryID:  2,
			Name:        "Ticket 2",
			Description: "Description 2",
			Price:       pointers.New[float32](100),
			Quantity:    2,
			Tags: []entities.Tag{
				{ID: 3, Name: "Tag3"},
				{ID: 4, Name: "Tag4"},
			},
		},
	}

	testCases := []struct {
		name       string
		args       args
		setupMocks func(
			ticketsService *mockservices.MockTicketsService,
			toysService *mockservices.MockToysService,
		)
		expected      []entities.Ticket
		errorExpected bool
		expectedError error
	}{
		{
			name: "successful get user tickets with tags",
			args: args{
				ctx:    context.Background(),
				userID: testUserID,
			},
			setupMocks: func(
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
			) {
				ticketsService.
					EXPECT().
					GetUserTickets(gomock.Any(), testUserID).
					Return(testRawTickets, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(testTags, nil).
					Times(1)
			},
			expected:      testProcessedTickets,
			errorExpected: false,
		},
		{
			name: "failed to get user tickets",
			args: args{
				ctx:    context.Background(),
				userID: testUserID,
			},
			setupMocks: func(
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
			) {
				ticketsService.
					EXPECT().
					GetUserTickets(gomock.Any(), testUserID).
					Return(nil, errors.New("service error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
			expectedError: errors.New("service error"),
		},
		{
			name: "failed to get tags but still returns tickets",
			args: args{
				ctx:    context.Background(),
				userID: testUserID,
			},
			setupMocks: func(
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
			) {
				ticketsService.
					EXPECT().
					GetUserTickets(gomock.Any(), testUserID).
					Return(testRawTickets, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(nil, errors.New("tags service error")).
					Times(1)
			},
			expected: []entities.Ticket{
				{
					ID:          1,
					UserID:      testUserID,
					CategoryID:  1,
					Name:        "Ticket 1",
					Description: "Description 1",
					Price:       pointers.New[float32](100),
					Quantity:    1,
					Tags: []entities.Tag{
						{ID: 1},
						{ID: 2},
					},
				},
				{
					ID:          2,
					UserID:      testUserID,
					CategoryID:  2,
					Name:        "Ticket 2",
					Description: "Description 2",
					Price:       pointers.New[float32](100),
					Quantity:    2,
					Tags: []entities.Tag{
						{ID: 3},
						{ID: 4},
					},
				},
			},
			errorExpected: false,
		},
		{
			name: "empty tickets list",
			args: args{
				ctx:    context.Background(),
				userID: testUserID,
			},
			setupMocks: func(
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
			) {
				ticketsService.
					EXPECT().
					GetUserTickets(gomock.Any(), testUserID).
					Return([]entities.RawTicket{}, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(testTags, nil).
					Times(1)
			},
			expected:      []entities.Ticket{},
			errorExpected: false,
		},
		{
			name: "tickets with no tags",
			args: args{
				ctx:    context.Background(),
				userID: testUserID,
			},
			setupMocks: func(
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
			) {
				rawTickets := []entities.RawTicket{
					{
						ID:          3,
						UserID:      testUserID,
						CategoryID:  3,
						Name:        "Ticket 3",
						Description: "Description 3",
						Price:       pointers.New[float32](100),
						Quantity:    3,
						TagIDs:      []uint32{},
					},
				}

				ticketsService.
					EXPECT().
					GetUserTickets(gomock.Any(), testUserID).
					Return(rawTickets, nil).
					Times(1)

				// Called but not used since no tags in tickets
				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(testTags, nil).
					Times(1)
			},
			expected: []entities.Ticket{
				{
					ID:          3,
					UserID:      testUserID,
					CategoryID:  3,
					Name:        "Ticket 3",
					Description: "Description 3",
					Price:       pointers.New[float32](100),
					Quantity:    3,
					Tags:        []entities.Tag{},
				},
			},
			errorExpected: false,
		},
	}

	ctrl := gomock.NewController(t)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	useCases := &UseCases{
		ticketsService: ticketsService,
		toysService:    toysService,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ticketsService, toysService)
			}

			actual, err := useCases.GetUserTickets(tc.args.ctx, tc.args.userID)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.EqualError(t, err, tc.expectedError.Error())
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetMyTickets(t *testing.T) {
	type args struct {
		ctx         context.Context
		accessToken string
	}

	// Test data
	testUser := &entities.User{
		ID:          1,
		DisplayName: "Test User",
		Email:       "test@example.com",
	}

	testTickets := []entities.Ticket{
		{
			ID:          1,
			UserID:      testUser.ID,
			CategoryID:  1,
			Name:        "My Ticket 1",
			Description: "Description 1",
			Price:       pointers.New[float32](100),
			Quantity:    1,
			Tags: []entities.Tag{
				{ID: 1, Name: "Tag1"},
			},
		},
		{
			ID:          2,
			UserID:      testUser.ID,
			CategoryID:  2,
			Name:        "My Ticket 2",
			Description: "Description 2",
			Price:       pointers.New[float32](100),
			Quantity:    2,
			Tags: []entities.Tag{
				{ID: 2, Name: "Tag2"},
			},
		},
	}

	testTags := []entities.Tag{
		{ID: 1, Name: "Tag1"},
		{ID: 2, Name: "Tag2"},
		{ID: 3, Name: "Tag3"},
		{ID: 4, Name: "Tag4"},
	}

	testCases := []struct {
		name       string
		args       args
		setupMocks func(
			ssoService *mockservices.MockSsoService,
			ticketsService *mockservices.MockTicketsService,
			toysService *mockservices.MockToysService,
		)
		expected      []entities.Ticket
		errorExpected bool
		expectedError error
	}{
		{
			name: "successful get my tickets",
			args: args{
				ctx:         context.Background(),
				accessToken: "valid_token",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				// Mock GetUserTickets behavior
				ticketsService.
					EXPECT().
					GetUserTickets(gomock.Any(), testUser.ID).
					Return(
						[]entities.RawTicket{
							{
								ID:          1,
								UserID:      testUser.ID,
								CategoryID:  1,
								Name:        "My Ticket 1",
								Description: "Description 1",
								Price:       pointers.New[float32](100),
								Quantity:    1,
								TagIDs:      []uint32{1},
							},
							{
								ID:          2,
								UserID:      testUser.ID,
								CategoryID:  2,
								Name:        "My Ticket 2",
								Description: "Description 2",
								Price:       pointers.New[float32](100),
								Quantity:    2,
								TagIDs:      []uint32{2},
							},
						},
						nil,
					).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(testTags, nil).
					Times(1)
			},
			expected:      testTickets,
			errorExpected: false,
		},
		{
			name: "invalid access token",
			args: args{
				ctx:         context.Background(),
				accessToken: "invalid_token",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "invalid_token").
					Return(nil, errors.New("invalid token")).
					Times(1)

				// Should not be called when auth fails
				ticketsService.
					EXPECT().
					GetUserTickets(gomock.Any(), gomock.Any()).
					Times(0)
			},
			expected:      nil,
			errorExpected: true,
			expectedError: errors.New("invalid token"),
		},
		{
			name: "failed to get user tickets",
			args: args{
				ctx:         context.Background(),
				accessToken: "valid_token",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetUserTickets(gomock.Any(), testUser.ID).
					Return(nil, errors.New("tickets service error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
			expectedError: errors.New("tickets service error"),
		},
		{
			name: "empty tickets list",
			args: args{
				ctx:         context.Background(),
				accessToken: "valid_token",
			},
			setupMocks: func(
				ssoService *mockservices.MockSsoService,
				ticketsService *mockservices.MockTicketsService,
				toysService *mockservices.MockToysService,
			) {
				ssoService.
					EXPECT().
					GetMe(gomock.Any(), "valid_token").
					Return(testUser, nil).
					Times(1)

				ticketsService.
					EXPECT().
					GetUserTickets(gomock.Any(), testUser.ID).
					Return([]entities.RawTicket{}, nil).
					Times(1)

				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(testTags, nil).
					Times(1)
			},
			expected:      []entities.Ticket{},
			errorExpected: false,
		},
	}

	ctrl := gomock.NewController(t)
	ssoService := mockservices.NewMockSsoService(ctrl)
	ticketsService := mockservices.NewMockTicketsService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	useCases := &UseCases{
		ssoService:     ssoService,
		ticketsService: ticketsService,
		toysService:    toysService,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoService, ticketsService, toysService)
			}

			actual, err := useCases.GetMyTickets(tc.args.ctx, tc.args.accessToken)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.EqualError(t, err, tc.expectedError.Error())
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetAllTags(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	// Test data
	testTags := []entities.Tag{
		{ID: 1, Name: "Tag1"},
		{ID: 2, Name: "Tag2"},
		{ID: 3, Name: "Tag3"},
	}

	testCases := []struct {
		name          string
		args          args
		setupMocks    func(toysService *mockservices.MockToysService)
		expected      []entities.Tag
		errorExpected bool
		expectedError error
	}{
		{
			name: "successful get all tags",
			args: args{
				ctx: context.Background(),
			},
			setupMocks: func(toysService *mockservices.MockToysService) {
				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(testTags, nil).
					Times(1)
			},
			expected:      testTags,
			errorExpected: false,
		},
		{
			name: "failed to get tags",
			args: args{
				ctx: context.Background(),
			},
			setupMocks: func(toysService *mockservices.MockToysService) {
				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(nil, errors.New("service error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
			expectedError: errors.New("service error"),
		},
		{
			name: "empty tags list",
			args: args{
				ctx: context.Background(),
			},
			setupMocks: func(toysService *mockservices.MockToysService) {
				toysService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{}, nil).
					Times(1)
			},
			expected:      []entities.Tag{},
			errorExpected: false,
		},
	}

	ctrl := gomock.NewController(t)
	toysService := mockservices.NewMockToysService(ctrl)
	useCases := &UseCases{
		toysService: toysService,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysService)
			}

			actual, err := useCases.GetAllTags(tc.args.ctx)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.EqualError(t, err, tc.expectedError.Error())
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetTagByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uint32
	}

	// Test data
	testTag := &entities.Tag{
		ID:   1,
		Name: "Test Tag",
	}

	testCases := []struct {
		name          string
		args          args
		setupMocks    func(toysService *mockservices.MockToysService)
		expected      *entities.Tag
		errorExpected bool
		expectedError error
	}{
		{
			name: "successful get tag by ID",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			setupMocks: func(toysService *mockservices.MockToysService) {
				toysService.
					EXPECT().
					GetTagByID(gomock.Any(), uint32(1)).
					Return(testTag, nil).
					Times(1)
			},
			expected:      testTag,
			errorExpected: false,
		},
		{
			name: "tag not found",
			args: args{
				ctx: context.Background(),
				id:  999,
			},
			setupMocks: func(toysService *mockservices.MockToysService) {
				toysService.
					EXPECT().
					GetTagByID(gomock.Any(), uint32(999)).
					Return(nil, errors.New("tag not found")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
			expectedError: errors.New("tag not found"),
		},
		{
			name: "service error",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			setupMocks: func(toysService *mockservices.MockToysService) {
				toysService.
					EXPECT().
					GetTagByID(gomock.Any(), uint32(1)).
					Return(nil, errors.New("service unavailable")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
			expectedError: errors.New("service unavailable"),
		},
	}

	ctrl := gomock.NewController(t)
	toysService := mockservices.NewMockToysService(ctrl)
	useCases := &UseCases{
		toysService: toysService,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysService)
			}

			actual, err := useCases.GetTagByID(tc.args.ctx, tc.args.id)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.EqualError(t, err, tc.expectedError.Error())
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetCategoryByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uint32
	}

	// Test data
	testCategory := &entities.Category{
		ID:   1,
		Name: "Test Category",
	}

	testCases := []struct {
		name          string
		args          args
		setupMocks    func(toysService *mockservices.MockToysService)
		expected      *entities.Category
		errorExpected bool
		expectedError error
	}{
		{
			name: "successful get category by ID",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			setupMocks: func(toysService *mockservices.MockToysService) {
				toysService.
					EXPECT().
					GetCategoryByID(gomock.Any(), uint32(1)).
					Return(testCategory, nil).
					Times(1)
			},
			expected:      testCategory,
			errorExpected: false,
		},
		{
			name: "category not found",
			args: args{
				ctx: context.Background(),
				id:  999,
			},
			setupMocks: func(toysService *mockservices.MockToysService) {
				toysService.
					EXPECT().
					GetCategoryByID(gomock.Any(), uint32(999)).
					Return(nil, errors.New("category not found")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
			expectedError: errors.New("category not found"),
		},
		{
			name: "service error",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			setupMocks: func(toysService *mockservices.MockToysService) {
				toysService.
					EXPECT().
					GetCategoryByID(gomock.Any(), uint32(1)).
					Return(nil, errors.New("service unavailable")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
			expectedError: errors.New("service unavailable"),
		},
	}

	ctrl := gomock.NewController(t)
	toysService := mockservices.NewMockToysService(ctrl)
	useCases := &UseCases{
		toysService: toysService,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysService)
			}

			actual, err := useCases.GetCategoryByID(tc.args.ctx, tc.args.id)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.EqualError(t, err, tc.expectedError.Error())
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetAllCategories(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	// Test data
	testCategories := []entities.Category{
		{
			ID:   1,
			Name: "Category 1",
		},
		{
			ID:   2,
			Name: "Category 2",
		},
	}

	testCases := []struct {
		name          string
		args          args
		setupMocks    func(toysService *mockservices.MockToysService)
		expected      []entities.Category
		errorExpected bool
		expectedError error
	}{
		{
			name: "successful get all categories",
			args: args{
				ctx: context.Background(),
			},
			setupMocks: func(toysService *mockservices.MockToysService) {
				toysService.
					EXPECT().
					GetAllCategories(gomock.Any()).
					Return(testCategories, nil).
					Times(1)
			},
			expected:      testCategories,
			errorExpected: false,
		},
		{
			name: "empty categories list",
			args: args{
				ctx: context.Background(),
			},
			setupMocks: func(toysService *mockservices.MockToysService) {
				toysService.
					EXPECT().
					GetAllCategories(gomock.Any()).
					Return([]entities.Category{}, nil).
					Times(1)
			},
			expected:      []entities.Category{},
			errorExpected: false,
		},
		{
			name: "service error",
			args: args{
				ctx: context.Background(),
			},
			setupMocks: func(toysService *mockservices.MockToysService) {
				toysService.
					EXPECT().
					GetAllCategories(gomock.Any()).
					Return(nil, errors.New("service error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
			expectedError: errors.New("service error"),
		},
	}

	ctrl := gomock.NewController(t)
	toysService := mockservices.NewMockToysService(ctrl)
	useCases := &UseCases{
		toysService: toysService,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysService)
			}

			actual, err := useCases.GetAllCategories(tc.args.ctx)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.EqualError(t, err, tc.expectedError.Error())
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}
