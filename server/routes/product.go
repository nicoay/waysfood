package routes

import (
	"mytask/handlers"
	"mytask/pkg/middleware"
	"mytask/pkg/mysql"
	repositories "mytask/repository"

	"github.com/labstack/echo/v4"
)

func RouteProduct(e *echo.Group) {
	ProductRepo := repositories.RepositoryProduct(mysql.DB)
	UserRepo := repositories.RepositoryUser(mysql.DB)
	h := handlers.HandlerProduct(ProductRepo, UserRepo)

	e.GET("/products", h.FindProduct)
	e.GET("/product/:id", h.GetProduct)
	e.GET("/product-partner/:id", h.FindProductPartner)
	e.POST("/product", middleware.Auth(middleware.UploadFile(h.CreateProduct)))
	// e.PATCH("/product/:id", middleware.Auth(h.UpdateProduct))
	// e.DELETE("/product/:id", middleware.Auth(h.DeleteProduct))

}
