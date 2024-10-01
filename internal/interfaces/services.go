package interfaces

import (
	ssoentities "github.com/DKhorkov/hmtm-sso/pkg/entities"
)

type SsoService interface {
	GetAllUsers() ([]*ssoentities.User, error)
	GetUserByID(int) (*ssoentities.User, error)
	LoginUser(userData ssoentities.LoginUserDTO) (string, error)
	RegisterUser(userData ssoentities.RegisterUserDTO) (int, error)
}
