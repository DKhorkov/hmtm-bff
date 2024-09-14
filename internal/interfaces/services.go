package interfaces

import (
	"github.com/DKhorkov/hmtm-bff/internal/entities"
)

type SsoService interface {
	GetAllUsers() ([]*entities.User, error)
	GetUserByID(int) (*entities.User, error)
	LoginUser(userData entities.LoginUserDTO) (string, error)
	RegisterUser(userData entities.RegisterUserDTO) (int, error)
}
