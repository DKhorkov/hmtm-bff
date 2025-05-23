package services

import (
	"context"
	"fmt"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

type SsoService struct {
	ssoRepository interfaces.SsoRepository
	logger        logging.Logger
}

func NewSsoService(ssoRepository interfaces.SsoRepository, logger logging.Logger) *SsoService {
	return &SsoService{
		ssoRepository: ssoRepository,
		logger:        logger,
	}
}

func (service *SsoService) GetUsers(ctx context.Context, pagination *entities.Pagination) ([]entities.User, error) {
	users, err := service.ssoRepository.GetUsers(ctx, pagination)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to get all Users",
			err,
		)
	}

	return users, err
}

func (service *SsoService) GetUserByID(ctx context.Context, id uint64) (*entities.User, error) {
	user, err := service.ssoRepository.GetUserByID(ctx, id)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get User with ID=%d", id),
			err,
		)
	}

	return user, err
}

func (service *SsoService) GetUserByEmail(
	ctx context.Context,
	email string,
) (*entities.User, error) {
	user, err := service.ssoRepository.GetUserByEmail(ctx, email)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to get User with Email="+email,
			err,
		)
	}

	return user, err
}

func (service *SsoService) RegisterUser(
	ctx context.Context,
	userData entities.RegisterUserDTO,
) (uint64, error) {
	userID, err := service.ssoRepository.RegisterUser(ctx, userData)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to register User",
			err,
		)
	}

	return userID, err
}

func (service *SsoService) LoginUser(
	ctx context.Context,
	userData entities.LoginUserDTO,
) (*entities.TokensDTO, error) {
	tokens, err := service.ssoRepository.LoginUser(ctx, userData)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to login User with email="+userData.Email,
			err,
		)
	}

	return tokens, err
}

func (service *SsoService) LogoutUser(ctx context.Context, accessToken string) error {
	err := service.ssoRepository.LogoutUser(ctx, accessToken)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to logout User with AccessToken="+accessToken,
			err,
		)
	}

	return err
}

func (service *SsoService) VerifyUserEmail(ctx context.Context, verifyEmailToken string) error {
	err := service.ssoRepository.VerifyUserEmail(ctx, verifyEmailToken)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to verify email for User with verify email token="+verifyEmailToken,
			err,
		)
	}

	return err
}

func (service *SsoService) ForgetPassword(ctx context.Context, forgetPasswordToken, newPassword string) error {
	err := service.ssoRepository.ForgetPassword(ctx, forgetPasswordToken, newPassword)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to forget password for User with ForgetPasswordToken="+forgetPasswordToken,
			err,
		)
	}

	return err
}

func (service *SsoService) ChangePassword(
	ctx context.Context,
	accessToken string,
	oldPassword string,
	newPassword string,
) error {
	err := service.ssoRepository.ChangePassword(ctx, accessToken, oldPassword, newPassword)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to change password for User with AccessToken="+accessToken,
			err,
		)
	}

	return err
}

func (service *SsoService) SendVerifyEmailMessage(ctx context.Context, email string) error {
	err := service.ssoRepository.SendVerifyEmailMessage(ctx, email)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to send verify-email message to User with email="+email,
			err,
		)
	}

	return err
}

func (service *SsoService) SendForgetPasswordMessage(ctx context.Context, email string) error {
	err := service.ssoRepository.SendForgetPasswordMessage(ctx, email)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to send forget-password message to User with email="+email,
			err,
		)
	}

	return err
}

func (service *SsoService) UpdateUserProfile(
	ctx context.Context,
	userProfileData entities.UpdateUserProfileDTO,
) error {
	err := service.ssoRepository.UpdateUserProfile(ctx, userProfileData)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to update profile for User with AccessToken="+userProfileData.AccessToken,
			err,
		)
	}

	return err
}

func (service *SsoService) GetMe(ctx context.Context, accessToken string) (*entities.User, error) {
	user, err := service.ssoRepository.GetMe(ctx, accessToken)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to get User with AccessToken="+accessToken,
			err,
		)
	}

	return user, err
}

func (service *SsoService) RefreshTokens(
	ctx context.Context,
	refreshToken string,
) (*entities.TokensDTO, error) {
	tokens, err := service.ssoRepository.RefreshTokens(ctx, refreshToken)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to refresh tokens with RefreshToken="+refreshToken,
			err,
		)
	}

	return tokens, err
}
