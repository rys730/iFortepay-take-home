package product

import "github.com/labstack/echo/v4"

func (h *ProductHandler) HandleCheckout(c echo.Context) error {
	return c.String(200, "Checkout processed")
}
