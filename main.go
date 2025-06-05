package main

import (
	"Savannah_Screening_Test/config"
	"Savannah_Screening_Test/migrations"
	"Savannah_Screening_Test/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"

	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "Savannah_Screening_Test/docs"
)

// @title Savannah Screening API
// @version 1.0
// @description API for handling customer signups, orders, and product data.
// @termsOfService http://swagger.io/terms/

// @contact.name Ian Wanyonyi Wanjala
// @contact.email your-email@example.com

// @host localhost:8088
// @BasePath /

// @securityDefinitions.oauth2.password OAuth2Password
// @tokenUrl http://localhost:8080/realms/master/protocol/openid-connect/token
// @scope profile email
func main() {

	db := config.ConnectDatabase()
	err := migrations.AutoMigrate(db)
	if err != nil {
		log.Fatal("Migration failed:", err)
	} else {
		log.Println("Migration succeeded")
	}

	r := gin.Default()
	r.GET("/swagger/*any" /*security middleware if I want*/, ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.RegisterRoutes(r, db)

	era := godotenv.Load(".env")
	if era != nil {
		log.Fatal("Error loading .env file:: Using os files instead")
	}

	server_port := os.Getenv("SERVER_PORT")
	if err := r.Run(":" + server_port); err != nil {
		log.Fatal("❌ Server error:", err)
	}

}
