package repositories

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

func NewToysRepository(client interfaces.ToysClient) *ToysRepository {
	return &ToysRepository{client: client}
}

type ToysRepository struct {
	client interfaces.ToysClient
}

func (repo *ToysRepository) AddToy(ctx context.Context, toyData entities.AddToyDTO) (uint64, error) {
	response, err := repo.client.AddToy(
		ctx,
		&toys.AddToyIn{
			UserID:      toyData.UserID,
			CategoryID:  toyData.CategoryID,
			Name:        toyData.Name,
			Description: toyData.Description,
			Price:       toyData.Price,
			Quantity:    toyData.Quantity,
			TagIDs:      toyData.TagIDs,
			Attachments: toyData.Attachments,
		},
	)

	if err != nil {
		return 0, err
	}

	return response.GetToyID(), nil
}

func (repo *ToysRepository) GetAllToys(ctx context.Context) ([]entities.Toy, error) {
	response, err := repo.client.GetToys(
		ctx,
		&emptypb.Empty{},
	)

	if err != nil {
		return nil, err
	}

	allToys := make([]entities.Toy, len(response.GetToys()))
	for i, toyResponse := range response.GetToys() {
		allToys[i] = *repo.processToyResponse(toyResponse)
	}

	return allToys, nil
}

func (repo *ToysRepository) GetMasterToys(ctx context.Context, masterID uint64) ([]entities.Toy, error) {
	response, err := repo.client.GetMasterToys(
		ctx,
		&toys.GetMasterToysIn{
			MasterID: masterID,
		},
	)

	if err != nil {
		return nil, err
	}

	masterToys := make([]entities.Toy, len(response.GetToys()))
	for i, toyResponse := range response.GetToys() {
		masterToys[i] = *repo.processToyResponse(toyResponse)
	}

	return masterToys, nil
}

func (repo *ToysRepository) GetUserToys(ctx context.Context, userID uint64) ([]entities.Toy, error) {
	response, err := repo.client.GetUserToys(
		ctx,
		&toys.GetUserToysIn{
			UserID: userID,
		},
	)

	if err != nil {
		return nil, err
	}

	userToys := make([]entities.Toy, len(response.GetToys()))
	for i, toyResponse := range response.GetToys() {
		userToys[i] = *repo.processToyResponse(toyResponse)
	}

	return userToys, nil
}

func (repo *ToysRepository) GetToyByID(ctx context.Context, id uint64) (*entities.Toy, error) {
	response, err := repo.client.GetToy(
		ctx,
		&toys.GetToyIn{
			ID: id,
		},
	)

	if err != nil {
		return nil, err
	}

	return repo.processToyResponse(response), nil
}

func (repo *ToysRepository) GetAllMasters(ctx context.Context) ([]entities.Master, error) {
	response, err := repo.client.GetMasters(
		ctx,
		&emptypb.Empty{},
	)

	if err != nil {
		return nil, err
	}

	masters := make([]entities.Master, len(response.GetMasters()))
	for i, masterResponse := range response.GetMasters() {
		masters[i] = entities.Master{
			ID:        masterResponse.GetID(),
			UserID:    masterResponse.GetUserID(),
			Info:      masterResponse.GetInfo(),
			CreatedAt: masterResponse.GetCreatedAt().AsTime(),
			UpdatedAt: masterResponse.GetUpdatedAt().AsTime(),
		}
	}

	return masters, nil
}

func (repo *ToysRepository) GetMasterByID(ctx context.Context, id uint64) (*entities.Master, error) {
	response, err := repo.client.GetMaster(
		ctx,
		&toys.GetMasterIn{
			ID: id,
		},
	)

	if err != nil {
		return nil, err
	}

	return &entities.Master{
		ID:        response.GetID(),
		UserID:    response.GetUserID(),
		Info:      response.GetInfo(),
		CreatedAt: response.GetCreatedAt().AsTime(),
		UpdatedAt: response.GetUpdatedAt().AsTime(),
	}, nil
}

func (repo *ToysRepository) RegisterMaster(
	ctx context.Context,
	masterData entities.RegisterMasterDTO,
) (uint64, error) {
	response, err := repo.client.RegisterMaster(
		ctx,
		&toys.RegisterMasterIn{
			UserID: masterData.UserID,
			Info:   masterData.Info,
		},
	)

	if err != nil {
		return 0, err
	}

	return response.GetMasterID(), nil
}

