package main

import (
	"log"
	"os"
	"personal-api/api/v1/handler"
	"personal-api/api/v1/routes"
	"personal-api/configs"
	"personal-api/internal/entity"
	"personal-api/internal/models"
	"personal-api/internal/repository"
	"personal-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	db := configs.InitDB()

	// Auto Migrate the schema
	db.AutoMigrate(
		&entity.User{},
		&entity.Role{},
		&entity.Post{},
		&entity.PersonalInfo{},
		&entity.ContactInfo{},
		&entity.SocialLink{},
		&models.PostImage{},
	)

	// Create default roles if they don't exist
	var adminRole entity.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		adminRole = entity.Role{
			Name:        "admin",
			Description: "Administrator role with full access",
		}
		db.Create(&adminRole)
	}

	var userRole entity.Role
	if err := db.Where("name = ?", "user").First(&userRole).Error; err != nil {
		userRole = entity.Role{
			Name:        "user",
			Description: "Regular user role",
		}
		db.Create(&userRole)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	personalInfoRepo := repository.NewPersonalInfoRepository(db)
	contactInfoRepo := repository.NewContactInfoRepository(db)
	socialLinkRepo := repository.NewSocialLinkRepository(db)
	postImageRepo := repository.NewPostImageRepository(db)

	// Initialize services
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(jwtKey) == 0 {
		jwtKey = []byte("your-256-bit-secret") // Default key for development
	}
	authService := service.NewAuthService(userRepo, jwtKey)
	postService := service.NewPostService(postRepo)
	personalInfoService := service.NewPersonalInfoService(personalInfoRepo)
	contactInfoService := service.NewContactInfoService(contactInfoRepo)
	socialLinkService := service.NewSocialLinkService(socialLinkRepo)
	postImageService := service.NewPostImageService(postImageRepo, postRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	postHandler := handler.NewPostHandler(postService)
	personalInfoHandler := handler.NewPersonalInfoHandler(personalInfoService)
	contactInfoHandler := handler.NewContactInfoHandler(contactInfoService)
	socialLinkHandler := handler.NewSocialLinkHandler(socialLinkService)
	postImageHandler := handler.NewPostImageHandler(postImageService)

	// Initialize Gin router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r, authHandler, postHandler, personalInfoHandler, contactInfoHandler, socialLinkHandler, postImageHandler)

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
