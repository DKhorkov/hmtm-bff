package interfaces

import (
	"context"

	"github.com/DKhorkov/hmtm-bff/internal/models"
)

type SsoRepository interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id uint64) (*models.User, error)
	RegisterUser(ctx context.Context, userData models.RegisterUserDTO) (userID uint64, err error)
	LoginUser(ctx context.Context, userData models.LoginUserDTO) (*models.TokensDTO, error)
	GetMe(ctx context.Context, accessToken string) (*models.User, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*models.TokensDTO, error)
}

type ToysRepository interface {
	AddToy(ctx context.Context, toyData models.AddToyDTO) (toyID uint64, err error)
	GetAllToys(ctx context.Context) ([]models.Toy, error)
	GetToyByID(ctx context.Context, id uint64) (*models.Toy, error)
	GetMasterToys(ctx context.Context, masterID uint64) ([]models.Toy, error)
	GetAllMasters(ctx context.Context) ([]models.Master, error)
	GetMasterByID(ctx context.Context, id uint64) (*models.Master, error)
	RegisterMaster(ctx context.Context, masterData models.RegisterMasterDTO) (masterID uint64, err error)
	GetAllCategories(ctx context.Context) ([]models.Category, error)
	GetCategoryByID(ctx context.Context, id uint32) (*models.Category, error)
	GetAllTags(ctx context.Context) ([]models.Tag, error)
	GetTagByID(ctx context.Context, id uint32) (*models.Tag, error)
}

type FileStorageRepository interface {
	Upload(ctx context.Context, key string, file []byte) (string, error)
}
