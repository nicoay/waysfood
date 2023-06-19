package routes

import (
	"mytask/handlers"
	"mytask/pkg/middleware"
	"mytask/pkg/mysql"
	repositories "mytask/repository"

	"github.com/labstack/echo/v4"
)

func RouteOrder(e *echo.Group) {
	OrderRepo := repositories.RepositoryOrder(mysql.DB)
	h := handlers.HandlerOrder(OrderRepo)

	e.GET("/orders", h.FindOrder)
	e.GET("/order/:id", h.GetOrder)
	e.GET("/order-product/:id", middleware.Auth(h.GetOrderByUserProduct))
	e.GET("/order-buyer/:id", middleware.Auth(h.GetOrderByUserSeller))
	e.GET("/order-user", middleware.Auth(h.GetOrderbyUser))
	e.POST("/order", h.CreateOrder)
	e.DELETE("/order/:id", middleware.Auth(h.DeleteOrder))
	e.DELETE("/delete-order", middleware.Auth(h.DeleteAllOrder))
}
