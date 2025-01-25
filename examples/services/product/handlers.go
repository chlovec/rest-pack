package product

import (
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct {
	logger *log.Logger
}

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price float64 `json:"price"`
	Qty   int     `json:"qty"`
}

func NewHandler(logger *log.Logger) *Handler {
	return &Handler{}
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products := []Product{
		{
			ID:    1,
			Name:  "Sample Product",
			Price: 99.99,
			Qty:   4,
		},
		{
			ID:    2,
			Name:  "Real Product",
			Price: 100.23,
			Qty:   22,
		},
	}

	// Serialize the product to JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		// Handle JSON encoding errors
		h.logger.Println("Error encoding products to JSON:", err)
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
		return
	}
}