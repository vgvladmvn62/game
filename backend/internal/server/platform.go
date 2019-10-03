package server

import (
	"log"
	"net/http"
)

func (s *Server) turnOffLightsHandler(w http.ResponseWriter, r *http.Request) {
	err := s.hardwareService.TurnOffLights()
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
}
