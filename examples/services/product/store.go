package product

import (
	"database/sql"
	"errors"

	"github.com/chlovec/rest-pack/examples/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProduct(id int) (*types.Product, error) {
	query := "SELECT * FROM products WHERE id = ? LIMIT 1"
	row := s.db.QueryRow(query, id)
	product, err := scanProductRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return product, nil
}

func (s *Store) ListProducts(limit int, offset int) ([]*types.Product, error) {
	if limit <= 0 || limit > 1000 {
		limit = 1000
	}

	query := "SELECT * FROM products ORDER BY id ASC LIMIT ? OFFSET ?"
	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan results
	products := []*types.Product{}
	for rows.Next() {
		product, err := scanProductRow(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (s *Store) CreateProduct(product types.CreateProductPayload) (int64, error) {
	query := "INSERT INTO products(name, description, ImageUrl, price, quantity) VALUES(?, ?, ?, ?, ?)"
	res, err := s.db.Exec(query, product.Name, product.Description, product.ImageUrl, product.Price, product.Quantity)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (s *Store) UpdateProduct(product types.UpdateProductPayload) error {
	query := "UPDATE products SET name = ?, description = ?, imageUrl = ?, price = ?, quantity = ? WHERE id = ?"
	_, err := s.db.Exec(query, product.Name, product.Description, product.ImageUrl, product.Price, product.Quantity, product.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteProduct(id int) error {
	query := "DELETE FROM products WHERE id = ?"
	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func scanProductRow(scanner interface{ Scan(dest ...interface{}) error }) (*types.Product, error) {
	var product types.Product
	err := scanner.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.ImageUrl,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &product, nil
}
