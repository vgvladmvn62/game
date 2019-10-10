package productcache

import (
	"fmt"
	"log"

	"github.com/kyma-incubator/bullseye-showcase/backend/internal/db/stands"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/ec"
)

// ProductService enables getting information about products.
type ProductService interface {
	GetProductDetailsByID(ID string) (ec.ProductDTO, error)
	GetProductJSONByID(ID string) (string, error)
	BuildProductImage(dto *ec.ProductDTO) *ec.ProductDTO
}

// StandsRepository enables getting information about stands.
type StandsRepository interface {
	GetAllStands() ([]stands.StandDTO, error)
}

// ProductsRepository enables managing products' data.
type ProductsRepository interface {
	CreateTable() error
	DropTable() error
	Exists() (bool, error)
	AddProduct(ID string, data string) error
	GetAllProducts() ([]ec.ProductDTO, error)
	GetProductByID(ID string) (ec.ProductDTO, error)
	UpdateProductDataByID(ID string, newData string) error
}

//go:generate mockery -name=ProductService -output=mocks -outpkg=mocks -case=underscore
//go:generate mockery -name=StandsRepository -output=mocks -outpkg=mocks -case=underscore
//go:generate mockery -name=ProductsRepository -output=mocks -outpkg=mocks -case=underscore

// Service is a Product Cache Service that fetches product information
// from database or from EC using Product Service.
type Service struct {
	productService     ProductService
	standsRepository   StandsRepository
	productsRepository ProductsRepository
}

// NewProductCacheService returns new instance of Product Cache Service.
func NewProductCacheService(productService ProductService, standsRepository StandsRepository, productsRepository ProductsRepository) *Service {
	return &Service{
		productService:     productService,
		standsRepository:   standsRepository,
		productsRepository: productsRepository,
	}
}

// GetProductDetailsByID returns information about the product searched by ID.
func (s *Service) GetProductDetailsByID(ID string) (ec.ProductDTO, error) {
	var product ec.ProductDTO
	var err error

	product, err = s.fetchProductData(ID)
	if err != nil {
		return ec.ProductDTO{}, err
	}

	s.productService.BuildProductImage(&product)

	if product.IsEmpty() {
		err = s.updateProductByID(ID)
		if err != nil {
			return ec.ProductDTO{}, err
		}

		product, err = s.fetchProductData(ID)
		if err != nil {
			return ec.ProductDTO{}, err
		}

		s.productService.BuildProductImage(&product)
	}

	return product, nil
}

// UpdateProducts updates information about the product.
func (s *Service) UpdateProducts() error {
	var err error

	err = s.createTableIfNotExists()
	if err != nil {
		log.Println("Error creating table in Product Cache")
		return err
	}

	existingStands, err := s.standsRepository.GetAllStands()
	if err != nil {
		log.Println("Error getting stands in Product Cache")
		return err
	}

	err = s.updateProductsFromStands(existingStands)
	if err != nil {
		log.Println("Could not force update products in Product Cache:")
		return err
	}

	return nil
}

// ForceUpdateProducts updates information about all products mapped to stands.
func (s *Service) ForceUpdateProducts() error {
	var err error

	err = s.createTableIfNotExists()
	if err != nil {
		log.Println("Error creating table in Product Cache")
		return err
	}

	existingStands, err := s.standsRepository.GetAllStands()
	if err != nil {
		log.Println("Error getting stands in Product Cache")
		return err
	}

	err = s.forceUpdateProductsFromStands(existingStands)
	if err != nil {
		log.Println("Could not force update products in Product Cache:", err)
	}

	return nil
}

// ForceUpdateProductByID updates information about the specific product.
func (s *Service) ForceUpdateProductByID(ID string) error {
	var err error

	err = s.createTableIfNotExists()
	if err != nil {
		return err
	}

	err = s.updateProductByID(ID)
	if err != nil {
		return err
	}

	return err
}

func (s *Service) fetchProductData(ID string) (ec.ProductDTO, error) {
	var product ec.ProductDTO
	var err error

	product, err = s.productsRepository.GetProductByID(ID)
	if err != nil {
		log.Println("Could not fetch data from table products: ", err)

		product, err = s.productService.GetProductDetailsByID(ID)
		if err != nil {
			log.Println("Could not fetch data from EC: ", err)
			return ec.ProductDTO{}, err
		}
	}

	return product, nil
}

func (s *Service) createTableIfNotExists() error {
	var exists bool
	var err error

	exists, err = s.productsRepository.Exists()
	if err != nil {
		return err
	}

	if !exists {
		err = s.productsRepository.CreateTable()
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func (s *Service) updateProductsFromStands(existingStands []stands.StandDTO) error {
	var productIsInDatabase bool
	var err error

	updated := false

	for i := range existingStands {
		productIsInDatabase, err = s.isProductInDatabase(existingStands[i].ProductID)
		if err != nil {
			return err
		}

		if !productIsInDatabase {
			err = s.updateProductByID(existingStands[i].ProductID)
			if err != nil {
				return err
			}
			updated = true
		}
	}

	if updated {
		fmt.Println("Updated products - updateProductsFromStands()")
	}

	return nil
}

func (s *Service) forceUpdateProductsFromStands(existingStands []stands.StandDTO) error {
	var err error

	for i := range existingStands {
		err = s.updateProductByID(existingStands[i].ProductID)
		if err != nil {
			return err
		}
	}

	fmt.Println("Force updated products")

	return nil
}

func (s *Service) isProductInDatabase(ID string) (bool, error) {
	product, err := s.productsRepository.GetProductByID(ID)
	if err != nil {
		return false, err
	}

	if product.IsEmpty() {
		return false, nil
	}

	return true, nil
}

func (s *Service) updateProductByID(ID string) error {
	var data string
	var productIsInDatabase bool
	var err error

	data, err = s.productService.GetProductJSONByID(ID)
	if err != nil {
		return err
	}

	productIsInDatabase, err = s.isProductInDatabase(ID)
	if err != nil {
		return err
	}

	if productIsInDatabase {
		err = s.productsRepository.UpdateProductDataByID(ID, data)
		if err != nil {
			return err
		}
	} else {
		err = s.productsRepository.AddProduct(ID, data)
		if err != nil {
			return err
		}
	}

	return nil
}
