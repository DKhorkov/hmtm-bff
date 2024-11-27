package interfaces

import (
	ssoentities "github.com/DKhorkov/hmtm-sso/pkg/entities"
)

type SsoRepository interface {
	GetAllUsers() ([]*ssoentities.User, error)
	GetUserByID(id int) (*ssoentities.User, error)
	RegisterUser(userData ssoentities.RegisterUserDTO) (userID int, err error)
	LoginUser(userData ssoentities.LoginUserDTO) (*ssoentities.TokensDTO, error)
	GetMe(accessToken string) (*ssoentities.User, error)
	RefreshTokens(refreshTokensData ssoentities.TokensDTO) (*ssoentities.TokensDTO, error)
}
