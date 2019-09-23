package server_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kyma-incubator/bullseye-showcase/backend/internal/server/mocks"

	"github.com/go-http-utils/logger"
	"github.com/stretchr/testify/assert"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/ec"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/server"
)

func prepareCorrectConfigForTests() *server.Config {
	var config = server.Config{}
	config.Port = 8080
	config.IP = "0.0.0.0"
	config.Logger.Type = logger.Type(3)
	return &config
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

func fixedAPIError() server.APIError {
	return server.NewAPIError("An OCC error occured", http.StatusNotFound)
}

func TestServer(t *testing.T) {
	cache := &mocks.ProductCacheService{}
	cache.On("GetProductDetailsByID", "1").Return(fixedProductDTO(), nil)
	cache.On("GetProductDetailsByID", "10").Return(ec.ProductDTO{}, errors.New("An OCC error occured"))

	srv := server.NewServer(prepareCorrectConfigForTests(), cache, nil, nil, nil, nil, nil, nil)

	type serverTest struct {
		title string
		path  string
		code  int
		body  string
	}

	tests := []serverTest{
		{title: "Good request",
			path: "/product/1",
			code: http.StatusOK,
			body: `{
					"id": "1",
					"description": "First product to test",
					"price": {
						"currency": "EUR",
						"value": 42000,
						"formatted_value": "420.0eur"
					},
					"name": "Product 1",
					"image": "http://url.to/img"
				}`,
		},
		{title: "Non existing product",
			path: "/product/10",
			code: http.StatusNotFound,
			body: `{
					"message": "An OCC error occured",
					"code": 404
				}`,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, err := http.NewRequest("GET", test.path, nil)
			srv.Router.ServeHTTP(rr, req)

			assert.NoError(t, err)
			assert.JSONEq(t, test.body, rr.Body.String())
			assert.Equal(t, test.code, rr.Code)
		})
	}
}
