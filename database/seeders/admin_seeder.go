package seeders

import (
	"log"

	"gorm.io/gorm"
	"personal-api/internal/models"
	"personal-api/pkg/utils"
)

func SeedAdminUser(db *gorm.DB) error {
	// Check if admin role exists
	var adminRole models.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		// Create admin role if it doesn't exist
		adminRole = models.Role{
			Name:        "admin",
			Description: "Administrator role with full access",
		}
		if err := db.Create(&adminRole).Error; err != nil {
			return err
		}
	}

	// Check if admin user exists
	var count int64
	db.Model(&models.User{}).Where("email = ?", "admin@example.com").Count(&count)
	if count == 0 {
		// Create admin user if it doesn't exist
		hashedPassword, err := utils.HashPassword("admin123")
		if err != nil {
			return err
		}

		adminUser := models.User{
			Username: "admin",
			Email:    "admin@example.com",
			Password: hashedPassword,
			RoleID:   adminRole.ID,
		}

		if err := db.Create(&adminUser).Error; err != nil {
			return err
		}
		log.Println("Admin user created successfully")
	} else {
		log.Println("Admin user already exists")
	}

	return nil
}
