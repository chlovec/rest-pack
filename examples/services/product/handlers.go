package product

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/chlovec/rest-pack/examples/config"
	"github.com/chlovec/rest-pack/examples/types"
	"github.com/chlovec/rest-pack/utils"
)

type Handler struct {
	logger *log.Logger
	store types.ProductStore
}

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price float64 `json:"price"`
	Qty   int     `json:"qty"`
}

func NewHandler(logger *log.Logger, store types.ProductStore) *Handler {
	return &Handler{
		logger: logger,
		store: store,
	}
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// Default values
	const defaultPageSize = 1000
	const defaultPageNum = 0

	query := r.URL.Query()

	// Parse page size
	pageSize, err := strconv.Atoi(query.Get("pageSize"))
	if err != nil || pageSize <= 0 {
		pageSize = defaultPageSize
	}

	// Parse page number
	pageNum, err := strconv.Atoi(query.Get("pageNumber"))
	if err != nil || pageNum < 1 {
		pageNum = defaultPageNum
	} else {
		pageNum-- // Convert to zero-based index
	}

	// Fetch products
	products, err := h.store.ListProducts(pageSize, pageNum)
	if err != nil {
		utils.WriteInternalServerError(w, "", nil)
		return
	}

	// Send response
	utils.WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var product types.CreateProductPayload
	if err := utils.ParseJSON(r, &product); err != nil {
		utils.WriteBadRequest(w, "missing request body", nil)
		return
	}

	// Validate payload
	if err := utils.Validate.Struct(product); err != nil {
		details := utils.GetValidationError(err)
		utils.WriteBadRequest(w, "Validation Error", details)
		return
	}

	// Create product
	productID, err := h.store.CreateProduct(product);
	if err != nil {
		utils.WriteInternalServerError(w, "", nil)
		return
	}

	// Write response
	newProductURL := fmt.Sprintf("%s%s/products/%d", config.Envs.BaseUrl, config.Envs.PathPrefix, productID)
	w.Header().Set("Location", newProductURL)

	response := map[string]interface{}{
		"message": "Product created successfully",
		"id":      productID,
		"url":     newProductURL,
	}
	utils.WriteJSON(w, http.StatusCreated, response)
}

// func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
// 	// Parse request body
// 	var product types.UpdateProductPayload
// 	if err := utils.ParseJSON(r, &product); err != nil {
// 		utils.WriteErrorJSON(w, http.StatusBadRequest, err, nil)
// 		return
// 	}

// 	// Validate payload
// 	if err := utils.Validate.Struct(product); err != nil {
// 		error := fmt.Errorf("Validation Error")
// 		details := utils.GetValidationError(err)
// 		utils.WriteErrorJSON(w, http.StatusBadRequest, error, details)
// 		return
// 	}

// 	// Create product
// 	err := h.store.UpdateProduct(product);
// 	if err != nil {
// 		h.logger.Printf("error creating product: %v", err)
// 		error := fmt.Errorf("Internal Server Error")
// 		utils.WriteErrorJSON(w, http.StatusInternalServerError, error, nil)
// 		return
// 	}

// 	// Write response
// 	response := map[string]interface{}{
// 		"message": "Product updated successfully",
// 	}
// 	utils.WriteJSON(w, http.StatusCreated, response)
// }

// func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
// 	// Parse request body
// 	productIDStr := r.PathValue("id")
// 	if productIDStr == "" {
// 		err := fmt.Errorf("id is required")
// 		utils.WriteErrorJSON(w, http.StatusBadRequest, err, nil)
// 		return
// 	}

// 	productID, err := strconv.Atoi(productIDStr)
// 	if err != nil {
// 		err := fmt.Errorf("Invalid productid: %s", productIDStr)
// 		utils.WriteErrorJSON(w, http.StatusBadRequest, err, nil)
// 		return
// 	}

// 	// Create product
// 	err = h.store.DeleteProduct(productID);
// 	if err != nil {
// 		h.logger.Printf("error deleting product with id `%d`: ", productID, err)
// 		error := fmt.Errorf("Internal Server Error")
// 		utils.WriteErrorJSON(w, http.StatusInternalServerError, error, nil)
// 		return
// 	}

// 	// Write response
// 	response := map[string]interface{}{
// 		"message": "Product updated successfully",
// 	}
// 	utils.WriteJSON(w, http.StatusCreated, response)
// }