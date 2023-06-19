package main

import (
	"fmt"
	"mytask/database"
	"os"

	"mytask/pkg/mysql"
	"mytask/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic("failed load env file")
	}

	e := echo.New()

	mysql.DatabaseInit()
	database.RunMigration()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PATCH, echo.DELETE},
		AllowHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
	}))

	e.Static("/uploads", "./uploads")

	routes.RouteInit(e.Group("/api/v1"))

	PORT := os.Getenv("PORT")

	fmt.Println("Server is runnning on localhost:" + PORT)
	e.Logger.Fatal(e.Start(":" + PORT))
	// e.Logger.Fatal(e.Start("localhost:8080"))
}
