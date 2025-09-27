package product

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/rys730/iFortepay-take-home/db"
	"github.com/rys730/iFortepay-take-home/internal/model/entity"
)

func (pr *ProductRepository) GetProductByID(ctx context.Context, productID int32) (entity.Product, error) {
	// Implementation to get product by ID from the database
	q := db.New(pr.db)
	row, err := q.GetProductByID(ctx, productID)
	if err != nil {
		log.Error().Err(err).Msg("failed getting product by id")
		return entity.Product{}, err
	}
	product, err := entity.ProductFromDB(&row)
	if err != nil {
		log.Error().Err(err).Msg("failed converting db product to entity product")
		return entity.Product{}, err
	}
	return product, nil
}

func (pr *ProductRepository) UpdateProductQuantityByID(ctx context.Context, productID int32, quantity int32) error {
	q := db.New(pr.db)
	_, err := q.UpdateProductQuantity(ctx, db.UpdateProductQuantityParams{
		Quantity: quantity,
		ID:       productID,
	})
	if err != nil {
		return err
	}
	return nil
}
