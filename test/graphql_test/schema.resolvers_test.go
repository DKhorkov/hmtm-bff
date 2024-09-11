package tests_test

import (
	"testing"
	"time"

	"github.com/DKhorkov/hmtm-bff/graph"
	"github.com/DKhorkov/hmtm-bff/graph/model"
	"github.com/DKhorkov/hmtm-bff/internal/mocks"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name     string
		input    model.NewUser
		expected *model.User
		message  string
		wantErr  bool
	}{
		{
			name: "should create a new user",
			input: model.NewUser{
				Email:    "test@example.com",
				Password: "password",
			},
			expected: &model.User{
				ID:        1,
				Email:     "test@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			message: "should return a new user with correct data",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &graph.Resolver{UsersService: &mocks.MockUsersService{}}
			actual, err := r.UsersService.CreateUser(tc.input)

			require.NoError(
				t,
				err,
				"%s - error: %v", tc.message, err)
			assert.Equal(
				t,
				tc.expected.ID,
				actual,
				"\n%s - actual: %v, expected: %v", tc.message, actual, tc.expected.ID)
		})
	}
}

func TestGetUsersWithoutExistingUsers(t *testing.T) {
	r := &graph.Resolver{UsersService: &mocks.MockUsersService{}}
	users, err := r.UsersService.GetUsers()

	require.NoError(t, err, "Should return no error")
	assert.Empty(t, users, "Should return an empty list")
}

func TestGetUsersWithExistingUsers(t *testing.T) {
	r := &graph.Resolver{UsersService: &mocks.MockUsersService{}}
	testUsers := []struct {
		Email    string
		Password string
	}{
		{"test1@hamster.com", "password1"},
		{"test2@wopwop.com", "password2"},
		{"test3@gogogo.com", "password3"},
	}

	for index, user := range testUsers {
		createdUserID, err := r.UsersService.CreateUser(model.NewUser{Email: user.Email, Password: user.Password})
		require.NoError(t, err, "Should create user without error")

		assert.Equal(t, createdUserID, index+1, "Should return correct ID for the user")
	}

	users, err := r.UsersService.GetUsers()
	require.NoError(t, err, "Should return no error")
	assert.Len(t, users, len(testUsers), "Should return the correct number of users")

	for index, user := range users {
		assert.Equal(t, user.Email, testUsers[index].Email, "Should return correct email for the user")
		assert.Equal(t, user.ID, index+1, "Should return correct ID for the user")
	}
}
