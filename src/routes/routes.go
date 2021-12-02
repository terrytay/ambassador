package routes

import (
	"github.com/terrytay/ambassador/src/controllers"
	"github.com/terrytay/ambassador/src/middlewares"

	"github.com/labstack/echo/v4"
)

func Setup(app *echo.Echo) {
	api := app.Group("/api")

	admin := api.Group("/admin")

	admin.POST("/register", controllers.Register)
	admin.POST("/login", controllers.Login)

	admin.POST("/logout", controllers.Logout, middlewares.IsAuthenticated)
	admin.GET("/user", controllers.User, middlewares.IsAuthenticated)
	admin.PUT("/users/info", controllers.UpdateInfo, middlewares.IsAuthenticated)
	admin.PUT("/users/password", controllers.UpdatePassword, middlewares.IsAuthenticated)

}
