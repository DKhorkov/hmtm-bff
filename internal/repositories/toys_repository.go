package repositories

import (
	"context"

	"github.com/DKhorkov/libs/contextlib"
	"github.com/DKhorkov/libs/requestid"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"
)

func NewGrpcToysRepository(grpcClient interfaces.ToysGrpcClient) *GrpcToysRepository {
	return &GrpcToysRepository{client: grpcClient}
}

type GrpcToysRepository struct {
	client interfaces.ToysGrpcClient
}

func (repo *GrpcToysRepository) AddToy(ctx context.Context, toyData entities.AddToyDTO) (uint64, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.AddToy(
		ctx,
		&toys.AddToyRequest{
			RequestID:   requestID,
			AccessToken: toyData.AccessToken,
			CategoryID:  toyData.CategoryID,
			Name:        toyData.Name,
			Description: toyData.Description,
			Price:       toyData.Price,
			Quantity:    toyData.Quantity,
			TagIDs:      toyData.TagsIDs,
		},
	)

	if err != nil {
		return 0, err
	}

	return response.GetToyID(), nil
}

func (repo *GrpcToysRepository) GetAllToys(ctx context.Context) ([]entities.Toy, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.GetToys(
		ctx,
		&toys.GetToysRequest{RequestID: requestID},
	)

	if err != nil {
		return nil, err
	}

	allToys := make([]entities.Toy, len(response.GetToys()))
	for index, toyResponse := range response.GetToys() {
		allToys[index] = *repo.processToyResponse(toyResponse)
	}

	return allToys, nil
}

func (repo *GrpcToysRepository) GetMasterToys(ctx context.Context, masterID uint64) ([]entities.Toy, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.GetMasterToys(
		ctx,
		&toys.GetMasterToysRequest{
			RequestID: requestID,
			MasterID:  masterID,
		},
	)

	if err != nil {
		return nil, err
	}

	masterToys := make([]entities.Toy, len(response.GetToys()))
	for index, toyResponse := range response.GetToys() {
		masterToys[index] = *repo.processToyResponse(toyResponse)
	}

	return masterToys, nil
}

func (repo *GrpcToysRepository) GetToyByID(ctx context.Context, id uint64) (*entities.Toy, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.GetToy(
		ctx,
		&toys.GetToyRequest{
			RequestID: requestID,
			ID:        id,
		},
	)

	if err != nil {
		return nil, err
	}

	return repo.processToyResponse(response), nil
}

func (repo *GrpcToysRepository) GetAllMasters(ctx context.Context) ([]entities.Master, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.GetMasters(
		ctx,
		&toys.GetMastersRequest{RequestID: requestID},
	)

	if err != nil {
		return nil, err
	}

	masters := make([]entities.Master, len(response.GetMasters()))
	for index, masterResponse := range response.GetMasters() {
		masters[index] = entities.Master{
			ID:        masterResponse.GetID(),
			UserID:    masterResponse.GetUserID(),
			Info:      masterResponse.GetInfo(),
			CreatedAt: masterResponse.GetCreatedAt().AsTime(),
			UpdatedAt: masterResponse.GetUpdatedAt().AsTime(),
		}
	}

	return masters, nil
}

func (repo *GrpcToysRepository) GetMasterByID(ctx context.Context, id uint64) (*entities.Master, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.GetMaster(
		ctx,
		&toys.GetMasterRequest{
			RequestID: requestID,
			ID:        id,
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

func (repo *GrpcToysRepository) RegisterMaster(
	ctx context.Context,
	masterData entities.RegisterMasterDTO,
) (uint64, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.RegisterMaster(
		ctx,
		&toys.RegisterMasterRequest{
			RequestID:   requestID,
			AccessToken: masterData.AccessToken,
			Info:        masterData.Info,
		},
	)

	if err != nil {
		return 0, err
	}

	return response.GetMasterID(), nil
}

func (repo *GrpcToysRepository) GetAllCategories(ctx context.Context) ([]entities.Category, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.GetCategories(
		ctx,
		&toys.GetCategoriesRequest{RequestID: requestID},
	)

	if err != nil {
		return nil, err
	}

	categories := make([]entities.Category, len(response.GetCategories()))
	for index, categoryResponse := range response.GetCategories() {
		categories[index] = *repo.processCategoryResponse(categoryResponse)
	}

	return categories, nil
}

func (repo *GrpcToysRepository) GetCategoryByID(ctx context.Context, id uint32) (*entities.Category, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.GetCategory(
		ctx,
		&toys.GetCategoryRequest{
			RequestID: requestID,
			ID:        id,
		},
	)

	if err != nil {
		return nil, err
	}

	return repo.processCategoryResponse(response), nil
}

func (repo *GrpcToysRepository) GetAllTags(ctx context.Context) ([]entities.Tag, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.GetTags(
		ctx,
		&toys.GetTagsRequest{RequestID: requestID},
	)

	if err != nil {
		return nil, err
	}

	tags := make([]entities.Tag, len(response.GetTags()))
	for index, tagResponse := range response.GetTags() {
		tags[index] = *repo.processTagResponse(tagResponse)
	}

	return tags, nil
}

func (repo *GrpcToysRepository) GetTagByID(ctx context.Context, id uint32) (*entities.Tag, error) {
	requestID, _ := contextlib.GetValue[string](ctx, requestid.Key)
	response, err := repo.client.GetTag(
		ctx,
		&toys.GetTagRequest{
			RequestID: requestID,
			ID:        id,
		},
	)

	if err != nil {
		return nil, err
	}

	return repo.processTagResponse(response), nil
}

func (repo *GrpcToysRepository) processTagResponse(tagResponse *toys.GetTagResponse) *entities.Tag {
	return &entities.Tag{
		ID:   tagResponse.GetID(),
		Name: tagResponse.GetName(),
	}
}

func (repo *GrpcToysRepository) processCategoryResponse(
	categoryResponse *toys.GetCategoryResponse,
) *entities.Category {
	return &entities.Category{
		ID:   categoryResponse.GetID(),
		Name: categoryResponse.GetName(),
	}
}

func (repo *GrpcToysRepository) processToyResponse(toyResponse *toys.GetToyResponse) *entities.Toy {
	tags := make([]entities.Tag, len(toyResponse.GetTags()))
	for index, tagResponse := range toyResponse.GetTags() {
		tags[index] = *repo.processTagResponse(tagResponse)
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
		CreatedAt:   toyResponse.GetCreatedAt().AsTime(),
		UpdatedAt:   toyResponse.GetUpdatedAt().AsTime(),
	}
}
