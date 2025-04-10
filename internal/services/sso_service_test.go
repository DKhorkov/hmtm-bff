package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/libs/pointers"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	mockrepositories "github.com/DKhorkov/hmtm-bff/mocks/repositories"
)

func TestSsoService_GetAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoRepository := mockrepositories.NewMockSsoRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewSsoService(ssoRepository, logger)

	testCases := []struct {
		name          string
		setupMocks    func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger)
		expectedUsers []entities.User
		errorExpected bool
	}{
		{
			name: "success",
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					GetAllUsers(gomock.Any()).
					Return([]entities.User{{ID: 1, Email: "user@example.com"}}, nil).
					Times(1)
			},
			expectedUsers: []entities.User{{ID: 1, Email: "user@example.com"}},
			errorExpected: false,
		},
		{
			name: "error",
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					GetAllUsers(gomock.Any()).
					Return(nil, errors.New("fetch failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedUsers: nil,
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoRepository, logger)
			}

			users, err := service.GetAllUsers(context.Background())
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, users)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedUsers, users)
			}
		})
	}
}

func TestSsoService_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoRepository := mockrepositories.NewMockSsoRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewSsoService(ssoRepository, logger)

	testCases := []struct {
		name          string
		id            uint64
		setupMocks    func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger)
		expectedUser  *entities.User
		errorExpected bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(&entities.User{ID: 1, Email: "user@example.com"}, nil).
					Times(1)
			},
			expectedUser:  &entities.User{ID: 1, Email: "user@example.com"},
			errorExpected: false,
		},
		{
			name: "error",
			id:   1,
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					GetUserByID(gomock.Any(), uint64(1)).
					Return(nil, errors.New("not found")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedUser:  nil,
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoRepository, logger)
			}

			user, err := service.GetUserByID(context.Background(), tc.id)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedUser, user)
			}
		})
	}
}

func TestSsoService_GetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoRepository := mockrepositories.NewMockSsoRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewSsoService(ssoRepository, logger)

	testCases := []struct {
		name          string
		email         string
		setupMocks    func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger)
		expectedUser  *entities.User
		errorExpected bool
	}{
		{
			name:  "success",
			email: "user@example.com",
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					GetUserByEmail(gomock.Any(), "user@example.com").
					Return(&entities.User{ID: 1, Email: "user@example.com"}, nil).
					Times(1)
			},
			expectedUser:  &entities.User{ID: 1, Email: "user@example.com"},
			errorExpected: false,
		},
		{
			name:  "error",
			email: "user@example.com",
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					GetUserByEmail(gomock.Any(), "user@example.com").
					Return(nil, errors.New("not found")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedUser:  nil,
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoRepository, logger)
			}

			user, err := service.GetUserByEmail(context.Background(), tc.email)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedUser, user)
			}
		})
	}
}

func TestSsoService_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoRepository := mockrepositories.NewMockSsoRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewSsoService(ssoRepository, logger)

	testCases := []struct {
		name           string
		userData       entities.RegisterUserDTO
		setupMocks     func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger)
		expectedUserID uint64
		errorExpected  bool
	}{
		{
			name:     "success",
			userData: entities.RegisterUserDTO{Email: "user@example.com", Password: "password"},
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					RegisterUser(gomock.Any(), entities.RegisterUserDTO{Email: "user@example.com", Password: "password"}).
					Return(uint64(1), nil).
					Times(1)
			},
			expectedUserID: 1,
			errorExpected:  false,
		},
		{
			name:     "error",
			userData: entities.RegisterUserDTO{Email: "user@example.com", Password: "password"},
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					RegisterUser(gomock.Any(), entities.RegisterUserDTO{Email: "user@example.com", Password: "password"}).
					Return(uint64(0), errors.New("registration failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedUserID: 0,
			errorExpected:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoRepository, logger)
			}

			userID, err := service.RegisterUser(context.Background(), tc.userData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedUserID, userID)
		})
	}
}

