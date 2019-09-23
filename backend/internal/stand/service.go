package stand

import (
	"fmt"
	"strconv"

	"github.com/kyma-incubator/bullseye-showcase/backend/internal/db/stands"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/ec"
)

// WithProductDetailsDTO represents mapping between shelf and product.
type WithProductDetailsDTO struct {
	ID      int           `json:"standNr"`
	Product ec.ProductDTO `json:"product"`
}

type productService interface {
	GetProductDetailsByID(ID string) (ec.ProductDTO, error)
}

type standsRepository interface {
	GetAllActiveStands() ([]stands.StandDTO, error)
	GetAllActiveStandsMap() (map[int]string, error)
}

// Service is a Stand Service that performs operations on stands.
type Service struct {
	productService   productService
	standsRepository standsRepository
}

// NewStandService returns new Stand Service.
func NewStandService(productService productService, standsRepository standsRepository) *Service {
	return &Service{
		productService:   productService,
		standsRepository: standsRepository,
	}
}

// GetActiveProductsWithData fetches information about products on stands are active.
func (s *Service) GetActiveProductsWithData() ([]WithProductDetailsDTO, error) {
	var existingStands []WithProductDetailsDTO

	data, err := s.standsRepository.GetAllActiveStands()
	if err != nil {
		return nil, err
	}

	for i := range data {
		stand := WithProductDetailsDTO{}
		stand.ID, err = strconv.Atoi(data[i].ID)
		if err != nil {
			fmt.Println(err)
			return []WithProductDetailsDTO{}, err
		}

		stand.Product, err = s.productService.GetProductDetailsByID(data[i].ProductID)
		if err != nil {
			return nil, err
		}

		existingStands = append(existingStands, stand)

	}

	return existingStands, nil
}

// GetAllProductsWithData fetches information about all products on stands
func (s *Service) GetAllProductsWithData() ([]WithProductDetailsDTO, error) {
	var existingStands []WithProductDetailsDTO

	data, err := s.standsRepository.GetAllActiveStandsMap()
	if err != nil {
		return nil, err
	}

	for ID, productID := range data {
		stand := WithProductDetailsDTO{}
		stand.ID = ID

		stand.Product, err = s.productService.GetProductDetailsByID(productID)
		if err != nil {
			return nil, err
		}

		existingStands = append(existingStands, stand)
	}

	return existingStands, nil
}
