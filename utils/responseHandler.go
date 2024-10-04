package utils

import (
	"net/http"
	"encoding/json"
)

// for HTTP Responses
func HandleResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
