package entity

import (
	"time"

	"gorm.io/gorm"
)

type PersonalInfo struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	FullName  string         `gorm:"not null" json:"full_name"`
	Title     string         `json:"title"`
	Bio       string         `json:"bio"`
	AvatarURL string         `json:"avatar_url"`
	ResumeURL string         `json:"resume_url"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
