package repositories

import (
	"context"

	"github.com/DKhorkov/libs/contextlib"
	"github.com/DKhorkov/libs/requestid"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
)

func NewGrpcSsoRepository(grpcClient interfaces.SsoGrpcClient) *GrpcSsoRepository {
	return &GrpcSsoRepository{client: grpcClient}
}

type GrpcSsoRepository struct {
	client interfaces.SsoGrpcClient
}

func (repo *GrpcSsoRepository) RegisterUser(ctx context.Context, userData entities.RegisterUserDTO) (uint64, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.Register(
		ctx,
		&sso.RegisterIn{
			RequestID: requestID,
			Email:     userData.Email,
			Password:  userData.Password,
		},
	)

	if err != nil {
		return 0, err
	}

	return response.GetUserID(), nil
}

func (repo *GrpcSsoRepository) GetUserByID(ctx context.Context, id uint64) (*entities.User, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.GetUser(
		ctx,
		&sso.GetUserIn{
			RequestID: requestID,
			ID:        id,
		},
	)

	if err != nil {
		return nil, err
	}

	return &entities.User{
		ID:        response.GetID(),
		Email:     response.GetEmail(),
		CreatedAt: response.GetCreatedAt().AsTime(),
		UpdatedAt: response.GetUpdatedAt().AsTime(),
	}, nil
}

func (repo *GrpcSsoRepository) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.GetUsers(
		ctx,
		&sso.GetUsersIn{
			RequestID: requestID,
		},
	)

	if err != nil {
		return nil, err
	}

	users := make([]entities.User, len(response.GetUsers()))
	for index, user := range response.GetUsers() {
		users[index] = entities.User{
			ID:        user.GetID(),
			Email:     user.GetEmail(),
			CreatedAt: user.GetCreatedAt().AsTime(),
			UpdatedAt: user.GetUpdatedAt().AsTime(),
		}
	}

	return users, nil
}

func (repo *GrpcSsoRepository) LoginUser(
	ctx context.Context,
	userData entities.LoginUserDTO,
) (*entities.TokensDTO, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.Login(
		ctx,
		&sso.LoginIn{
			RequestID: requestID,
			Email:     userData.Email,
			Password:  userData.Password,
		},
	)

	if err != nil {
		return nil, err
	}

	return &entities.TokensDTO{
		AccessToken:  response.GetAccessToken(),
		RefreshToken: response.GetRefreshToken(),
	}, nil
}

func (repo *GrpcSsoRepository) GetMe(ctx context.Context, accessToken string) (*entities.User, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.GetMe(
		ctx,
		&sso.GetMeIn{
			RequestID:   requestID,
			AccessToken: accessToken,
		},
	)

	if err != nil {
		return nil, err
	}

	return &entities.User{
		ID:        response.GetID(),
		Email:     response.GetEmail(),
		CreatedAt: response.GetCreatedAt().AsTime(),
		UpdatedAt: response.GetUpdatedAt().AsTime(),
	}, nil
}

func (repo *GrpcSsoRepository) RefreshTokens(ctx context.Context, refreshToken string) (*entities.TokensDTO, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.RefreshTokens(
		ctx,
		&sso.RefreshTokensIn{
			RequestID:    requestID,
			RefreshToken: refreshToken,
		},
	)

	if err != nil {
		return nil, err
	}

	return &entities.TokensDTO{
		AccessToken:  response.GetAccessToken(),
		RefreshToken: response.GetRefreshToken(),
	}, nil
}
