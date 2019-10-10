package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/kyma-incubator/bullseye-showcase/backend/internal/db/attributes"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/matching"
)

// AnswersDTO contains answers sent to the server.
type AnswersDTO struct {
	Tags []string `json:"answers"`
}

// ResultsDTO contains information about products matched
// by the service and index highlighted product.
type ResultsDTO struct {
	Matched            []matching.MatchedProductDTO `json:"matched"`
	HighlightedProduct int                          `json:"highlightedProduct"`
}

func setupResponse(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func (s *Server) resultsPOSTHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Body == nil {
		_ = NewAPIError("No body", http.StatusBadRequest).Send(w)
		return

	}

	defer func() { _ = r.Body.Close() }()

	answers := AnswersDTO{}

	err := json.NewDecoder(r.Body).Decode(&answers)
	if err != nil {
		_ = NewAPIError(err.Error(), http.StatusBadRequest).Send(w)
		return
	}

	var toMatch []attributes.Attribute

	for i := range answers.Tags {
		temp := answers.Tags[i]
		toMatch = append(toMatch, attributes.Attribute(temp))
	}

	matched, err := s.matchingService.MatchProducts(toMatch)
	if err != nil {
		_ = NewAPIError(err.Error(), http.StatusInternalServerError).Send(w)
		return
	}

	results := ResultsDTO{
		Matched:            matched,
		HighlightedProduct: 0,
	}

	body, err := json.Marshal(results)
	if err != nil {
		_ = NewAPIError(err.Error(), http.StatusInternalServerError).Send(w)
		return
	}

	_, _ = w.Write(body)

	err = HighlightMatchedProducts(s.hardwareService, matched)
	if err != nil {
		_ = NewAPIError(err.Error(), http.StatusInternalServerError).Send(w)
		return
	}
}

type hardwarer interface {
	TurnOffLights() error
	TurnOnGreenLight(platformID byte) error
}

func HighlightMatchedProducts(hw hardwarer, matches []matching.MatchedProductDTO) error {
	err := hw.TurnOffLights()
	if err != nil {
		log.Println(err)
	}

	for _, id := range selectPlatformsToHighlightFromSortedMatches(matches) {
		idAsInt, err := strconv.Atoi(id)
		if err != nil {
			return err
		}

		err = hw.TurnOnGreenLight(byte(idAsInt))
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func selectPlatformsToHighlightFromSortedMatches(matches []matching.MatchedProductDTO) []string {
	// function requires matches to be sorted (which they are and there is test for it in the matchingService)
	var platforms []string

	if len(matches) == 0 {
		return platforms
	}

	highest := matches[0].Score

	if highest == 0 {
		return platforms
	}

	for _, match := range matches {
		if match.Score == highest {
			platforms = append(platforms, match.StandID)
		} else {
			break
		}
	}

	return platforms
}
