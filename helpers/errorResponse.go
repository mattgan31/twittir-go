package helpers

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Status  string            `json:"status"`
	Message map[string]string `json:"message"`
}

func sendErrorResponse(w http.ResponseWriter, status int, errors map[string]string) {
	errorResponse := ErrorResponse{
		Status:  "BAD_REQUEST",
		Message: errors,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		// Handle JSON encoding error if necessary
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