func (repo *ToysRepository) GetAllCategories(ctx context.Context) ([]entities.Category, error) {
	response, err := repo.client.GetCategories(
		ctx,
		&emptypb.Empty{},
	)

	if err != nil {
		return nil, err
	}

	categories := make([]entities.Category, len(response.GetCategories()))
	for i, categoryResponse := range response.GetCategories() {
		categories[i] = *repo.processCategoryResponse(categoryResponse)
	}

	return categories, nil
}

func (repo *ToysRepository) GetCategoryByID(ctx context.Context, id uint32) (*entities.Category, error) {
	response, err := repo.client.GetCategory(
		ctx,
		&toys.GetCategoryIn{
			ID: id,
		},
	)

	if err != nil {
		return nil, err
	}

	return repo.processCategoryResponse(response), nil
}

func (repo *ToysRepository) GetAllTags(ctx context.Context) ([]entities.Tag, error) {
	response, err := repo.client.GetTags(
		ctx,
		&emptypb.Empty{},
	)

	if err != nil {
		return nil, err
	}

	tags := make([]entities.Tag, len(response.GetTags()))
	for i, tagResponse := range response.GetTags() {
		tags[i] = *repo.processTagResponse(tagResponse)
	}

	return tags, nil
}

func (repo *ToysRepository) GetTagByID(ctx context.Context, id uint32) (*entities.Tag, error) {
	response, err := repo.client.GetTag(
		ctx,
		&toys.GetTagIn{
			ID: id,
		},
	)

	if err != nil {
		return nil, err
	}

	return repo.processTagResponse(response), nil
}

func (repo *ToysRepository) processTagResponse(tagResponse *toys.GetTagOut) *entities.Tag {
	return &entities.Tag{
		ID:   tagResponse.GetID(),
		Name: tagResponse.GetName(),
	}
}

func (repo *ToysRepository) processCategoryResponse(categoryResponse *toys.GetCategoryOut) *entities.Category {
	return &entities.Category{
		ID:   categoryResponse.GetID(),
		Name: categoryResponse.GetName(),
	}
}

func (repo *ToysRepository) processToyResponse(toyResponse *toys.GetToyOut) *entities.Toy {
	tags := make([]entities.Tag, len(toyResponse.GetTags()))
	for i, tagResponse := range toyResponse.GetTags() {
		tags[i] = *repo.processTagResponse(tagResponse)
	}

	attachments := make([]entities.ToyAttachment, len(toyResponse.GetAttachments()))
	for i, attachment := range toyResponse.GetAttachments() {
		attachments[i] = entities.ToyAttachment{
			ID:        attachment.GetID(),
			ToyID:     attachment.GetToyID(),
			Link:      attachment.GetLink(),
			CreatedAt: attachment.GetCreatedAt().AsTime(),
			UpdatedAt: attachment.GetUpdatedAt().AsTime(),
		}
	}

	return &entities.Toy{
		ID:          toyResponse.GetID(),
		MasterID:    toyResponse.GetMasterID(),
		CategoryID:  toyResponse.GetCategoryID(),
		Name:        toyResponse.GetName(),
		Description: toyResponse.GetDescription(),
		Price:       toyResponse.GetPrice(),
		Quantity:    toyResponse.GetQuantity(),
		Tags:        tags,
		Attachments: attachments,
		CreatedAt:   toyResponse.GetCreatedAt().AsTime(),
		UpdatedAt:   toyResponse.GetUpdatedAt().AsTime(),
	}
}

func (repo *ToysRepository) CreateTags(ctx context.Context, tagsData []entities.CreateTagDTO) ([]uint32, error) {
	tagsRequest := make([]*toys.CreateTagIn, len(tagsData))
	for i, tag := range tagsData {
		tagsRequest[i] = &toys.CreateTagIn{
			Name: tag.Name,
		}
	}

	response, err := repo.client.CreateTags(ctx, &toys.CreateTagsIn{Tags: tagsRequest})
	if err != nil {
		return nil, err
	}

	tagIDs := make([]uint32, len(response.GetTags()))
	for i, tagResponse := range response.GetTags() {
		tagIDs[i] = tagResponse.GetID()
	}

	return tagIDs, nil
}
