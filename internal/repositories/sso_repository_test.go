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

	"github.com/DKhorkov/libs/pointers"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	mockclients "github.com/DKhorkov/hmtm-bff/mocks/clients"
	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
)

func TestSsoRepository_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoClient := mockclients.NewMockSsoClient(ctrl)
	repo := NewSsoRepository(ssoClient)

	testCases := []struct {
		name           string
		userData       entities.RegisterUserDTO
		setupMocks     func(ssoClient *mockclients.MockSsoClient)
		expectedUserID uint64
		errorExpected  bool
	}{
		{
			name: "success",
			userData: entities.RegisterUserDTO{
				DisplayName: "Test User",
				Email:       "test@example.com",
				Password:    "password123",
			},
			setupMocks: func(ssoClient *mockclients.MockSsoClient) {
				ssoClient.
					EXPECT().
					Register(
						gomock.Any(),
						&sso.RegisterIn{
							DisplayName: "Test User",
							Email:       "test@example.com",
							Password:    "password123",
						},
					).
					Return(&sso.RegisterOut{UserID: 1}, nil).
					Times(1)
			},
			expectedUserID: 1,
			errorExpected:  false,
		},
		{
			name: "error",
			userData: entities.RegisterUserDTO{
				Email: "test@example.com",
			},
			setupMocks: func(ssoClient *mockclients.MockSsoClient) {
				ssoClient.
					EXPECT().
					Register(
						gomock.Any(),
						&sso.RegisterIn{
							Email: "test@example.com",
						},
					).
					Return(nil, errors.New("registration failed")).
					Times(1)
			},
			expectedUserID: 0,
			errorExpected:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoClient)
			}

			userID, err := repo.RegisterUser(context.Background(), tc.userData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedUserID, userID)
		})
	}
}

