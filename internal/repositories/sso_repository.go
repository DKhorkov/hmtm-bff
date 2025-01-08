package repositories

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
)

func NewGrpcSsoRepository(client interfaces.SsoGrpcClient) *GrpcSsoRepository {
	return &GrpcSsoRepository{client: client}
}

type GrpcSsoRepository struct {
	client interfaces.SsoGrpcClient
}

func (repo *GrpcSsoRepository) RegisterUser(ctx context.Context, userData entities.RegisterUserDTO) (uint64, error) {
	response, err := repo.client.Register(
		ctx,
		&sso.RegisterIn{
			DisplayName: userData.DisplayName,
			Email:       userData.Email,
			Password:    userData.Password,
		},
	)

	if err != nil {
		return 0, err
	}

	return response.GetUserID(), nil
}

func (repo *GrpcSsoRepository) GetUserByID(ctx context.Context, id uint64) (*entities.User, error) {
	response, err := repo.client.GetUser(
		ctx,
		&sso.GetUserIn{
			ID: id,
		},
	)

	if err != nil {
		return nil, err
	}

	return repo.processUserResponse(response), nil
}

func (repo *GrpcSsoRepository) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	response, err := repo.client.GetUsers(
		ctx,
		&emptypb.Empty{},
	)

	if err != nil {
		return nil, err
	}

	users := make([]entities.User, len(response.GetUsers()))
	for index, userResponse := range response.GetUsers() {
		users[index] = *repo.processUserResponse(userResponse)
	}

	return users, nil
}

func (repo *GrpcSsoRepository) LoginUser(
	ctx context.Context,
	userData entities.LoginUserDTO,
) (*entities.TokensDTO, error) {
	response, err := repo.client.Login(
		ctx,
		&sso.LoginIn{
			Email:    userData.Email,
			Password: userData.Password,
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
	response, err := repo.client.GetMe(
		ctx,
		&sso.GetMeIn{
			AccessToken: accessToken,
		},
	)

	if err != nil {
		return nil, err
	}

	return repo.processUserResponse(response), nil
}

func (repo *GrpcSsoRepository) RefreshTokens(ctx context.Context, refreshToken string) (*entities.TokensDTO, error) {
	response, err := repo.client.RefreshTokens(
		ctx,
		&sso.RefreshTokensIn{
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

func (repo *GrpcSsoRepository) processUserResponse(userResponse *sso.GetUserOut) *entities.User {
	return &entities.User{
		ID:                userResponse.GetID(),
		DisplayName:       userResponse.GetDisplayName(),
		Email:             userResponse.GetEmail(),
		EmailConfirmed:    userResponse.GetEmailConfirmed(),
		Phone:             userResponse.Phone,
		PhoneConfirmed:    userResponse.GetPhoneConfirmed(),
		Telegram:          userResponse.Telegram,
		TelegramConfirmed: userResponse.GetTelegramConfirmed(),
		Avatar:            userResponse.Avatar,
		CreatedAt:         userResponse.GetCreatedAt().AsTime(),
		UpdatedAt:         userResponse.GetUpdatedAt().AsTime(),
	}
}
