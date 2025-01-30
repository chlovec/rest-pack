package product

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/chlovec/rest-pack/examples/types"
	"github.com/stretchr/testify/assert"
)

const DbError string = "db error"
const One int64 = 1

func TestGetProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := NewStore(db)

	t.Run("should get product", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "imageUrl", "price", "quantity"}).
			AddRow(1, "Test Product", "New product for testing", "test/image-url", 22.20, 20)

		mock.ExpectQuery("SELECT \\* FROM products WHERE id = \\? LIMIT 1").
			WithArgs(1).
			WillReturnRows(rows)

		product, err := store.GetProduct(1)

		assert.NoError(t, err)
		assert.NotNil(t, product)
	})

	t.Run("should return nil if product does not exist", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "imageUrl", "price", "quantity"})

		mock.ExpectQuery("SELECT \\* FROM products WHERE id = \\? LIMIT 1").
			WithArgs(1).
			WillReturnRows(rows)

		product, err := store.GetProduct(1)

		assert.NoError(t, err)
		assert.Nil(t, product)
	})

	t.Run("should return db error", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM products WHERE id = \\? LIMIT 1").
			WithArgs(1).
			WillReturnError(errors.New(DbError))

		product, err := store.GetProduct(1)

		assert.Error(t, err)
		assert.Equal(t, DbError, err.Error())
		assert.Nil(t, product)
	})
}

func TestCreateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := NewStore(db)
	product := types.CreateProductPayload{
		Name:        "Test Product",
		Description: "New product for testing",
		ImageUrl:    "test/image-url",
		Price:       22.45,
		Quantity:    20,
	}

	t.Run("should create product", func(t *testing.T) {
		// Define expected behavior
		mock.ExpectExec("INSERT INTO products").
			WithArgs(product.Name, product.Description, product.ImageUrl, product.Price, product.Quantity).WillReturnResult(sqlmock.NewResult(One, One))

		res, err := store.CreateProduct(product)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, One, res)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should fail to create product with db err", func(t *testing.T) {
		// Define expected behavior
		mock.ExpectExec("INSERT INTO products").
			WithArgs(product.Name, product.Description, product.ImageUrl, product.Price, product.Quantity).WillReturnError(errors.New(DbError))

		res, err := store.CreateProduct(product)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, DbError, err.Error())
		assert.Equal(t, int64(0), res)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := NewStore(db)
	product := types.UpdateProductPayload{
		ID:          1,
		Name:        "Test Product",
		Description: "New product for testing",
		ImageUrl:    "test/image-url",
		Price:       22.45,
		Quantity:    20,
	}

	t.Run("should update product", func(t *testing.T) {
		// Expect the query to be executed
		mock.ExpectExec(regexp.QuoteMeta("UPDATE products SET name = ?, description = ?, imageUrl = ?, price = ?, quantity = ? WHERE id = ?")).
			WithArgs(product.Name, product.Description, product.ImageUrl, product.Price, product.Quantity, product.ID).
			WillReturnResult(sqlmock.NewResult(One, One))

		err := store.UpdateProduct(product)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should fail with db error", func(t *testing.T) {
		// Expect the query to be executed
		mock.ExpectExec(regexp.QuoteMeta("UPDATE products SET name = ?, description = ?, imageUrl = ?, price = ?, quantity = ? WHERE id = ?")).
			WithArgs(product.Name, product.Description, product.ImageUrl, product.Price, product.Quantity, product.ID).
			WillReturnError(errors.New(DbError))

		err := store.UpdateProduct(product)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, DbError, err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := NewStore(db)
	t.Run("should delete product", func(t *testing.T) {
		// Expect the query to be executed
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM products WHERE id = ?")).
			WithArgs(One).
			WillReturnResult(sqlmock.NewResult(One, One))

		err := store.DeleteProduct(1)

		// Assert
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should fail to delete product with db error", func(t *testing.T) {
		// Expect the query to be executed
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM products WHERE id = ?")).
			WithArgs(One).
			WillReturnError(errors.New(DbError))

		err := store.DeleteProduct(1)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, DbError, err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
