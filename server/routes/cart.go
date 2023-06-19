package routes

import (
	"mytask/handlers"
	"mytask/pkg/middleware"
	"mytask/pkg/mysql"
	repositories "mytask/repository"

	"github.com/labstack/echo/v4"
)

func RouteCart(e *echo.Group) {
	CartRepo := repositories.RepositoryCart(mysql.DB)
	h := handlers.HandlerCart(CartRepo)

	e.GET("/carts", h.FindCarts)
	e.GET("/cart/:id", h.GetCarts)
	e.POST("/cart", middleware.Auth(h.CreateCarts))
	// e.DELETE("/cart/:id", middleware.Auth(h.DeleteCarts))
}
