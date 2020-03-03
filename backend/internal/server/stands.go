package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kyma-incubator/bullseye-showcase/backend/internal/db/stands"
)

// StandsDTO stores list of stands for unmarshaling
type StandsDTO struct {
	Stands []stands.StandDTO `json:"stands"`
}

func (s *Server) standsPOSTHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		_ = NewAPIError("No body", http.StatusBadRequest).Send(w)
		return
	}

	defer func() { _ = r.Body.Close() }()

	standsConfig := StandsDTO{}

	err := json.NewDecoder(r.Body).Decode(&standsConfig)
	if err != nil {
		_ = NewAPIError(err.Error(), http.StatusBadRequest).Send(w)
		return
	}

	err = s.standsRepository.DropTable()
	if err != nil {
		log.Println("Not dropping shelves: ", err)
	}

	err = s.standsRepository.CreateTable()
	if err != nil {
		_ = NewAPIError(err.Error(), http.StatusInternalServerError).Send(w)
		return
	}

	for _, stand := range standsConfig.Stands {
		err = s.standsRepository.AddStand(stand.ID, stand.ProductID, stand.Active)
		if err != nil {
			_ = NewAPIError(err.Error(), http.StatusInternalServerError).Send(w)
			return
		}
	}

	err = s.productCacheService.ForceUpdateProducts()
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) standsGETHandler(w http.ResponseWriter, r *http.Request) {
	config := StandsDTO{}

	var err error

	config.Stands, err = s.standsRepository.GetAllStands()
	if err != nil {
		_ = NewAPIError(err.Error(), http.StatusInternalServerError).Send(w)
		return
	}

	body, err := json.Marshal(config)
	if err != nil {
		_ = NewAPIError(err.Error(), http.StatusInternalServerError).Send(w)
		return
	}

	_, _ = w.Write(body)
}
