package entity

import (
	"time"

	"gorm.io/gorm"
)

type ContactInfo struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Email     string         `json:"email"`
	Phone     string         `json:"phone"`
	Address   string         `json:"address"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
