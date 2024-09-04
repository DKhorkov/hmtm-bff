package services

import "hmtmbff/graph/model"

type UsersService interface {
	GetUsers() ([]*model.User, error)
	CreateUser(newUser model.NewUser) (*model.User, error)
}
