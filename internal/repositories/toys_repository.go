package repositories

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"
	toysentities "github.com/DKhorkov/hmtm-toys/pkg/entities"
)

type GrpcToysRepository struct {
	client interfaces.ToysGrpcClient
}

func (repo *GrpcToysRepository) AddToy(toyData toysentities.RawAddToyDTO) (uint64, error) {
	response, err := repo.client.AddToy(
		context.Background(),
		&toys.AddToyRequest{
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

func (repo *GrpcToysRepository) GetAllToys() ([]*toysentities.Toy, error) {
	response, err := repo.client.GetToys(
		context.Background(),
		&emptypb.Empty{},
	)

	if err != nil {
		return nil, err
	}

	allToys := make([]*toysentities.Toy, len(response.GetToys()))
	for index, toyResponse := range response.GetToys() {
		allToys[index] = repo.processToyResponse(toyResponse)
	}

	return allToys, nil
}

func (repo *GrpcToysRepository) GetToyByID(id uint64) (*toysentities.Toy, error) {
	response, err := repo.client.GetToy(
		context.Background(),
		&toys.GetToyRequest{
			ID: id,
		},
	)

	if err != nil {
		return nil, err
	}

	return repo.processToyResponse(response), nil
}

func (repo *GrpcToysRepository) GetAllMasters() ([]*toysentities.Master, error) {
	response, err := repo.client.GetMasters(
		context.Background(),
		&emptypb.Empty{},
	)

	if err != nil {
		return nil, err
	}

	masters := make([]*toysentities.Master, len(response.GetMasters()))
	for index, masterResponse := range response.GetMasters() {
		masters[index] = repo.processMasterResponse(masterResponse)
	}

	return masters, nil
}

func (repo *GrpcToysRepository) GetMasterByID(id uint64) (*toysentities.Master, error) {
	response, err := repo.client.GetMaster(
		context.Background(),
		&toys.GetMasterRequest{
			ID: id,
		},
	)

	if err != nil {
		return nil, err
	}

	return repo.processMasterResponse(response), nil
}

func (repo *GrpcToysRepository) RegisterMaster(masterData toysentities.RawRegisterMasterDTO) (uint64, error) {
	response, err := repo.client.RegisterMaster(
		context.Background(),
		&toys.RegisterMasterRequest{
			AccessToken: masterData.AccessToken,
			Info:        masterData.Info,
		},
	)

	if err != nil {
		return 0, err
	}

	return response.GetMasterID(), nil
}

func (repo *GrpcToysRepository) GetAllCategories() ([]*toysentities.Category, error) {
	response, err := repo.client.GetCategories(
		context.Background(),
		&emptypb.Empty{},
	)

	if err != nil {
		return nil, err
	}

	categories := make([]*toysentities.Category, len(response.GetCategories()))
	for index, categoryResponse := range response.GetCategories() {
		categories[index] = repo.processCategoryResponse(categoryResponse)
	}

	return categories, nil
}

func (repo *GrpcToysRepository) GetCategoryByID(id uint32) (*toysentities.Category, error) {
	response, err := repo.client.GetCategory(
		context.Background(),
		&toys.GetCategoryRequest{
			ID: id,
		},
	)

	if err != nil {
		return nil, err
	}

	return repo.processCategoryResponse(response), nil
}

func (repo *GrpcToysRepository) GetAllTags() ([]*toysentities.Tag, error) {
	response, err := repo.client.GetTags(
		context.Background(),
		&emptypb.Empty{},
	)

	if err != nil {
		return nil, err
	}

	tags := make([]*toysentities.Tag, len(response.GetTags()))
	for index, tagResponse := range response.GetTags() {
		tags[index] = repo.processTagResponse(tagResponse)
	}

	return tags, nil
}

func (repo *GrpcToysRepository) GetTagByID(id uint32) (*toysentities.Tag, error) {
	response, err := repo.client.GetTag(
		context.Background(),
		&toys.GetTagRequest{
			ID: id,
		},
	)

	if err != nil {
		return nil, err
	}

	return repo.processTagResponse(response), nil
}

func (repo *GrpcToysRepository) processTagResponse(tagResponse *toys.GetTagResponse) *toysentities.Tag {
	return &toysentities.Tag{
		ID:        tagResponse.GetID(),
		Name:      tagResponse.GetName(),
		CreatedAt: tagResponse.GetCreatedAt().AsTime(),
		UpdatedAt: tagResponse.GetUpdatedAt().AsTime(),
	}
}

func (repo *GrpcToysRepository) processCategoryResponse(
	categoryResponse *toys.GetCategoryResponse,
) *toysentities.Category {
	return &toysentities.Category{
		ID:        categoryResponse.GetID(),
		Name:      categoryResponse.GetName(),
		CreatedAt: categoryResponse.GetCreatedAt().AsTime(),
		UpdatedAt: categoryResponse.GetUpdatedAt().AsTime(),
	}
}

func (repo *GrpcToysRepository) processToyResponse(toyResponse *toys.GetToyResponse) *toysentities.Toy {
	tags := make([]*toysentities.Tag, len(toyResponse.GetTags()))
	for index, tagResponse := range toyResponse.GetTags() {
		tags[index] = repo.processTagResponse(tagResponse)
	}

	return &toysentities.Toy{
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

func (repo *GrpcToysRepository) processMasterResponse(masterResponse *toys.GetMasterResponse) *toysentities.Master {
	masterToys := make([]*toysentities.Toy, len(masterResponse.GetToys()))
	for index, toyResponse := range masterResponse.GetToys() {
		masterToys[index] = repo.processToyResponse(toyResponse)
	}

	return &toysentities.Master{
		ID:        masterResponse.GetID(),
		UserID:    masterResponse.GetUserID(),
		Info:      masterResponse.GetInfo(),
		CreatedAt: masterResponse.GetCreatedAt().AsTime(),
		UpdatedAt: masterResponse.GetUpdatedAt().AsTime(),
		Toys:      masterToys,
	}
}

func NewGrpcToysRepository(grpcClient interfaces.ToysGrpcClient) *GrpcToysRepository {
	return &GrpcToysRepository{client: grpcClient}
}
