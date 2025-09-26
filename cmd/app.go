package main

import (
	"github.com/labstack/echo/v4"
	"github.com/rys730/iFortepay-take-home/infrastructure/postgres"
	"github.com/rys730/iFortepay-take-home/internal/common/config"
	"github.com/rys730/iFortepay-take-home/internal/handler"
	"github.com/rys730/iFortepay-take-home/internal/handler/product"
	productRepo "github.com/rys730/iFortepay-take-home/internal/repository/product"
	productUsecase "github.com/rys730/iFortepay-take-home/internal/usecase/product"
)

func CreateApp(cfg *config.Config) *echo.Echo {
	db := postgres.NewPostgres(&cfg.DB)

	productRepo := productRepo.NewProductRepository(db.Conn)
	productUc := productUsecase.NewProductUsecase(productRepo)
	productHandler := product.NewProductHandler(productUc)

	e := echo.New()
	handler.RegisterHandlers(e, []handler.APIHandler{
		productHandler,
	})
	return e
}
