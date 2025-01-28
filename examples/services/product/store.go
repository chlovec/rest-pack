package product

import (
	"database/sql"

	"github.com/chlovec/rest-pack/examples/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateProduct(product types.CreateProductPayload) (int64, error) {
	query := "INSERT INTO products(name, description, ImageUrl, price, quantity) VALUES(?, ?, ?, ?, ?)"
	res, err := s.db.Exec(query, product.Name, product.Description, product.ImageUrl, product.Price, product.Quantity)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}