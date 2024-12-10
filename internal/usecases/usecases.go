package usecases

import (
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
	"github.com/DKhorkov/hmtm-bff/internal/models"
)

type CommonUseCases struct {
	ssoService  interfaces.SsoService
	toysService interfaces.ToysService
}

func (useCases *CommonUseCases) RegisterUser(userData models.RegisterUserDTO) (uint64, error) {
	return useCases.ssoService.RegisterUser(userData)
}

func (useCases *CommonUseCases) LoginUser(userData models.LoginUserDTO) (*models.TokensDTO, error) {
	return useCases.ssoService.LoginUser(userData)
}

func (useCases *CommonUseCases) GetMe(accessToken string) (*models.User, error) {
	return useCases.ssoService.GetMe(accessToken)
}

func (useCases *CommonUseCases) RefreshTokens(refreshToken string) (*models.TokensDTO, error) {
	return useCases.ssoService.RefreshTokens(refreshToken)
}

func (useCases *CommonUseCases) GetUserByID(id uint64) (*models.User, error) {
	return useCases.ssoService.GetUserByID(id)
}

func (useCases *CommonUseCases) GetAllUsers() ([]models.User, error) {
	return useCases.ssoService.GetAllUsers()
}

func (useCases *CommonUseCases) AddToy(toyData models.AddToyDTO) (uint64, error) {
	return useCases.toysService.AddToy(toyData)
}

func (useCases *CommonUseCases) GetAllToys() ([]models.Toy, error) {
	return useCases.toysService.GetAllToys()
}

func (useCases *CommonUseCases) GetMasterToys(masterID uint64) ([]models.Toy, error) {
	return useCases.toysService.GetMasterToys(masterID)
}

func (useCases *CommonUseCases) GetToyByID(id uint64) (*models.Toy, error) {
	return useCases.toysService.GetToyByID(id)
}

func (useCases *CommonUseCases) GetAllMasters() ([]models.Master, error) {
	return useCases.toysService.GetAllMasters()
}

func (useCases *CommonUseCases) GetMasterByID(id uint64) (*models.Master, error) {
	return useCases.toysService.GetMasterByID(id)
}

func (useCases *CommonUseCases) RegisterMaster(masterData models.RegisterMasterDTO) (uint64, error) {
	return useCases.toysService.RegisterMaster(masterData)
}

func (useCases *CommonUseCases) GetAllCategories() ([]models.Category, error) {
	return useCases.toysService.GetAllCategories()
}

func (useCases *CommonUseCases) GetCategoryByID(id uint32) (*models.Category, error) {
	return useCases.toysService.GetCategoryByID(id)
}

func (useCases *CommonUseCases) GetAllTags() ([]models.Tag, error) {
	return useCases.toysService.GetAllTags()
}

func (useCases *CommonUseCases) GetTagByID(id uint32) (*models.Tag, error) {
	return useCases.toysService.GetTagByID(id)
}

func NewCommonUseCases(
	ssoService interfaces.SsoService,
	toysService interfaces.ToysService,
) *CommonUseCases {
	return &CommonUseCases{
		ssoService:  ssoService,
		toysService: toysService,
	}
}
