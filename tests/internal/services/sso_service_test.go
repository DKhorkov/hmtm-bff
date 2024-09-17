package services__test

import (
	"sort"
	"testing"
	"time"

	"github.com/DKhorkov/hmtm-bff/internal/errors"

	ssoentities "github.com/DKhorkov/hmtm-sso/entities"

	mocks "github.com/DKhorkov/hmtm-bff/internal/mocks/repositories"
	"github.com/DKhorkov/hmtm-bff/internal/services"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	testCases := []struct {
		name     string
		input    ssoentities.RegisterUserDTO
		expected int
		message  string
	}{
		{
			name: "should register a new user",
			input: ssoentities.RegisterUserDTO{
				Credentials: ssoentities.LoginUserDTO{
					Email:    "tests@example.com",
					Password: "password",
				},
			},
			expected: 1,
			message:  "should return a new user id",
		},
	}

	ssoRepository := &mocks.MockedSsoRepository{UsersStorage: map[int]*ssoentities.User{}}
	ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ssoService.RegisterUser(tc.input)

			require.NoError(
				t,
				err,
				"%s - error: %v", tc.message, err)
			assert.Equal(
				t,
				tc.expected,
				actual,
				"\n%s - actual: %v, expected: %v", tc.message, actual, tc.expected)
		})
	}
}

func TestGetAllUsersWithoutExistingUsers(t *testing.T) {
	ssoRepository := &mocks.MockedSsoRepository{UsersStorage: map[int]*ssoentities.User{}}
	ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
	users, err := ssoService.GetAllUsers()

	require.NoError(t, err, "Should return no error")
	assert.Empty(t, users, "Should return an empty list")
}

func TestGetAllUsersWithExistingUsers(t *testing.T) {
	testUsers := [3]ssoentities.RegisterUserDTO{
		{
			Credentials: ssoentities.LoginUserDTO{
				Email:    "test1@example.com",
				Password: "password1",
			},
		},
		{
			Credentials: ssoentities.LoginUserDTO{
				Email:    "test2@example.com",
				Password: "password2",
			},
		},
		{
			Credentials: ssoentities.LoginUserDTO{
				Email:    "test3@example.com",
				Password: "password3",
			},
		},
	}

	ssoRepository := &mocks.MockedSsoRepository{UsersStorage: map[int]*ssoentities.User{}}
	ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}
	for index, userData := range testUsers {
		registeredUserID, err := ssoService.RegisterUser(userData)
		require.NoError(t, err, "Should create user without error")
		assert.Equal(t, registeredUserID, index+1, "Should return correct ID for registered user")
	}

	users, err := ssoService.GetAllUsers()
	require.NoError(t, err, "Should return no error")
	assert.Len(t, users, len(testUsers), "Should return correct number of users")

	// Sorting slice of users to avoid IDs and Emails mismatch errors due to slice structure:
	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})

	for index, user := range users {
		assert.Equal(t, user.Email, testUsers[index].Credentials.Email, "Should return correct email for user")
		assert.Equal(t, user.ID, index+1, "Should return correct ID for user")
	}
}

func TestGetUserByID(t *testing.T) {
	testCases := []struct {
		name     string
		input    int
		expected *ssoentities.User
	}{
		{
			name:     "should find user by ID",
			input:    4,
			expected: &ssoentities.User{ID: 4, Email: "test@example4.com"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ssoRepository := &mocks.MockedSsoRepository{
				UsersStorage: map[int]*ssoentities.User{
					1: {ID: 1, Email: "test@example.com"},
					2: {ID: 2, Email: "test@example2.com"},
					3: {ID: 3, Email: "test@example3.com"},
					4: {ID: 4, Email: "test@example4.com"},
				},
			}
			ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}

			actual, err := ssoService.GetUserByID(tc.input)

			require.NoError(t, err, "Should return no error")
			assert.Equal(t, tc.expected, actual, "\n%s - actual: %v, expected: %v", tc.name, actual, tc.expected)
		})
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	ssoService := &services.CommonSsoService{
		SsoRepository: &mocks.MockedSsoRepository{
			UsersStorage: map[int]*ssoentities.User{
				1: {ID: 1, Email: "test@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
				2: {ID: 2, Email: "test@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			},
		},
	}

	id := 3

	_, err := ssoService.GetUserByID(id)

	assert.IsType(t, &errors.UserNotFoundError{}, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestGetUserByID_Nil(t *testing.T) {
	ssoService := &services.CommonSsoService{
		SsoRepository: &mocks.MockedSsoRepository{
			UsersStorage: map[int]*ssoentities.User{
				1: nil,
			},
		},
	}

	id := 1

	_, err := ssoService.GetUserByID(id)

	assert.IsType(t, &errors.UserNotFoundError{}, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestLoginUser(t *testing.T) {
	testCases := []struct {
		name     string
		input    ssoentities.LoginUserDTO
		expected string
	}{
		{
			name:     "should return token",
			input:    ssoentities.LoginUserDTO{Email: "test@example.com", Password: "password"},
			expected: "test@example.com_password",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ssoRepository := &mocks.MockedSsoRepository{}
			ssoService := &services.CommonSsoService{SsoRepository: ssoRepository}

			actual, err := ssoService.LoginUser(tc.input)

			require.NoError(t, err, "Should return no error")
			assert.Equal(t, tc.expected, actual, "\n%s - actual: %v, expected: %v", tc.name, actual, tc.expected)
		})
	}
}
