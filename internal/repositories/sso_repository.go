package repositories

import (
	"context"
	ssogrpcclient "github.com/DKhorkov/hmtm-bff/internal/clients/sso/grpc"
	ssoentities "github.com/DKhorkov/hmtm-sso/entities"
	"github.com/DKhorkov/hmtm-sso/protobuf/generated/go/sso"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcSsoRepository struct {
	Client *ssogrpcclient.Client
}

func (repo *GrpcSsoRepository) RegisterUser(userData ssoentities.RegisterUserDTO) (int, error) {
	response, err := repo.Client.Auth.Register(
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

	return int(response.UserID), nil
}

func (repo *GrpcSsoRepository) GetUserByID(id int) (*ssoentities.User, error) {
	response, err := repo.Client.Users.GetUser(
		context.Background(),
		&sso.GetUserRequest{
			UserID: int64(id),
		},
	)

	if err != nil {
		return nil, err
	}

	return &ssoentities.User{
		ID:        int(response.UserID),
		Email:     response.Email,
		CreatedAt: response.CreatedAt.AsTime(),
		UpdatedAt: response.UpdatedAt.AsTime(),
	}, nil
}

func (repo *GrpcSsoRepository) GetAllUsers() ([]*ssoentities.User, error) {
	response, err := repo.Client.Users.GetUsers(
		context.Background(),
		&emptypb.Empty{},
	)

	if err != nil {
		return nil, err
	}

	users := make([]*ssoentities.User, len(response.Users))
	for index, user := range response.Users {
		users[index] = &ssoentities.User{
			ID:        int(user.UserID),
			Email:     user.Email,
			CreatedAt: user.CreatedAt.AsTime(),
			UpdatedAt: user.UpdatedAt.AsTime(),
		}
	}

	return users, nil
}

func (repo *GrpcSsoRepository) LoginUser(userData ssoentities.LoginUserDTO) (string, error) {
	response, err := repo.Client.Auth.Login(
		context.Background(),
		&sso.LoginRequest{
			Email:    userData.Email,
			Password: userData.Password,
		},
	)

	if err != nil {
		return "", err
	}

	return response.Token, nil
}
