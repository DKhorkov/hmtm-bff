package interfaces

import "github.com/DKhorkov/hmtm-bff/internal/entities"

type SsoRepository interface {
	GetUserByID(id int) (*entities.User, error)
	GetAllUsers() ([]*entities.User, error)
	RegisterUser(user entities.RegisterUserDTO) (int, error)
	GetUserByEmail(email string) (*entities.User, error)
	LoginUser(userData entities.LoginUserDTO) (string, error)
}
