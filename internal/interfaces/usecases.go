package interfaces

import (
	ssoentities "github.com/DKhorkov/hmtm-sso/pkg/entities"
)

type UseCases interface {
	GetUserByID(id int) (*ssoentities.User, error)
	GetAllUsers() ([]*ssoentities.User, error)
	RegisterUser(userData ssoentities.RegisterUserDTO) (int, error)
	LoginUser(userData ssoentities.LoginUserDTO) (*ssoentities.TokensDTO, error)
	GetMe(accessToken string) (*ssoentities.User, error)
	RefreshTokens(refreshTokensData ssoentities.TokensDTO) (*ssoentities.TokensDTO, error)
}
