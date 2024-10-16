package services

import (
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
	ssoentities "github.com/DKhorkov/hmtm-sso/pkg/entities"
)

type CommonSsoService struct {
	SsoRepository interfaces.SsoRepository
}

func (service *CommonSsoService) GetAllUsers() ([]*ssoentities.User, error) {
	return service.SsoRepository.GetAllUsers()
}

func (service *CommonSsoService) GetUserByID(id int) (*ssoentities.User, error) {
	return service.SsoRepository.GetUserByID(id)
}

func (service *CommonSsoService) RegisterUser(userData ssoentities.RegisterUserDTO) (int, error) {
	return service.SsoRepository.RegisterUser(userData)
}

func (service *CommonSsoService) LoginUser(userData ssoentities.LoginUserDTO) (*ssoentities.TokensDTO, error) {
	return service.SsoRepository.LoginUser(userData)
}

func (service *CommonSsoService) GetMe(accessToken string) (*ssoentities.User, error) {
	return service.SsoRepository.GetMe(accessToken)
}

func (service *CommonSsoService) RefreshTokens(
	refreshTokensData ssoentities.TokensDTO,
) (*ssoentities.TokensDTO, error) {
	return service.SsoRepository.RefreshTokens(refreshTokensData)
}
