package product

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/chlovec/rest-pack/examples/config"
	"github.com/chlovec/rest-pack/examples/types"
	"github.com/chlovec/rest-pack/utils"
	"github.com/gorilla/mux"
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
	pageSize, err := strconv.Atoi(query.Get("pagesize"))
	if err != nil || pageSize <= 0 {
		pageSize = defaultPageSize
	}

	// Parse page number
	pageNum, err := strconv.Atoi(query.Get("pagenumber"))
	if err != nil || pageNum < 1 {
		pageNum = defaultPageNum
	} else {
		pageNum-- // Convert to zero-based index
	}

	// Fetch products
	products, err := h.store.ListProducts(pageSize, pageNum * pageSize)
	if err != nil {
		utils.WriteInternalServerError(w, "", nil)
		return
	}

	// Send response
	utils.WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteBadRequest(w, "", nil)
		return
	}

	product, err := h.store.GetProduct(productId)
	if err != nil {
		utils.WriteInternalServerError(w, "", nil)
		return
	}

	if product == nil {
		utils.WriteNotFound(w, "", nil)
		return
	}
	utils.WriteJSON(w, http.StatusOK, product)
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var product types.CreateProductPayload
	if err := utils.ParseJSON(r, &product); err != nil {
		utils.WriteBadRequest(w, "", nil)
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

func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Parse request param
	productID, err := GetProductId(r)
	if err != nil {
		utils.WriteBadRequest(w, "", nil)
		return
	}

	// Parse request body
	var product types.UpdateProductPayload
	err = utils.ParseJSON(r, &product)
	if err != nil || product.ID != productID {
		utils.WriteBadRequest(w, "", nil)
		return
	}

	// Validate payload
	if err := utils.Validate.Struct(product); err != nil {
		details := utils.GetValidationError(err)
		utils.WriteBadRequest(w, "Validation Error", details)
		return
	}

	// Check that product exists
	existingProduct, err := h.store.GetProduct(productID)
	if err != nil {
		utils.WriteInternalServerError(w, "", nil)
		return
	} else if existingProduct == nil {
		utils.WriteBadRequest(w, "", nil)
		return
	}

	// Update product
	err = h.store.UpdateProduct(product);
	if err != nil {
		utils.WriteInternalServerError(w, "", nil)
		return
	}

	// Write response
	response := map[string]interface{}{
		"message": "Product updated successfully",
	}
	utils.WriteJSON(w, http.StatusNoContent, response)
}

func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Parse request param
	productID, err := GetProductId(r)
	if err != nil {
		utils.WriteBadRequest(w, "", nil)
		return
	}

	// Check that product exists
	product, err := h.store.GetProduct(productID)
	if err != nil {
		utils.WriteInternalServerError(w, "", nil)
		return
	} else if product == nil {
		utils.WriteBadRequest(w, "", nil)
		return
	}

	// Create product
	err = h.store.DeleteProduct(productID);
	if err != nil {
		utils.WriteInternalServerError(w, "", nil)
		return
	}

	// Write response
	response := map[string]interface{}{
		"message": "Product deleted successfully",
	}
	utils.WriteJSON(w, http.StatusNoContent, response)
}

func GetProductId(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	return strconv.Atoi(vars["id"])
}