func TestSsoService_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoRepository := mockrepositories.NewMockSsoRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewSsoService(ssoRepository, logger)

	testCases := []struct {
		name           string
		userData       entities.LoginUserDTO
		setupMocks     func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger)
		expectedTokens *entities.TokensDTO
		errorExpected  bool
	}{
		{
			name:     "success",
			userData: entities.LoginUserDTO{Email: "user@example.com", Password: "password"},
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					LoginUser(gomock.Any(), entities.LoginUserDTO{Email: "user@example.com", Password: "password"}).
					Return(&entities.TokensDTO{AccessToken: "access", RefreshToken: "refresh"}, nil).
					Times(1)
			},
			expectedTokens: &entities.TokensDTO{AccessToken: "access", RefreshToken: "refresh"},
			errorExpected:  false,
		},
		{
			name:     "error",
			userData: entities.LoginUserDTO{Email: "user@example.com", Password: "password"},
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					LoginUser(gomock.Any(), entities.LoginUserDTO{Email: "user@example.com", Password: "password"}).
					Return(nil, errors.New("login failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedTokens: nil,
			errorExpected:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoRepository, logger)
			}

			tokens, err := service.LoginUser(context.Background(), tc.userData)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, tokens)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedTokens, tokens)
			}
		})
	}
}

func TestSsoService_LogoutUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoRepository := mockrepositories.NewMockSsoRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewSsoService(ssoRepository, logger)

	testCases := []struct {
		name          string
		accessToken   string
		setupMocks    func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger)
		errorExpected bool
	}{
		{
			name:        "success",
			accessToken: "access-token",
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					LogoutUser(gomock.Any(), "access-token").
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:        "error",
			accessToken: "access-token",
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					LogoutUser(gomock.Any(), "access-token").
					Return(errors.New("logout failed")).
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
				tc.setupMocks(ssoRepository, logger)
			}

			err := service.LogoutUser(context.Background(), tc.accessToken)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSsoService_VerifyUserEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoRepository := mockrepositories.NewMockSsoRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewSsoService(ssoRepository, logger)

	testCases := []struct {
		name             string
		verifyEmailToken string
		setupMocks       func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger)
		errorExpected    bool
	}{
		{
			name:             "success",
			verifyEmailToken: "verify-token",
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					VerifyUserEmail(gomock.Any(), "verify-token").
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:             "error",
			verifyEmailToken: "verify-token",
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					VerifyUserEmail(gomock.Any(), "verify-token").
					Return(errors.New("verify failed")).
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
				tc.setupMocks(ssoRepository, logger)
			}

			err := service.VerifyUserEmail(context.Background(), tc.verifyEmailToken)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSsoService_ForgetPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoRepository := mockrepositories.NewMockSsoRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewSsoService(ssoRepository, logger)

	testCases := []struct {
		name                string
		forgetPasswordToken string
		newPassword         string
		setupMocks          func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger)
		errorExpected       bool
	}{
		{
			name:                "success",
			forgetPasswordToken: "forget-token",
			newPassword:         "newpass",
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					ForgetPassword(gomock.Any(), "forget-token", "newpass").
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:                "error",
			forgetPasswordToken: "forget-token",
			newPassword:         "newpass",
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					ForgetPassword(gomock.Any(), "forget-token", "newpass").
					Return(errors.New("forget failed")).
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
				tc.setupMocks(ssoRepository, logger)
			}

			err := service.ForgetPassword(context.Background(), tc.forgetPasswordToken, tc.newPassword)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSsoService_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoRepository := mockrepositories.NewMockSsoRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewSsoService(ssoRepository, logger)

	testCases := []struct {
		name          string
		accessToken   string
		oldPassword   string
		newPassword   string
		setupMocks    func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger)
		errorExpected bool
	}{
		{
			name:        "success",
			accessToken: "access-token",
			oldPassword: "oldpass",
			newPassword: "newpass",
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					ChangePassword(gomock.Any(), "access-token", "oldpass", "newpass").
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:        "error",
			accessToken: "access-token",
			oldPassword: "oldpass",
			newPassword: "newpass",
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					ChangePassword(gomock.Any(), "access-token", "oldpass", "newpass").
					Return(errors.New("change failed")).
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
				tc.setupMocks(ssoRepository, logger)
			}

			err := service.ChangePassword(context.Background(), tc.accessToken, tc.oldPassword, tc.newPassword)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSsoService_SendVerifyEmailMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoRepository := mockrepositories.NewMockSsoRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewSsoService(ssoRepository, logger)

	testCases := []struct {
		name          string
		email         string
		setupMocks    func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger)
		errorExpected bool
	}{
		{
			name:  "success",
			email: "user@example.com",
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					SendVerifyEmailMessage(gomock.Any(), "user@example.com").
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:  "error",
			email: "user@example.com",
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					SendVerifyEmailMessage(gomock.Any(), "user@example.com").
					Return(errors.New("send failed")).
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
				tc.setupMocks(ssoRepository, logger)
			}

			err := service.SendVerifyEmailMessage(context.Background(), tc.email)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSsoService_SendForgetPasswordMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoRepository := mockrepositories.NewMockSsoRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewSsoService(ssoRepository, logger)

	testCases := []struct {
		name          string
		email         string
		setupMocks    func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger)
		errorExpected bool
	}{
		{
			name:  "success",
			email: "user@example.com",
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					SendForgetPasswordMessage(gomock.Any(), "user@example.com").
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:  "error",
			email: "user@example.com",
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					SendForgetPasswordMessage(gomock.Any(), "user@example.com").
					Return(errors.New("send failed")).
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
				tc.setupMocks(ssoRepository, logger)
			}

			err := service.SendForgetPasswordMessage(context.Background(), tc.email)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSsoService_UpdateUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoRepository := mockrepositories.NewMockSsoRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewSsoService(ssoRepository, logger)

	testCases := []struct {
		name            string
		userProfileData entities.UpdateUserProfileDTO
		setupMocks      func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger)
		errorExpected   bool
	}{
		{
			name: "success",
			userProfileData: entities.UpdateUserProfileDTO{
				AccessToken: "access-token",
				Phone:       pointers.New("89612245678"),
				Telegram:    pointers.New("@test"),
				Avatar:      pointers.New("link"),
			},
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					UpdateUserProfile(
						gomock.Any(),
						entities.UpdateUserProfileDTO{
							AccessToken: "access-token",
							Phone:       pointers.New("89612245678"),
							Telegram:    pointers.New("@test"),
							Avatar:      pointers.New("link"),
						},
					).
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "error",
			userProfileData: entities.UpdateUserProfileDTO{
				AccessToken: "access-token",
				Phone:       pointers.New("89612245678"),
				Telegram:    pointers.New("@test"),
				Avatar:      pointers.New("link"),
			},
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					UpdateUserProfile(
						gomock.Any(),
						entities.UpdateUserProfileDTO{
							AccessToken: "access-token",
							Phone:       pointers.New("89612245678"),
							Telegram:    pointers.New("@test"),
							Avatar:      pointers.New("link"),
						},
					).
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
				tc.setupMocks(ssoRepository, logger)
			}

			err := service.UpdateUserProfile(context.Background(), tc.userProfileData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSsoService_GetMe(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoRepository := mockrepositories.NewMockSsoRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewSsoService(ssoRepository, logger)

	testCases := []struct {
		name          string
		accessToken   string
		setupMocks    func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger)
		expectedUser  *entities.User
		errorExpected bool
	}{
		{
			name:        "success",
			accessToken: "access-token",
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					GetMe(gomock.Any(), "access-token").
					Return(&entities.User{ID: 1, Email: "user@example.com"}, nil).
					Times(1)
			},
			expectedUser:  &entities.User{ID: 1, Email: "user@example.com"},
			errorExpected: false,
		},
		{
			name:        "error",
			accessToken: "access-token",
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					GetMe(gomock.Any(), "access-token").
					Return(nil, errors.New("not found")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedUser:  nil,
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoRepository, logger)
			}

			user, err := service.GetMe(context.Background(), tc.accessToken)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedUser, user)
			}
		})
	}
}

