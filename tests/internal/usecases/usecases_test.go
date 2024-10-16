package usecases__test

import (
	"testing"
	"time"

	ssoerrors "github.com/DKhorkov/hmtm-sso/pkg/errors"

	"github.com/DKhorkov/hmtm-bff/internal/services"

	mocks "github.com/DKhorkov/hmtm-bff/internal/mocks/repositories"
	"github.com/DKhorkov/hmtm-bff/internal/usecases"
	ssoentities "github.com/DKhorkov/hmtm-sso/pkg/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser(t *testing.T) {
	const (
		testUserID    = 1
		testUserEmail = "test@example.com"
	)

	t.Run("Success", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: make(map[int]*ssoentities.User),
		}

		ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
		useCases := &usecases.CommonUseCases{SsoService: ssoService}

		userData := ssoentities.RegisterUserDTO{
			Credentials: ssoentities.LoginUserDTO{
				Email:    testUserEmail,
				Password: "password",
			},
		}

		userID, err := useCases.RegisterUser(userData)
		require.NoError(t, err)
		assert.Equal(t, testUserID, userID)
	})

	t.Run("UserAlreadyExists", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: map[int]*ssoentities.User{
				testUserID: {Email: testUserEmail},
			},
		}

		ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
		useCases := &usecases.CommonUseCases{SsoService: ssoService}

		userData := ssoentities.RegisterUserDTO{
			Credentials: ssoentities.LoginUserDTO{
				Email:    testUserEmail,
				Password: "password",
			},
		}

		userID, err := useCases.RegisterUser(userData)
		require.Error(t, err)
		assert.Equal(t, 0, userID)
		assert.IsType(t, &ssoerrors.UserAlreadyExistsError{}, err)
	})
}

func TestLoginUser(t *testing.T) {
	const (
		testUserID       = 1
		testUserEmail    = "test@example.com"
		testUserPassword = "password"
	)

	t.Run("Success", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: map[int]*ssoentities.User{
				testUserID: {
					Email:    testUserEmail,
					Password: testUserPassword,
				},
			},
		}

		ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
		useCases := &usecases.CommonUseCases{SsoService: ssoService}

		userData := ssoentities.LoginUserDTO{
			Email:    testUserEmail,
			Password: testUserPassword,
		}

		expected := &ssoentities.TokensDTO{
			AccessToken:  "AccessToken",
			RefreshToken: "RefreshToken",
		}

		tokens, err := useCases.LoginUser(userData)
		require.NoError(t, err)
		assert.Equal(t, expected, tokens)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: make(map[int]*ssoentities.User),
		}

		ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
		useCases := &usecases.CommonUseCases{SsoService: ssoService}

		userData := ssoentities.LoginUserDTO{
			Email:    testUserEmail,
			Password: testUserPassword,
		}

		tokens, err := useCases.LoginUser(userData)
		require.Error(t, err)
		assert.Nil(t, tokens)
		assert.IsType(t, &ssoerrors.UserNotFoundError{}, err)
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: map[int]*ssoentities.User{
				testUserID: {
					Email:    testUserEmail,
					Password: testUserPassword,
				},
			},
		}

		ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
		useCases := &usecases.CommonUseCases{SsoService: ssoService}

		userData := ssoentities.LoginUserDTO{
			Email:    testUserEmail,
			Password: "wrongPassword",
		}

		tokens, err := useCases.LoginUser(userData)
		require.Error(t, err)
		assert.Nil(t, tokens)
		assert.IsType(t, &ssoerrors.InvalidPasswordError{}, err)
	})
}

func TestGetUserByID(t *testing.T) {
	const (
		testUserID    = 1
		testUserEmail = "test@example.com"
	)

	t.Run("Success", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: map[int]*ssoentities.User{
				testUserID: {
					ID:        testUserID,
					Email:     testUserEmail,
					Password:  "password",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
		}

		ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
		useCases := &usecases.CommonUseCases{SsoService: ssoService}

		user, err := useCases.GetUserByID(testUserID)
		require.NoError(t, err)
		assert.Equal(t, testUserEmail, user.Email)
		assert.Equal(t, testUserID, user.ID)
	})
}

func TestGetUserByIDNotFound(t *testing.T) {
	ssoRepository := &mocks.MockedSsoRepository{
		UsersStorage: map[int]*ssoentities.User{
			1: {
				ID:        1,
				Email:     "test@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			2: {
				ID:        2,
				Email:     "test@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
	useCases := &usecases.CommonUseCases{SsoService: ssoService}

	userID := 3
	user, err := useCases.GetUserByID(userID)
	assert.IsType(t, &ssoerrors.UserNotFoundError{}, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Nil(t, user)
}

func TestGetAllUsersWithExistingUsers(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: map[int]*ssoentities.User{
				1: {
					ID:        1,
					Email:     "test@example.com",
					Password:  "password",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				2: {
					ID:        2,
					Email:     "test2@example.com",
					Password:  "password2",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
		}

		ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
		useCases := &usecases.CommonUseCases{SsoService: ssoService}

		users, err := useCases.GetAllUsers()
		require.NoError(t, err)
		assert.Len(
			t,
			users,
			len(ssoRepository.UsersStorage),
			"expected to get %d users, got %d", len(ssoRepository.UsersStorage), len(users))
	})
}

func TestGetAllUsersWithoutExistingUsers(t *testing.T) {
	ssoRepository := &mocks.MockedSsoRepository{
		UsersStorage: make(map[int]*ssoentities.User),
	}

	ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
	useCases := &usecases.CommonUseCases{SsoService: ssoService}

	users, err := useCases.GetAllUsers()
	require.NoError(t, err)
	assert.Empty(t, users)
}
