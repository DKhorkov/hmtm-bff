package mocks

import (
	"time"

	"github.com/DKhorkov/hmtm-bff/graph/model"
)

type MockUsersService struct {
	UsersStorage []*model.User
}

func (service *MockUsersService) CreateUser(newUser model.NewUser) (*model.User, error) {
	var user model.User
	user.Email = newUser.Email
	user.ID = 1
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	service.UsersStorage = append(service.UsersStorage, &user)
	return &user, nil
}

func (service *MockUsersService) GetUsers() ([]*model.User, error) {
	return service.UsersStorage, nil
}
