package utils

import (
	"encoding/json"
	"net/http"
)

// WriteJSONResponse - writes json response
func WriteJSONResponse(w http.ResponseWriter, obj interface{}) {
	// Convert person struct to JSON
	jsonData, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set JSON response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write JSON data to response body
	w.Write(jsonData)
}
