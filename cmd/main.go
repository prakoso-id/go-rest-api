package main

import (
	"log"
	"os"
	"personal-api/config"
	"personal-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	config.InitDB()

	// Auto Migrate the schema
	config.DB.AutoMigrate(&models.User{}, &models.Role{}, &models.Post{})

	// Create default admin role if it doesn't exist
	var adminRole models.Role
	if err := config.DB.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		adminRole = models.Role{
			Name:        "admin",
			Description: "Administrator role with full access",
		}
		config.DB.Create(&adminRole)
	}

	// Initialize Gin router
	r := gin.Default()

	// TODO: Add routes here

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
