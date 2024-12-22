package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

func NewCommonToysService(toysRepository interfaces.ToysRepository, logger *slog.Logger) *CommonToysService {
	return &CommonToysService{
		toysRepository: toysRepository,
		logger:         logger,
	}
}

type CommonToysService struct {
	toysRepository interfaces.ToysRepository
	logger         *slog.Logger
}

func (service *CommonToysService) AddToy(ctx context.Context, toyData entities.AddToyDTO) (uint64, error) {
	toyID, err := service.toysRepository.AddToy(ctx, toyData)
	if err != nil {
		logging.LogErrorContext(ctx, service.logger, "Error occurred while trying to add new Toy", err)
	}

	return toyID, err
}

func (service *CommonToysService) GetAllToys(ctx context.Context) ([]entities.Toy, error) {
	toys, err := service.toysRepository.GetAllToys(ctx)
	if err != nil {
		logging.LogErrorContext(ctx, service.logger, "Error occurred while trying to get all Toys", err)
	}

	return toys, err
}

func (service *CommonToysService) GetMasterToys(ctx context.Context, masterID uint64) ([]entities.Toy, error) {
	toys, err := service.toysRepository.GetMasterToys(ctx, masterID)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get all Toys for Master with ID=%d", masterID),
			err,
		)
	}

	return toys, err
}

func (service *CommonToysService) GetToyByID(ctx context.Context, id uint64) (*entities.Toy, error) {
	toy, err := service.toysRepository.GetToyByID(ctx, id)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get Toy with ID=%v", id),
			err,
		)
	}

	return toy, err
}

func (service *CommonToysService) GetAllMasters(ctx context.Context) ([]entities.Master, error) {
	masters, err := service.toysRepository.GetAllMasters(ctx)
	if err != nil {
		logging.LogErrorContext(ctx, service.logger, "Error occurred while trying to get all Masters", err)
	}

	return masters, err
}

func (service *CommonToysService) GetMasterByID(ctx context.Context, id uint64) (*entities.Master, error) {
	master, err := service.toysRepository.GetMasterByID(ctx, id)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get Master with ID=%v", id),
			err,
		)
	}

	return master, err
}

func (service *CommonToysService) RegisterMaster(
	ctx context.Context,
	masterData entities.RegisterMasterDTO,
) (uint64, error) {
	masterID, err := service.toysRepository.RegisterMaster(ctx, masterData)
	if err != nil {
		logging.LogErrorContext(ctx, service.logger, "Error occurred while trying to register Master", err)
	}

	return masterID, err
}

func (service *CommonToysService) GetAllCategories(ctx context.Context) ([]entities.Category, error) {
	categories, err := service.toysRepository.GetAllCategories(ctx)
	if err != nil {
		logging.LogErrorContext(ctx, service.logger, "Error occurred while trying to get all Categories", err)
	}

	return categories, err
}

func (service *CommonToysService) GetCategoryByID(ctx context.Context, id uint32) (*entities.Category, error) {
	category, err := service.toysRepository.GetCategoryByID(ctx, id)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get Category with ID=%v", id),
			err,
		)
	}

	return category, err
}

func (service *CommonToysService) GetAllTags(ctx context.Context) ([]entities.Tag, error) {
	tags, err := service.toysRepository.GetAllTags(ctx)
	if err != nil {
		logging.LogErrorContext(ctx, service.logger, "Error occurred while trying to get all Tags", err)
	}

	return tags, err
}

func (service *CommonToysService) GetTagByID(ctx context.Context, id uint32) (*entities.Tag, error) {
	tag, err := service.toysRepository.GetTagByID(ctx, id)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get Tag with ID=%v", id),
			err,
		)
	}

	return tag, err
}
