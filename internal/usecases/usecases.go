package usecases

import (
	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

type CommonUseCases struct {
	SsoService interfaces.SsoService
}

func (useCases *CommonUseCases) RegisterUser(userData entities.RegisterUserDTO) (int, error) {
	return useCases.SsoService.RegisterUser(userData)
}

func (useCases *CommonUseCases) LoginUser(userData entities.LoginUserDTO) (string, error) {
	return useCases.SsoService.LoginUser(userData)
}

func (useCases *CommonUseCases) GetUserByID(id int) (*entities.User, error) {
	return useCases.SsoService.GetUserByID(id)
}

func (useCases *CommonUseCases) GetAllUsers() ([]*entities.User, error) {
	return useCases.SsoService.GetAllUsers()
}
