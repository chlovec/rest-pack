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

func (s *Store) CreateProduct(types.CreateProductPayload) (int, error) {
	return 0, nil
}