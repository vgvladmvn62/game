package ec

import (
	"encoding/json"
		"io/ioutil"
	"net/http"
		"fmt"
)

// ProductService provides methods operating on products.
// Uses HTTPClient for HTTP communication.
type ProductService struct {
	Cfg  *Config
	Doer httpDoer
}

// NewProductService returns new product service using passed config.
func NewProductService(config *Config) *ProductService {
	return &ProductService{
		Cfg:  config,
		Doer: NewHTTPClient(config.HTTP.Client.Timeout.Seconds),
	}
}

// GetProductDetailsByID returns product details from EC by product ID.
// If http request fails function returns 404 when not found, otherwise returns RequestFailedError.
func (p *ProductService) GetProductDetailsByID(ID string) (ProductDTO, error) {
	product := ProductDTO{}

	url := fmt.Sprintf("%s/%s/products/%s?fields=FULL",
		p.Cfg.Host.API,
		p.Cfg.Products.Site,
		ID,
	)

	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return ProductDTO{}, err
	}

	response, err := p.Doer.Do(request)

	if err == nil {
		if response.StatusCode >= 400 && response.StatusCode < 500 {
			return ProductDTO{}, RequestFailedError
		}
	} else {
		return ProductDTO{}, err
	}

	defer func() { _ = response.Body.Close() }()

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return ProductDTO{}, ReadDataFailedError
	}

	err = json.Unmarshal(data, &product)

	if err != nil {
		return ProductDTO{}, UnmarshalDataFailedError
	}

	p.BuildProductImage(&product)

	return product, nil
}

// GetProductJSONByID returns product details from EC by product ID in JSON format as string.
// If http request fails function returns 404 when not found, otherwise returns RequestFailedError.
func (p *ProductService) GetProductJSONByID(ID string) (string, error) {
	url := fmt.Sprintf("%s/%s/products/%s?fields=FULL",
		p.Cfg.Host.API,
		p.Cfg.Products.Site,
		ID,
	)

	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return "", err
	}

	response, err := p.Doer.Do(request)

	if err == nil {
		if response.StatusCode >= 400 && response.StatusCode < 500 {
			return "", RequestFailedError
		}
	} else {
		return "", err
	}

	defer func() { _ = response.Body.Close() }()

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", ReadDataFailedError
	}

	dataText := string(data)

	return dataText, nil
}

// BuildProductImage fill's products Image field with server host.
// It allows getting a working URL.
func (p *ProductService) BuildProductImage(dto *ProductDTO) *ProductDTO {
	dto.Image = p.Cfg.Host.Static + dto.Image

	return dto
}
