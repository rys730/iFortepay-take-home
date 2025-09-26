package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type APIHandler interface {
	RegisterRoutes(e *echo.Group)
}

func RegisterHandlers(e *echo.Echo, handlers []APIHandler) {
	e.Use(middleware.Recover())
	for _, h  := range handlers {
		h.RegisterRoutes(e.Group("/api"))
	}
}