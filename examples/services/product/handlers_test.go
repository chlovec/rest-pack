package product

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chlovec/rest-pack/examples/services/mocks"
	"github.com/chlovec/rest-pack/examples/types"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCreateProductHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := prod_mocks.NewMockProductStore(ctrl)
	handler := NewHandler(log.Default(), mockStore)

	t.Run("should create new product", func(t *testing.T) {
		product := types.CreateProductPayload{
			Name:        "Test Product",
			Description: "New product for testing",
			ImageUrl:    "test/image-url",
			Price:       22.45,
			Quantity:    20,
		}

		mockStore.EXPECT().CreateProduct(product).Return(int64(1), nil)

		// Create http request
		body, _ := json.Marshal(product)
		req, err := http.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))

		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products", handler.CreateProduct).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		expectedResponse := `{
			"id": 1,
			"message": "Product created successfully",
			"url": "http://example.com/api/v1/products/1"
		}`
		assert.JSONEq(t, expectedResponse, rr.Body.String())
		assert.Equal(t, "http://example.com/api/v1/products/1", rr.Header().Get("Location"))
	})

	t.Run("should fail if DB error", func(t *testing.T) {
		product := types.CreateProductPayload{
			Name:        "Test Product",
			Description: "New product for testing",
			ImageUrl:    "test/image-url",
			Price:       22.45,
			Quantity:    20,
		}

		mockStore.EXPECT().CreateProduct(product).Return(int64(0), errors.New("DB error"))

		// Create http request
		body, _ := json.Marshal(product)
		req, err := http.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))

		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products", handler.CreateProduct).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		expectedResponse := `{
			"error": "Internal Server Error"
		}`
		assert.JSONEq(t, expectedResponse, rr.Body.String())
	})

	t.Run("should fail request has no payload", func(t *testing.T) {
		// Create http request
		req, err := http.NewRequest(http.MethodPost, "/products", nil)

		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products", handler.CreateProduct).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		expectedResponse := `{
			"error": "missing request body"
		}`
		assert.JSONEq(t, expectedResponse, rr.Body.String())
	})

	t.Run("should fail if payload is invalid", func(t *testing.T) {
		product := types.CreateProductPayload{}

		// Create http request
		body, _ := json.Marshal(product)
		req, err := http.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))

		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products", handler.CreateProduct).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		expectedResponse := `{
			"details": {
				"Name": "'Name' is required",
				"Price": "'Price' is required", 
				"Quantity": "'Quantity' is required"
			}, 
			"error":"Validation Error"
		}`
		assert.JSONEq(t, expectedResponse, rr.Body.String())
	})
}
