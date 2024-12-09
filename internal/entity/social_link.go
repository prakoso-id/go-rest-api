package entity

import (
	"time"

	"gorm.io/gorm"
)

type SocialLink struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Platform  string         `gorm:"not null" json:"platform"`
	URL       string         `gorm:"not null" json:"url"`
	Icon      string         `json:"icon"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
