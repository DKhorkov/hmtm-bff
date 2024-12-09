package repositories

import (
	"context"

	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
	"github.com/DKhorkov/hmtm-bff/internal/models"
	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcSsoRepository struct {
	client interfaces.SsoGrpcClient
}

func (repo *GrpcSsoRepository) RegisterUser(userData models.RegisterUserDTO) (uint64, error) {
	response, err := repo.client.Register(
		context.Background(),
		&sso.RegisterRequest{
			Credentials: &sso.LoginRequest{
				Email:    userData.Email,
				Password: userData.Password,
			},
		},
	)

	if err != nil {
		return 0, err
	}

	return response.GetUserID(), nil
}

func (repo *GrpcSsoRepository) GetUserByID(id uint64) (*models.User, error) {
	response, err := repo.client.GetUser(
		context.Background(),
		&sso.GetUserRequest{
			ID: id,
		},
	)

	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:        response.GetID(),
		Email:     response.GetEmail(),
		CreatedAt: response.GetCreatedAt().AsTime(),
		UpdatedAt: response.GetUpdatedAt().AsTime(),
	}, nil
}

func (repo *GrpcSsoRepository) GetAllUsers() ([]models.User, error) {
	response, err := repo.client.GetUsers(
		context.Background(),
		&emptypb.Empty{},
	)

	if err != nil {
		return nil, err
	}

	users := make([]models.User, len(response.GetUsers()))
	for index, user := range response.GetUsers() {
		users[index] = models.User{
			ID:        user.GetID(),
			Email:     user.GetEmail(),
			CreatedAt: user.GetCreatedAt().AsTime(),
			UpdatedAt: user.GetUpdatedAt().AsTime(),
		}
	}

	return users, nil
}

func (repo *GrpcSsoRepository) LoginUser(userData models.LoginUserDTO) (*models.TokensDTO, error) {
	response, err := repo.client.Login(
		context.Background(),
		&sso.LoginRequest{
			Email:    userData.Email,
			Password: userData.Password,
		},
	)

	if err != nil {
		return nil, err
	}

	return &models.TokensDTO{
		AccessToken:  response.GetAccessToken(),
		RefreshToken: response.GetRefreshToken(),
	}, nil
}

func (repo *GrpcSsoRepository) GetMe(accessToken string) (*models.User, error) {
	response, err := repo.client.GetMe(
		context.Background(),
		&sso.GetMeRequest{AccessToken: accessToken},
	)

	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:        response.GetID(),
		Email:     response.GetEmail(),
		CreatedAt: response.GetCreatedAt().AsTime(),
		UpdatedAt: response.GetUpdatedAt().AsTime(),
	}, nil
}

func (repo *GrpcSsoRepository) RefreshTokens(refreshToken string) (*models.TokensDTO, error) {
	response, err := repo.client.RefreshTokens(
		context.Background(),
		&sso.RefreshTokensRequest{
			RefreshToken: refreshToken,
		},
	)

	if err != nil {
		return nil, err
	}

	return &models.TokensDTO{
		AccessToken:  response.GetAccessToken(),
		RefreshToken: response.GetRefreshToken(),
	}, nil
}

func NewGrpcSsoRepository(grpcClient interfaces.SsoGrpcClient) *GrpcSsoRepository {
	return &GrpcSsoRepository{client: grpcClient}
}
