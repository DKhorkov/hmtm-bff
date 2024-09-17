package usecases_test

import (
	"testing"
	"time"

	"github.com/DKhorkov/hmtm-bff/internal/errors"
	mocks "github.com/DKhorkov/hmtm-bff/internal/mocks/repositories"
	"github.com/DKhorkov/hmtm-bff/internal/usecases"
	ssoentities "github.com/DKhorkov/hmtm-sso/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: make(map[int]*ssoentities.User),
		}

		useCases := &usecases.CommonUseCases{SsoService: ssoRepository}
		userData := ssoentities.RegisterUserDTO{
			Credentials: ssoentities.LoginUserDTO{
				Email:    "test@example.com",
				Password: "password",
			},
		}

		userID, err := useCases.RegisterUser(userData)
		require.NoError(t, err)

		assert.Equal(t, 1, userID)
		assert.NotNil(t, ssoRepository.UsersStorage[userID])
		assert.Equal(t, "test@example.com", ssoRepository.UsersStorage[userID].Email)
	})
}

func TestLoginUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: make(map[int]*ssoentities.User),
		}

		useCases := &usecases.CommonUseCases{SsoService: ssoRepository}
		userData := ssoentities.LoginUserDTO{
			Email:    "test@example.com",
			Password: "password",
		}

		token, err := useCases.LoginUser(userData)
		require.NoError(t, err)

		assert.Equal(t, "test@example.com_password", token)
	})
}

func TestGetUserByID(t *testing.T) {
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
			},
		}

		useCases := &usecases.CommonUseCases{SsoService: ssoRepository}

		userResult, err := useCases.GetUserByID(1)
		require.NoError(t, err)

		assert.Equal(t, "test@example.com", userResult.Email)
	})
}

func TestGetUserByID_NotFound(t *testing.T) {
	repo := &mocks.MockedSsoRepository{
		UsersStorage: map[int]*ssoentities.User{
			1: {ID: 1, Email: "test@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			2: {ID: 2, Email: "test@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
	}

	useCases := &usecases.CommonUseCases{SsoService: repo}

	id := 3

	userResult, err := useCases.GetUserByID(id)

	assert.IsType(t, &errors.UserNotFoundError{}, err)
	assert.Equal(t, "user not found", err.Error())

	assert.Nil(t, userResult)
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

		useCases := &usecases.CommonUseCases{SsoService: ssoRepository}

		usersResult, err := useCases.GetAllUsers()
		require.NoError(t, err)

		assert.Len(t, usersResult, 2, "expected to get 2 users")
	})
}

func TestGetAllUsersWithoutExistingUsers(t *testing.T) {
	ssoRepository := &mocks.MockedSsoRepository{
		UsersStorage: make(map[int]*ssoentities.User),
	}

	useCases := &usecases.CommonUseCases{SsoService: ssoRepository}

	usersResult, err := useCases.GetAllUsers()
	require.NoError(t, err)

	assert.Empty(t, usersResult)
}
