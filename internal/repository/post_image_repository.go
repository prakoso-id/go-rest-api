package repository

import (
	"context"

	"gorm.io/gorm"
	"personal-api/internal/models"
)

type PostImageRepository interface {
	Create(ctx context.Context, postImage *models.PostImage) error
	GetByID(ctx context.Context, id uint) (*models.PostImage, error)
	GetByPostID(ctx context.Context, postID uint) ([]models.PostImage, error)
	Update(ctx context.Context, postImage *models.PostImage) error
	Delete(ctx context.Context, id uint) error
}

type postImageRepository struct {
	db *gorm.DB
}

func NewPostImageRepository(db *gorm.DB) PostImageRepository {
	return &postImageRepository{db: db}
}

func (r *postImageRepository) Create(ctx context.Context, postImage *models.PostImage) error {
	return r.db.WithContext(ctx).Create(postImage).Error
}

func (r *postImageRepository) GetByID(ctx context.Context, id uint) (*models.PostImage, error) {
	var postImage models.PostImage
	err := r.db.WithContext(ctx).First(&postImage, id).Error
	if err != nil {
		return nil, err
	}
	return &postImage, nil
}

func (r *postImageRepository) GetByPostID(ctx context.Context, postID uint) ([]models.PostImage, error) {
	var postImages []models.PostImage
	err := r.db.WithContext(ctx).Where("post_id = ?", postID).Find(&postImages).Error
	if err != nil {
		return nil, err
	}
	return postImages, nil
}

func (r *postImageRepository) Update(ctx context.Context, postImage *models.PostImage) error {
	return r.db.WithContext(ctx).Save(postImage).Error
}

func (r *postImageRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.PostImage{}, id).Error
}
