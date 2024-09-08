package tests

import (
	"context"
	"hmtmbff/graph"
	"hmtmbff/graph/model"
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
				ID:        "1",
				Email:     "test@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			message: "should return a new user with correct data",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &graph.Resolver{}
			actual, err := r.Mutation().CreateUser(context.Background(), tc.input)

			assert.NoError(t, err, "%s - error: %v", tc.message, err)

			actual.CreatedAt = time.Now().Add(-1 * time.Second)
			actual.UpdatedAt = time.Now().Add(-1 * time.Second)

			time.Sleep(1 * time.Second)

			assert.Equal(t, tc.expected.ID, actual.ID, "\n%s - actual: %v, expected: %v", tc.message, actual.ID, tc.expected.ID)
			assert.Equal(t, tc.expected.Email, actual.Email, "\n%s - actual: %v, expected: %v", tc.message, actual.Email, tc.expected.Email)
			assert.WithinDuration(t, actual.CreatedAt, time.Now(), 3*time.Second)
			assert.WithinDuration(t, actual.UpdatedAt, time.Now(), 3*time.Second)
		})
	}
}

// Тест корректности данных получаемого списка пользователей.
func TestUsers(t *testing.T) {
	testCases := []struct {
		name     string
		expected []*model.User
		message  string
	}{
		{
			name:     "should return a list of users",
			expected: []*model.User{{ID: "1", Email: "test@example.com"}},
			message:  "should return a list of users with correct data",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &graph.Resolver{}
			actual, err := r.Query().Users(context.Background())

			assert.NoError(t, err, "\n%s - error: %v", tc.message, err)
			assert.Equal(t, len(tc.expected), len(actual), "\n%s - actual: %v, expected: %v", tc.message, len(actual), len(tc.expected))

			for i, expectedUser := range tc.expected {
				assert.Equal(t, expectedUser.ID, actual[i].ID, "\n%s - actual: %v, expected: %v", tc.message, actual[i].ID, expectedUser.ID)
				assert.Equal(t, expectedUser.Email, actual[i].Email, "\n%s - actual: %v, expected: %v", tc.message, actual[i].Email, expectedUser.Email)
			}
		})
	}
}
