package repository

import (
	"personal-api/internal/entity"

	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *entity.Post) error
	Update(post *entity.Post) error
	Delete(id uint) error
	FindByID(id uint) (*entity.Post, error)
	FindAll() ([]entity.Post, error)
	FindByUserID(userID uint) ([]entity.Post, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *entity.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) Update(post *entity.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Post{}, id).Error
}

func (r *postRepository) FindByID(id uint) (*entity.Post, error) {
	var post entity.Post
	err := r.db.Preload("User").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) FindAll() ([]entity.Post, error) {
	var posts []entity.Post
	err := r.db.Preload("User").Find(&posts).Error
	return posts, err
}

func (r *postRepository) FindByUserID(userID uint) ([]entity.Post, error) {
	var posts []entity.Post
	err := r.db.Preload("User").Where("user_id = ?", userID).Find(&posts).Error
	return posts, err
}
