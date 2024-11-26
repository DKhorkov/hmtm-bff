package graphqlcontoller__test

import (
	"testing"
	"time"

	graphqlcore "github.com/DKhorkov/hmtm-bff/internal/controllers/graph/core"
	"github.com/DKhorkov/hmtm-bff/internal/services"
	"github.com/DKhorkov/hmtm-bff/internal/usecases"
	mocks "github.com/DKhorkov/hmtm-bff/mocks/repositories"
	ssoentities "github.com/DKhorkov/hmtm-sso/pkg/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestControllersRegisterUserResolverWithoutExistingUsers(t *testing.T) {
	const (
		testUserID    = 1
		testUserEmail = "test@example.com"
	)

	ssoRepository := &mocks.MockedSsoRepository{
		UsersStorage: make(map[int]*ssoentities.User),
	}

	ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
	useCases := &usecases.CommonUseCases{SsoService: ssoService}
	resolvers := &graphqlcore.Resolver{UseCases: useCases}

	userData := ssoentities.RegisterUserDTO{
		Credentials: ssoentities.LoginUserDTO{
			Email:    testUserEmail,
			Password: "testPassword",
		},
	}

	result, err := resolvers.UseCases.RegisterUser(userData)
	require.NoError(
		t,
		err,
		result,
		"Error registering user")
	assert.Equal(
		t,
		testUserID,
		result,
		"should return userID=%d", testUserID)
}

func TestControllersRegisterUserResolverWithExistingUsers(t *testing.T) {
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

		currentUsersLength := len(ssoRepository.UsersStorage)
		userID, err := resolvers.UseCases.RegisterUser(userData)
		require.NoError(
			t,
			err,
			"unexpected error: %v", err)
		assert.Len(
			t,
			ssoRepository.UsersStorage,
			currentUsersLength+1,
			"expected userID to be %d, got %d", currentUsersLength, userID)
	})
}

func TestControllersLoginUserResolver(t *testing.T) {
	const (
		testUserID    = 1
		testUserEmail = "test@example.com"
	)

	t.Run("should return a valid token when login is successful", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: map[int]*ssoentities.User{
				testUserID: {
					ID:        testUserID,
					Email:     testUserEmail,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Password:  "testPassword",
				},
			},
		}

		ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
		useCases := &usecases.CommonUseCases{SsoService: ssoService}
		resolvers := &graphqlcore.Resolver{UseCases: useCases}

		userData := ssoentities.LoginUserDTO{
			Email:    testUserEmail,
			Password: "testPassword",
		}

		expected := &ssoentities.TokensDTO{
			AccessToken:  "AccessToken",
			RefreshToken: "RefreshToken",
		}

		tokens, err := resolvers.UseCases.LoginUser(userData)
		require.NoError(
			t,
			err,
			"unexpected error during user login")
		assert.Equal(
			t,
			expected,
			tokens)
	})

	t.Run("should return an error when user not found", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: map[int]*ssoentities.User{},
		}

		ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
		useCases := &usecases.CommonUseCases{SsoService: ssoService}
		resolvers := &graphqlcore.Resolver{UseCases: useCases}

		userData := ssoentities.LoginUserDTO{
			Email:    testUserEmail,
			Password: "password",
		}

		tokens, err := resolvers.UseCases.LoginUser(userData)
		require.Error(
			t,
			err,
			"should return an error")
		assert.Nil(t, tokens)
	})

	t.Run("should return error when login fails", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: map[int]*ssoentities.User{
				testUserID: {
					ID:        testUserID,
					Email:     testUserEmail,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Password:  "testPassword",
				},
			},
		}

		ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
		useCases := &usecases.CommonUseCases{SsoService: ssoService}
		resolvers := &graphqlcore.Resolver{UseCases: useCases}

		userData := ssoentities.LoginUserDTO{
			Email:    testUserEmail,
			Password: "wrongPassword",
		}

		tokens, err := resolvers.UseCases.LoginUser(userData)
		require.Error(
			t,
			err,
			"should return an error")
		assert.Nil(t, tokens)
	})
}

func TestControllersGetUserResolver(t *testing.T) {
	const (
		testUserID    = 1
		testUserEmail = "test@example.com"
	)

	t.Run("should return a valid user when user exists", func(t *testing.T) {
		ssoRepository := &mocks.MockedSsoRepository{
			UsersStorage: map[int]*ssoentities.User{
				testUserID: {
					ID:        testUserID,
					Email:     testUserEmail,
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
			testUserID,
			user.ID,
			"expected user ID to be '%d', got '%d'", testUserID, user.ID)
		assert.Equal(
			t,
			testUserEmail,
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

		user, err := resolvers.UseCases.GetUserByID(testUserID)
		require.Error(
			t,
			err,
			"expected error, got nil")
		assert.Nil(
			t,
			user,
			"should return nul if user doesn't exist")
	})
}

func TestControllersGetAllUsersResolver(t *testing.T) {
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
			len(ssoRepository.UsersStorage),
			"expected to get %d users, got %d", len(ssoRepository.UsersStorage), len(users))
	})
}
