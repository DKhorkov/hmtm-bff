package usecases

import (
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
	ssoentities "github.com/DKhorkov/hmtm-sso/pkg/entities"
)

type CommonUseCases struct {
	SsoService interfaces.SsoService
}

func (useCases *CommonUseCases) RegisterUser(userData ssoentities.RegisterUserDTO) (int, error) {
	return useCases.SsoService.RegisterUser(userData)
}

func (useCases *CommonUseCases) LoginUser(userData ssoentities.LoginUserDTO) (*ssoentities.TokensDTO, error) {
	return useCases.SsoService.LoginUser(userData)
}

func (useCases *CommonUseCases) GetMe(accessToken string) (*ssoentities.User, error) {
	return useCases.SsoService.GetMe(accessToken)
}

func (useCases *CommonUseCases) RefreshTokens(refreshTokensData ssoentities.TokensDTO) (*ssoentities.TokensDTO, error) {
	return useCases.SsoService.RefreshTokens(refreshTokensData)
}

func (useCases *CommonUseCases) GetUserByID(id int) (*ssoentities.User, error) {
	return useCases.SsoService.GetUserByID(id)
}

func (useCases *CommonUseCases) GetAllUsers() ([]*ssoentities.User, error) {
	return useCases.SsoService.GetAllUsers()
}
