package routes

import (
	"mytask/handlers"
	"mytask/pkg/middleware"
	"mytask/pkg/mysql"
	repositories "mytask/repository"

	"github.com/labstack/echo/v4"
)

func AuthRoutes(e *echo.Group) {
	AuthRepository := repositories.RepositoryAuth(mysql.DB)
	h := handlers.HandlerAuth(AuthRepository)

	e.POST("/register", middleware.UploadFile(h.Register))
	e.POST("/login", h.Login)
	e.GET("/check-auth", middleware.Auth(h.CheckAuth)) // add this code
}
