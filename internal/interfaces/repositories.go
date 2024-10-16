package interfaces

import (
	ssoentities "github.com/DKhorkov/hmtm-sso/pkg/entities"
)

type SsoRepository interface {
	GetUserByID(id int) (*ssoentities.User, error)
	GetAllUsers() ([]*ssoentities.User, error)
	RegisterUser(user ssoentities.RegisterUserDTO) (int, error)
	LoginUser(userData ssoentities.LoginUserDTO) (*ssoentities.TokensDTO, error)
	GetMe(accessToken string) (*ssoentities.User, error)
	RefreshTokens(refreshTokensData ssoentities.TokensDTO) (*ssoentities.TokensDTO, error)
}
