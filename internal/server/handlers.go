package server

import (
	"encoding/json"
	"net/http"

	rmodel "github.com/dormitory-life/gateway/internal/server/request_models"
)

func (s *Server) pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func writeErrorResponse(w http.ResponseWriter, err error, code int, details ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := rmodel.ErrorResponse{
		Error:   err.Error(),
		Details: details,
	}

	_ = json.NewEncoder(w).Encode(response)
}
