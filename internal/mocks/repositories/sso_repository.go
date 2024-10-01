package mocks

import (
	"time"

	ssoerrors "github.com/DKhorkov/hmtm-sso/pkg/errors"

	ssoentities "github.com/DKhorkov/hmtm-sso/pkg/entities"
)

type MockedSsoRepository struct {
	UsersStorage map[int]*ssoentities.User
}

func (repo *MockedSsoRepository) RegisterUser(userData ssoentities.RegisterUserDTO) (int, error) {
	for _, user := range repo.UsersStorage {
		if user.Email == userData.Credentials.Email {
			return 0, &ssoerrors.UserAlreadyExistsError{}
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
				return "", &ssoerrors.InvalidPasswordError{}
			}

			return "someToken", nil
		}
	}

	return "", &ssoerrors.UserNotFoundError{}
}

func (repo *MockedSsoRepository) GetUserByID(id int) (*ssoentities.User, error) {
	user := repo.UsersStorage[id]
	if user != nil {
		return user, nil
	}

	return nil, &ssoerrors.UserNotFoundError{}
}

func (repo *MockedSsoRepository) GetAllUsers() ([]*ssoentities.User, error) {
	var users []*ssoentities.User
	for _, user := range repo.UsersStorage {
		users = append(users, user)
	}

	return users, nil
}
