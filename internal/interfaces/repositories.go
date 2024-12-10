package interfaces

import "github.com/DKhorkov/hmtm-bff/internal/models"

type SsoRepository interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id uint64) (*models.User, error)
	RegisterUser(userData models.RegisterUserDTO) (userID uint64, err error)
	LoginUser(userData models.LoginUserDTO) (*models.TokensDTO, error)
	GetMe(accessToken string) (*models.User, error)
	RefreshTokens(refreshToken string) (*models.TokensDTO, error)
}

type ToysRepository interface {
	AddToy(toyData models.AddToyDTO) (toyID uint64, err error)
	GetAllToys() ([]models.Toy, error)
	GetToyByID(id uint64) (*models.Toy, error)
	GetMasterToys(masterID uint64) ([]models.Toy, error)
	GetAllMasters() ([]models.Master, error)
	GetMasterByID(id uint64) (*models.Master, error)
	RegisterMaster(masterData models.RegisterMasterDTO) (masterID uint64, err error)
	GetAllCategories() ([]models.Category, error)
	GetCategoryByID(id uint32) (*models.Category, error)
	GetAllTags() ([]models.Tag, error)
	GetTagByID(id uint32) (*models.Tag, error)
}
