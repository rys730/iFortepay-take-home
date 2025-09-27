package product

import (
	"context"

	"github.com/rys730/iFortepay-take-home/internal/model/dto"
	"github.com/rys730/iFortepay-take-home/internal/repository/product"
	"github.com/rys730/iFortepay-take-home/internal/repository/promotion"
)

type Usecase interface {
	Checkout(ctx context.Context, req dto.CheckoutRequest) (dto.CheckoutResponse, error)
}

type ProductUsecase struct {
	productRepo product.Repository
	promoRepo   promotion.Repository
}

func NewProductUsecase(
	productRepo product.Repository,
	promoRepo promotion.Repository,
) Usecase {
	return &ProductUsecase{
		productRepo: productRepo,
		promoRepo:   promoRepo,
	}
}
