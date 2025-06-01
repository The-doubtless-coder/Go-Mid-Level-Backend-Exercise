package main

import (
	"Savannah_Screening_Test/config"
	"Savannah_Screening_Test/migrations"
	"Savannah_Screening_Test/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	db := config.ConnectDatabase()
	err := migrations.AutoMigrate(db)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Migrations complete")

	r := gin.Default()
	routes.RegisterRoutes(r)

	era := godotenv.Load(".env")
	if era != nil {
		log.Fatal("Error loading .env file:: Using os files instead")
	}

	server_port := os.Getenv("SERVER_PORT")
	if err := r.Run(":" + server_port); err != nil {
		log.Fatal("❌ Server error:", err)
	}

}
