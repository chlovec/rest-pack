package types

import "time"

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageUrl    string    `json:"image"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"createdAt"`
}

// Stores
type ProductStore interface {
	CreateProduct(product CreateProductPayload) (int64, error)
	UpdateProduct(product UpdateProductPayload) error
	DeleteProduct(id int) error
	GetProduct(id int) (*Product, error)
	ListProducts(limit int, offset int) ([]*Product, error)
}

// Payloads
type CreateProductPayload struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	ImageUrl    string  `json:"image"`
	Price       float64 `json:"price" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required"`
}

type UpdateProductPayload struct {
	ID          int     `json:"id" validate:"required,gt=0"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	ImageUrl    string  `json:"image"`
	Price       float64 `json:"price" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required"`
}