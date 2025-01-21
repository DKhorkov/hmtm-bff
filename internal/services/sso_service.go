package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

func NewCommonSsoService(ssoRepository interfaces.SsoRepository, logger *slog.Logger) *CommonSsoService {
	return &CommonSsoService{
		ssoRepository: ssoRepository,
		logger:        logger,
	}
}

type CommonSsoService struct {
	ssoRepository interfaces.SsoRepository
	logger        *slog.Logger
}

func (service *CommonSsoService) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	users, err := service.ssoRepository.GetAllUsers(ctx)
	if err != nil {
		logging.LogErrorContext(ctx, service.logger, "Error occurred while trying to get all Users", err)
	}

	return users, err
}

func (service *CommonSsoService) GetUserByID(ctx context.Context, id uint64) (*entities.User, error) {
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

func (service *CommonSsoService) RegisterUser(ctx context.Context, userData entities.RegisterUserDTO) (uint64, error) {
	userID, err := service.ssoRepository.RegisterUser(ctx, userData)
	if err != nil {
		logging.LogErrorContext(ctx, service.logger, "Error occurred while trying to register User", err)
	}

	return userID, err
}

func (service *CommonSsoService) LoginUser(
	ctx context.Context,
	userData entities.LoginUserDTO,
) (*entities.TokensDTO, error) {
	tokens, err := service.ssoRepository.LoginUser(ctx, userData)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to login User with email=%s", userData.Email),
			err,
		)
	}

	return tokens, err
}

func (service *CommonSsoService) LogoutUser(ctx context.Context, accessToken string) error {
	err := service.ssoRepository.LogoutUser(ctx, accessToken)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to logout User with AccessToken=%s", accessToken),
			err,
		)
	}

	return err
}

func (service *CommonSsoService) GetMe(ctx context.Context, accessToken string) (*entities.User, error) {
	user, err := service.ssoRepository.GetMe(ctx, accessToken)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get User with AccessToken=%s", accessToken),
			err,
		)
	}

	return user, err
}

func (service *CommonSsoService) RefreshTokens(ctx context.Context, refreshToken string) (*entities.TokensDTO, error) {
	tokens, err := service.ssoRepository.RefreshTokens(ctx, refreshToken)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to refresh tokens with RefreshToken=%s", refreshToken),
			err,
		)
	}

	return tokens, err
}
