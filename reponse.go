package controller

import (
	"encoding/json"
	m "uts/model"
	"net/http"
)

func sendErrorResponse(w http.ResponseWriter, message string) {
	var response m.ErrorResponse
	response.Status = 400
	response.Message = message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
