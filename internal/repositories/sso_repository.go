package repositories

import (
	customerrors "github.com/DKhorkov/hmtm-bff/internal/errors"
	ssoentities "github.com/DKhorkov/hmtm-sso/entities"
)

type GrpcSsoRepository struct {
}

func (repo *GrpcSsoRepository) RegisterUser(userData ssoentities.RegisterUserDTO) (int, error) {
	return 0, nil
}

func (repo *GrpcSsoRepository) GetUserByID(id int) (*ssoentities.User, error) {
	return nil, &customerrors.UserNotFoundError{}
}

func (repo *GrpcSsoRepository) GetAllUsers() ([]*ssoentities.User, error) {
	return []*ssoentities.User{}, nil
}

func (repo *GrpcSsoRepository) LoginUser(userData ssoentities.LoginUserDTO) (string, error) {
	return "", nil
}
