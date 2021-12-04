package main

import (
	"github.com/terrytay/ambassador/src/database"
	"github.com/terrytay/ambassador/src/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	database.Connect()
	database.AutoMigrate()
	database.SetupRedis()
	database.SetupCacheChannel()

	app := echo.New()

	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	app.Start(":8000")

}
