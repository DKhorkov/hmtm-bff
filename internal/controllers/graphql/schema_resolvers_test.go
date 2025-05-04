package graphqlcontroller

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/DKhorkov/libs/contextlib"
	"github.com/DKhorkov/libs/cookies"
	mocklogger "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/libs/middlewares"
	"github.com/DKhorkov/libs/pointers"

	graphqlapi "github.com/DKhorkov/hmtm-bff/api/graphql"
	"github.com/DKhorkov/hmtm-bff/internal/config"
	"github.com/DKhorkov/hmtm-bff/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-bff/internal/errors"
	mockusecases "github.com/DKhorkov/hmtm-bff/mocks/usecases"
)

var (
	ctx               = context.Background()
	userID     uint64 = 1
	masterID   uint64 = 1
	categoryID uint32 = 1
	tagID      uint32 = 1
	ticketID   uint64 = 1
	respondID  uint64 = 1
	toyID      uint64 = 1
	user              = &entities.User{
		ID:          userID,
		Email:       "user@example.com",
		DisplayName: "test",
	}

	now = time.Now()
)

func TestEmailResolver_User(t *testing.T) {
	testCases := []struct {
		name          string
		obj           *entities.Email
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *entities.User
		errorExpected bool
	}{
		{
			name: "success",
			obj: &entities.Email{
				UserID: userID,
				Email:  user.Email,
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetUserByID(gomock.Any(), userID).
					Return(user, nil).
					Times(1)
			},
			expected: user,
		},
		{
			name: "error",
			obj: &entities.Email{
				UserID: userID,
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetUserByID(gomock.Any(), userID).
					Return(nil, errors.New("error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name:     "nil obj",
			obj:      nil,
			expected: nil,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	mainResolver := NewResolver(useCases, logger, config.CookiesConfig{})
	testedResolver := &emailResolver{Resolver: mainResolver}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := testedResolver.User(ctx, tc.obj)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestMasterResolver_User(t *testing.T) {
	testCases := []struct {
		name          string
		obj           *entities.Master
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *entities.User
		errorExpected bool
	}{
		{
			name: "success",
			obj: &entities.Master{
				ID:     masterID,
				UserID: userID,
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetUserByID(gomock.Any(), userID).
					Return(user, nil).
					Times(1)
			},
			expected: user,
		},
		{
			name: "error",
			obj: &entities.Master{
				ID:     masterID,
				UserID: userID,
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetUserByID(gomock.Any(), userID).
					Return(nil, errors.New("error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name:     "nil obj",
			obj:      nil,
			expected: nil,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	mainResolver := NewResolver(useCases, logger, config.CookiesConfig{})
	testedResolver := &masterResolver{Resolver: mainResolver}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := testedResolver.User(ctx, tc.obj)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestMutationResolver_RegisterUser(t *testing.T) {
	validInput := graphqlapi.RegisterUserInput{
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Password:    "securepassword123",
	}

	testCases := []struct {
		name          string
		input         graphqlapi.RegisterUserInput
		setupMocks    func(useCases *mockusecases.MockUseCases)
		expected      string
		errorExpected bool
	}{
		{
			name:  "success",
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					RegisterUser(
						gomock.Any(),
						entities.RegisterUserDTO{
							DisplayName: validInput.DisplayName,
							Email:       validInput.Email,
							Password:    validInput.Password,
						},
					).
					Return(userID, nil).
					Times(1)
			},
			expected: strconv.FormatUint(userID, 10),
		},
		{
			name:  "error from useCases",
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					RegisterUser(
						gomock.Any(),
						gomock.Any(),
					).
					Return(uint64(0), errors.New("registration error")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	mainResolver := NewResolver(useCases, logger, config.CookiesConfig{})
	testedResolver := &mutationResolver{Resolver: mainResolver}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := testedResolver.RegisterUser(ctx, tc.input)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestMutationResolver_LoginUser(t *testing.T) {
	validInput := graphqlapi.LoginUserInput{
		Email:    user.Email,
		Password: "validpassword123",
	}

	confirmedUser := &entities.User{
		Email:          validInput.Email,
		EmailConfirmed: true,
	}

	unconfirmedUser := &entities.User{
		Email:          validInput.Email,
		EmailConfirmed: false,
	}

	validTokens := &entities.TokensDTO{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}

	testCases := []struct {
		name           string
		prepareContext func() context.Context
		input          graphqlapi.LoginUserInput
		setupMocks     func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected       bool
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful login",
			prepareContext: func() context.Context {
				recorder := httptest.NewRecorder()
				return contextlib.WithValue(ctx, middlewares.CookiesWriterName, recorder)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetUserByEmail(gomock.Any(), validInput.Email).
					Return(confirmedUser, nil).
					Times(1)

				useCases.
					EXPECT().
					LoginUser(
						gomock.Any(),
						entities.LoginUserDTO{
							Email:    validInput.Email,
							Password: validInput.Password,
						},
					).
					Return(validTokens, nil).
					Times(1)
			},
			expected: true,
		},
		{
			name:  "user not found",
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetUserByEmail(gomock.Any(), validInput.Email).
					Return(nil, errors.New("user not found")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name:  "email not confirmed",
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetUserByEmail(gomock.Any(), validInput.Email).
					Return(unconfirmedUser, nil).
					Times(1)
			},
			expectedError: &customerrors.PermissionDeniedError{
				Message: fmt.Sprintf("User with Email=%s has not confirmed it", validInput.Email),
			},
			errorExpected: true,
		},
		{
			name:  "invalid credentials",
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetUserByEmail(gomock.Any(), validInput.Email).
					Return(confirmedUser, nil).
					Times(1)

				useCases.
					EXPECT().
					LoginUser(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("invalid credentials")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name:  "missing cookies writer in context",
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetUserByEmail(gomock.Any(), validInput.Email).
					Return(confirmedUser, nil).
					Times(1)

				useCases.
					EXPECT().
					LoginUser(gomock.Any(), gomock.Any()).
					Return(validTokens, nil).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedError: &contextlib.ValueNotFoundError{Message: middlewares.CookiesWriterName},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	cookiesConfig := config.CookiesConfig{
		AccessToken:  cookies.Config{},
		RefreshToken: cookies.Config{},
	}

	mainResolver := NewResolver(useCases, logger, cookiesConfig)
	testedResolver := &mutationResolver{Resolver: mainResolver}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := ctx
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext()
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := testedResolver.LoginUser(testCtx, tc.input)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, tc.expectedError, err)
				}
			} else {
				require.NoError(t, err)

				// Verify cookies were set
				if recorder, ok := testCtx.Value(middlewares.CookiesWriterName).(*httptest.ResponseRecorder); ok {
					setCookies := recorder.Result().Cookies()
					require.Len(t, setCookies, 2)

					var accessTokenFound, refreshTokenFound bool
					for _, cookie := range setCookies {
						if cookie.Name == accessTokenCookieName {
							accessTokenFound = true
						}
						if cookie.Name == refreshTokenCookieName {
							refreshTokenFound = true
						}
					}

					require.True(t, accessTokenFound)
					require.True(t, refreshTokenFound)
				}
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestMutationResolver_LogoutUser(t *testing.T) {
	validAccessToken := &http.Cookie{
		Name:  accessTokenCookieName,
		Value: "valid_access_token",
	}

	testCases := []struct {
		name           string
		prepareContext func(ctx context.Context) context.Context
		setupMocks     func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected       bool
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful logout",
			prepareContext: func(ctx context.Context) context.Context {
				recorder := httptest.NewRecorder()
				ctx = contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
				return contextlib.WithValue(ctx, middlewares.CookiesWriterName, recorder)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					LogoutUser(gomock.Any(), validAccessToken.Value).
					Return(nil).
					Times(1)
			},
			expected: true,
		},
		{
			name:          "access token not found in context",
			expectedError: &cookies.NotFoundError{Message: accessTokenCookieName},
			errorExpected: true,
		},
		{
			name: "logout use case error",
			prepareContext: func(ctx context.Context) context.Context {
				ctx = contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
				return contextlib.WithValue(ctx, middlewares.CookiesWriterName, httptest.NewRecorder())
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					LogoutUser(gomock.Any(), validAccessToken.Value).
					Return(errors.New("logout error")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "cookies writer not found in context",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					LogoutUser(gomock.Any(), validAccessToken.Value).
					Return(nil).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedError: &contextlib.ValueNotFoundError{Message: middlewares.CookiesWriterName},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	cookiesConfig := config.CookiesConfig{
		AccessToken:  cookies.Config{MaxAge: -1},
		RefreshToken: cookies.Config{MaxAge: -1},
	}

	mainResolver := NewResolver(useCases, logger, cookiesConfig)
	testedResolver := &mutationResolver{Resolver: mainResolver}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := ctx
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := testedResolver.LogoutUser(testCtx)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, tc.expectedError, err)
				}
			} else {
				require.NoError(t, err)

				// Verify cookies were deleted
				if recorder, ok := testCtx.Value(middlewares.CookiesWriterName).(*httptest.ResponseRecorder); ok {
					setCookies := recorder.Result().Cookies()
					require.Len(t, setCookies, 2)

					var accessTokenDeleted, refreshTokenDeleted bool
					for _, cookie := range setCookies {
						if cookie.Name == accessTokenCookieName && cookie.MaxAge == -1 {
							accessTokenDeleted = true
						}
						if cookie.Name == refreshTokenCookieName && cookie.MaxAge == -1 {
							refreshTokenDeleted = true
						}
					}

					require.True(t, accessTokenDeleted)
					require.True(t, refreshTokenDeleted)
				}
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestMutationResolver_RefreshTokens(t *testing.T) {
	validRefreshToken := &http.Cookie{
		Name:  refreshTokenCookieName,
		Value: "valid_refresh_token",
	}

	newTokens := &entities.TokensDTO{
		AccessToken:  "new_access_token",
		RefreshToken: "new_refresh_token",
	}

	testCases := []struct {
		name           string
		prepareContext func(ctx context.Context) context.Context
		setupMocks     func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected       bool
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful refresh",
			prepareContext: func(ctx context.Context) context.Context {
				recorder := httptest.NewRecorder()
				ctx = contextlib.WithValue(ctx, refreshTokenCookieName, validRefreshToken)
				return contextlib.WithValue(ctx, middlewares.CookiesWriterName, recorder)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					RefreshTokens(gomock.Any(), validRefreshToken.Value).
					Return(newTokens, nil).
					Times(1)
			},
			expected: true,
		},
		{
			name:          "refresh token not found in context",
			expectedError: &cookies.NotFoundError{Message: refreshTokenCookieName},
			errorExpected: true,
		},
		{
			name: "refresh tokens use case error",
			prepareContext: func(ctx context.Context) context.Context {
				recorder := httptest.NewRecorder()
				ctx = contextlib.WithValue(ctx, refreshTokenCookieName, validRefreshToken)
				return contextlib.WithValue(ctx, middlewares.CookiesWriterName, recorder)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					RefreshTokens(gomock.Any(), validRefreshToken.Value).
					Return(nil, errors.New("refresh error")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "cookies writer not found in context",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, refreshTokenCookieName, validRefreshToken)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					RefreshTokens(gomock.Any(), validRefreshToken.Value).
					Return(newTokens, nil).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedError: &contextlib.ValueNotFoundError{Message: middlewares.CookiesWriterName},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	cookiesConfig := config.CookiesConfig{
		AccessToken:  cookies.Config{},
		RefreshToken: cookies.Config{},
	}

	mainResolver := NewResolver(useCases, logger, cookiesConfig)
	testedResolver := &mutationResolver{Resolver: mainResolver}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := ctx
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := testedResolver.RefreshTokens(testCtx)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, tc.expectedError, err)
				}
			} else {
				require.NoError(t, err)

				// Verify new cookies were set
				if recorder, ok := testCtx.Value(middlewares.CookiesWriterName).(*httptest.ResponseRecorder); ok {
					setCookies := recorder.Result().Cookies()
					require.Len(t, setCookies, 2)

					var accessTokenFound, refreshTokenFound bool
					for _, cookie := range setCookies {
						if cookie.Name == accessTokenCookieName && cookie.Value == newTokens.AccessToken {
							accessTokenFound = true
						}
						if cookie.Name == refreshTokenCookieName && cookie.Value == newTokens.RefreshToken {
							refreshTokenFound = true
						}
					}

					require.True(t, accessTokenFound)
					require.True(t, refreshTokenFound)
				}
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestMutationResolver_VerifyUserEmail(t *testing.T) {
	validInput := graphqlapi.VerifyUserEmailInput{
		VerifyEmailToken: "valid_token",
	}

	testCases := []struct {
		name          string
		input         graphqlapi.VerifyUserEmailInput
		setupMocks    func(useCases *mockusecases.MockUseCases)
		expected      bool
		errorExpected bool
	}{
		{
			name:  "successful verification",
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					VerifyUserEmail(gomock.Any(), validInput.VerifyEmailToken).
					Return(nil).
					Times(1)
			},
			expected: true,
		},
		{
			name:  "error",
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					VerifyUserEmail(gomock.Any(), validInput.VerifyEmailToken).
					Return(errors.New("invalid token")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	mainResolver := NewResolver(useCases, logger, config.CookiesConfig{})
	testedResolver := &mutationResolver{Resolver: mainResolver}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := testedResolver.VerifyUserEmail(context.Background(), tc.input)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestMutationResolver_SendVerifyEmailMessage(t *testing.T) {
	validInput := graphqlapi.SendVerifyEmailMessageInput{
		Email: "user@example.com",
	}

	testCases := []struct {
		name          string
		input         graphqlapi.SendVerifyEmailMessageInput
		setupMocks    func(useCases *mockusecases.MockUseCases)
		expected      bool
		errorExpected bool
	}{
		{
			name:  "successful email sending",
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					SendVerifyEmailMessage(gomock.Any(), validInput.Email).
					Return(nil).
					Times(1)
			},
			expected: true,
		},
		{
			name:  "error",
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					SendVerifyEmailMessage(gomock.Any(), validInput.Email).
					Return(errors.New("user not found")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	mainResolver := NewResolver(useCases, logger, config.CookiesConfig{})
	testedResolver := &mutationResolver{Resolver: mainResolver}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := testedResolver.SendVerifyEmailMessage(context.Background(), tc.input)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestMutationResolver_ForgetPassword(t *testing.T) {
	validInput := graphqlapi.ForgetPasswordInput{
		ForgetPasswordToken: "valid_token",
		NewPassword:         "new_secure_password123",
	}

	testCases := []struct {
		name          string
		input         graphqlapi.ForgetPasswordInput
		setupMocks    func(useCases *mockusecases.MockUseCases)
		expected      bool
		errorExpected bool
	}{
		{
			name:  "successful password reset",
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					ForgetPassword(
						gomock.Any(),
						validInput.ForgetPasswordToken,
						validInput.NewPassword,
					).
					Return(nil).
					Times(1)
			},
			expected: true,
		},
		{
			name:  "error",
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					ForgetPassword(
						gomock.Any(),
						validInput.ForgetPasswordToken,
						validInput.NewPassword,
					).
					Return(errors.New("test")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	mainResolver := NewResolver(useCases, logger, config.CookiesConfig{})
	testedResolver := &mutationResolver{Resolver: mainResolver}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := testedResolver.ForgetPassword(context.Background(), tc.input)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestMutationResolver_SendForgetPasswordMessage(t *testing.T) {
	validInput := graphqlapi.SendForgetPasswordMessageInput{
		Email: "user@example.com",
	}

	testCases := []struct {
		name          string
		input         graphqlapi.SendForgetPasswordMessageInput
		setupMocks    func(useCases *mockusecases.MockUseCases)
		expected      bool
		errorExpected bool
	}{
		{
			name:  "successful message sending",
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					SendForgetPasswordMessage(gomock.Any(), validInput.Email).
					Return(nil).
					Times(1)
			},
			expected: true,
		},
		{
			name:  "error",
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					SendForgetPasswordMessage(gomock.Any(), validInput.Email).
					Return(errors.New("test")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	mainResolver := NewResolver(useCases, logger, config.CookiesConfig{})
	testedResolver := &mutationResolver{Resolver: mainResolver}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := testedResolver.SendForgetPasswordMessage(context.Background(), tc.input)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestMutationResolver_ChangePassword(t *testing.T) {
	validInput := graphqlapi.ChangePasswordInput{
		OldPassword: "old_password123",
		NewPassword: "new_password456",
	}

	validAccessToken := &http.Cookie{
		Name:  accessTokenCookieName,
		Value: "valid_access_token",
	}

	testCases := []struct {
		name           string
		prepareContext func(ctx context.Context) context.Context
		input          graphqlapi.ChangePasswordInput
		setupMocks     func(useCases *mockusecases.MockUseCases)
		expected       bool
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful password change",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					ChangePassword(
						gomock.Any(),
						validAccessToken.Value,
						validInput.OldPassword,
						validInput.NewPassword,
					).
					Return(nil).
					Times(1)
			},
			expected: true,
		},
		{
			name:          "access token not found in context",
			input:         validInput,
			expectedError: &cookies.NotFoundError{Message: accessTokenCookieName},
			errorExpected: true,
		},
		{
			name: "invalid old password",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					ChangePassword(
						gomock.Any(),
						validAccessToken.Value,
						validInput.OldPassword,
						validInput.NewPassword,
					).
					Return(errors.New("error")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	mainResolver := NewResolver(useCases, logger, config.CookiesConfig{})
	testedResolver := &mutationResolver{Resolver: mainResolver}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := ctx
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := testedResolver.ChangePassword(testCtx, tc.input)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, tc.expectedError, err)
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestMutationResolver_UpdateUserProfile(t *testing.T) {
	validAccessToken := &http.Cookie{
		Name:  accessTokenCookieName,
		Value: "valid_access_token",
	}

	validInput := graphqlapi.UpdateUserProfileInput{
		DisplayName: pointers.New("New Name"),
		Phone:       pointers.New("+1234567890"),
		Telegram:    pointers.New("@telegram"),
	}

	testCases := []struct {
		name           string
		prepareContext func(ctx context.Context) context.Context
		input          graphqlapi.UpdateUserProfileInput
		setupMocks     func(useCases *mockusecases.MockUseCases)
		expected       bool
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful profile update",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					UpdateUserProfile(
						gomock.Any(),
						entities.RawUpdateUserProfileDTO{
							AccessToken: validAccessToken.Value,
							DisplayName: validInput.DisplayName,
							Phone:       validInput.Phone,
							Telegram:    validInput.Telegram,
							Avatar:      validInput.Avatar,
						},
					).
					Return(nil).
					Times(1)
			},
			expected: true,
		},
		{
			name:          "access token not found",
			input:         validInput,
			expectedError: &cookies.NotFoundError{Message: accessTokenCookieName},
			errorExpected: true,
		},
		{
			name: "useCases error",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: graphqlapi.UpdateUserProfileInput{
				DisplayName: validInput.DisplayName,
				Phone:       pointers.New("+1234567890"),
				Telegram:    validInput.Telegram,
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					UpdateUserProfile(gomock.Any(), gomock.Any()).
					Return(errors.New("error")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &mutationResolver{
		Resolver: NewResolver(
			useCases,
			logger, config.
				CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := ctx
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.UpdateUserProfile(testCtx, tc.input)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, err, tc.expectedError)
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestMutationResolver_RegisterMaster(t *testing.T) {
	validAccessToken := &http.Cookie{
		Name:  accessTokenCookieName,
		Value: "valid_access_token",
	}

	validInput := graphqlapi.RegisterMasterInput{
		Info: pointers.New("Master info text"),
	}

	testCases := []struct {
		name           string
		prepareContext func(ctx context.Context) context.Context
		input          graphqlapi.RegisterMasterInput
		setupMocks     func(useCases *mockusecases.MockUseCases)
		expected       string
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful master registration",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					RegisterMaster(
						gomock.Any(),
						entities.RawRegisterMasterDTO{
							AccessToken: validAccessToken.Value,
							Info:        validInput.Info,
						},
					).
					Return(masterID, nil).
					Times(1)
			},
			expected: strconv.FormatUint(masterID, 10),
		},
		{
			name:          "access token not found",
			input:         validInput,
			expectedError: &cookies.NotFoundError{Message: accessTokenCookieName},
			errorExpected: true,
		},
		{
			name: "use case error",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					RegisterMaster(gomock.Any(), gomock.Any()).
					Return(uint64(0), errors.New("registration error")).
					Times(1)
			},
			expected:      "0",
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &mutationResolver{
		Resolver: NewResolver(
			useCases,
			logger, config.
				CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := ctx
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.RegisterMaster(testCtx, tc.input)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, err, tc.expectedError)
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestMutationResolver_UpdateMaster(t *testing.T) {
	validAccessToken := &http.Cookie{
		Name:  accessTokenCookieName,
		Value: "valid_access_token",
	}

	validInput := graphqlapi.UpdateMasterInput{
		ID:   strconv.FormatUint(masterID, 10),
		Info: pointers.New("Updated master info"),
	}

	testCases := []struct {
		name           string
		prepareContext func(ctx context.Context) context.Context
		input          graphqlapi.UpdateMasterInput
		setupMocks     func(useCases *mockusecases.MockUseCases)
		expected       bool
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful master update",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					UpdateMaster(
						gomock.Any(),
						entities.RawUpdateMasterDTO{
							AccessToken: validAccessToken.Value,
							ID:          masterID,
							Info:        validInput.Info,
						},
					).
					Return(nil).
					Times(1)
			},
			expected: true,
		},
		{
			name:          "access token not found",
			input:         validInput,
			expectedError: &cookies.NotFoundError{Message: accessTokenCookieName},
			errorExpected: true,
		},
		{
			name: "invalid master ID format",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: graphqlapi.UpdateMasterInput{
				ID:   "invalid_id",
				Info: validInput.Info,
			},
			errorExpected: true,
		},
		{
			name: "use case error",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					UpdateMaster(gomock.Any(), gomock.Any()).
					Return(errors.New("update error")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &mutationResolver{
		Resolver: NewResolver(
			useCases,
			logger, config.
				CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := ctx
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.UpdateMaster(testCtx, tc.input)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, err, tc.expectedError)
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestMutationResolver_AddToy(t *testing.T) {
	validAccessToken := &http.Cookie{
		Name:  accessTokenCookieName,
		Value: "valid_access_token",
	}

	validInput := graphqlapi.AddToyInput{
		CategoryID:  strconv.Itoa(int(categoryID)),
		Name:        "New Toy",
		Description: "Toy description",
		Price:       19.99,
		Quantity:    10,
		Tags:        []string{"educational", "wooden"},
	}

	testCases := []struct {
		name           string
		prepareContext func(ctx context.Context) context.Context
		input          graphqlapi.AddToyInput
		setupMocks     func(useCases *mockusecases.MockUseCases)
		expected       string
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful toy addition",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					AddToy(
						gomock.Any(),
						entities.RawAddToyDTO{
							AccessToken: validAccessToken.Value,
							CategoryID:  categoryID,
							Name:        validInput.Name,
							Description: validInput.Description,
							Price:       19.99,
							Quantity:    10,
							Tags:        validInput.Tags,
							Attachments: nil,
						},
					).
					Return(toyID, nil).
					Times(1)
			},
			expected: strconv.FormatUint(toyID, 10),
		},
		{
			name:          "access token not found",
			input:         validInput,
			expectedError: &cookies.NotFoundError{Message: accessTokenCookieName},
			errorExpected: true,
		},
		{
			name: "invalid category ID format",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: graphqlapi.AddToyInput{
				CategoryID:  "invalid",
				Name:        validInput.Name,
				Description: validInput.Description,
				Price:       validInput.Price,
				Quantity:    validInput.Quantity,
			},
			errorExpected: true,
		},
		{
			name: "use case error",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					AddToy(gomock.Any(), gomock.Any()).
					Return(uint64(0), errors.New("add toy error")).
					Times(1)
			},
			expected:      "0",
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &mutationResolver{
		Resolver: NewResolver(
			useCases,
			logger, config.
				CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := ctx
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.AddToy(testCtx, tc.input)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, err, tc.expectedError)
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestMutationResolver_UpdateToy(t *testing.T) {
	validAccessToken := &http.Cookie{
		Name:  accessTokenCookieName,
		Value: "valid_access_token",
	}

	validName := "Updated Toy"
	validDescription := "Updated description"
	validPrice := 29.99
	validQuantity := 5
	validTags := []string{"educational", "wooden"}

	testCases := []struct {
		name           string
		prepareContext func(ctx context.Context) context.Context
		input          graphqlapi.UpdateToyInput
		setupMocks     func(useCases *mockusecases.MockUseCases)
		expected       bool
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful full update",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: graphqlapi.UpdateToyInput{
				ID:          strconv.FormatUint(toyID, 10),
				CategoryID:  pointers.New(strconv.FormatUint(uint64(categoryID), 10)),
				Name:        &validName,
				Description: &validDescription,
				Price:       &validPrice,
				Quantity:    &validQuantity,
				Tags:        validTags,
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				expectedDTO := entities.RawUpdateToyDTO{
					AccessToken: validAccessToken.Value,
					ID:          toyID,
					CategoryID:  &categoryID,
					Name:        &validName,
					Description: &validDescription,
					Price:       pointers.New[float32](29.99),
					Quantity:    pointers.New[uint32](5),
					Tags:        validTags,
					Attachments: nil,
				}
				useCases.
					EXPECT().
					UpdateToy(gomock.Any(), expectedDTO).
					Return(nil).
					Times(1)
			},
			expected: true,
		},
		{
			name: "successful partial update",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: graphqlapi.UpdateToyInput{
				ID:   strconv.FormatUint(toyID, 10),
				Name: &validName,
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				expectedDTO := entities.RawUpdateToyDTO{
					AccessToken: validAccessToken.Value,
					ID:          toyID,
					Name:        &validName,
				}
				useCases.
					EXPECT().
					UpdateToy(gomock.Any(), expectedDTO).
					Return(nil).
					Times(1)
			},
			expected: true,
		},
		{
			name:          "access token not found",
			expectedError: &cookies.NotFoundError{Message: accessTokenCookieName},
			errorExpected: true,
		},
		{
			name: "invalid toy ID format",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: graphqlapi.UpdateToyInput{
				ID: "invalid_id",
			},
			setupMocks:    nil,
			errorExpected: true,
		},
		{
			name: "invalid category ID format",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: graphqlapi.UpdateToyInput{
				ID:         strconv.FormatUint(toyID, 10),
				CategoryID: pointers.New("invalid"),
			},
			setupMocks:    nil,
			errorExpected: true,
		},
		{
			name: "use case error",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: graphqlapi.UpdateToyInput{
				ID:   strconv.FormatUint(toyID, 10),
				Name: &validName,
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					UpdateToy(gomock.Any(), gomock.Any()).
					Return(errors.New("update error")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &mutationResolver{
		Resolver: NewResolver(
			useCases,
			logger, config.
				CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := ctx
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.UpdateToy(testCtx, tc.input)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, err, tc.expectedError)
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestMutationResolver_DeleteToy(t *testing.T) {
	validAccessToken := &http.Cookie{
		Name:  accessTokenCookieName,
		Value: "valid_access_token",
	}

	testCases := []struct {
		name           string
		prepareContext func(ctx context.Context) context.Context
		input          graphqlapi.DeleteToyInput
		setupMocks     func(useCases *mockusecases.MockUseCases)
		expected       bool
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful toy deletion",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: graphqlapi.DeleteToyInput{ID: strconv.FormatUint(toyID, 10)},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					DeleteToy(
						gomock.Any(),
						validAccessToken.Value,
						toyID,
					).
					Return(nil).
					Times(1)
			},
			expected: true,
		},
		{
			name:          "access token not found",
			expectedError: &cookies.NotFoundError{Message: accessTokenCookieName},
			errorExpected: true,
		},
		{
			name: "invalid toy ID format",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: graphqlapi.DeleteToyInput{
				ID: "invalid_id",
			},
			setupMocks:    nil,
			errorExpected: true,
		},
		{
			name: "use case error",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: graphqlapi.DeleteToyInput{ID: strconv.FormatUint(toyID, 10)},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					DeleteToy(
						gomock.Any(),
						validAccessToken.Value,
						toyID,
					).
					Return(errors.New("test")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &mutationResolver{
		Resolver: NewResolver(
			useCases,
			logger, config.
				CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := ctx
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.DeleteToy(testCtx, tc.input)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, err, tc.expectedError)
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestMutationResolver_CreateTicket(t *testing.T) {
	validAccessToken := &http.Cookie{
		Name:  accessTokenCookieName,
		Value: "valid_access_token",
	}

	validInput := graphqlapi.CreateTicketInput{
		CategoryID:  strconv.FormatUint(uint64(categoryID), 10),
		Name:        "Test Ticket",
		Description: "Ticket description",
		Price:       pointers.New(29.99),
		Quantity:    5,
		Tags:        []string{"event", "concert"},
		Attachments: nil,
	}

	testCases := []struct {
		name           string
		prepareContext func(ctx context.Context) context.Context
		input          graphqlapi.CreateTicketInput
		setupMocks     func(useCases *mockusecases.MockUseCases)
		expected       string
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful ticket creation",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					CreateTicket(
						gomock.Any(),
						entities.RawCreateTicketDTO{
							AccessToken: validAccessToken.Value,
							CategoryID:  categoryID,
							Name:        validInput.Name,
							Description: validInput.Description,
							Price:       pointers.New[float32](float32(29.99)),
							Quantity:    uint32(validInput.Quantity),
							Tags:        validInput.Tags,
							Attachments: validInput.Attachments,
						},
					).
					Return(ticketID, nil).
					Times(1)
			},
			expected: strconv.FormatUint(ticketID, 10),
		},
		{
			name:          "access token not found",
			input:         validInput,
			expectedError: &cookies.NotFoundError{Message: accessTokenCookieName},
			errorExpected: true,
		},
		{
			name: "invalid category ID format",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: graphqlapi.CreateTicketInput{
				CategoryID:  "invalid",
				Name:        validInput.Name,
				Description: validInput.Description,
				Price:       validInput.Price,
				Quantity:    validInput.Quantity,
			},
			errorExpected: true,
		},
		{
			name: "use case error",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					CreateTicket(gomock.Any(), gomock.Any()).
					Return(uint64(0), errors.New("create ticket error")).
					Times(1)
			},
			expected:      "0",
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &mutationResolver{
		Resolver: NewResolver(
			useCases,
			logger, config.
				CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.CreateTicket(testCtx, tc.input)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, tc.expectedError, err)
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestMutationResolver_RespondToTicket(t *testing.T) {
	validAccessToken := &http.Cookie{
		Name:  accessTokenCookieName,
		Value: "valid_access_token",
	}

	validComment := "Response comment"
	validInput := graphqlapi.RespondToTicketInput{
		TicketID: strconv.FormatUint(ticketID, 10),
		Price:    45.99,
		Comment:  &validComment,
	}

	testCases := []struct {
		name           string
		prepareContext func(ctx context.Context) context.Context
		input          graphqlapi.RespondToTicketInput
		setupMocks     func(useCases *mockusecases.MockUseCases)
		expected       string
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful ticket response",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					RespondToTicket(
						gomock.Any(),
						entities.RawRespondToTicketDTO{
							AccessToken: validAccessToken.Value,
							TicketID:    ticketID,
							Price:       float32(validInput.Price),
							Comment:     validInput.Comment,
						},
					).
					Return(respondID, nil).
					Times(1)
			},
			expected: strconv.FormatUint(ticketID, 10),
		},
		{
			name:          "access token not found",
			input:         validInput,
			expectedError: &cookies.NotFoundError{Message: accessTokenCookieName},
			errorExpected: true,
		},
		{
			name: "invalid ticket ID format",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: graphqlapi.RespondToTicketInput{
				TicketID: "invalid",
				Price:    validInput.Price,
				Comment:  validInput.Comment,
			},
			errorExpected: true,
		},
		{
			name: "use case error",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					RespondToTicket(gomock.Any(), gomock.Any()).
					Return(uint64(0), errors.New("respond to ticket error")).
					Times(1)
			},
			expected:      "0",
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &mutationResolver{
		Resolver: NewResolver(
			useCases,
			logger, config.
				CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.RespondToTicket(testCtx, tc.input)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, tc.expectedError, err)
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestMutationResolver_UpdateRespond(t *testing.T) {
	validAccessToken := &http.Cookie{
		Name:  accessTokenCookieName,
		Value: "valid_access_token",
	}

	validPrice := 55.99
	validComment := "Updated comment"
	validInput := graphqlapi.UpdateRespondInput{
		ID:      strconv.FormatUint(respondID, 10),
		Price:   &validPrice,
		Comment: &validComment,
	}

	testCases := []struct {
		name           string
		prepareContext func(ctx context.Context) context.Context
		input          graphqlapi.UpdateRespondInput
		setupMocks     func(useCases *mockusecases.MockUseCases)
		expected       bool
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful respond update",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					UpdateRespond(
						gomock.Any(),
						entities.RawUpdateRespondDTO{
							AccessToken: validAccessToken.Value,
							ID:          respondID,
							Price:       pointers.New[float32](float32(validPrice)),
							Comment:     validInput.Comment,
						},
					).
					Return(nil).
					Times(1)
			},
			expected: true,
		},
		{
			name:          "access token not found",
			input:         validInput,
			expectedError: &cookies.NotFoundError{Message: accessTokenCookieName},
			errorExpected: true,
		},
		{
			name: "invalid respond ID format",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: graphqlapi.UpdateRespondInput{
				ID:      "invalid",
				Price:   validInput.Price,
				Comment: validInput.Comment,
			},
			errorExpected: true,
		},
		{
			name: "use case error",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					UpdateRespond(gomock.Any(), gomock.Any()).
					Return(errors.New("update respond error")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "update with nil price",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: graphqlapi.UpdateRespondInput{
				ID:      strconv.FormatUint(respondID, 10),
				Price:   nil,
				Comment: &validComment,
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					UpdateRespond(
						gomock.Any(),
						entities.RawUpdateRespondDTO{
							AccessToken: validAccessToken.Value,
							ID:          respondID,
							Price:       nil,
							Comment:     &validComment,
						},
					).
					Return(nil).
					Times(1)
			},
			expected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &mutationResolver{
		Resolver: NewResolver(
			useCases,
			logger, config.
				CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.UpdateRespond(testCtx, tc.input)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, tc.expectedError, err)
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestMutationResolver_DeleteRespond(t *testing.T) {
	validAccessToken := &http.Cookie{
		Name:  accessTokenCookieName,
		Value: "valid_access_token",
	}

	validInput := graphqlapi.DeleteRespondInput{
		ID: strconv.FormatUint(respondID, 10),
	}

	testCases := []struct {
		name           string
		prepareContext func(ctx context.Context) context.Context
		input          graphqlapi.DeleteRespondInput
		setupMocks     func(useCases *mockusecases.MockUseCases)
		expected       bool
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful respond deletion",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					DeleteRespond(
						gomock.Any(),
						validAccessToken.Value,
						respondID,
					).
					Return(nil).
					Times(1)
			},
			expected: true,
		},
		{
			name:          "access token not found",
			input:         validInput,
			expectedError: &cookies.NotFoundError{Message: accessTokenCookieName},
			errorExpected: true,
		},
		{
			name: "invalid respond ID format",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: graphqlapi.DeleteRespondInput{
				ID: "invalid",
			},
			errorExpected: true,
		},
		{
			name: "use case error",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					DeleteRespond(
						gomock.Any(),
						validAccessToken.Value,
						respondID,
					).
					Return(errors.New("delete respond error")).
					Times(1)
			},
			expected:      false,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &mutationResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.DeleteRespond(testCtx, tc.input)
			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, tc.expectedError, err)
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestMutationResolver_UpdateTicket(t *testing.T) {
	validAccessToken := &http.Cookie{
		Name:  accessTokenCookieName,
		Value: "valid_access_token",
	}

	validName := "Updated Ticket"
	validDescription := "Updated description"
	validPrice := 39.99
	validQuantity := 15
	validInput := graphqlapi.UpdateTicketInput{
		ID:          strconv.FormatUint(ticketID, 10),
		CategoryID:  pointers.New(strconv.FormatUint(uint64(categoryID), 10)),
		Name:        &validName,
		Description: &validDescription,
		Price:       &validPrice,
		Quantity:    &validQuantity,
		Tags:        []string{"updated", "ticket"},
		Attachments: nil,
	}

	testCases := []struct {
		name           string
		prepareContext func(ctx context.Context) context.Context
		input          graphqlapi.UpdateTicketInput
		setupMocks     func(useCases *mockusecases.MockUseCases)
		expected       bool
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful ticket update",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					UpdateTicket(
						gomock.Any(),
						entities.RawUpdateTicketDTO{
							AccessToken: validAccessToken.Value,
							ID:          ticketID,
							CategoryID:  &categoryID,
							Name:        &validName,
							Description: &validDescription,
							Price:       pointers.New[float32](float32(validPrice)),
							Quantity:    pointers.New[uint32](uint32(validQuantity)),
							Tags:        validInput.Tags,
							Attachments: nil,
						},
					).
					Return(nil).
					Times(1)
			},
			expected: true,
		},
		{
			name:          "access token not found",
			input:         validInput,
			expectedError: &cookies.NotFoundError{Message: accessTokenCookieName},
			errorExpected: true,
		},
		{
			name: "invalid ticket ID format",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: graphqlapi.UpdateTicketInput{
				ID: "invalid",
			},
			errorExpected: true,
		},
		{
			name: "invalid category ID format",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: graphqlapi.UpdateTicketInput{
				ID:         strconv.FormatUint(ticketID, 10),
				CategoryID: pointers.New("invalid"),
			},
			errorExpected: true,
		},
		{
			name: "use case error",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					UpdateTicket(gomock.Any(), gomock.Any()).
					Return(errors.New("update ticket error")).
					Times(1)
			},
			expected:      false,
			errorExpected: true,
		},
		{
			name: "update with minimal fields",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: graphqlapi.UpdateTicketInput{
				ID: strconv.FormatUint(ticketID, 10),
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					UpdateTicket(
						gomock.Any(),
						entities.RawUpdateTicketDTO{
							AccessToken: validAccessToken.Value,
							ID:          ticketID,
						},
					).
					Return(nil).
					Times(1)
			},
			expected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &mutationResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.UpdateTicket(testCtx, tc.input)
			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, tc.expectedError, err)
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestMutationResolver_DeleteTicket(t *testing.T) {
	validAccessToken := &http.Cookie{
		Name:  accessTokenCookieName,
		Value: "valid_access_token",
	}

	validInput := graphqlapi.DeleteTicketInput{
		ID: strconv.FormatUint(ticketID, 10),
	}

	testCases := []struct {
		name           string
		prepareContext func(ctx context.Context) context.Context
		input          graphqlapi.DeleteTicketInput
		setupMocks     func(useCases *mockusecases.MockUseCases)
		expected       bool
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful ticket deletion",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					DeleteTicket(
						gomock.Any(),
						validAccessToken.Value,
						ticketID,
					).
					Return(nil).
					Times(1)
			},
			expected: true,
		},
		{
			name:          "access token not found",
			input:         validInput,
			expectedError: &cookies.NotFoundError{Message: accessTokenCookieName},
			errorExpected: true,
		},
		{
			name: "invalid ticket ID format",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: graphqlapi.DeleteTicketInput{
				ID: "invalid",
			},
			errorExpected: true,
		},
		{
			name: "use case error",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			input: validInput,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					DeleteTicket(
						gomock.Any(),
						validAccessToken.Value,
						ticketID,
					).
					Return(errors.New("delete ticket error")).
					Times(1)
			},
			expected:      false,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &mutationResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.DeleteTicket(testCtx, tc.input)
			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, tc.expectedError, err)
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_Users(t *testing.T) {
	testUsers := []entities.User{
		*user,
	}

	testCases := []struct {
		name          string
		input         *graphqlapi.UsersInput
		setupMocks    func(useCases *mockusecases.MockUseCases)
		expected      []*entities.User
		errorExpected bool
	}{
		{
			name: "successful users retrieval",
			input: &graphqlapi.UsersInput{
				Pagination: &entities.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetUsers(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return(testUsers, nil).
					Times(1)
			},
			expected: []*entities.User{
				user,
			},
		},
		{
			name: "use case error",
			input: &graphqlapi.UsersInput{
				Pagination: &entities.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetUsers(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return(nil, errors.New("get users error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.Users(testCtx, tc.input)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_User(t *testing.T) {
	testCases := []struct {
		name          string
		id            string
		setupMocks    func(useCases *mockusecases.MockUseCases)
		expected      *entities.User
		errorExpected bool
	}{
		{
			name: "successful user retrieval",
			id:   strconv.FormatUint(userID, 10),
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetUserByID(gomock.Any(), userID).
					Return(user, nil).
					Times(1)
			},
			expected: user,
		},
		{
			name:          "invalid id format",
			id:            "invalid",
			errorExpected: true,
		},
		{
			name: "use case error",
			id:   strconv.FormatUint(userID, 10),
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetUserByID(gomock.Any(), userID).
					Return(nil, errors.New("get user error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.User(testCtx, tc.id)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_UserByEmail(t *testing.T) {
	testCases := []struct {
		name          string
		email         string
		setupMocks    func(useCases *mockusecases.MockUseCases)
		expected      *entities.User
		errorExpected bool
	}{
		{
			name:  "successful user retrieval",
			email: user.Email,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetUserByEmail(gomock.Any(), user.Email).
					Return(user, nil).
					Times(1)
			},
			expected: user,
		},
		{
			name:  "use case error",
			email: user.Email,
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetUserByEmail(gomock.Any(), user.Email).
					Return(nil, errors.New("get user error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.UserByEmail(testCtx, tc.email)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_Me(t *testing.T) {
	validAccessToken := &http.Cookie{
		Name:  accessTokenCookieName,
		Value: "valid_access_token",
	}

	testCases := []struct {
		name           string
		prepareContext func(ctx context.Context) context.Context
		setupMocks     func(useCases *mockusecases.MockUseCases)
		expected       *entities.User
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful get current user",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMe(gomock.Any(), validAccessToken.Value).
					Return(user, nil).
					Times(1)
			},
			expected: user,
		},
		{
			name:          "access token not found",
			expectedError: &cookies.NotFoundError{Message: accessTokenCookieName},
			errorExpected: true,
		},
		{
			name: "use case error",
			prepareContext: func(ctx context.Context) context.Context {
				invalidToken := &http.Cookie{
					Name:  accessTokenCookieName,
					Value: "invalid_token",
				}
				return contextlib.WithValue(ctx, accessTokenCookieName, invalidToken)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMe(gomock.Any(), "invalid_token").
					Return(nil, errors.New("test")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := ctx
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.Me(testCtx)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, err, tc.expectedError)
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_MyToys(t *testing.T) {
	validAccessToken := &http.Cookie{
		Name:  accessTokenCookieName,
		Value: "valid_access_token",
	}

	mockToys := []entities.Toy{
		{
			ID:          toyID,
			Name:        "Toy 1",
			Description: "Description 1",
			Price:       19.99,
		},
	}

	testCases := []struct {
		name           string
		input          *graphqlapi.MyToysInput
		prepareContext func(ctx context.Context) context.Context
		setupMocks     func(useCases *mockusecases.MockUseCases)
		expected       []*entities.Toy
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful get my toys",
			input: &graphqlapi.MyToysInput{
				Pagination: &entities.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
			},
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMyToys(
						gomock.Any(),
						validAccessToken.Value,
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return(mockToys, nil).
					Times(1)
			},
			expected: []*entities.Toy{&mockToys[0]},
		},
		{
			name: "empty toys list",
			input: &graphqlapi.MyToysInput{
				Pagination: &entities.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
			},
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMyToys(
						gomock.Any(),
						validAccessToken.Value,
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return([]entities.Toy{}, nil).
					Times(1)
			},
			expected: []*entities.Toy{},
		},
		{
			name:          "access token not found",
			expectedError: &cookies.NotFoundError{Message: accessTokenCookieName},
			errorExpected: true,
		},
		{
			name: "use case error",
			input: &graphqlapi.MyToysInput{
				Pagination: &entities.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
			},
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMyToys(
						gomock.Any(),
						validAccessToken.Value,
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return(nil, errors.New("error")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := ctx
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.MyToys(testCtx, tc.input)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, err, tc.expectedError)
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_MyTickets(t *testing.T) {
	validAccessToken := &http.Cookie{
		Name:  accessTokenCookieName,
		Value: "valid_access_token",
	}

	mockTickets := []entities.Ticket{
		{
			ID:          ticketID,
			Name:        "Ticket 1",
			Description: "open",
		},
	}

	testCases := []struct {
		name           string
		prepareContext func(ctx context.Context) context.Context
		setupMocks     func(useCases *mockusecases.MockUseCases)
		expected       []*entities.Ticket
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful get my tickets",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMyTickets(gomock.Any(), validAccessToken.Value).
					Return(mockTickets, nil).
					Times(1)
			},
			expected: []*entities.Ticket{&mockTickets[0]},
		},
		{
			name: "empty tickets list",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMyTickets(gomock.Any(), validAccessToken.Value).
					Return([]entities.Ticket{}, nil).
					Times(1)
			},
			expected: []*entities.Ticket{},
		},
		{
			name:          "access token not found",
			expectedError: &cookies.NotFoundError{Message: accessTokenCookieName},
			errorExpected: true,
		},
		{
			name: "use case error",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMyTickets(gomock.Any(), validAccessToken.Value).
					Return(nil, errors.New("error")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := ctx
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.MyTickets(testCtx)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, err, tc.expectedError)
				}
			} else {
				require.NoError(t, err)
			}

			require.Len(t, actual, len(tc.expected))
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_MyResponds(t *testing.T) {
	validAccessToken := &http.Cookie{
		Name:  accessTokenCookieName,
		Value: "valid_access_token",
	}

	mockResponds := []entities.Respond{
		{
			ID:        respondID,
			TicketID:  ticketID,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	testCases := []struct {
		name           string
		prepareContext func(ctx context.Context) context.Context
		setupMocks     func(useCases *mockusecases.MockUseCases)
		expected       []*entities.Respond
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful get my responds",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMyResponds(gomock.Any(), validAccessToken.Value).
					Return(mockResponds, nil).
					Times(1)
			},
			expected: []*entities.Respond{&mockResponds[0]},
		},
		{
			name: "empty responds list",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMyResponds(gomock.Any(), validAccessToken.Value).
					Return([]entities.Respond{}, nil).
					Times(1)
			},
			expected: []*entities.Respond{},
		},
		{
			name:          "access token not found",
			expectedError: &cookies.NotFoundError{Message: accessTokenCookieName},
			errorExpected: true,
		},
		{
			name: "use case error",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMyResponds(gomock.Any(), validAccessToken.Value).
					Return(nil, errors.New("test")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := ctx
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.MyResponds(testCtx)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, err, tc.expectedError)
				}
			} else {
				require.NoError(t, err)
			}

			require.Len(t, actual, len(tc.expected))
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_MyEmailCommunications(t *testing.T) {
	validAccessToken := &http.Cookie{
		Name:  accessTokenCookieName,
		Value: "valid_access_token",
	}

	mockEmails := []entities.Email{
		{
			ID:     1,
			UserID: userID,
			SentAt: now,
		},
	}

	testCases := []struct {
		name           string
		prepareContext func(ctx context.Context) context.Context
		setupMocks     func(useCases *mockusecases.MockUseCases)
		expected       []*entities.Email
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful get email communications",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMyEmailCommunications(gomock.Any(), validAccessToken.Value).
					Return(mockEmails, nil).
					Times(1)
			},
			expected: []*entities.Email{&mockEmails[0]},
		},
		{
			name: "empty email communications list",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMyEmailCommunications(gomock.Any(), validAccessToken.Value).
					Return([]entities.Email{}, nil).
					Times(1)
			},
			expected: []*entities.Email{},
		},
		{
			name:          "access token not found",
			expectedError: &cookies.NotFoundError{Message: accessTokenCookieName},
			errorExpected: true,
		},
		{
			name: "use case error",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMyEmailCommunications(gomock.Any(), validAccessToken.Value).
					Return(nil, errors.New("test")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := ctx
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.MyEmailCommunications(testCtx)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, err, tc.expectedError)
				}
			} else {
				require.NoError(t, err)
			}

			require.Len(t, actual, len(tc.expected))
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_Master(t *testing.T) {
	testCases := []struct {
		name          string
		id            string
		setupMocks    func(useCases *mockusecases.MockUseCases)
		expected      *entities.Master
		errorExpected bool
	}{
		{
			name: "successful get master by ID",
			id:   strconv.FormatUint(masterID, 10),
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMasterByID(gomock.Any(), masterID).
					Return(
						&entities.Master{
							ID:        masterID,
							UserID:    userID,
							Info:      pointers.New("test"),
							CreatedAt: now,
						},
						nil,
					).
					Times(1)
			},
			expected: &entities.Master{
				ID:        masterID,
				UserID:    userID,
				Info:      pointers.New("test"),
				CreatedAt: now,
			},
		},
		{
			name:          "invalid master ID format",
			id:            "invalid",
			errorExpected: true,
		},
		{
			name: "use case error",
			id:   strconv.FormatUint(masterID, 10),
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMasterByID(gomock.Any(), masterID).
					Return(nil, errors.New("error")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.Master(ctx, tc.id)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_Masters(t *testing.T) {
	mockMasters := []entities.Master{
		{
			ID:        masterID,
			UserID:    userID,
			Info:      pointers.New("test"),
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	testCases := []struct {
		name          string
		input         *graphqlapi.MastersInput
		setupMocks    func(useCases *mockusecases.MockUseCases)
		expected      []*entities.Master
		errorExpected bool
	}{
		{
			name: "successful get all masters",
			input: &graphqlapi.MastersInput{
				Pagination: &entities.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMasters(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return(mockMasters, nil).
					Times(1)
			},
			expected: []*entities.Master{&mockMasters[0]},
		},
		{
			name: "empty masters list",
			input: &graphqlapi.MastersInput{
				Pagination: &entities.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMasters(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return([]entities.Master{}, nil).
					Times(1)
			},
			expected: []*entities.Master{},
		},
		{
			name: "use case error",
			input: &graphqlapi.MastersInput{
				Pagination: &entities.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMasters(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return(nil, errors.New("test")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.Masters(ctx, tc.input)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Len(t, actual, len(tc.expected))
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_MasterToys(t *testing.T) {
	mockToys := []entities.Toy{
		{
			ID:          toyID,
			Name:        "Test Toy",
			Description: "Test description",
			Price:       19.99,
			MasterID:    masterID,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	testCases := []struct {
		name          string
		input         graphqlapi.MasterToysInput
		setupMocks    func(useCases *mockusecases.MockUseCases)
		expected      []*entities.Toy
		errorExpected bool
	}{
		{
			name: "successful get master toys",
			input: graphqlapi.MasterToysInput{
				MasterID: strconv.FormatUint(masterID, 10),
				Pagination: &entities.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMasterToys(
						gomock.Any(),
						masterID,
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return(mockToys, nil).
					Times(1)
			},
			expected: []*entities.Toy{&mockToys[0]},
		},
		{
			name: "empty toys list",
			input: graphqlapi.MasterToysInput{
				MasterID: strconv.FormatUint(masterID, 10),
				Pagination: &entities.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMasterToys(
						gomock.Any(),
						masterID,
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return([]entities.Toy{}, nil).
					Times(1)
			},
			expected: []*entities.Toy{},
		},
		{
			name: "invalid master ID",
			input: graphqlapi.MasterToysInput{
				MasterID: "invalid",
				Pagination: &entities.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
			},
			setupMocks:    func(useCases *mockusecases.MockUseCases) {},
			errorExpected: true,
		},
		{
			name: "use case error",
			input: graphqlapi.MasterToysInput{
				MasterID: strconv.FormatUint(masterID, 10),
				Pagination: &entities.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMasterToys(
						gomock.Any(),
						masterID,
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return(nil, errors.New("test")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.MasterToys(ctx, tc.input)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Len(t, actual, len(tc.expected))
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_Toy(t *testing.T) {
	validToy := entities.Toy{
		ID:          1,
		Name:        "Test Toy",
		Description: "Toy description",
		Price:       19.99,
		Quantity:    10,
	}

	testCases := []struct {
		name          string
		id            string
		setupMocks    func(useCases *mockusecases.MockUseCases)
		expected      *entities.Toy
		expectedError error
		errorExpected bool
	}{
		{
			name: "successful toy retrieval",
			id:   strconv.FormatUint(toyID, 10),
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetToyByID(gomock.Any(), toyID).
					Return(&validToy, nil).
					Times(1)
			},
			expected: &validToy,
		},
		{
			name:          "invalid id format",
			id:            "invalid",
			errorExpected: true,
		},
		{
			name: "use case error",
			id:   strconv.FormatUint(toyID, 10),
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetToyByID(gomock.Any(), toyID).
					Return(nil, errors.New("get toy error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.Toy(testCtx, tc.id)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_Toys(t *testing.T) {
	testToys := []entities.Toy{
		{
			ID:          toyID,
			Name:        "Toy1",
			Description: "Description1",
			Price:       19.99,
			Quantity:    10,
		},
	}

	testCases := []struct {
		name          string
		input         *graphqlapi.ToysInput
		setupMocks    func(useCases *mockusecases.MockUseCases)
		expected      []*entities.Toy
		expectedError error
		errorExpected bool
	}{
		{
			name: "successful toys retrieval",
			input: &graphqlapi.ToysInput{
				Pagination: &entities.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetToys(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return(testToys, nil).
					Times(1)
			},
			expected: []*entities.Toy{
				&testToys[0],
			},
		},
		{
			name: "use case error",
			input: &graphqlapi.ToysInput{
				Pagination: &entities.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetToys(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return(nil, errors.New("get toys error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.Toys(testCtx, tc.input)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_Tag(t *testing.T) {
	validTag := entities.Tag{
		ID:   tagID,
		Name: "TestTag",
	}

	testCases := []struct {
		name          string
		id            string
		setupMocks    func(useCases *mockusecases.MockUseCases)
		expected      *entities.Tag
		expectedError error
		errorExpected bool
	}{
		{
			name: "successful tag retrieval",
			id:   strconv.Itoa(int(tagID)),
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetTagByID(gomock.Any(), tagID).
					Return(&validTag, nil).
					Times(1)
			},
			expected: &validTag,
		},
		{
			name:          "invalid id format",
			id:            "invalid",
			errorExpected: true,
		},
		{
			name: "use case error",
			id:   strconv.Itoa(int(tagID)),
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetTagByID(gomock.Any(), tagID).
					Return(nil, errors.New("get tag error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.Tag(testCtx, tc.id)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_Tags(t *testing.T) {
	testTags := []entities.Tag{
		{
			ID:   tagID,
			Name: "Tag1",
		},
	}

	testCases := []struct {
		name          string
		setupMocks    func(useCases *mockusecases.MockUseCases)
		expected      []*entities.Tag
		expectedError error
		errorExpected bool
	}{
		{
			name: "successful tags retrieval",
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(testTags, nil).
					Times(1)
			},
			expected: []*entities.Tag{
				&testTags[0],
			},
		},
		{
			name: "use case error",
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(nil, errors.New("get tags error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.Tags(testCtx)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_Category(t *testing.T) {
	validCategory := entities.Category{
		ID:   categoryID,
		Name: "TestCategory",
	}

	testCases := []struct {
		name          string
		id            string
		setupMocks    func(useCases *mockusecases.MockUseCases)
		expected      *entities.Category
		expectedError error
		errorExpected bool
	}{
		{
			name: "successful category retrieval",
			id:   strconv.Itoa(int(categoryID)),
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetCategoryByID(gomock.Any(), categoryID).
					Return(&validCategory, nil).
					Times(1)
			},
			expected: &validCategory,
		},
		{
			name:          "invalid id format",
			id:            "invalid",
			errorExpected: true,
		},
		{
			name: "use case error",
			id:   strconv.Itoa(int(categoryID)),
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetCategoryByID(gomock.Any(), categoryID).
					Return(nil, errors.New("get category error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.Category(testCtx, tc.id)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_Categories(t *testing.T) {
	testCategories := []entities.Category{
		{
			ID:   categoryID,
			Name: "Category1",
		},
	}

	testCases := []struct {
		name          string
		setupMocks    func(useCases *mockusecases.MockUseCases)
		expected      []*entities.Category
		expectedError error
		errorExpected bool
	}{
		{
			name: "successful categories retrieval",
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetAllCategories(gomock.Any()).
					Return(testCategories, nil).
					Times(1)
			},
			expected: []*entities.Category{
				&testCategories[0],
			},
		},
		{
			name: "use case error",
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetAllCategories(gomock.Any()).
					Return(nil, errors.New("get categories error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.Categories(testCtx)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_Ticket(t *testing.T) {
	validTicket := entities.Ticket{
		ID:          ticketID,
		Name:        "Test Ticket",
		Description: "Ticket description",
		Price:       pointers.New[float32](29.99),
		Quantity:    5,
	}

	testCases := []struct {
		name          string
		id            string
		setupMocks    func(useCases *mockusecases.MockUseCases)
		expected      *entities.Ticket
		expectedError error
		errorExpected bool
	}{
		{
			name: "successful ticket retrieval",
			id:   strconv.FormatUint(ticketID, 10),
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetTicketByID(gomock.Any(), ticketID).
					Return(&validTicket, nil).
					Times(1)
			},
			expected: &validTicket,
		},
		{
			name:          "invalid id format",
			id:            "invalid",
			errorExpected: true,
		},
		{
			name: "use case error",
			id:   strconv.FormatUint(ticketID, 10),
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetTicketByID(gomock.Any(), ticketID).
					Return(nil, errors.New("get ticket error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.Ticket(testCtx, tc.id)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_Tickets(t *testing.T) {
	testTickets := []entities.Ticket{
		{
			ID:          ticketID,
			Name:        "Ticket2",
			Description: "Description2",
			Price:       pointers.New[float32](39.99),
			Quantity:    3,
		},
	}

	testCases := []struct {
		name          string
		setupMocks    func(useCases *mockusecases.MockUseCases)
		expected      []*entities.Ticket
		expectedError error
		errorExpected bool
	}{
		{
			name: "successful tickets retrieval",
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetAllTickets(gomock.Any()).
					Return(testTickets, nil).
					Times(1)
			},
			expected: []*entities.Ticket{
				&testTickets[0],
			},
		},
		{
			name: "use case error",
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetAllTickets(gomock.Any()).
					Return(nil, errors.New("get tickets error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.Tickets(testCtx)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_UserTickets(t *testing.T) {
	testTickets := []entities.Ticket{
		{
			ID:          ticketID,
			UserID:      userID,
			Name:        "Ticket1",
			Description: "Description1",
			Price:       pointers.New[float32](39.99),
			Quantity:    5,
		},
	}

	testCases := []struct {
		name          string
		userID        string
		setupMocks    func(useCases *mockusecases.MockUseCases)
		expected      []*entities.Ticket
		expectedError error
		errorExpected bool
	}{
		{
			name:   "successful user tickets retrieval",
			userID: strconv.FormatUint(userID, 10),
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetUserTickets(gomock.Any(), userID).
					Return(testTickets, nil).
					Times(1)
			},
			expected: []*entities.Ticket{
				&testTickets[0],
			},
		},
		{
			name:          "invalid user ID format",
			userID:        "invalid",
			errorExpected: true,
		},
		{
			name:   "use case error",
			userID: strconv.FormatUint(userID, 10),
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetUserTickets(gomock.Any(), userID).
					Return(nil, errors.New("get user tickets error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.UserTickets(testCtx, tc.userID)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_Respond(t *testing.T) {
	validAccessToken := &http.Cookie{
		Name:  accessTokenCookieName,
		Value: "valid_access_token",
	}

	validRespond := entities.Respond{
		ID:       respondID,
		TicketID: ticketID,
		Price:    45.99,
		Comment:  pointers.New("Test comment"),
	}

	testCases := []struct {
		name           string
		id             string
		prepareContext func(ctx context.Context) context.Context
		setupMocks     func(useCases *mockusecases.MockUseCases)
		expected       *entities.Respond
		expectedError  error
		errorExpected  bool
	}{
		{
			name: "successful respond retrieval",
			id:   strconv.FormatUint(respondID, 10),
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetRespondByID(
						gomock.Any(),
						respondID,
						validAccessToken.Value,
					).
					Return(&validRespond, nil).
					Times(1)
			},
			expected: &validRespond,
		},
		{
			name:          "access token not found",
			id:            strconv.FormatUint(respondID, 10),
			expectedError: &cookies.NotFoundError{Message: accessTokenCookieName},
			errorExpected: true,
		},
		{
			name: "invalid id format",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			id:            "invalid",
			errorExpected: true,
		},
		{
			name: "use case error",
			id:   strconv.FormatUint(respondID, 10),
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetRespondByID(
						gomock.Any(),
						respondID,
						validAccessToken.Value,
					).
					Return(nil, errors.New("get respond error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.Respond(testCtx, tc.id)
			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, tc.expectedError, err)
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_TicketResponds(t *testing.T) {
	validAccessToken := &http.Cookie{
		Name:  accessTokenCookieName,
		Value: "valid_access_token",
	}

	testResponds := []entities.Respond{
		{
			ID:       respondID,
			TicketID: ticketID,
			Price:    45.99,
			Comment:  pointers.New("Test comment"),
		},
	}

	testCases := []struct {
		name           string
		ticketID       string
		prepareContext func(ctx context.Context) context.Context
		setupMocks     func(useCases *mockusecases.MockUseCases)
		expected       []*entities.Respond
		expectedError  error
		errorExpected  bool
	}{
		{
			name:     "successful ticket responds retrieval",
			ticketID: strconv.FormatUint(ticketID, 10),
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetTicketResponds(
						gomock.Any(),
						ticketID,
						validAccessToken.Value,
					).
					Return(testResponds, nil).
					Times(1)
			},
			expected: []*entities.Respond{
				&testResponds[0],
			},
		},
		{
			name:          "access token not found",
			ticketID:      strconv.FormatUint(ticketID, 10),
			expectedError: &cookies.NotFoundError{Message: accessTokenCookieName},
			errorExpected: true,
		},
		{
			name: "invalid ticket ID format",
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			ticketID:      "invalid",
			errorExpected: true,
		},
		{
			name:     "use case error",
			ticketID: strconv.FormatUint(ticketID, 10),
			prepareContext: func(ctx context.Context) context.Context {
				return contextlib.WithValue(ctx, accessTokenCookieName, validAccessToken)
			},
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetTicketResponds(
						gomock.Any(),
						ticketID,
						validAccessToken.Value,
					).
					Return(nil, errors.New("get ticket responds error")).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.TicketResponds(testCtx, tc.ticketID)
			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, tc.expectedError, err)
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestRespondResolver_Ticket(t *testing.T) {
	validRespond := &entities.Respond{
		ID:       respondID,
		TicketID: ticketID,
		Price:    45.99,
		Comment:  pointers.New("Test comment"),
	}

	validTicket := &entities.Ticket{
		ID:          ticketID,
		UserID:      userID,
		Name:        "Ticket1",
		Description: "Description1",
		Price:       pointers.New[float32](39.99),
		Quantity:    5,
	}

	testCases := []struct {
		name          string
		obj           *entities.Respond
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *entities.Ticket
		expectedError error
		errorExpected bool
	}{
		{
			name: "successful ticket retrieval",
			obj:  validRespond,
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetTicketByID(gomock.Any(), ticketID).
					Return(validTicket, nil).
					Times(1)
			},
			expected: validTicket,
		},
		{
			name: "use case error with logging",
			obj:  validRespond,
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetTicketByID(gomock.Any(), ticketID).
					Return(nil, errors.New("get ticket error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expected:      nil,
			errorExpected: true,
		},
		{
			name:     "nil obj",
			obj:      nil,
			expected: nil,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &respondResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()

			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := resolver.Ticket(testCtx, tc.obj)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestRespondResolver_Master(t *testing.T) {
	mockRespond := &entities.Respond{
		ID:        respondID,
		MasterID:  masterID,
		CreatedAt: now,
	}

	mockMaster := &entities.Master{
		ID:        masterID,
		UserID:    userID,
		Info:      pointers.New("Professional master"),
		CreatedAt: now,
	}

	testCases := []struct {
		name          string
		obj           *entities.Respond
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *entities.Master
		errorExpected bool
	}{
		{
			name: "successful get master for respond",
			obj:  mockRespond,
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetMasterByID(gomock.Any(), mockRespond.MasterID).
					Return(mockMaster, nil).
					Times(1)
			},
			expected: mockMaster,
		},
		{
			name: "master not found",
			obj:  mockRespond,
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetMasterByID(gomock.Any(), mockRespond.MasterID).
					Return(nil, errors.New("test")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "database error",
			obj:  mockRespond,
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetMasterByID(gomock.Any(), mockRespond.MasterID).
					Return(nil, errors.New("database error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name:     "nil obj",
			obj:      nil,
			expected: nil,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &respondResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := resolver.Master(ctx, tc.obj)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestRespondResolver_Price(t *testing.T) {
	testCases := []struct {
		name          string
		obj           *entities.Respond
		expected      float64
		expectedError error
	}{
		{
			name: "successful get price",
			obj: &entities.Respond{
				ID:    respondID,
				Price: 19.99,
			},
			expected: 19.99,
		},
		{
			name: "zero price",
			obj: &entities.Respond{
				ID:    respondID,
				Price: 0,
			},
			expected: 0.00,
		},
		{
			name:     "nil obj",
			obj:      nil,
			expected: 0,
		},
	}

	resolver := &respondResolver{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := resolver.Price(ctx, tc.obj)
			rounded := math.Round(actual*100) / 100

			require.NoError(t, err)
			require.Equal(t, tc.expected, rounded)
		})
	}
}

func TestTicketResolver_User(t *testing.T) {
	mockTicket := &entities.Ticket{
		ID:        ticketID,
		UserID:    userID,
		CreatedAt: now,
	}

	testCases := []struct {
		name          string
		obj           *entities.Ticket
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *entities.User
		errorExpected bool
	}{
		{
			name: "successful get user for ticket",
			obj:  mockTicket,
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetUserByID(gomock.Any(), mockTicket.UserID).
					Return(user, nil).
					Times(1)
			},
			expected: user,
		},
		{
			name: "use case error",
			obj:  mockTicket,
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetUserByID(gomock.Any(), mockTicket.UserID).
					Return(nil, errors.New("test")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name:     "nil obj",
			obj:      nil,
			expected: nil,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &ticketResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := resolver.User(ctx, tc.obj)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestTicketResolver_Category(t *testing.T) {
	mockTicket := &entities.Ticket{
		ID:         ticketID,
		CategoryID: categoryID,
		CreatedAt:  now,
	}

	mockCategory := &entities.Category{
		ID:   categoryID,
		Name: "Test Category",
	}

	testCases := []struct {
		name          string
		obj           *entities.Ticket
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *entities.Category
		errorExpected bool
	}{
		{
			name: "successful get category for ticket",
			obj:  mockTicket,
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetCategoryByID(gomock.Any(), mockTicket.CategoryID).
					Return(mockCategory, nil).
					Times(1)
			},
			expected: mockCategory,
		},
		{
			name: "use case error",
			obj:  mockTicket,
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetCategoryByID(gomock.Any(), mockTicket.CategoryID).
					Return(nil, errors.New("test")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name:     "nil obj",
			obj:      nil,
			expected: nil,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &ticketResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := resolver.Category(ctx, tc.obj)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestTicketResolver_Price(t *testing.T) {
	testCases := []struct {
		name          string
		obj           *entities.Ticket
		expected      *float64
		errorExpected bool
	}{
		{
			name: "ticket with price",
			obj: &entities.Ticket{
				ID:    ticketID,
				Price: pointers.New[float32](19.99),
			},
			expected: pointers.New(19.99),
		},

		{
			name: "ticket without price",
			obj: &entities.Ticket{
				ID:    ticketID,
				Price: nil,
			},
			expected: nil,
		},
		{
			name:     "nil obj",
			obj:      nil,
			expected: nil,
		},
	}

	resolver := &ticketResolver{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := resolver.Price(ctx, tc.obj)
			require.NoError(t, err)

			if tc.expected == nil {
				require.Nil(t, actual)
			} else {
				require.NotNil(t, actual)
				rounded := math.Round(*actual*float64(100)) / 100
				require.Equal(t, *tc.expected, rounded)
			}
		})
	}
}

func TestTicketResolver_Quantity(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name          string
		obj           *entities.Ticket
		expected      int
		errorExpected bool
	}{
		{
			name: "positive quantity",
			obj: &entities.Ticket{
				ID:       ticketID,
				Quantity: 5,
			},
			expected: 5,
		},
		{
			name: "zero quantity",
			obj: &entities.Ticket{
				ID:       ticketID,
				Quantity: 0,
			},
			expected: 0,
		},
		{
			name: "max quantity",
			obj: &entities.Ticket{
				ID:       ticketID,
				Quantity: math.MaxInt32,
			},
			expected: math.MaxInt32,
		},
		{
			name:     "nil obj",
			obj:      nil,
			expected: 0,
		},
	}

	resolver := &ticketResolver{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := resolver.Quantity(ctx, tc.obj)

			require.NoError(t, err)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestToyResolver_Master(t *testing.T) {
	mockToy := &entities.Toy{
		ID:       toyID,
		MasterID: masterID,
		Name:     "Test Toy",
	}

	mockMaster := &entities.Master{
		ID:        masterID,
		UserID:    userID,
		Info:      pointers.New("Professional master"),
		CreatedAt: now,
	}

	testCases := []struct {
		name          string
		obj           *entities.Toy
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *entities.Master
		errorExpected bool
	}{
		{
			name: "successful get master for toy",
			obj:  mockToy,
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetMasterByID(gomock.Any(), mockToy.MasterID).
					Return(mockMaster, nil).
					Times(1)
			},
			expected: mockMaster,
		},
		{
			name: "use case error",
			obj:  mockToy,
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetMasterByID(gomock.Any(), mockToy.MasterID).
					Return(nil, errors.New("error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name:     "nil obj",
			obj:      nil,
			expected: nil,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &toyResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := resolver.Master(ctx, tc.obj)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestToyResolver_Category(t *testing.T) {
	mockToy := &entities.Toy{
		ID:         toyID,
		CategoryID: categoryID,
		Name:       "Test Toy",
		CreatedAt:  now,
	}

	mockCategory := &entities.Category{
		ID:   categoryID,
		Name: "Test Category",
	}

	testCases := []struct {
		name          string
		obj           *entities.Toy
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *entities.Category
		errorExpected bool
	}{
		{
			name: "successful get category for toy",
			obj:  mockToy,
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetCategoryByID(gomock.Any(), mockToy.CategoryID).
					Return(mockCategory, nil).
					Times(1)
			},
			expected: mockCategory,
		},
		{
			name: "use case error",
			obj:  mockToy,
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetCategoryByID(gomock.Any(), mockToy.CategoryID).
					Return(nil, errors.New("test")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name:     "nil obj",
			obj:      nil,
			expected: nil,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &toyResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := resolver.Category(ctx, tc.obj)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestToyResolver_Price(t *testing.T) {
	testCases := []struct {
		name          string
		obj           *entities.Toy
		expected      float64
		errorExpected bool
	}{
		{
			name: "ticket with price",
			obj: &entities.Toy{
				ID:    ticketID,
				Price: 19.99,
			},
			expected: 19.99,
		},
		{
			name:     "nil obj",
			obj:      nil,
			expected: 0,
		},
	}

	resolver := &toyResolver{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := resolver.Price(ctx, tc.obj)
			rounded := math.Round(actual*100) / 100
			require.NoError(t, err)
			require.Equal(t, tc.expected, rounded)
		})
	}
}

func TestToyResolver_Quantity(t *testing.T) {
	testCases := []struct {
		name          string
		obj           *entities.Toy
		expected      int
		errorExpected bool
	}{
		{
			name: "positive quantity",
			obj: &entities.Toy{
				ID:       ticketID,
				Quantity: 5,
			},
			expected: 5,
		},
		{
			name: "zero quantity",
			obj: &entities.Toy{
				ID:       ticketID,
				Quantity: 0,
			},
			expected: 0,
		},
		{
			name:     "nil obj",
			obj:      nil,
			expected: 0,
		},
	}

	resolver := &toyResolver{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := resolver.Quantity(ctx, tc.obj)

			require.NoError(t, err)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestQueryResolver_MasterByUserID(t *testing.T) {
	master := &entities.Master{
		ID:        masterID,
		UserID:    user.ID,
		Info:      pointers.New("Tes master"),
		CreatedAt: now,
		UpdatedAt: now,
	}

	testCases := []struct {
		name           string
		userID         string
		prepareContext func(ctx context.Context) context.Context
		setupMocks     func(useCases *mockusecases.MockUseCases)
		expected       *entities.Master
		expectedError  error
		errorExpected  bool
	}{
		{
			name:   "successful get master by user id",
			userID: "1",
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMasterByUserID(gomock.Any(), uint64(1)).
					Return(master, nil).
					Times(1)
			},
			expected: master,
		},
		{
			name:   "use case error",
			userID: "1",
			setupMocks: func(useCases *mockusecases.MockUseCases) {
				useCases.
					EXPECT().
					GetMasterByUserID(gomock.Any(), uint64(1)).
					Return(nil, errors.New("test")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	resolver := &queryResolver{
		Resolver: NewResolver(
			useCases,
			logger,
			config.CookiesConfig{},
		),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := ctx
			if tc.prepareContext != nil {
				testCtx = tc.prepareContext(testCtx)
			}

			if tc.setupMocks != nil {
				tc.setupMocks(useCases)
			}

			actual, err := resolver.MasterByUser(testCtx, tc.userID)

			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, err, tc.expectedError)
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestPaginationResolver_Limit(t *testing.T) {
	tests := []struct {
		name     string
		obj      *entities.Pagination
		data     *int
		expected *uint64
	}{
		{
			name:     "nil obj and nil data",
			obj:      nil,
			data:     nil,
			expected: nil,
		},
		{
			name:     "nil obj with data",
			obj:      nil,
			data:     pointers.New(10),
			expected: nil,
		},
		{
			name:     "obj with nil data",
			obj:      &entities.Pagination{},
			data:     nil,
			expected: nil,
		},
		{
			name:     "set positive limit",
			obj:      &entities.Pagination{},
			data:     pointers.New(25),
			expected: pointers.New(uint64(25)),
		},
		{
			name:     "set zero limit",
			obj:      &entities.Pagination{},
			data:     pointers.New(0),
			expected: pointers.New(uint64(0)),
		},
		{
			name:     "overwrite existing limit",
			obj:      &entities.Pagination{Limit: pointers.New(uint64(50))},
			data:     pointers.New(30),
			expected: pointers.New(uint64(30)),
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	mainResolver := NewResolver(useCases, logger, config.CookiesConfig{})
	testedResolver := &paginationResolver{Resolver: mainResolver}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := testedResolver.Limit(context.Background(), tc.obj, tc.data)
			require.NoError(t, err)

			if tc.obj != nil {
				require.Equal(t, tc.expected, tc.obj.Limit)
			}
		})
	}
}

func TestPaginationResolver_Offset(t *testing.T) {
	tests := []struct {
		name     string
		obj      *entities.Pagination
		data     *int
		expected *uint64
	}{
		{
			name:     "nil obj and nil data",
			obj:      nil,
			data:     nil,
			expected: nil,
		},
		{
			name:     "nil obj with data",
			obj:      nil,
			data:     pointers.New(100),
			expected: nil,
		},
		{
			name:     "obj with nil data",
			obj:      &entities.Pagination{},
			data:     nil,
			expected: nil,
		},
		{
			name:     "set positive offset",
			obj:      &entities.Pagination{},
			data:     pointers.New(50),
			expected: pointers.New(uint64(50)),
		},
		{
			name:     "set zero offset",
			obj:      &entities.Pagination{},
			data:     pointers.New(0),
			expected: pointers.New(uint64(0)),
		},
		{
			name:     "overwrite existing offset",
			obj:      &entities.Pagination{Offset: pointers.New(uint64(200))},
			data:     pointers.New(150),
			expected: pointers.New(uint64(150)),
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	mainResolver := NewResolver(useCases, logger, config.CookiesConfig{})
	testedResolver := &paginationResolver{Resolver: mainResolver}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := testedResolver.Offset(context.Background(), tc.obj, tc.data)
			require.NoError(t, err)

			if tc.obj != nil {
				require.Equal(t, tc.expected, tc.obj.Offset)
			}
		})
	}
}
