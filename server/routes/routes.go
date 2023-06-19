package routes

import "github.com/labstack/echo/v4"

func RouteInit(e *echo.Group) {
	RouteProduct(e)
	AuthRoutes(e)
	RouteCart(e)
	RouteOrder(e)
	RouteUser(e)
	TransactionRoutes(e)

}
