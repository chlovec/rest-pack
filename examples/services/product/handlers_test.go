package product

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/chlovec/rest-pack/examples/config"
	"github.com/chlovec/rest-pack/examples/services/mocks"
	"github.com/chlovec/rest-pack/examples/types"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestListProductsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockProductStore(ctrl)
	handler := NewHandler(log.Default(), mockStore)

	t.Run("should list products", func(t *testing.T) {
		expectedProducts := []*types.Product{&prodA, &prodB}
		mockStore.EXPECT().ListProducts(1000, 0).Return(expectedProducts, nil)

		req, err := http.NewRequest(http.MethodGet, "/products", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products", handler.ListProducts).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var actualProducts []*types.Product
		err = json.Unmarshal(rr.Body.Bytes(), &actualProducts)
		assert.NoError(t, err)
		assert.Equal(t, expectedProducts, actualProducts)
	})

	t.Run("should return empty list", func(t *testing.T) {
		expectedProducts := []*types.Product{}
		mockStore.EXPECT().ListProducts(1000, 0).Return(expectedProducts, nil)

		req, err := http.NewRequest(http.MethodGet, "/products", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products", handler.ListProducts).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var actualProducts []*types.Product
		err = json.Unmarshal(rr.Body.Bytes(), &actualProducts)
		assert.NoError(t, err)
		assert.Equal(t, expectedProducts, actualProducts)
	})

	t.Run("should handle page size and number", func(t *testing.T) {
		expectedProducts := []*types.Product{&prodB}
		mockStore.EXPECT().ListProducts(100, 800).Return(expectedProducts, nil)

		req, err := http.NewRequest(http.MethodGet, "/products?pagesize=100&pagenumber=9", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products", handler.ListProducts).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var actualProducts []*types.Product
		err = json.Unmarshal(rr.Body.Bytes(), &actualProducts)
		assert.NoError(t, err)
		assert.Equal(t, expectedProducts, actualProducts)
	})

	t.Run("should return internal server error", func(t *testing.T) {
		mockStore.EXPECT().ListProducts(100, 800).Return(nil, errors.New(DbError))

		req, err := http.NewRequest(http.MethodGet, "/products?pagesize=100&pagenumber=9", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products", handler.ListProducts).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		expectedResponse := `{
			"error": "Internal Server Error"
		}`
		assert.JSONEq(t, expectedResponse, rr.Body.String())
	})
}

func TestGetProductHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockProductStore(ctrl)
	handler := NewHandler(log.Default(), mockStore)

	t.Run("should return product", func(t *testing.T) {
		mockStore.EXPECT().GetProduct(1).Return(&prodA, nil)

		req, err := http.NewRequest(http.MethodGet, "/products/1", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products/{id}", handler.GetProduct).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var actualProducts *types.Product
		err = json.Unmarshal(rr.Body.Bytes(), &actualProducts)
		assert.NoError(t, err)
		assert.Equal(t, &prodA, actualProducts)
	})

	t.Run("should return not found", func(t *testing.T) {
		mockStore.EXPECT().GetProduct(1).Return(nil, nil)

		req, err := http.NewRequest(http.MethodGet, "/products/1", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products/{id}", handler.GetProduct).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		expectedResponse := `{
			"error": "Not Found"
		}`
		assert.JSONEq(t, expectedResponse, rr.Body.String())
	})

	t.Run("should internal server error", func(t *testing.T) {
		mockStore.EXPECT().GetProduct(1).Return(nil, errors.New(DbError))

		req, err := http.NewRequest(http.MethodGet, "/products/1", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products/{id}", handler.GetProduct).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		expectedResponse := `{
			"error": "Internal Server Error"
		}`
		assert.JSONEq(t, expectedResponse, rr.Body.String())
	})

	t.Run("should return bad request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/products/1AbC", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products/{id}", handler.GetProduct).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		expectedResponse := `{
			"error": "Bad Request"
		}`
		assert.JSONEq(t, expectedResponse, rr.Body.String())
	})
}

func TestCreateProductHandler(t *testing.T) {
	// Set up test environment variables
	os.Setenv("BASE_URL", "http://example.com")
	os.Setenv("PATH_PREFIX", "/api/v1")
	config.InitConfig()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockProductStore(ctrl)
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
