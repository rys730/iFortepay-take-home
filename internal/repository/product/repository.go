package product

import (
	"context"

	"github.com/rys730/iFortepay-take-home/infrastructure/postgres"
	"github.com/rys730/iFortepay-take-home/internal/model/entity"
)

type Repository interface {
	GetProductByID(ctx context.Context, productID int32) (entity.Product, error)
	UpdateProductQuantityByID(ctx context.Context, productID int32, quantity int32) error
}

type ProductRepository struct {
	db postgres.Pool
}

func NewProductRepository(db postgres.Pool) Repository {
	return &ProductRepository{
		db: db,
	}
}
