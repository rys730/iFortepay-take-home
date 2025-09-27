package product

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/rys730/iFortepay-take-home/db"
	"github.com/rys730/iFortepay-take-home/internal/model/entity"
	"github.com/stretchr/testify/require"
)

func TestProductRepository_GetProductByID(t *testing.T) {
	ctx := context.Background()

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Failed to create mock pool: %v", err)
	}
	defer mock.Close()

	repo := &ProductRepository{db: mock}

	productID := int32(1)
	var num pgtype.Numeric
	num.Scan("10.00")
	insertedProduct := db.Product{
		ID:        1,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		Sku:       "test",
		Name:      "test",
		Price:     num,
		Quantity:  10,
	}

	expectedProduct, _ := entity.ProductFromDB(&insertedProduct)

	rows := pgxmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "sku", "name", "price", "quantity"}).
		AddRow(insertedProduct.ID, insertedProduct.CreatedAt, insertedProduct.UpdatedAt, insertedProduct.DeletedAt, insertedProduct.Sku, insertedProduct.Name, insertedProduct.Price, insertedProduct.Quantity)
	mock.ExpectQuery(regexp.QuoteMeta(
		`select id, created_at, updated_at, deleted_at, sku, name, price, quantity 
     from products 
     where id = $1 and deleted_at is null and quantity > 0`)).WithArgs(productID).WillReturnRows(rows)

	product, err := repo.GetProductByID(ctx, productID)
	require.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if !reflect.DeepEqual(product, expectedProduct) {
		t.Errorf("Product is not as expected: \nExpected: %v\nGot: %v", insertedProduct, product)
	}
}

func TestProductRepository_GetProductByID_ProductNotFound(t *testing.T) {
	ctx := context.Background()

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Failed to create mock pool: %v", err)
	}
	defer mock.Close()

	repo := &ProductRepository{db: mock}

	productID := int32(1)

	mock.ExpectQuery(regexp.QuoteMeta(
		`select id, created_at, updated_at, deleted_at, sku, name, price, quantity 
		from products where id = $1 and deleted_at is null and quantity > 0`)).
		WithArgs(productID).WillReturnError(errors.New("product not found"))

	_, err = repo.GetProductByID(ctx, productID)
	require.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if err.Error() != "product not found" {
		t.Errorf("Error is not as expected: \nExpected: %v\nGot: %v", "product not found", err.Error())
	}
}

func TestProductRepository_UpdateProductQuantityByID(t *testing.T) {
	ctx := context.Background()

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Failed to create mock pool: %v", err)
	}
	defer mock.Close()

	repo := &ProductRepository{db: mock}

	productID := int32(1)
	quantity := int32(5)

	var num pgtype.Numeric
	num.Scan("10.00")
	insertedProduct := db.Product{
		ID:        1,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		Sku:       "test",
		Name:      "test",
		Price:     num,
		Quantity:  10,
	}

	// expectedProduct, _ := entity.ProductFromDB(&insertedProduct)

	rows := pgxmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "sku", "name", "price", "quantity"}).
		AddRow(insertedProduct.ID, insertedProduct.CreatedAt, insertedProduct.UpdatedAt, insertedProduct.DeletedAt, insertedProduct.Sku, insertedProduct.Name, insertedProduct.Price, insertedProduct.Quantity)
	mock.ExpectQuery(regexp.QuoteMeta(`update products set quantity = quantity - $1, updated_at = CURRENT_TIMESTAMP 
	where id = $2 and quantity >= $1 and deleted_at is null 
	returning id, created_at, updated_at, deleted_at, sku, name, price, quantity`)).
		WithArgs(quantity, productID).WillReturnRows(rows)

	err = repo.UpdateProductQuantityByID(ctx, productID, quantity)
	require.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestProductRepository_UpdateProductQuantityByID_InternalServerError(t *testing.T) {
	ctx := context.Background()

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Failed to create mock pool: %v", err)
	}
	defer mock.Close()

	repo := &ProductRepository{db: mock}

	productID := int32(1)
	quantity := int32(5)

	mock.ExpectQuery(regexp.QuoteMeta(`update products set quantity = quantity - $1, updated_at = CURRENT_TIMESTAMP 
	where id = $2 and quantity >= $1 and deleted_at is null 
	returning id, created_at, updated_at, deleted_at, sku, name, price, quantity`)).
		WithArgs(quantity, productID).WillReturnError(errors.New("internal server error"))

	err = repo.UpdateProductQuantityByID(ctx, productID, quantity)
	require.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if err.Error() != "internal server error" {
		t.Errorf("Error is not as expected: \nExpected: %v\nGot: %v", "internal server error", err.Error())
	}
}
