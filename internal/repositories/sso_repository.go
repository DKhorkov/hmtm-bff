package repositories

import (
	"github.com/DKhorkov/hmtm-bff/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-bff/internal/errors"
)

type GrpcSsoRepository struct {
}

func (repo *GrpcSsoRepository) RegisterUser(userData entities.RegisterUserDTO) (int, error) {
	return 0, nil
}

func (repo *GrpcSsoRepository) GetUserByID(id int) (*entities.User, error) {
	return nil, &customerrors.UserNotFoundError{}
}

func (repo *GrpcSsoRepository) GetAllUsers() ([]*entities.User, error) {
	return []*entities.User{}, nil
}

func (repo *GrpcSsoRepository) LoginUser(userData entities.LoginUserDTO) (string, error) {
	return "", nil
}
