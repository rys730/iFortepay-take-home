package app

import (
	"github.com/labstack/echo/v4"
	_ "github.com/rys730/iFortepay-take-home/docs"
	"github.com/rys730/iFortepay-take-home/infrastructure/postgres"
	"github.com/rys730/iFortepay-take-home/internal/common/config"
	"github.com/rys730/iFortepay-take-home/internal/handler"
	"github.com/rys730/iFortepay-take-home/internal/handler/product"
	productRepo "github.com/rys730/iFortepay-take-home/internal/repository/product"
	promotionRepo "github.com/rys730/iFortepay-take-home/internal/repository/promotion"
	productUsecase "github.com/rys730/iFortepay-take-home/internal/usecase/product"
)

// @title Swagger Example API
// @version 1.0
func CreateApp(cfg *config.Config) *echo.Echo {
	db := postgres.NewPostgres(&cfg.DB)

	productRepo := productRepo.NewProductRepository(db.Conn)
	promotionRepo := promotionRepo.NewPromotionRepository(db.Conn)
	productUc := productUsecase.NewProductUsecase(productRepo, promotionRepo)
	productHandler := product.NewProductHandler(productUc)

	e := echo.New()
	handler.RegisterHandlers(e, []handler.APIHandler{
		productHandler,
	})
	return e
}
