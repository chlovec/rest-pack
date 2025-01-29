package utils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestPayload struct {
	Name  string `json:"name" validate:"required,min=3,max=10"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"gte=5"`
}

func TestParseJSON(t *testing.T) {
	t.Run("Valid JSON", func(t *testing.T) {
		payload := TestPayload{}
		body := `{"name": "John", "email": "john@example.com"}`
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		err := ParseJSON(r, &payload)
		assert.NoError(t, err)
		assert.Equal(t, "John", payload.Name)
		assert.Equal(t, "john@example.com", payload.Email)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		payload := TestPayload{}
		body := `{name: "John", "email": "john@example.com"}` // Invalid JSON
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		err := ParseJSON(r, &payload)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid JSON")
	})

	t.Run("Missing Request Body", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/", nil)
		r.Body = nil
		err := ParseJSON(r, &TestPayload{})
		assert.Error(t, err)
		assert.Equal(t, "missing request body", err.Error())
	})

	t.Run("Missing Body", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/", nil)
		err := ParseJSON(r, &TestPayload{})
		assert.Error(t, err)
		assert.Equal(t, "missing request body", err.Error())
	})
}

func TestWriteErrorJSON(t *testing.T) {
	t.Run("Test Error", func(t *testing.T) {
		rw := httptest.NewRecorder()
		err := errors.New("test error")
		details := map[string]string{"field": "error detail"}

		WriteErrorJSON(rw, http.StatusBadRequest, err, details)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		assert.JSONEq(t, `{"error": "test error", "details": {"field": "error detail"}}`, rw.Body.String())
	})

	t.Run("JSON Marshal Error", func(t *testing.T) {
		rw := httptest.NewRecorder()
		err := errors.New("test error")

		// Pass an unsupported type (function) as `details` to force Marshal to fail
		WriteErrorJSON(rw, http.StatusInternalServerError, err, func() {})

		// The response should fall back to the hardcoded internal server error response
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.JSONEq(t, `{"error": "internal server error"}`, rw.Body.String())
	})

}

func TestWriteJSON(t *testing.T) {
	rw := httptest.NewRecorder()
	data := map[string]string{"message": "success"}

	err := WriteJSON(rw, http.StatusOK, data)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rw.Code)
	assert.JSONEq(t, `{"message": "success"}`, rw.Body.String())
}

func TestGetValidationError(t *testing.T) {
	t.Run("Required", func(t *testing.T) {
		testPayload := TestPayload{}
		err := Validate.Struct(testPayload)
		assert.Error(t, err)

		validationErrors := GetValidationError(err)
		assert.Contains(t, validationErrors, "Name")
		assert.Contains(t, validationErrors, "Email")
		assert.Equal(t, "'Name' is required", validationErrors["Name"])
		assert.Equal(t, "'Email' is required", validationErrors["Email"])
	})

	t.Run("Invalid", func(t *testing.T) {
		testPayload := TestPayload{Name: "John", Email: "invalid-email"}
		err := Validate.Struct(testPayload)
		assert.Error(t, err)

		msg := GetValidationError(err)
		assert.Contains(t, msg, "Email")
		assert.Equal(t, "'Email' is invalid", msg["Email"])
	})

	t.Run("Minimum Length", func(t *testing.T) {
		testPayload := TestPayload{Name: "Jo", Email: "alice.doe@example.com"}
		err := Validate.Struct(testPayload)
		assert.Error(t, err)

		msg := GetValidationError(err)
		assert.Contains(t, msg, "Name")
		assert.Equal(t, "'Name' is less than the required minimum length", msg["Name"])
	})

	t.Run("Maximum Length", func(t *testing.T) {
		testPayload := TestPayload{Name: "TestJohnDoe", Email: "alice.doe@example.com"}
		err := Validate.Struct(testPayload)
		assert.Error(t, err)

		msg := GetValidationError(err)
		assert.Contains(t, msg, "Name")
		assert.Equal(t, "'Name' is greater than the required maximum length", msg["Name"])
	})

	t.Run("Unknown Error", func(t *testing.T) {
		testPayload := TestPayload{Name: "John", Email: "john.doe@example.com", Age: 4}
		err := Validate.Struct(testPayload)
		assert.Error(t, err)

		msg := GetValidationError(err)
		assert.Contains(t, msg, "Age")
		assert.Equal(t, "Key: 'TestPayload.Age' Error:Field validation for 'Age' failed on the 'gte' tag", msg["Age"])
	})
}
