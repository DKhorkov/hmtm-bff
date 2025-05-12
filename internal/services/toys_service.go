package services

import (
	"context"
	"fmt"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

type ToysService struct {
	toysRepository interfaces.ToysRepository
	logger         logging.Logger
}

func NewToysService(toysRepository interfaces.ToysRepository, logger logging.Logger) *ToysService {
	return &ToysService{
		toysRepository: toysRepository,
		logger:         logger,
	}
}

func (service *ToysService) AddToy(
	ctx context.Context,
	toyData entities.AddToyDTO,
) (uint64, error) {
	toyID, err := service.toysRepository.AddToy(ctx, toyData)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to add new Toy",
			err,
		)
	}

	return toyID, err
}

func (service *ToysService) GetToys(ctx context.Context, pagination *entities.Pagination) ([]entities.Toy, error) {
	toys, err := service.toysRepository.GetToys(ctx, pagination)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to get all Toys",
			err,
		)
	}

	return toys, err
}

func (service *ToysService) CountToys(ctx context.Context) (uint64, error) {
	count, err := service.toysRepository.CountToys(ctx)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to count all Toys",
			err,
		)
	}

	return count, err
}

func (service *ToysService) GetMasterToys(
	ctx context.Context,
	masterID uint64,
	pagination *entities.Pagination,
) ([]entities.Toy, error) {
	toys, err := service.toysRepository.GetMasterToys(ctx, masterID, pagination)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf(
				"Error occurred while trying to get all Toys for Master with ID=%d",
				masterID,
			),
			err,
		)
	}

	return toys, err
}

func (service *ToysService) GetUserToys(
	ctx context.Context,
	userID uint64,
	pagination *entities.Pagination,
) ([]entities.Toy, error) {
	toys, err := service.toysRepository.GetUserToys(ctx, userID, pagination)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get all Toys for User with ID=%d", userID),
			err,
		)
	}

	return toys, err
}

func (service *ToysService) GetToyByID(ctx context.Context, id uint64) (*entities.Toy, error) {
	toy, err := service.toysRepository.GetToyByID(ctx, id)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get Toy with ID=%d", id),
			err,
		)
	}

	return toy, err
}

func (service *ToysService) GetMasters(
	ctx context.Context,
	pagination *entities.Pagination,
) ([]entities.Master, error) {
	masters, err := service.toysRepository.GetMasters(ctx, pagination)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to get all Masters",
			err,
		)
	}

	return masters, err
}

func (service *ToysService) GetMasterByID(
	ctx context.Context,
	id uint64,
) (*entities.Master, error) {
	master, err := service.toysRepository.GetMasterByID(ctx, id)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get Master with ID=%d", id),
			err,
		)
	}

	return master, err
}

func (service *ToysService) RegisterMaster(
	ctx context.Context,
	masterData entities.RegisterMasterDTO,
) (uint64, error) {
	masterID, err := service.toysRepository.RegisterMaster(ctx, masterData)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to register Master",
			err,
		)
	}

	return masterID, err
}

func (service *ToysService) GetAllCategories(ctx context.Context) ([]entities.Category, error) {
	categories, err := service.toysRepository.GetAllCategories(ctx)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to get all Categories",
			err,
		)
	}

	return categories, err
}

func (service *ToysService) GetCategoryByID(
	ctx context.Context,
	id uint32,
) (*entities.Category, error) {
	category, err := service.toysRepository.GetCategoryByID(ctx, id)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get Category with ID=%d", id),
			err,
		)
	}

	return category, err
}

func (service *ToysService) GetAllTags(ctx context.Context) ([]entities.Tag, error) {
	tags, err := service.toysRepository.GetAllTags(ctx)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to get all Tags",
			err,
		)
	}

	return tags, err
}

func (service *ToysService) GetTagByID(ctx context.Context, id uint32) (*entities.Tag, error) {
	tag, err := service.toysRepository.GetTagByID(ctx, id)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get Tag with ID=%d", id),
			err,
		)
	}

	return tag, err
}

func (service *ToysService) CreateTags(
	ctx context.Context,
	tagsData []entities.CreateTagDTO,
) ([]uint32, error) {
	tagIDs, err := service.toysRepository.CreateTags(ctx, tagsData)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to create Tags",
			err,
		)
	}

	return tagIDs, err
}

func (service *ToysService) UpdateToy(ctx context.Context, toyData entities.UpdateToyDTO) error {
	err := service.toysRepository.UpdateToy(ctx, toyData)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to update Toy with ID=%d", toyData.ID),
			err,
		)
	}

	return err
}

func (service *ToysService) DeleteToy(ctx context.Context, id uint64) error {
	err := service.toysRepository.DeleteToy(ctx, id)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to delete Toy with ID=%d", id),
			err,
		)
	}

	return err
}

func (service *ToysService) GetMasterByUserID(
	ctx context.Context,
	userID uint64,
) (*entities.Master, error) {
	master, err := service.toysRepository.GetMasterByUserID(ctx, userID)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get Master with userID=%d", userID),
			err,
		)
	}

	return master, err
}

func (service *ToysService) UpdateMaster(
	ctx context.Context,
	masterData entities.UpdateMasterDTO,
) error {
	err := service.toysRepository.UpdateMaster(ctx, masterData)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to update Mastrer with ID=%d", masterData.ID),
			err,
		)
	}

	return err
}
