package services

import (
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
	toysentities "github.com/DKhorkov/hmtm-toys/pkg/entities"
)

type CommonToysService struct {
	toysRepository interfaces.ToysRepository
}

func (service *CommonToysService) AddToy(toyData toysentities.RawAddToyDTO) (uint64, error) {
	return service.toysRepository.AddToy(toyData)
}

func (service *CommonToysService) GetAllToys() ([]*toysentities.Toy, error) {
	return service.toysRepository.GetAllToys()
}

func (service *CommonToysService) GetToyByID(id uint64) (*toysentities.Toy, error) {
	return service.toysRepository.GetToyByID(id)
}

func (service *CommonToysService) GetAllMasters() ([]*toysentities.Master, error) {
	return service.toysRepository.GetAllMasters()
}

func (service *CommonToysService) GetMasterByID(id uint64) (*toysentities.Master, error) {
	return service.toysRepository.GetMasterByID(id)
}

func (service *CommonToysService) RegisterMaster(masterData toysentities.RawRegisterMasterDTO) (uint64, error) {
	return service.toysRepository.RegisterMaster(masterData)
}

func (service *CommonToysService) GetAllCategories() ([]*toysentities.Category, error) {
	return service.toysRepository.GetAllCategories()
}

func (service *CommonToysService) GetCategoryByID(id uint32) (*toysentities.Category, error) {
	return service.toysRepository.GetCategoryByID(id)
}

func (service *CommonToysService) GetAllTags() ([]*toysentities.Tag, error) {
	return service.toysRepository.GetAllTags()
}

func (service *CommonToysService) GetTagByID(id uint32) (*toysentities.Tag, error) {
	return service.toysRepository.GetTagByID(id)
}

func NewCommonToysService(toysRepository interfaces.ToysRepository) *CommonToysService {
	return &CommonToysService{toysRepository: toysRepository}
}
