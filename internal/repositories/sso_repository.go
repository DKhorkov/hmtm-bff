package repositories

import (
	"context"

	"github.com/DKhorkov/hmtm-bff/internal/interfaces"

	ssoentities "github.com/DKhorkov/hmtm-sso/pkg/entities"
	"github.com/DKhorkov/hmtm-sso/protobuf/generated/go/sso"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcSsoRepository struct {
	client interfaces.SsoGrpcClient
}

func (repo *GrpcSsoRepository) RegisterUser(userData ssoentities.RegisterUserDTO) (int, error) {
	response, err := repo.client.Register(
		context.Background(),
		&sso.RegisterRequest{
			Credentials: &sso.LoginRequest{
				Email:    userData.Credentials.Email,
				Password: userData.Credentials.Password,
			},
		},
	)

	if err != nil {
		return 0, err
	}

	return int(response.GetUserID()), nil
}

func (repo *GrpcSsoRepository) GetUserByID(id int) (*ssoentities.User, error) {
	response, err := repo.client.GetUser(
		context.Background(),
		&sso.GetUserRequest{
			UserID: int64(id),
		},
	)

	if err != nil {
		return nil, err
	}

	return &ssoentities.User{
		ID:        int(response.GetUserID()),
		Email:     response.GetEmail(),
		CreatedAt: response.GetCreatedAt().AsTime(),
		UpdatedAt: response.GetUpdatedAt().AsTime(),
	}, nil
}

func (repo *GrpcSsoRepository) GetAllUsers() ([]*ssoentities.User, error) {
	response, err := repo.client.GetUsers(
		context.Background(),
		&emptypb.Empty{},
	)

	if err != nil {
		return nil, err
	}

	users := make([]*ssoentities.User, len(response.GetUsers()))
	for index, user := range response.GetUsers() {
		users[index] = &ssoentities.User{
			ID:        int(user.GetUserID()),
			Email:     user.GetEmail(),
			CreatedAt: user.GetCreatedAt().AsTime(),
			UpdatedAt: user.GetUpdatedAt().AsTime(),
		}
	}

	return users, nil
}

func (repo *GrpcSsoRepository) LoginUser(userData ssoentities.LoginUserDTO) (*ssoentities.TokensDTO, error) {
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

	return &ssoentities.TokensDTO{
		AccessToken:  response.GetAccessToken(),
		RefreshToken: response.GetRefreshToken(),
	}, nil
}

func (repo *GrpcSsoRepository) GetMe(accessToken string) (*ssoentities.User, error) {
	response, err := repo.client.GetMe(
		context.Background(),
		&sso.GetMeRequest{AccessToken: accessToken},
	)

	if err != nil {
		return nil, err
	}

	return &ssoentities.User{
		ID:        int(response.GetUserID()),
		Email:     response.GetEmail(),
		CreatedAt: response.GetCreatedAt().AsTime(),
		UpdatedAt: response.GetUpdatedAt().AsTime(),
	}, nil
}

func (repo *GrpcSsoRepository) RefreshTokens(refreshTokensData ssoentities.TokensDTO) (*ssoentities.TokensDTO, error) {
	response, err := repo.client.RefreshTokens(
		context.Background(),
		&sso.RefreshTokensRequest{
			AccessToken:  refreshTokensData.AccessToken,
			RefreshToken: refreshTokensData.RefreshToken,
		},
	)

	if err != nil {
		return nil, err
	}

	return &ssoentities.TokensDTO{
		AccessToken:  response.GetAccessToken(),
		RefreshToken: response.GetRefreshToken(),
	}, nil
}

func NewGrpcSsoRepository(grpcClient interfaces.SsoGrpcClient) *GrpcSsoRepository {
	return &GrpcSsoRepository{client: grpcClient}
}
