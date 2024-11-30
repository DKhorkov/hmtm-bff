package usecases

import (
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
	ssoentities "github.com/DKhorkov/hmtm-sso/pkg/entities"
	toysentities "github.com/DKhorkov/hmtm-toys/pkg/entities"
)

type CommonUseCases struct {
	ssoService  interfaces.SsoService
	toysService interfaces.ToysService
}

func (useCases *CommonUseCases) RegisterUser(userData ssoentities.RegisterUserDTO) (uint64, error) {
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

func (useCases *CommonUseCases) GetUserByID(id uint64) (*ssoentities.User, error) {
	return useCases.ssoService.GetUserByID(id)
}

func (useCases *CommonUseCases) GetAllUsers() ([]*ssoentities.User, error) {
	return useCases.ssoService.GetAllUsers()
}

func (useCases *CommonUseCases) AddToy(toyData toysentities.RawAddToyDTO) (uint64, error) {
	return useCases.toysService.AddToy(toyData)
}

func (useCases *CommonUseCases) GetAllToys() ([]*toysentities.Toy, error) {
	return useCases.toysService.GetAllToys()
}

func (useCases *CommonUseCases) GetToyByID(id uint64) (*toysentities.Toy, error) {
	return useCases.toysService.GetToyByID(id)
}

func (useCases *CommonUseCases) GetAllMasters() ([]*toysentities.Master, error) {
	return useCases.toysService.GetAllMasters()
}

func (useCases *CommonUseCases) GetMasterByID(id uint64) (*toysentities.Master, error) {
	return useCases.toysService.GetMasterByID(id)
}

func (useCases *CommonUseCases) RegisterMaster(masterData toysentities.RawRegisterMasterDTO) (uint64, error) {
	return useCases.toysService.RegisterMaster(masterData)
}

func (useCases *CommonUseCases) GetAllCategories() ([]*toysentities.Category, error) {
	return useCases.toysService.GetAllCategories()
}

func (useCases *CommonUseCases) GetCategoryByID(id uint32) (*toysentities.Category, error) {
	return useCases.toysService.GetCategoryByID(id)
}

func (useCases *CommonUseCases) GetAllTags() ([]*toysentities.Tag, error) {
	return useCases.toysService.GetAllTags()
}

func (useCases *CommonUseCases) GetTagByID(id uint32) (*toysentities.Tag, error) {
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
