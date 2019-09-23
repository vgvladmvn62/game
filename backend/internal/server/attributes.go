package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kyma-incubator/bullseye-showcase/backend/internal/db/attributes"
)

// ProductWithAttributesDTO contains information about product's ID and
// it's assigned attributes.
type ProductWithAttributesDTO struct {
	ID         string                 `json:"id"`
	Attributes []attributes.Attribute `json:"attributes"`
}

// ProductsWithAttributesDTO stores collection of ProductWithAttributesDTO.
type ProductsWithAttributesDTO struct {
	Products []ProductWithAttributesDTO `json:"mapping"`
}

func (s *Server) attributesPOSTHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		_ = NewAPIError("No body", http.StatusBadRequest).Send(w)
		return
	}

	defer func() { _ = r.Body.Close() }()

	config := ProductsWithAttributesDTO{}

	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		_ = NewAPIError(err.Error(), http.StatusBadRequest).Send(w)
		return
	}

	err = s.attributesRepository.DropTable()
	if err != nil {
		log.Println("Unable to drop attributes: ", err)
	}

	err = s.attributesRepository.CreateTable()
	if err != nil {
		_ = NewAPIError(err.Error(), http.StatusInternalServerError).Send(w)
		return
	}

	for _, productWithAttributes := range config.Products {
		err = s.attributesRepository.AddAttributes(productWithAttributes.ID, productWithAttributes.Attributes)
		if err != nil {
			_ = NewAPIError(err.Error(), http.StatusInternalServerError).Send(w)
			return
		}
	}

	err = s.productCacheService.ForceUpdateProducts()
	if err != nil {
		log.Println("Could not force update products info: ", err)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) attributesGETHandler(w http.ResponseWriter, r *http.Request) {
	config := ProductsWithAttributesDTO{}

	var err error

	activeStands, err := s.standsRepository.GetAllStands()

	if err != nil {
		_ = NewAPIError(err.Error(), 500).Send(w)
		return
	}

	for _, stand := range activeStands {
		productAttributes, err := s.attributesRepository.GetAttributes(stand.ProductID)
		if err != nil {
			_ = NewAPIError(err.Error(), 500).Send(w)
			return
		}

		if productAttributes == nil {
			productAttributes = []attributes.Attribute{}
		}

		config.Products = append(config.Products, ProductWithAttributesDTO{ID: stand.ProductID, Attributes: productAttributes})
	}

	body, err := json.Marshal(config)
	if err != nil {
		_ = NewAPIError(err.Error(), http.StatusInternalServerError).Send(w)
		return
	}

	_, _ = w.Write(body)
}
