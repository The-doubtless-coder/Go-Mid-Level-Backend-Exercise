package routes

import (
	"Savannah_Screening_Test/controllers"
	"Savannah_Screening_Test/handler"
	"Savannah_Screening_Test/middlewares"
	"Savannah_Screening_Test/repository"
	"Savannah_Screening_Test/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	api := r.Group("/api/v1")
	//api.Use(middlewares.RequireAuth("customer"))
	//todo: fetch product list - customer
	//todo: fetch orders by custID - customer
	//todo: fetch products by category - customer

	//api.POST("/products", middlewares.RequireAuth("backend_admin"), controllers.CreateProduct)
	api.POST("/products", middlewares.RequireAuth("backend_admin"), productHandler.CreateProduct)
	api.GET("/products", middlewares.RequireAuth("customer"), productHandler.GetProducts) //get_by_category_id

	//api.POST("/orders", middlewares.RequireAuth("customer"), controllers.CreateOrder)
	api.POST("/orders", middlewares.RequireAuth("customer"), orderHandler.CreateOrder)
	api.GET("/orders/by-customer", middlewares.RequireAuth("customer"), orderHandler.GetOrdersByCustomer)

	//api.GET("/categories/average-prices", middlewares.RequireAuth("backend_admin"), controllers.GetAveragePricePerCategory)
	api.GET("/categories/average-prices", middlewares.RequireAuth("backend_admin"), productHandler.GetAveragePricePerCategoryHandler)
	api.POST("/categories", middlewares.RequireAuth("backend_admin"), categoryHandler.CreateCategory)

	api.POST("/users/signup", controllers.SignUpHandler)
	api.POST("/users/signin", controllers.LoginHandler)

	api.POST("/send-message", controllers.SendTestMessage) //test, use goroutine once order is placed
	api.POST("/send-email", controllers.SendTestEmail)     //test, use goroutine once order is placed

}
