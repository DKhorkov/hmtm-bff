package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/DKhorkov/libs/logging"

	customerrors "github.com/DKhorkov/hmtm-bff/internal/errors"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

func NewCommonFileStorageService(
	fileStorageRepository interfaces.FileStorageRepository,
	logger *slog.Logger,
) *CommonFileStorageService {
	return &CommonFileStorageService{
		fileStorageRepository: fileStorageRepository,
		logger:                logger,
	}
}

type CommonFileStorageService struct {
	fileStorageRepository interfaces.FileStorageRepository
	logger                *slog.Logger
}

func (service *CommonFileStorageService) Upload(ctx context.Context, key string, data []byte) (string, error) {
	url, err := service.fileStorageRepository.Upload(ctx, key, data)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to upload File with key=%s", key),
			err,
		)

		return "", customerrors.UploadFileError{Message: key}
	}

	return url, nil
}