func TestSsoRepository_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoClient := mockclients.NewMockSsoClient(ctrl)
	repo := NewSsoRepository(ssoClient)

	now := time.Now().UTC().Truncate(time.Second)

	testCases := []struct {
		name          string
		id            uint64
		setupMocks    func(ssoClient *mockclients.MockSsoClient)
		expectedUser  *entities.User
		errorExpected bool
	}{
		{
			name: "success",
			id:   1,
			setupMocks: func(ssoClient *mockclients.MockSsoClient) {
				ssoClient.
					EXPECT().
					GetUser(
						gomock.Any(),
						&sso.GetUserIn{ID: 1},
					).
					Return(&sso.GetUserOut{
						ID:          1,
						DisplayName: "Test User",
						Email:       "test@example.com",
						CreatedAt:   timestamppb.New(now),
						UpdatedAt:   timestamppb.New(now),
					}, nil).
					Times(1)
			},
			expectedUser: &entities.User{
				ID:          1,
				DisplayName: "Test User",
				Email:       "test@example.com",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			errorExpected: false,
		},
		{
			name: "error",
			id:   1,
			setupMocks: func(ssoClient *mockclients.MockSsoClient) {
				ssoClient.
					EXPECT().
					GetUser(
						gomock.Any(),
						&sso.GetUserIn{ID: 1},
					).
					Return(nil, errors.New("get user failed")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoClient)
			}

			user, err := repo.GetUserByID(context.Background(), tc.id)
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

func TestSsoRepository_GetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoClient := mockclients.NewMockSsoClient(ctrl)
	repo := NewSsoRepository(ssoClient)

	now := time.Now().UTC().Truncate(time.Second)

	testCases := []struct {
		name          string
		email         string
		setupMocks    func(ssoClient *mockclients.MockSsoClient)
		expectedUser  *entities.User
		errorExpected bool
	}{
		{
			name:  "success",
			email: "test@example.com",
			setupMocks: func(ssoClient *mockclients.MockSsoClient) {
				ssoClient.
					EXPECT().
					GetUserByEmail(
						gomock.Any(),
						&sso.GetUserByEmailIn{Email: "test@example.com"},
					).
					Return(&sso.GetUserOut{
						ID:          1,
						DisplayName: "Test User",
						Email:       "test@example.com",
						CreatedAt:   timestamppb.New(now),
						UpdatedAt:   timestamppb.New(now),
					}, nil).
					Times(1)
			},
			expectedUser: &entities.User{
				ID:          1,
				DisplayName: "Test User",
				Email:       "test@example.com",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			errorExpected: false,
		},
		{
			name:  "error",
			email: "test@example.com",
			setupMocks: func(ssoClient *mockclients.MockSsoClient) {
				ssoClient.
					EXPECT().
					GetUserByEmail(
						gomock.Any(),
						&sso.GetUserByEmailIn{Email: "test@example.com"},
					).
					Return(nil, errors.New("get user failed")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoClient)
			}

			user, err := repo.GetUserByEmail(context.Background(), tc.email)
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

func TestSsoRepository_GetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoClient := mockclients.NewMockSsoClient(ctrl)
	repo := NewSsoRepository(ssoClient)

	now := time.Now().UTC().Truncate(time.Second)

	testCases := []struct {
		name          string
		pagination    *entities.Pagination
		setupMocks    func(ssoClient *mockclients.MockSsoClient)
		expectedUsers []entities.User
		errorExpected bool
	}{
		{
			name: "success",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			setupMocks: func(ssoClient *mockclients.MockSsoClient) {
				ssoClient.
					EXPECT().
					GetUsers(
						gomock.Any(),
						&sso.GetUsersIn{
							Pagination: &sso.Pagination{
								Limit:  pointers.New[uint64](1),
								Offset: pointers.New[uint64](1),
							},
						},
					).
					Return(&sso.GetUsersOut{
						Users: []*sso.GetUserOut{
							{ID: 1, DisplayName: "User1", Email: "user1@example.com", CreatedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)},
							{ID: 2, DisplayName: "User2", Email: "user2@example.com", CreatedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)},
						},
					}, nil).
					Times(1)
			},
			expectedUsers: []entities.User{
				{ID: 1, DisplayName: "User1", Email: "user1@example.com", CreatedAt: now, UpdatedAt: now},
				{ID: 2, DisplayName: "User2", Email: "user2@example.com", CreatedAt: now, UpdatedAt: now},
			},
			errorExpected: false,
		},
		{
			name: "error",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			setupMocks: func(ssoClient *mockclients.MockSsoClient) {
				ssoClient.
					EXPECT().
					GetUsers(
						gomock.Any(),
						&sso.GetUsersIn{
							Pagination: &sso.Pagination{
								Limit:  pointers.New[uint64](1),
								Offset: pointers.New[uint64](1),
							},
						},
					).
					Return(nil, errors.New("get users failed")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoClient)
			}

			users, err := repo.GetUsers(context.Background(), tc.pagination)
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

func TestSsoRepository_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoClient := mockclients.NewMockSsoClient(ctrl)
	repo := NewSsoRepository(ssoClient)

	testCases := []struct {
		name           string
		userData       entities.LoginUserDTO
		setupMocks     func(ssoClient *mockclients.MockSsoClient)
		expectedTokens *entities.TokensDTO
		errorExpected  bool
	}{
		{
			name: "success",
			userData: entities.LoginUserDTO{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(ssoClient *mockclients.MockSsoClient) {
				ssoClient.
					EXPECT().
					Login(
						gomock.Any(),
						&sso.LoginIn{
							Email:    "test@example.com",
							Password: "password123",
						},
					).
					Return(&sso.LoginOut{
						AccessToken:  "access_token",
						RefreshToken: "refresh_token",
					}, nil).
					Times(1)
			},
			expectedTokens: &entities.TokensDTO{
				AccessToken:  "access_token",
				RefreshToken: "refresh_token",
			},
			errorExpected: false,
		},
		{
			name: "error",
			userData: entities.LoginUserDTO{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(ssoClient *mockclients.MockSsoClient) {
				ssoClient.
					EXPECT().
					Login(
						gomock.Any(),
						&sso.LoginIn{
							Email:    "test@example.com",
							Password: "password123",
						},
					).
					Return(nil, errors.New("login failed")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoClient)
			}

			tokens, err := repo.LoginUser(context.Background(), tc.userData)
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

func TestSsoRepository_LogoutUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoClient := mockclients.NewMockSsoClient(ctrl)
	repo := NewSsoRepository(ssoClient)

	testCases := []struct {
		name          string
		accessToken   string
		setupMocks    func(ssoClient *mockclients.MockSsoClient)
		errorExpected bool
	}{
		{
			name:        "success",
			accessToken: "access_token",
			setupMocks: func(ssoClient *mockclients.MockSsoClient) {
				ssoClient.
					EXPECT().
					Logout(
						gomock.Any(),
						&sso.LogoutIn{AccessToken: "access_token"},
					).
					Return(&emptypb.Empty{}, nil).
					Times(1)
			},
			errorExpected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoClient)
			}

			err := repo.LogoutUser(context.Background(), tc.accessToken)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSsoRepository_VerifyUserEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoClient := mockclients.NewMockSsoClient(ctrl)
	repo := NewSsoRepository(ssoClient)

	testCases := []struct {
		name          string
		token         string
		setupMocks    func(ssoClient *mockclients.MockSsoClient)
		errorExpected bool
	}{
		{
			name:  "success",
			token: "verify_token",
			setupMocks: func(ssoClient *mockclients.MockSsoClient) {
				ssoClient.
					EXPECT().
					VerifyEmail(
						gomock.Any(),
						&sso.VerifyEmailIn{VerifyEmailToken: "verify_token"},
					).
					Return(&emptypb.Empty{}, nil).
					Times(1)
			},
			errorExpected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoClient)
			}

			err := repo.VerifyUserEmail(context.Background(), tc.token)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSsoRepository_ForgetPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoClient := mockclients.NewMockSsoClient(ctrl)
	repo := NewSsoRepository(ssoClient)

	testCases := []struct {
		name          string
		token         string
		newPassword   string
		setupMocks    func(ssoClient *mockclients.MockSsoClient)
		errorExpected bool
	}{
		{
			name:        "success",
			token:       "forget_token",
			newPassword: "newpassword123",
			setupMocks: func(ssoClient *mockclients.MockSsoClient) {
				ssoClient.
					EXPECT().
					ForgetPassword(
						gomock.Any(),
						&sso.ForgetPasswordIn{
							ForgetPasswordToken: "forget_token",
							NewPassword:         "newpassword123",
						},
					).
					Return(&emptypb.Empty{}, nil).
					Times(1)
			},
			errorExpected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoClient)
			}

			err := repo.ForgetPassword(context.Background(), tc.token, tc.newPassword)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSsoRepository_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoClient := mockclients.NewMockSsoClient(ctrl)
	repo := NewSsoRepository(ssoClient)

	testCases := []struct {
		name          string
		accessToken   string
		oldPassword   string
		newPassword   string
		setupMocks    func(ssoClient *mockclients.MockSsoClient)
		errorExpected bool
	}{
		{
			name:        "success",
			accessToken: "access_token",
			oldPassword: "oldpassword",
			newPassword: "newpassword",
			setupMocks: func(ssoClient *mockclients.MockSsoClient) {
				ssoClient.
					EXPECT().
					ChangePassword(
						gomock.Any(),
						&sso.ChangePasswordIn{
							AccessToken: "access_token",
							OldPassword: "oldpassword",
							NewPassword: "newpassword",
						},
					).
					Return(&emptypb.Empty{}, nil).
					Times(1)
			},
			errorExpected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoClient)
			}

			err := repo.ChangePassword(context.Background(), tc.accessToken, tc.oldPassword, tc.newPassword)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSsoRepository_SendVerifyEmailMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoClient := mockclients.NewMockSsoClient(ctrl)
	repo := NewSsoRepository(ssoClient)

	testCases := []struct {
		name          string
		email         string
		setupMocks    func(ssoClient *mockclients.MockSsoClient)
		errorExpected bool
	}{
		{
			name:  "success",
			email: "test@example.com",
			setupMocks: func(ssoClient *mockclients.MockSsoClient) {
				ssoClient.
					EXPECT().
					SendVerifyEmailMessage(
						gomock.Any(),
						&sso.SendVerifyEmailMessageIn{Email: "test@example.com"},
					).
					Return(&emptypb.Empty{}, nil).
					Times(1)
			},
			errorExpected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoClient)
			}

			err := repo.SendVerifyEmailMessage(context.Background(), tc.email)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSsoRepository_SendForgetPasswordMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoClient := mockclients.NewMockSsoClient(ctrl)
	repo := NewSsoRepository(ssoClient)

	testCases := []struct {
		name          string
		email         string
		setupMocks    func(ssoClient *mockclients.MockSsoClient)
		errorExpected bool
	}{
		{
			name:  "success",
			email: "test@example.com",
			setupMocks: func(ssoClient *mockclients.MockSsoClient) {
				ssoClient.
					EXPECT().
					SendForgetPasswordMessage(
						gomock.Any(),
						&sso.SendForgetPasswordMessageIn{Email: "test@example.com"},
					).
					Return(&emptypb.Empty{}, nil).
					Times(1)
			},
			errorExpected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoClient)
			}

			err := repo.SendForgetPasswordMessage(context.Background(), tc.email)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSsoRepository_GetMe(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoClient := mockclients.NewMockSsoClient(ctrl)
	repo := NewSsoRepository(ssoClient)

	now := time.Now().UTC().Truncate(time.Second)

	testCases := []struct {
		name          string
		accessToken   string
		setupMocks    func(ssoClient *mockclients.MockSsoClient)
		expectedUser  *entities.User
		errorExpected bool
	}{
		{
			name:        "success",
			accessToken: "access_token",
			setupMocks: func(ssoClient *mockclients.MockSsoClient) {
				ssoClient.
					EXPECT().
					GetMe(
						gomock.Any(),
						&sso.GetMeIn{AccessToken: "access_token"},
					).
					Return(&sso.GetUserOut{
						ID:          1,
						DisplayName: "Test User",
						Email:       "test@example.com",
						CreatedAt:   timestamppb.New(now),
						UpdatedAt:   timestamppb.New(now),
					}, nil).
					Times(1)
			},
			expectedUser: &entities.User{
				ID:          1,
				DisplayName: "Test User",
				Email:       "test@example.com",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			errorExpected: false,
		},
		{
			name:        "error",
			accessToken: "access_token",
			setupMocks: func(ssoClient *mockclients.MockSsoClient) {
				ssoClient.
					EXPECT().
					GetMe(
						gomock.Any(),
						&sso.GetMeIn{AccessToken: "access_token"},
					).
					Return(nil, errors.New("test")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoClient)
			}

			user, err := repo.GetMe(context.Background(), tc.accessToken)
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

func TestSsoRepository_RefreshTokens(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoClient := mockclients.NewMockSsoClient(ctrl)
	repo := NewSsoRepository(ssoClient)

	testCases := []struct {
		name           string
		refreshToken   string
		setupMocks     func(ssoClient *mockclients.MockSsoClient)
		expectedTokens *entities.TokensDTO
		errorExpected  bool
	}{
		{
			name:         "success",
			refreshToken: "refresh_token",
			setupMocks: func(ssoClient *mockclients.MockSsoClient) {
				ssoClient.
					EXPECT().
					RefreshTokens(
						gomock.Any(),
						&sso.RefreshTokensIn{RefreshToken: "refresh_token"},
					).
					Return(
						&sso.LoginOut{
							AccessToken:  "new_access_token",
							RefreshToken: "new_refresh_token",
						},
						nil,
					).
					Times(1)
			},
			expectedTokens: &entities.TokensDTO{
				AccessToken:  "new_access_token",
				RefreshToken: "new_refresh_token",
			},
			errorExpected: false,
		},
		{
			name:         "error",
			refreshToken: "refresh_token",
			setupMocks: func(ssoClient *mockclients.MockSsoClient) {
				ssoClient.
					EXPECT().
					RefreshTokens(
						gomock.Any(),
						&sso.RefreshTokensIn{RefreshToken: "refresh_token"},
					).
					Return(nil, errors.New("test")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoClient)
			}

			tokens, err := repo.RefreshTokens(context.Background(), tc.refreshToken)
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

func TestSsoRepository_UpdateUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	ssoClient := mockclients.NewMockSsoClient(ctrl)
	repo := NewSsoRepository(ssoClient)

	testCases := []struct {
		name          string
		userData      entities.UpdateUserProfileDTO
		setupMocks    func(ssoClient *mockclients.MockSsoClient)
		errorExpected bool
	}{
		{
			name: "success",
			userData: entities.UpdateUserProfileDTO{
				AccessToken: "access_token",
				DisplayName: pointers.New("Updated User"),
				Phone:       pointers.New("1234567890"),
				Telegram:    pointers.New("@updateduser"),
				Avatar:      pointers.New("avatar_url"),
			},
			setupMocks: func(ssoClient *mockclients.MockSsoClient) {
				ssoClient.
					EXPECT().
					UpdateUserProfile(
						gomock.Any(),
						&sso.UpdateUserProfileIn{
							AccessToken: "access_token",
							DisplayName: pointers.New("Updated User"),
							Phone:       pointers.New("1234567890"),
							Telegram:    pointers.New("@updateduser"),
							Avatar:      pointers.New("avatar_url"),
						},
					).
					Return(&emptypb.Empty{}, nil).
					Times(1)
			},
			errorExpected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoClient)
			}

			err := repo.UpdateUserProfile(context.Background(), tc.userData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
