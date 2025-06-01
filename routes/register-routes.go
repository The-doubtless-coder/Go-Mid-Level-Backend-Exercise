package routes

import (
	"Savannah_Screening_Test/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	api := r.Group("/api/v1")

	api.POST("/products", controllers.CreateProduct)
	api.POST("/orders", controllers.CreateOrder)
	api.GET("/categories/average-prices", controllers.GetAveragePricePerCategory)
	api.POST("/categories", controllers.CreateCategory)
	api.POST("/customers", controllers.CreateCustomer)
}
