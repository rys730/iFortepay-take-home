package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type APIHandler interface {
	RegisterRoutes(e *echo.Group)
}

func RegisterHandlers(e *echo.Echo, handlers []APIHandler) {
	e.Use(middleware.Recover())
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
		})
	})
	e.GET("/docs/*", echoSwagger.WrapHandler)
	for _, h := range handlers {
		h.RegisterRoutes(e.Group("/api"))
	}
}
