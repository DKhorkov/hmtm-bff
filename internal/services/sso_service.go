package services

import (
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
	ssoentities "github.com/DKhorkov/hmtm-sso/pkg/entities"
)

type CommonSsoService struct {
	ssoRepository interfaces.SsoRepository
}

func (service *CommonSsoService) GetAllUsers() ([]*ssoentities.User, error) {
	return service.ssoRepository.GetAllUsers()
}

func (service *CommonSsoService) GetUserByID(id uint64) (*ssoentities.User, error) {
	return service.ssoRepository.GetUserByID(id)
}

func (service *CommonSsoService) RegisterUser(userData ssoentities.RegisterUserDTO) (uint64, error) {
	return service.ssoRepository.RegisterUser(userData)
}

func (service *CommonSsoService) LoginUser(userData ssoentities.LoginUserDTO) (*ssoentities.TokensDTO, error) {
	return service.ssoRepository.LoginUser(userData)
}

func (service *CommonSsoService) GetMe(accessToken string) (*ssoentities.User, error) {
	return service.ssoRepository.GetMe(accessToken)
}

func (service *CommonSsoService) RefreshTokens(
	refreshTokensData ssoentities.TokensDTO,
) (*ssoentities.TokensDTO, error) {
	return service.ssoRepository.RefreshTokens(refreshTokensData)
}

func NewCommonSsoService(ssoRepository interfaces.SsoRepository) *CommonSsoService {
	return &CommonSsoService{ssoRepository: ssoRepository}
}
