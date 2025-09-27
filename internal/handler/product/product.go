package product

import (
	"github.com/labstack/echo/v4"
	"github.com/rys730/iFortepay-take-home/internal/model/dto"
)

// Checkout will process the checkout of products
//
//	@Summary	Checkout products
//	@Tags		api
//	@Accept		json
//	@Produce	json
//	@Param		request		body		dto.CheckoutRequest					true	"request body"
//	@Success	200			{object}	dto.CheckoutResponse	"successful response message"
//	@Router		/api/checkout [post]
func (h *ProductHandler) HandleCheckout(c echo.Context) error {
	var req dto.CheckoutRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	res, err := h.pu.Checkout(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return c.JSON(200, res)
}
