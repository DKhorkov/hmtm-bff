package services

import (
	"context"
	"fmt"

	"github.com/DKhorkov/libs/logging"

	customerrors "github.com/DKhorkov/hmtm-bff/internal/errors"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

func NewFileStorageService(
	fileStorageRepository interfaces.FileStorageRepository,
	logger logging.Logger,
) *FileStorageService {
	return &FileStorageService{
		fileStorageRepository: fileStorageRepository,
		logger:                logger,
	}
}

type FileStorageService struct {
	fileStorageRepository interfaces.FileStorageRepository
	logger                logging.Logger
}

func (service *FileStorageService) Upload(
	ctx context.Context,
	key string,
	data []byte,
) (string, error) {
	url, err := service.fileStorageRepository.Upload(ctx, key, data)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to upload File with key=%s", key),
			err,
		)

		return "", &customerrors.UploadFileError{Message: key}
	}

	return url, nil
}

func (service *FileStorageService) Delete(ctx context.Context, key string) error {
	err := service.fileStorageRepository.Delete(ctx, key)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to delete File with key=%s", key),
			err,
		)
	}

	return err
}

func (service *FileStorageService) DeleteMany(ctx context.Context, keys []string) []error {
	deleteErrors := service.fileStorageRepository.DeleteMany(ctx, keys)
	if deleteErrors != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Errors occurred while trying to delete Files with keys=%s", keys),
			fmt.Errorf("%v", deleteErrors),
		)
	}

	return deleteErrors
}
