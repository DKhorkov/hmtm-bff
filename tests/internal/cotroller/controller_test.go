package cotroller_test

import (
	"testing"
	"time"

	graphqlcore "github.com/DKhorkov/hmtm-bff/internal/controllers/graph/core"
	mocks "github.com/DKhorkov/hmtm-bff/internal/mocks/repositories"
	"github.com/DKhorkov/hmtm-bff/internal/services"
	"github.com/DKhorkov/hmtm-bff/internal/usecases"
	ssoentities "github.com/DKhorkov/hmtm-sso/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testUserID = 1

func TestRegisterUserResolverWithoutExistingUsers(t *testing.T) {
	ssoRepository := &mocks.MockedSsoRepository{
		UsersStorage: make(map[int]*ssoentities.User),
	}

	ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
	useCases := &usecases.CommonUseCases{SsoService: ssoService}
	resolvers := &graphqlcore.Resolver{UseCases: useCases}

	userData := ssoentities.RegisterUserDTO{
		Credentials: ssoentities.LoginUserDTO{
			Email:    "test@example.com",
			Password: "testPassword",
		},
	}

	result, err := resolvers.UseCases.RegisterUser(userData)
	require.NoError(
		t,
		err,
		"Error registering user")

	assert.Equal(t, testUserID, result, "should return user ID as 1")
}

func TestRegisterUserResolverWithExistingUsers(t *testing.T) {
	t.Run("should return the correct user ID when there are existing users", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: map[int]*ssoentities.User{
				1: {
					ID:        1,
					Email:     "existing@example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				2: {
					ID:        2,
					Email:     "another@example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				3: {
					ID:        3,
					Email:     "another@example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
		}

		ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
		useCases := &usecases.CommonUseCases{SsoService: ssoService}
		resolvers := &graphqlcore.Resolver{UseCases: useCases}

		userData := ssoentities.RegisterUserDTO{
			Credentials: ssoentities.LoginUserDTO{
				Email:    "new@example.com",
				Password: "password",
			},
		}
		userID, err := resolvers.UseCases.RegisterUser(userData)
		require.NoError(
			t,
			err,
			"unexpected error: %v", err)

		assert.Equal(t,
			4,
			userID,
			"expected user ID to be 4, got %d", userID)
	})
}

func TestLoginUserResolver(t *testing.T) {
	t.Run("should return a valid token when login is successful", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: map[int]*ssoentities.User{
				1: {
					ID:        1,
					Email:     "test@example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
		}

		ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
		useCases := &usecases.CommonUseCases{SsoService: ssoService}
		resolvers := &graphqlcore.Resolver{UseCases: useCases}

		userData := ssoentities.LoginUserDTO{
			Email:    "test@example.com",
			Password: "testPassword",
		}

		token, err := resolvers.UseCases.LoginUser(userData)
		require.NoError(
			t,
			err,
			"unexpected error during user login")

		assert.Equal(
			t,
			"test@example.com_testPassword",
			token,
			"expected token to be 'test@example.com_password', got '%s'", token)
	})

	t.Run("should return an error when login fails", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: map[int]*ssoentities.User{},
		}

		ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
		useCases := &usecases.CommonUseCases{SsoService: ssoService}
		resolvers := &graphqlcore.Resolver{UseCases: useCases}

		userData := ssoentities.LoginUserDTO{
			Email:    "test@example.com",
			Password: "password",
		}

		user := ssoRepository.UsersStorage[testUserID]
		result, err := resolvers.UseCases.LoginUser(userData)
		require.NoError(
			t,
			err,
			"should return an error")
		assert.Nil(
			t,
			user,
			"should return an nil, got '%s'", result)
	})
}

func TestGetUserResolver(t *testing.T) {
	t.Run("should return a valid user when user exists", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: map[int]*ssoentities.User{
				1: {
					ID:        1,
					Email:     "test@example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
		}

		ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
		useCases := &usecases.CommonUseCases{SsoService: ssoService}
		resolvers := &graphqlcore.Resolver{UseCases: useCases}

		user, err := resolvers.UseCases.GetUserByID(testUserID)
		require.NoError(
			t,
			err,
			"unexpected error during user retrieval")
		assert.Equal(
			t,
			"test@example.com",
			user.Email,
			"expected user email to be 'test@example.com', got '%s'", user.Email)
	})

	t.Run("should return an error when user doesn't exist", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: map[int]*ssoentities.User{},
		}

		ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
		useCases := &usecases.CommonUseCases{SsoService: ssoService}
		resolvers := &graphqlcore.Resolver{UseCases: useCases}

		_, err := resolvers.UseCases.GetUserByID(testUserID)
		assert.Error(
			t,
			err,
			"expected error, got nil")
	})
}

func TestGetAllUsersResolver(t *testing.T) {
	t.Run("should return all users", func(t *testing.T) {
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
					Email:     "another@example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
		}

		ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
		useCases := &usecases.CommonUseCases{SsoService: ssoService}
		resolvers := &graphqlcore.Resolver{UseCases: useCases}

		users, err := resolvers.UseCases.GetAllUsers()
		require.NoError(
			t,
			err,
			"unexpected error during user retrieval")

		assert.Len(
			t,
			users,
			2,
			"expected to get 2 users, got %d", len(users))
	})
}
