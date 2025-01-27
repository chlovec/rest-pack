package product

import (
	"bytes"
	"encoding/json"
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
}