func TestSsoService_RefreshTokens(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoRepository := mockrepositories.NewMockSsoRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewSsoService(ssoRepository, logger)

	testCases := []struct {
		name           string
		refreshToken   string
		setupMocks     func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger)
		expectedTokens *entities.TokensDTO
		errorExpected  bool
	}{
		{
			name:         "success",
			refreshToken: "refresh-token",
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					RefreshTokens(gomock.Any(), "refresh-token").
					Return(&entities.TokensDTO{AccessToken: "new-access", RefreshToken: "new-refresh"}, nil).
					Times(1)
			},
			expectedTokens: &entities.TokensDTO{AccessToken: "new-access", RefreshToken: "new-refresh"},
			errorExpected:  false,
		},
		{
			name:         "error",
			refreshToken: "refresh-token",
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *mocklogging.MockLogger) {
				ssoRepository.
					EXPECT().
					RefreshTokens(gomock.Any(), "refresh-token").
					Return(nil, errors.New("refresh failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedTokens: nil,
			errorExpected:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoRepository, logger)
			}

			tokens, err := service.RefreshTokens(context.Background(), tc.refreshToken)
			if tc.errorExpected {
				require.Error(t, err)
				require.Nil(t, tokens)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedTokens, tokens)
			}
		})
	}
}
