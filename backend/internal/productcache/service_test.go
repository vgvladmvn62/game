package productcache_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/ec"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/productcache"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/productcache/mocks"
)

func TestGetProductDetailsByIDFromRepo(t *testing.T) {
	// Given

	fixed := fixedProductDTO()

	productService := &mocks.ProductService{}
	productService.On("BuildProductImage", &fixed).Return(&fixed, nil)

	productsRepository := &mocks.ProductsRepository{}
	productsRepository.On("GetProductByID", "1").Return(fixed, nil)

	standsRepository := &mocks.StandsRepository{}

	productCacheService := productcache.NewProductCacheService(productService, standsRepository, productsRepository)

	// When

	output, _ := productCacheService.GetProductDetailsByID("1")

	// Then

	assert.Equal(t, output, fixed)
}

func TestGetProductDetailsByIDStraightFromEC(t *testing.T) {
	// Given

	fixed := fixedProductDTO()

	productService := &mocks.ProductService{}
	productService.On("GetProductDetailsByID", "1").Return(fixedProductDTO(), nil)
	productService.On("BuildProductImage", &fixed).Return(&fixed, nil)

	productsRepository := &mocks.ProductsRepository{}
	productsRepository.On("GetProductByID", "1").Return(ec.ProductDTO{},
		errors.New("Unable to fetch data from database"))

	standsRepository := &mocks.StandsRepository{}

	productCacheService := productcache.NewProductCacheService(productService, standsRepository, productsRepository)

	// When

	output, _ := productCacheService.GetProductDetailsByID("1")

	// Then

	assert.Equal(t, output, fixed)
}

func fixedProductDTO() ec.ProductDTO {
	return ec.ProductDTO{
		ID:          "1",
		Description: "First product to test",
		Price: ec.Price{
			Currency:       "EUR",
			Value:          42000,
			FormattedValue: "420.0eur",
		},
		Name:  "Product 1",
		Image: "http://url.to/img",
	}
}
