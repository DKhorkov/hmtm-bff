package services

import "github.com/DKhorkov/hmtm-bff/graph/model"

type UsersService interface {
	GetUsers() ([]*model.User, error)
	CreateUser(newUser model.NewUser) (*model.User, error)
}
