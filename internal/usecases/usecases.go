package usecases

import (
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
	ssoentities "github.com/DKhorkov/hmtm-sso/pkg/entities"
)

type CommonUseCases struct {
	ssoService interfaces.SsoService
}

func (useCases *CommonUseCases) RegisterUser(userData ssoentities.RegisterUserDTO) (int, error) {
	return useCases.ssoService.RegisterUser(userData)
}

func (useCases *CommonUseCases) LoginUser(userData ssoentities.LoginUserDTO) (*ssoentities.TokensDTO, error) {
	return useCases.ssoService.LoginUser(userData)
}

func (useCases *CommonUseCases) GetMe(accessToken string) (*ssoentities.User, error) {
	return useCases.ssoService.GetMe(accessToken)
}

func (useCases *CommonUseCases) RefreshTokens(refreshTokensData ssoentities.TokensDTO) (*ssoentities.TokensDTO, error) {
	return useCases.ssoService.RefreshTokens(refreshTokensData)
}

func (useCases *CommonUseCases) GetUserByID(id int) (*ssoentities.User, error) {
	return useCases.ssoService.GetUserByID(id)
}

func (useCases *CommonUseCases) GetAllUsers() ([]*ssoentities.User, error) {
	return useCases.ssoService.GetAllUsers()
}

func NewCommonUseCases(ssoService interfaces.SsoService) *CommonUseCases {
	return &CommonUseCases{ssoService: ssoService}
}
