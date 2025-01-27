package product

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/chlovec/rest-pack/examples/types"
	"github.com/stretchr/testify/assert"
)

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

	t.Run("should create user", func(t *testing.T) {
		// Define expected behavior
		mock.ExpectExec("INSERT INTO products").
			WithArgs("Test Product", "New product for testing", "test/image-url", 22.45, 20).WillReturnResult(sqlmock.NewResult(1, 1))

		res, err := store.CreateProduct(product)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, int64(1), res)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should fail to create user with db err", func(t *testing.T) {
		// Define expected behavior
		mock.ExpectExec("INSERT INTO products").
			WithArgs("Test Product", "New product for testing", "test/image-url", 22.45, 20).WillReturnError(errors.New("db error"))

		res, err := store.CreateProduct(product)
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, int64(0), res)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
