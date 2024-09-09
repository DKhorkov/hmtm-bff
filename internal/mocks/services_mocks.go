package mocks

import (
	"time"

	"github.com/DKhorkov/hmtm-bff/graph/model"
)

type MockUsersService struct {
	usersStorage []*model.User
}

func (service *MockUsersService) CreateUser(newUser model.NewUser) (*model.User, error) {
	var user model.User
	user.Email = newUser.Email
	user.ID = len(service.usersStorage) + 1
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	service.usersStorage = append(service.usersStorage, &user)
	return &user, nil
}

func (service *MockUsersService) GetUsers() ([]*model.User, error) {
	return service.usersStorage, nil
}
