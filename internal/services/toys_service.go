package services

import (
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
	"github.com/DKhorkov/hmtm-bff/internal/models"
)

type CommonToysService struct {
	toysRepository interfaces.ToysRepository
}

func (service *CommonToysService) AddToy(toyData models.AddToyDTO) (uint64, error) {
	return service.toysRepository.AddToy(toyData)
}

func (service *CommonToysService) GetAllToys() ([]models.Toy, error) {
	return service.toysRepository.GetAllToys()
}

func (service *CommonToysService) GetMasterToys(masterID uint64) ([]models.Toy, error) {
	return service.toysRepository.GetMasterToys(masterID)
}

func (service *CommonToysService) GetToyByID(id uint64) (*models.Toy, error) {
	return service.toysRepository.GetToyByID(id)
}

func (service *CommonToysService) GetAllMasters() ([]models.Master, error) {
	return service.toysRepository.GetAllMasters()
}

func (service *CommonToysService) GetMasterByID(id uint64) (*models.Master, error) {
	return service.toysRepository.GetMasterByID(id)
}

func (service *CommonToysService) RegisterMaster(masterData models.RegisterMasterDTO) (uint64, error) {
	return service.toysRepository.RegisterMaster(masterData)
}

func (service *CommonToysService) GetAllCategories() ([]models.Category, error) {
	return service.toysRepository.GetAllCategories()
}

func (service *CommonToysService) GetCategoryByID(id uint32) (*models.Category, error) {
	return service.toysRepository.GetCategoryByID(id)
}

func (service *CommonToysService) GetAllTags() ([]models.Tag, error) {
	return service.toysRepository.GetAllTags()
}

func (service *CommonToysService) GetTagByID(id uint32) (*models.Tag, error) {
	return service.toysRepository.GetTagByID(id)
}

func NewCommonToysService(toysRepository interfaces.ToysRepository) *CommonToysService {
	return &CommonToysService{toysRepository: toysRepository}
}
