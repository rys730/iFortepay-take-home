package promotion

import (
	"context"

	"github.com/rys730/iFortepay-take-home/infrastructure/postgres"
	"github.com/rys730/iFortepay-take-home/internal/model/entity"
)

type Repository interface {
	GetProductPromotionsByProductID(ctx context.Context, productID int32) ([]entity.ProductPromotion, error)
}

type PromotionRepository struct {
	db postgres.Pool
}

func NewPromotionRepository(db postgres.Pool) Repository {
	return &PromotionRepository{
		db: db,
	}
}
