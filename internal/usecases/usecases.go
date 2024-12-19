package usecases

import (
	"context"
	"path"

	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
	"github.com/DKhorkov/hmtm-bff/internal/models"
	"github.com/DKhorkov/libs/security"
)

func NewCommonUseCases(
	ssoService interfaces.SsoService,
	toysService interfaces.ToysService,
	fileStorageService interfaces.FileStorageService,
) *CommonUseCases {
	return &CommonUseCases{
		ssoService:         ssoService,
		toysService:        toysService,
		fileStorageService: fileStorageService,
	}
}

type CommonUseCases struct {
	ssoService         interfaces.SsoService
	toysService        interfaces.ToysService
	fileStorageService interfaces.FileStorageService
}

func (useCases *CommonUseCases) RegisterUser(ctx context.Context, userData models.RegisterUserDTO) (uint64, error) {
	return useCases.ssoService.RegisterUser(ctx, userData)
}

func (useCases *CommonUseCases) LoginUser(
	ctx context.Context,
	userData models.LoginUserDTO,
) (*models.TokensDTO, error) {
	return useCases.ssoService.LoginUser(ctx, userData)
}

func (useCases *CommonUseCases) GetMe(ctx context.Context, accessToken string) (*models.User, error) {
	return useCases.ssoService.GetMe(ctx, accessToken)
}

func (useCases *CommonUseCases) RefreshTokens(ctx context.Context, refreshToken string) (*models.TokensDTO, error) {
	return useCases.ssoService.RefreshTokens(ctx, refreshToken)
}

func (useCases *CommonUseCases) GetUserByID(ctx context.Context, id uint64) (*models.User, error) {
	return useCases.ssoService.GetUserByID(ctx, id)
}

func (useCases *CommonUseCases) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return useCases.ssoService.GetAllUsers(ctx)
}

func (useCases *CommonUseCases) AddToy(ctx context.Context, toyData models.AddToyDTO) (uint64, error) {
	return useCases.toysService.AddToy(ctx, toyData)
}

func (useCases *CommonUseCases) GetAllToys(ctx context.Context) ([]models.Toy, error) {
	return useCases.toysService.GetAllToys(ctx)
}

func (useCases *CommonUseCases) GetMasterToys(ctx context.Context, masterID uint64) ([]models.Toy, error) {
	return useCases.toysService.GetMasterToys(ctx, masterID)
}

func (useCases *CommonUseCases) GetToyByID(ctx context.Context, id uint64) (*models.Toy, error) {
	return useCases.toysService.GetToyByID(ctx, id)
}

func (useCases *CommonUseCases) GetAllMasters(ctx context.Context) ([]models.Master, error) {
	return useCases.toysService.GetAllMasters(ctx)
}

func (useCases *CommonUseCases) GetMasterByID(ctx context.Context, id uint64) (*models.Master, error) {
	return useCases.toysService.GetMasterByID(ctx, id)
}

func (useCases *CommonUseCases) RegisterMaster(
	ctx context.Context,
	masterData models.RegisterMasterDTO,
) (uint64, error) {
	return useCases.toysService.RegisterMaster(ctx, masterData)
}

func (useCases *CommonUseCases) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	return useCases.toysService.GetAllCategories(ctx)
}

func (useCases *CommonUseCases) GetCategoryByID(ctx context.Context, id uint32) (*models.Category, error) {
	return useCases.toysService.GetCategoryByID(ctx, id)
}

func (useCases *CommonUseCases) GetAllTags(ctx context.Context) ([]models.Tag, error) {
	return useCases.toysService.GetAllTags(ctx)
}

func (useCases *CommonUseCases) GetTagByID(ctx context.Context, id uint32) (*models.Tag, error) {
	return useCases.toysService.GetTagByID(ctx, id)
}

func (useCases *CommonUseCases) UploadFile(ctx context.Context, filename string, file []byte) (string, error) {
	key := security.RawEncode([]byte(filename)) + path.Ext(filename)
	return useCases.fileStorageService.Upload(ctx, key, file)
}
