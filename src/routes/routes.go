package routes

import (
	"github.com/terrytay/ambassador/src/controllers"
	"github.com/terrytay/ambassador/src/middlewares"

	"github.com/labstack/echo/v4"
)

func Setup(app *echo.Echo) {
	api := app.Group("/api")

	// ADMIN
	admin := api.Group("/admin")

	admin.POST("/register", controllers.Register)
	admin.POST("/login", controllers.Login)

	adminAuthenticated := admin.Group("", middlewares.IsAuthenticated)
	adminAuthenticated.POST("/logout", controllers.Logout)
	adminAuthenticated.GET("/user", controllers.User)
	adminAuthenticated.PUT("/users/info", controllers.UpdateInfo)
	adminAuthenticated.PUT("/users/password", controllers.UpdatePassword)
	adminAuthenticated.GET("/ambassadors", controllers.Ambassador)
	adminAuthenticated.GET("/users/:id/links", controllers.Link)

	products := adminAuthenticated.Group("/products")
	products.GET("", controllers.Products)
	products.GET("/:id", controllers.GetProduct)
	products.POST("", controllers.CreateProduct)
	products.PUT("/:id", controllers.UpdateProduct)
	products.DELETE("/:id", controllers.DeleteProduct)

	orders := adminAuthenticated.Group("/orders")
	orders.GET("", controllers.Orders)

	// AMBASSADORS
	ambassador := api.Group("/ambassador")

	ambassador.POST("/register", controllers.Register)
	ambassador.POST("/login", controllers.Login)
	ambassador.GET("/products/frontend", controllers.ProductsFrontend)
	ambassador.GET("/products/backend", controllers.ProductsBackend)

	ambassadorAuthenticated := ambassador.Group("", middlewares.IsAuthenticated)
	ambassadorAuthenticated.POST("/logout", controllers.Logout)
	ambassadorAuthenticated.GET("/user", controllers.User)
	ambassadorAuthenticated.PUT("/users/info", controllers.UpdateInfo)
	ambassadorAuthenticated.PUT("/users/password", controllers.UpdatePassword)
	ambassadorAuthenticated.POST("/links", controllers.CreateLink)
	ambassadorAuthenticated.GET("/stats", controllers.Stats)
	ambassadorAuthenticated.GET("/rankings", controllers.Rankings)
}
