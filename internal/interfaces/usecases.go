package interfaces

import (
	ssoentities "github.com/DKhorkov/hmtm-sso/entities"
)

type UseCases interface {
	GetUserByID(id int) (*ssoentities.User, error)
	GetAllUsers() ([]*ssoentities.User, error)
	RegisterUser(userData ssoentities.RegisterUserDTO) (int, error)
	LoginUser(userData ssoentities.LoginUserDTO) (string, error)
}
