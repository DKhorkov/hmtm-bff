package interfaces

import (
	ssoentities "github.com/DKhorkov/hmtm-sso/pkg/entities"
	toysentities "github.com/DKhorkov/hmtm-toys/pkg/entities"
)

type SsoRepository interface {
	GetAllUsers() ([]*ssoentities.User, error)
	GetUserByID(id uint64) (*ssoentities.User, error)
	RegisterUser(userData ssoentities.RegisterUserDTO) (userID uint64, err error)
	LoginUser(userData ssoentities.LoginUserDTO) (*ssoentities.TokensDTO, error)
	GetMe(accessToken string) (*ssoentities.User, error)
	RefreshTokens(refreshTokensData ssoentities.TokensDTO) (*ssoentities.TokensDTO, error)
}

type ToysRepository interface {
	AddToy(toyData toysentities.RawAddToyDTO) (toyID uint64, err error)
	GetAllToys() ([]*toysentities.Toy, error)
	GetToyByID(id uint64) (*toysentities.Toy, error)
	GetAllMasters() ([]*toysentities.Master, error)
	GetMasterByID(id uint64) (*toysentities.Master, error)
	RegisterMaster(masterData toysentities.RawRegisterMasterDTO) (masterID uint64, err error)
	GetAllCategories() ([]*toysentities.Category, error)
	GetCategoryByID(id uint32) (*toysentities.Category, error)
	GetAllTags() ([]*toysentities.Tag, error)
	GetTagByID(id uint32) (*toysentities.Tag, error)
}
