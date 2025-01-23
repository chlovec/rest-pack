package api

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIServer(t *testing.T) {
	// Set up a logger
	var logBuffer bytes.Buffer
	logger := log.New(&logBuffer, "", log.LstdFlags)

	// Create a new APIServer instance
	server := NewAPIServer(":8080", "/api", logger)

	t.Run("should handle request if route is valid", func(t *testing.T) {
		// Define a simple handler
		testHandler := func(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Test route hit"))
		}

		// Register the route
		server.RegisterRoute("/test", testHandler, "GET")

		// Create a test request
		req, err := http.NewRequest("GET", "/api/test", nil)
		assert.NoError(t, err)

		// Create a test response recorder
		rec := httptest.NewRecorder()

		// Serve the request
		server.apiRouter.ServeHTTP(rec, req)

		// Check the response
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, rec.Body.String(), "Test route hit")
	})

	t.Run("should fail if path is empty", func(t *testing.T) {
		// Define a simple handler
		testHandler := func(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Test route hit"))
		}

		// Try to register a route with an empty path
		server.RegisterRoute("", testHandler, "GET")

		// Check the log output for the expected error message
		expectedLogMessage := "Cannot register a route with an empty path"
		assert.Contains(t, logBuffer.String(), expectedLogMessage)
	})

	t.Run("should fail if handler is nil", func(t *testing.T) {
		// Try to register a route with an empty path
		server.RegisterRoute("/testPath", nil, "POST")

		// Check the log output for the expected error message
		expectedLogMessage := "Cannot register a route with a nil handler"
		assert.Contains(t, logBuffer.String(), expectedLogMessage)
	})
}