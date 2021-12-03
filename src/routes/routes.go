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

	adminAuthenticated := admin.Group("", middlewares.IsAuthenticated)
	adminAuthenticated.POST("/logout", controllers.Logout)
	adminAuthenticated.GET("/user", controllers.User)
	adminAuthenticated.PUT("/users/info", controllers.UpdateInfo)
	adminAuthenticated.PUT("/users/password", controllers.UpdatePassword)
	adminAuthenticated.GET("/ambassadors", controllers.Ambassador)

	products := adminAuthenticated.Group("/products", middlewares.IsAuthenticated)
	products.GET("", controllers.Products)
	products.GET("/:id", controllers.GetProduct)
	products.POST("", controllers.CreateProduct)
}
