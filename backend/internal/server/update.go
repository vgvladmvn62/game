package server

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (s *Server) updateGETHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	err = s.productCacheService.ForceUpdateProducts()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) updateByIDGETHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	err := s.productCacheService.ForceUpdateProductByID(id)
	if err != nil {
		_ = NewAPIError(err.Error(), http.StatusNotFound).Send(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
