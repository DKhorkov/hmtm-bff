package tests

import (
	"context"
	"github.com/DKhorkov/hmtm-bff/graph"
	"github.com/DKhorkov/hmtm-bff/graph/model"
	"github.com/DKhorkov/hmtm-bff/internal/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Тест создания пользователя.
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
			actual, err := r.Mutation().CreateUser(context.Background(), tc.input)

			assert.NoError(
				t,
				err,
				"%s - error: %v", tc.message, err)
			assert.Equal(
				t,
				tc.expected.ID,
				actual.ID,
				"\n%s - actual: %v, expected: %v", tc.message, actual.ID, tc.expected.ID)
			assert.Equal(
				t,
				tc.expected.Email,
				actual.Email,
				"\n%s - actual: %v, expected: %v", tc.message, actual.Email, tc.expected.Email)
		})
	}
}

// Тесты на получение списка с пользователями и пустого списка.
func TestGetUsers(t *testing.T) {
	testCases := []struct {
		name     string
		storage  []*model.User
		expected int
		isEmpty  bool
		message  string
	}{
		{
			name: "GetUsers returns all users",
			storage: []*model.User{
				{Email: "test1@hamster.com", ID: 1},
				{Email: "test2@wopwop.com", ID: 2},
				{Email: "test3@gogogo.com", ID: 3},
			},
			expected: 3,
			isEmpty:  false,
			message:  "Should return three users",
		},
		{
			name:     "GetUsers returns empty list",
			storage:  []*model.User{},
			expected: 0,
			isEmpty:  true,
			message:  "Should return an empty list",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := &mocks.MockUsersService{
				UsersStorage: tc.storage,
			}

			r := &graph.Resolver{
				UsersService: s,
			}

			users, err := r.Query().Users(context.Background())

			assert.NoError(
				t,
				err,
				"%s - error: %v", tc.message, err)
			assert.Len(
				t,
				users,
				tc.expected,
				"%s - actual: %v, expected: %v", tc.message, len(users), tc.expected)
			if tc.isEmpty {
				assert.Empty(
					t,
					users,
					"%s - users shouldn't be empty", tc.message)
			} else {
				assert.Equal(
					t,
					users[0].Email,
					"test1@hamster.com",
					"%s - actual: %v, expected: %v", tc.message, users[0].Email, "test1@hamster.com")
				assert.Equal(
					t,
					users[1].Email,
					"test2@wopwop.com",
					"%s - actual: %v, expected: %v", tc.message, users[1].Email, "test2@wopwop.com")
				assert.Equal(
					t,
					users[2].Email,
					"test3@gogogo.com",
					"%s - actual: %v, expected: %v", tc.message, users[2].Email, "test3@gogogo.com")
			}
		})
	}
}
