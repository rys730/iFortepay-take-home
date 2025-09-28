package product

import (
	"github.com/labstack/echo/v4"
	"github.com/rys730/iFortepay-take-home/internal/handler"
	productUsecase "github.com/rys730/iFortepay-take-home/internal/usecase/product"
)

type ProductHandler struct {
	pu productUsecase.Usecase
}

func NewProductHandler(pu productUsecase.Usecase) handler.APIHandler {
	return &ProductHandler{
		pu: pu,
	}
}

func (h *ProductHandler) RegisterRoutes(e *echo.Group) {
	e.POST("/checkout", h.HandleCheckout)
}
