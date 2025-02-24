package repositories

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

func NewSsoRepository(client interfaces.SsoClient) *SsoRepository {
	return &SsoRepository{client: client}
}

type SsoRepository struct {
	client interfaces.SsoClient
}

func (repo *SsoRepository) RegisterUser(ctx context.Context, userData entities.RegisterUserDTO) (uint64, error) {
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

func (repo *SsoRepository) GetUserByID(ctx context.Context, id uint64) (*entities.User, error) {
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

func (repo *SsoRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	response, err := repo.client.GetUserByEmail(
		ctx,
		&sso.GetUserByEmailIn{
			Email: email,
		},
	)

	if err != nil {
		return nil, err
	}

	return repo.processUserResponse(response), nil
}

func (repo *SsoRepository) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	response, err := repo.client.GetUsers(
		ctx,
		&emptypb.Empty{},
	)

	if err != nil {
		return nil, err
	}

	users := make([]entities.User, len(response.GetUsers()))
	for i, userResponse := range response.GetUsers() {
		users[i] = *repo.processUserResponse(userResponse)
	}

	return users, nil
}

func (repo *SsoRepository) LoginUser(
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

func (repo *SsoRepository) LogoutUser(ctx context.Context, accessToken string) error {
	_, err := repo.client.Logout(
		ctx,
		&sso.LogoutIn{
			AccessToken: accessToken,
		},
	)

	return err
}

func (repo *SsoRepository) VerifyUserEmail(ctx context.Context, verifyEmailToken string) error {
	_, err := repo.client.VerifyEmail(
		ctx,
		&sso.VerifyEmailIn{
			VerifyEmailToken: verifyEmailToken,
		},
	)

	return err
}

func (repo *SsoRepository) ForgetPassword(ctx context.Context, accessToken string) error {
	_, err := repo.client.ForgetPassword(
		ctx,
		&sso.ForgetPasswordIn{
			AccessToken: accessToken,
		},
	)

	return err
}

func (repo *SsoRepository) ChangePassword(
	ctx context.Context,
	accessToken string,
	oldPassword string,
	newPassword string,
) error {
	_, err := repo.client.ChangePassword(
		ctx,
		&sso.ChangePasswordIn{
			AccessToken: accessToken,
			OldPassword: oldPassword,
			NewPassword: newPassword,
		},
	)

	return err
}

func (repo *SsoRepository) SendVerifyEmailMessage(ctx context.Context, email string) error {
	_, err := repo.client.SendVerifyEmailMessage(
		ctx,
		&sso.SendVerifyEmailMessageIn{
			Email: email,
		},
	)

	return err
}

func (repo *SsoRepository) GetMe(ctx context.Context, accessToken string) (*entities.User, error) {
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

func (repo *SsoRepository) RefreshTokens(ctx context.Context, refreshToken string) (*entities.TokensDTO, error) {
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

func (repo *SsoRepository) UpdateUserProfile(ctx context.Context, userProfileData entities.UpdateUserProfileDTO) error {
	_, err := repo.client.UpdateUserProfile(
		ctx,
		&sso.UpdateUserProfileIn{
			AccessToken: userProfileData.AccessToken,
			DisplayName: userProfileData.DisplayName,
			Phone:       userProfileData.Phone,
			Telegram:    userProfileData.Telegram,
			Avatar:      userProfileData.Avatar,
		},
	)

	return err
}

func (repo *SsoRepository) processUserResponse(userResponse *sso.GetUserOut) *entities.User {
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
