package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func GetValidationError(err error) map[string]string {
	validationErrors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		validationErrors[err.Field()] = getValidationMessage(err)
	}
	return validationErrors
}

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	defer r.Body.Close() // Always close the body to avoid resource leaks
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("missing request body")
		}
		return fmt.Errorf("invalid JSON: %w", err)
	}

	return nil
}

func WriteErrorJSON(w http.ResponseWriter, status int, err error, details any) {
	// Marshal the response to JSON
	responseJSON, marshalErr := createJsonMessage("error", err.Error(), "details", details)
	if marshalErr != nil {
		// Fallback to http.Error if JSON marshaling fails
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Use http.Error to send the response as JSON
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, string(responseJSON), status)
}

func WriteBadRequest(w http.ResponseWriter, errorMessage string, details any) {
	if errorMessage == "" {
		errorMessage = "Bad Request"
	}
	writeError(w, http.StatusBadRequest, errorMessage, details)
}

func WriteInternalServerError(w http.ResponseWriter, errorMessage string, details any) {
	if errorMessage == "" {
		errorMessage = "Internal Server Error"
	}
	writeError(w, http.StatusInternalServerError, errorMessage, details)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteNotFound(w http.ResponseWriter, errorMessage string, details any) {
	if errorMessage == "" {
		errorMessage = "Not Found"
	}
	writeError(w, http.StatusNotFound, errorMessage, details)
}

func WriteLog(logger *log.Logger, category string, details any) {
	log, marshalErr := createJsonMessage("category", category, "message", details)
	if marshalErr != nil {
		logger.Printf("log error: %v", marshalErr)
		return
	}

	logger.Println(log)
}

func createJsonMessage(messageKey string, message string, detailsKey string, details any) ([]byte, error) {
	errorResponse := map[string]interface{}{
		messageKey: message,
	}

	if details != nil {
		errorResponse[detailsKey] = details
	}

	// Marshal the response to JSON
	msg, marshalErr := json.Marshal(errorResponse)
	if marshalErr != nil {
		return nil, marshalErr
	}

	return msg, nil
}

func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("'%s' is required", fe.Field())
	case "email":
		return fmt.Sprintf("'%s' is invalid", fe.Field())
	case "min":
		return fmt.Sprintf("'%s' is less than the required minimum length", fe.Field())
	case "max":
		return fmt.Sprintf("'%s' is greater than the required maximum length", fe.Field())
	}
	return fmt.Sprintf("%v", fe.Error())
}

func writeError(w http.ResponseWriter, httpErrorStatusCode int, errorMessage string, details any) {
	err := errors.New(errorMessage)
	WriteErrorJSON(w, httpErrorStatusCode, err, details)
}
