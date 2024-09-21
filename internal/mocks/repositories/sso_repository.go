package mocks

import (
	"time"

	customerrors "github.com/DKhorkov/hmtm-bff/internal/errors"
	ssoentities "github.com/DKhorkov/hmtm-sso/entities"
)

type MockedSsoRepository struct {
	UsersStorage map[int]*ssoentities.User
}

func (repo *MockedSsoRepository) RegisterUser(userData ssoentities.RegisterUserDTO) (int, error) {
	for _, user := range repo.UsersStorage {
		if user.Email == userData.Credentials.Email {
			return 0, &customerrors.UserAlreadyExistsError{}
		}
	}

	var user ssoentities.User
	user.Email = userData.Credentials.Email
	user.ID = len(repo.UsersStorage) + 1
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	repo.UsersStorage[user.ID] = &user
	return user.ID, nil
}

func (repo *MockedSsoRepository) LoginUser(userData ssoentities.LoginUserDTO) (string, error) {
	for _, user := range repo.UsersStorage {
		if user.Email == userData.Email {
			if user.Password != userData.Password {
				return "", &customerrors.InvalidPasswordError{}
			}
			return "someToken", nil
		}
	}
	return "", &customerrors.UserNotFoundError{}
}

func (repo *MockedSsoRepository) GetUserByID(id int) (*ssoentities.User, error) {
	user := repo.UsersStorage[id]
	if user != nil {
		return user, nil
	}

	return nil, &customerrors.UserNotFoundError{}
}

func (repo *MockedSsoRepository) GetAllUsers() ([]*ssoentities.User, error) {
	var users []*ssoentities.User
	for _, user := range repo.UsersStorage {
		users = append(users, user)
	}

	return users, nil
}
