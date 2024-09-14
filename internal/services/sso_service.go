package services

import (
	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

type CommonSsoService struct {
	SsoRepository interfaces.SsoRepository
}

func (service *CommonSsoService) GetAllUsers() ([]*entities.User, error) {
	return service.SsoRepository.GetAllUsers()
}

func (service *CommonSsoService) GetUserByID(id int) (*entities.User, error) {
	return service.SsoRepository.GetUserByID(id)
}

func (service *CommonSsoService) LoginUser(userData entities.LoginUserDTO) (string, error) {
	token, err := service.SsoRepository.LoginUser(userData)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (service *CommonSsoService) RegisterUser(userData entities.RegisterUserDTO) (int, error) {
	return service.SsoRepository.RegisterUser(userData)
}
