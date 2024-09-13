package services

import (
	"errors"

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
	user, err := service.SsoRepository.GetUserByEmail(userData.Email)
	if err != nil {
		return "", err
	}

	if user.Password != userData.Password {
		return "", errors.New("incorrect password")
	}

	return "someToken", nil
}

func (service *CommonSsoService) RegisterUser(userData entities.RegisterUserDTO) (int, error) {
	return service.SsoRepository.RegisterUser(userData)
}
