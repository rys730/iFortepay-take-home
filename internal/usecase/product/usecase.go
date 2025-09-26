package product

import "github.com/rys730/iFortepay-take-home/internal/repository/product"

type Usecase interface {}

type ProductUsecase struct {
	pr product.Repository
}

func NewProductUsecase(pr product.Repository) Usecase {
	return &ProductUsecase{
		pr: pr,
	}
}