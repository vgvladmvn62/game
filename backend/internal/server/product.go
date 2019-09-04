package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) productGETHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	product, err := s.productCacheService.GetProductDetailsByID(id)
	if err != nil {
		_ = NewAPIError(err.Error(), http.StatusNotFound).Send(w)
		return
	}

	body, err := json.Marshal(product)
	if err != nil {
		_ = NewAPIError(err.Error(), http.StatusInternalServerError).Send(w)
		return
	}

	_, _ = w.Write(body)
}
