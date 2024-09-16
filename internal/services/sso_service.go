package services

import (
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
	ssoentities "github.com/DKhorkov/hmtm-sso/entities"
)

type CommonSsoService struct {
	SsoRepository interfaces.SsoRepository
}

func (service *CommonSsoService) GetAllUsers() ([]*ssoentities.User, error) {
	return service.SsoRepository.GetAllUsers()
}

func (service *CommonSsoService) GetUserByID(id int) (*ssoentities.User, error) {
	return service.SsoRepository.GetUserByID(id)
}

func (service *CommonSsoService) LoginUser(userData ssoentities.LoginUserDTO) (string, error) {
	token, err := service.SsoRepository.LoginUser(userData)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (service *CommonSsoService) RegisterUser(userData ssoentities.RegisterUserDTO) (int, error) {
	return service.SsoRepository.RegisterUser(userData)
}
