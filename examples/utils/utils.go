package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	defer r.Body.Close() // Always close the body to avoid resource leaks
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	return nil
}

func WriteErrorJSON(w http.ResponseWriter, status int, err error, details any) {
	errorResponse := map[string]interface{}{
		"error": err.Error(),
	}

	if details != nil {
		errorResponse["details"] = details
	}

	// Marshal the response to JSON
	responseJSON, marshalErr := json.Marshal(errorResponse)
	if marshalErr != nil {
		// Fallback to http.Error if JSON marshaling fails
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Use http.Error to send the response as JSON
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, string(responseJSON), status)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}