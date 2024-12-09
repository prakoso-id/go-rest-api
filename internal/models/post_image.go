package models

import "time"

type PostImage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PostID    uint      `json:"post_id" gorm:"not null"`
	ImageURL  string    `json:"image_url" gorm:"not null"`
	ImageAlt  string    `json:"image_alt"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePostImageRequest struct {
	PostID   uint   `json:"post_id" binding:"required"`
	ImageURL string `json:"image_url" binding:"required"`
	ImageAlt string `json:"image_alt"`
}

type UpdatePostImageRequest struct {
	ImageURL string `json:"image_url"`
	ImageAlt string `json:"image_alt"`
}
