package product

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/rys730/iFortepay-take-home/internal/model/dto"
	"github.com/rys730/iFortepay-take-home/internal/model/entity"
)

func (pu *ProductUsecase) Checkout(ctx context.Context, req dto.CheckoutRequest) (dto.CheckoutResponse, error) {
	finalPrice := float64(0)
	checkedoutItems := []entity.CheckoutItem{}

	// check item exists
	for _, item := range req.Items {
		product, err := pu.productRepo.GetProductByID(ctx, item.ID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return dto.CheckoutResponse{}, echo.NewHTTPError(404, fmt.Sprintf("product with id %d not found or out of stock", item.ID))
			}
			return dto.CheckoutResponse{}, echo.NewHTTPError(500, "internal server error")
		}
		if product.Quantity < item.Quantity {
			return dto.CheckoutResponse{}, echo.NewHTTPError(400, fmt.Sprintf("product with id %d has only %d items left in stock", item.ID, product.Quantity))
		}
		// get item promotions by product id
		promotion, err := pu.promoRepo.GetProductPromotionsByProductID(ctx, req.Items[0].ID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return dto.CheckoutResponse{}, echo.ErrInternalServerError
		}
		// apply promotion rules
		checkoutItems, err := pu.ApplyPromotionRules(ctx, product, item.Quantity, promotion)
		if err != nil {
			return dto.CheckoutResponse{}, echo.ErrInternalServerError
		}
		checkedoutItems = append(checkedoutItems, checkoutItems...)
	}

	// recalculate total price
	for _, ci := range checkedoutItems {
		finalPrice += ci.TotalPrice
		err := pu.productRepo.UpdateProductQuantityByID(ctx, ci.ProductID, ci.Quantity)
		if err != nil {
			log.Error().Err(err).Msgf("failed updating product quantity for product id %d", ci.ProductID)
			return dto.CheckoutResponse{}, echo.ErrInternalServerError
		}
	}

	return dto.CheckoutResponse{
		Message:    "Checkout successful",
		TotalPrice: fmt.Sprintf("$%.2f", finalPrice),
		Items:      checkedoutItems,
	}, nil
}

func (pu *ProductUsecase) ApplyPromotionRules(ctx context.Context, product entity.Product, buyQuantity int32, promotions []entity.ProductPromotion) ([]entity.CheckoutItem, error) {
	checkoutItems := []entity.CheckoutItem{}
	if len(promotions) == 0 {
		checkoutItems = append(checkoutItems, entity.CheckoutItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			Quantity:    buyQuantity,
			TotalPrice:  product.Price * float64(buyQuantity),
		})
		return checkoutItems, nil
	}
	for _, promo := range promotions {
		switch promo.PromotionType {
		case "FREE_ITEM":
			checkoutItems = append(checkoutItems, entity.CheckoutItem{
				ProductID:   product.ID,
				ProductName: product.Name,
				Quantity:    buyQuantity,
				TotalPrice:  product.Price * float64(buyQuantity),
			})

			if buyQuantity >= promo.MinQuantity && promo.FreeProductID != nil && promo.FreeQuantity != nil {
				freeItem, err := pu.productRepo.GetProductByID(ctx, *promo.FreeProductID)
				if err != nil {
					if errors.Is(err, pgx.ErrNoRows) {
						log.Info().Msgf("%v/n", err)
						log.Info().Msgf("free item with id %d not found or out of stock, skipping promotion", *promo.FreeProductID)
					} else {
						return nil, err
					}
				} else {
					freeItemQuantity := buyQuantity / promo.MinQuantity * (*promo.FreeQuantity)
					if freeItemQuantity > freeItem.Quantity {
						freeItemQuantity = freeItem.Quantity
					}
					checkoutItems = append(checkoutItems, entity.CheckoutItem{
						ProductID:   freeItem.ID,
						ProductName: freeItem.Name,
						Quantity:    freeItemQuantity,
						TotalPrice:  0,
					})
				}
			}
		case "BUY_X_PAY_Y":
			if buyQuantity >= promo.MinQuantity && promo.PayY != nil {
				payQuantity := (buyQuantity / promo.MinQuantity * (*promo.PayY)) + (buyQuantity % promo.MinQuantity)
				checkoutItems = append(checkoutItems, entity.CheckoutItem{
					ProductID:   product.ID,
					ProductName: product.Name,
					Quantity:    buyQuantity,
					TotalPrice:  float64(payQuantity) * product.Price,
				})
			} else {
				checkoutItems = append(checkoutItems, entity.CheckoutItem{
					ProductID:   product.ID,
					ProductName: product.Name,
					Quantity:    buyQuantity,
					TotalPrice:  product.Price * float64(buyQuantity),
				})
			}
		case "BULK_DISCOUNT":
			if buyQuantity >= promo.MinQuantity && promo.Discount != nil {
				totalPrice := float64(buyQuantity) * product.Price
				discount := totalPrice * *promo.Discount
				discountedPrice := totalPrice - discount
				checkoutItems = append(checkoutItems, entity.CheckoutItem{
					ProductID:   product.ID,
					ProductName: product.Name,
					Quantity:    buyQuantity,
					TotalPrice:  discountedPrice,
				})
			} else {
				checkoutItems = append(checkoutItems, entity.CheckoutItem{
					ProductID:   product.ID,
					ProductName: product.Name,
					Quantity:    buyQuantity,
					TotalPrice:  product.Price * float64(buyQuantity),
				})
			}
		default:
			checkoutItems = append(checkoutItems, entity.CheckoutItem{
				ProductID:   product.ID,
				ProductName: product.Name,
				Quantity:    buyQuantity,
				TotalPrice:  product.Price * float64(buyQuantity),
			})
		}
	}

	return checkoutItems, nil
}
