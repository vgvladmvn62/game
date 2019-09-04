package server

import (
	"net/http"
)

func (s *Server) turnOffLightsHandler(w http.ResponseWriter, r *http.Request) {
	err := s.hardwareService.TurnOffLights()
	if err != nil {
		_ = NewAPIError(err.Error(), http.StatusInternalServerError).Send(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
