package product

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetProductByID(ctx context.Context, productID int32) error
	UpdateProductQuantityByID(ctx context.Context, productID int32, quantity int32) error
}

type ProductRepository struct {
	db *pgxpool.Pool	
}

func NewProductRepository(db *pgxpool.Pool) Repository {
	return &ProductRepository{
		db: db,
	}
}