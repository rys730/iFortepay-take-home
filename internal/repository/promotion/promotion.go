package promotion

import (
	"context"

	"github.com/rys730/iFortepay-take-home/db"
	"github.com/rys730/iFortepay-take-home/internal/model/entity"
)

func (pr *PromotionRepository) GetProductPromotionsByProductID(ctx context.Context, productID int32) ([]entity.ProductPromotion, error) {
	q := db.New(pr.db)
	rows, err := q.GetProductPromotionsByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	var promotions []entity.ProductPromotion
	for _, r := range rows {
		pricePg, err := r.Discount.Float64Value()
		if err != nil {
			return promotions, err
		}
		var discount *float64
		if pricePg.Valid {
			discount = &pricePg.Float64
		}
		pp := entity.ProductPromotion{
			PromotionType: r.PromotionType,
			ID:            r.ID,
			CreatedAt:     r.CreatedAt.Time,
			UpdatedAt:     r.UpdatedAt.Time,
			DeletedAt:     r.DeletedAt.Time,
			PromotionID:   r.PromotionID,
			ProductID:     r.ProductID,
			MinQuantity:   r.MinQuantity,
			FreeProductID: r.FreeProductID,
			Discount:      discount,
			FreeQuantity:  r.FreeQuantity,
			PayY:          r.PayY,
		}
		promotions = append(promotions, pp)
	}
	return promotions, nil
}
