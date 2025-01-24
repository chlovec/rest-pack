package api

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegisterRoute(t *testing.T) {
	t.Run("should register route with path prefix", func(t *testing.T) {
		// Create a new APIServer
		logger, _ := initLog()
		server := NewAPIServer(":8080", "/api", logger)

		// Define a test handler
		testHandler := func(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello, World!"))
		}

		// Register the test route
		server.RegisterRoute("/test", testHandler, http.MethodGet)

		// Create a test server
		testServer := httptest.NewServer(server.apiRouter)
		defer testServer.Close()

		// Send a request to the registered route
		resp, err := http.Get(testServer.URL + "/api/test")
		assert.NoError(t, err, "Failed to send request: %v", err)
		defer resp.Body.Close()

		// Read the response body
		body := make([]byte, resp.ContentLength)
		resp.Body.Read(body)

		// Assert the response status code
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status %d, got %d", http.StatusOK, resp.StatusCode)

		// Assert the response body
		expectedBody := "Hello, World!"
		actualBody := string(body)
		assert.Equal(t, expectedBody, actualBody, "Expected body %q, got %q", expectedBody, body)
	})

	t.Run("should register route with empty path prefix", func(t *testing.T) {
		// Create a new APIServer
		logger, _ := initLog()
		server := NewAPIServer(":8080", "", logger)

		// Define a test handler
		testHandler := func(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("Hello, World!"))
		}

		// Register the test route
		server.RegisterRoute("/test", testHandler, http.MethodPost)

		// Create a test server
		testServer := httptest.NewServer(server.apiRouter)
		defer testServer.Close()

		// Send a request to the registered route
		resp, err := http.Post(testServer.URL+"/test", "application/json", nil)
		assert.NoError(t, err, "Failed to send request: %v", err)
		defer resp.Body.Close()

		// Read the response body
		body := make([]byte, resp.ContentLength)
		resp.Body.Read(body)

		// Assert the response status code
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status %d, got %d", http.StatusCreated, resp.StatusCode)

		// Assert the response body
		expectedBody := "Hello, World!"
		actualBody := string(body)
		assert.Equal(t, expectedBody, actualBody, "Expected body %q, got %q", expectedBody, body)
	})

	t.Run("should handle request with invalid http method", func(t *testing.T) {
		// Create a new APIServer
		logger, _ := initLog()
		server := NewAPIServer(":8080", "", logger)

		// Define a test handler
		testHandler := func(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello, World!"))
		}

		// Register the test route
		server.RegisterRoute("/test", testHandler, http.MethodGet)

		// Create a test server
		testServer := httptest.NewServer(server.apiRouter)
		defer testServer.Close()

		// Send a request to the registered route
		resp, err := http.Post(testServer.URL+"/test", "application/json", nil)
		assert.NoError(t, err, "Failed to send request: %v", err)
		defer resp.Body.Close()

		// Read the response body
		body := make([]byte, resp.ContentLength)
		resp.Body.Read(body)

		// Assert the response status code
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode, "Expected status %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("should handle request to invalid route", func(t *testing.T) {
		// Create a new APIServer
		logger, _ := initLog()
		server := NewAPIServer(":8080", "", logger)

		// Define a test handler
		testHandler := func(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello, World!"))
		}

		// Register the test route
		server.RegisterRoute("/test", testHandler, http.MethodGet)

		// Create a test server
		testServer := httptest.NewServer(server.apiRouter)
		defer testServer.Close()

		// Send a request to the registered route
		resp, err := http.Get(testServer.URL + "/api/test")
		assert.NoError(t, err, "Failed to send request: %v", err)
		defer resp.Body.Close()

		// Read the response body
		body := make([]byte, resp.ContentLength)
		resp.Body.Read(body)

		// Assert the response status code
		assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Expected status %d, got %d", http.StatusNotFound, resp.StatusCode)

		// Assert the response body
		expectedBody := "404 page not found\n"
		actualBody := string(body)
		assert.Equal(t, expectedBody, actualBody, "Expected body %q, got %q", expectedBody, body)
	})

	t.Run("should fail to register route with empty path", func(t *testing.T) {
		// Create a new APIServer
		logger, logBuffer := initLog()
		server := NewAPIServer(":8080", "/api", logger)
		// Attempt to register a route with an empty path
		server.RegisterRoute("", func(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
			w.WriteHeader(http.StatusOK)
		}, http.MethodGet)

		// Assert that an appropriate log message was written
		logContent := logBuffer.String()
		expectedLog := "Cannot register a route with an empty path"
		assert.Contains(t, logContent, expectedLog, "Expected log message %q, but it was not found. Actual log message %q", expectedLog, logContent)
	})

	t.Run("should fail to register route with nil handler", func(t *testing.T) {
		// Create a new APIServer
		logger, logBuffer := initLog()
		server := NewAPIServer(":8080", "/api", logger)
		// Attempt to register a route with an empty path
		server.RegisterRoute("/test", nil, http.MethodGet)

		// Assert that an appropriate log message was written
		logContent := logBuffer.String()
		expectedLog := "Cannot register a route with a nil handler"
		assert.Contains(t, logContent, expectedLog, "Expected log message %q, but it was not found. Actual log message %q", expectedLog, logContent)
	})
}

func TestAPIServerStart(t *testing.T) {
	// Mock logger
	var logBuffer bytes.Buffer
	logger := log.New(&logBuffer, "", log.LstdFlags)

	// Create a new API server
	serverAddr := "127.0.0.1:0" // Use a random available port
	apiServer := NewAPIServer(serverAddr, "", logger)

	// Mock route to ensure the server is running
	apiServer.RegisterRoute("/health", func(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
		w.WriteHeader(http.StatusOK)
	}, "GET")

	// Channel to capture the Start function's return value
	errChan := make(chan error, 1)

	// Simulate sending an interrupt signal after a delay
	go func() {
		time.Sleep(1 * time.Second)
		p, _ := os.FindProcess(os.Getpid()) // Get the current process
		_ = p.Signal(os.Interrupt)         // Send interrupt signal
	}()

	// Run the Start function in a goroutine
	go func() {
		errChan <- apiServer.Start(2 * time.Second)
	}()

	// Wait for the server to shut down
	select {
	case err := <-errChan:
		if err != nil {
			t.Fatalf("Server failed to shut down gracefully: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Test timed out waiting for server to shut down")
	}

	// Check if the correct log messages were written
	if !bytes.Contains(logBuffer.Bytes(), []byte("Starting server on")) {
		t.Errorf("Expected log message 'Starting server on', but not found")
	}
	if !bytes.Contains(logBuffer.Bytes(), []byte("Shutting down gracefully...")) {
		t.Errorf("Expected log message 'Shutting down gracefully...', but not found")
	}
	if !bytes.Contains(logBuffer.Bytes(), []byte("Server stopped gracefully.")) {
		t.Errorf("Expected log message 'Server stopped gracefully.', but not found")
	}
}

func initLog() (*log.Logger, *bytes.Buffer) {
	logBuffer := &bytes.Buffer{}
	logger := log.New(logBuffer, "TEST: ", log.LstdFlags)
	return logger, logBuffer
}
