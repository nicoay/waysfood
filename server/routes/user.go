package routes

import (
	"mytask/handlers"
	"mytask/pkg/middleware"
	"mytask/pkg/mysql"
	repositories "mytask/repository"

	"github.com/labstack/echo/v4"
)

func RouteUser(e *echo.Group) {
	UserRepo := repositories.RepositoryUser(mysql.DB)
	h := handlers.UserHandler(UserRepo)

	e.GET("/users", h.FindUsers)
	e.GET("/user", middleware.Auth(h.GetUser))
	e.GET("/partner", h.FindPartner)
	e.PATCH("/update-user", middleware.Auth(middleware.UploadFile(h.UpdateUser)))
	e.DELETE("/user/:id", h.DeleteUser)
}
