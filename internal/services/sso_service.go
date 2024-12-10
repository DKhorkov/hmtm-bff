package services

import (
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
	"github.com/DKhorkov/hmtm-bff/internal/models"
)

type CommonSsoService struct {
	ssoRepository interfaces.SsoRepository
}

func (service *CommonSsoService) GetAllUsers() ([]models.User, error) {
	return service.ssoRepository.GetAllUsers()
}

func (service *CommonSsoService) GetUserByID(id uint64) (*models.User, error) {
	return service.ssoRepository.GetUserByID(id)
}

func (service *CommonSsoService) RegisterUser(userData models.RegisterUserDTO) (uint64, error) {
	return service.ssoRepository.RegisterUser(userData)
}

func (service *CommonSsoService) LoginUser(userData models.LoginUserDTO) (*models.TokensDTO, error) {
	return service.ssoRepository.LoginUser(userData)
}

func (service *CommonSsoService) GetMe(accessToken string) (*models.User, error) {
	return service.ssoRepository.GetMe(accessToken)
}

func (service *CommonSsoService) RefreshTokens(refreshToken string) (*models.TokensDTO, error) {
	return service.ssoRepository.RefreshTokens(refreshToken)
}

func NewCommonSsoService(ssoRepository interfaces.SsoRepository) *CommonSsoService {
	return &CommonSsoService{ssoRepository: ssoRepository}
}
