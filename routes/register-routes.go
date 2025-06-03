package routes

import (
	"Savannah_Screening_Test/controllers"
	"Savannah_Screening_Test/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	api := r.Group("/api/v1")
	//api.Use(middlewares.RequireAuth("customer"))

	//todo: fetch product list - customer
	//todo: fetch orders by custID - customer
	//todo: fetch products by category - customer

	api.POST("/products", middlewares.RequireAuth("backend_admin"), controllers.CreateProduct)
	api.GET("/products", middlewares.RequireAuth("customer"), controllers.GetProducts)
	api.POST("/orders", middlewares.RequireAuth("customer"), controllers.CreateOrder)
	api.GET("/orders/by-customer", middlewares.RequireAuth("customer"), controllers.GetOrdersByCustomer)
	api.GET("/categories/average-prices", middlewares.RequireAuth("backend_admin"), controllers.GetAveragePricePerCategory)
	api.POST("/categories", middlewares.RequireAuth("backend_admin"), controllers.CreateCategory)
	api.POST("/users/signup", controllers.SignUpHandler)
	api.POST("/users/signin", controllers.LoginHandler)
}
