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
	
	// Transaction methods
	BeginTx() (*gorm.DB, error)
	CreateTx(tx *gorm.DB, post *entity.Post) error
	UpdateTx(tx *gorm.DB, post *entity.Post) error
	CreateImagesTx(tx *gorm.DB, images []entity.PostImage) error
	DeleteImagesByPostIDTx(tx *gorm.DB, postID uint) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{
		db: db,
	}
}

func (r *postRepository) BeginTx() (*gorm.DB, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (r *postRepository) CreateTx(tx *gorm.DB, post *entity.Post) error {
	return tx.Create(post).Error
}

func (r *postRepository) UpdateTx(tx *gorm.DB, post *entity.Post) error {
	return tx.Save(post).Error
}

func (r *postRepository) CreateImagesTx(tx *gorm.DB, images []entity.PostImage) error {
	if len(images) == 0 {
		return nil
	}
	return tx.Create(&images).Error
}

func (r *postRepository) DeleteImagesByPostIDTx(tx *gorm.DB, postID uint) error {
	return tx.Where("post_id = ?", postID).Delete(&entity.PostImage{}).Error
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
	if err := r.db.Preload("User").Preload("Images").First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) FindAll() ([]entity.Post, error) {
	var posts []entity.Post
	if err := r.db.Preload("User").Preload("Images").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) FindByUserID(userID uint) ([]entity.Post, error) {
	var posts []entity.Post
	if err := r.db.Preload("User").Preload("Images").Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
