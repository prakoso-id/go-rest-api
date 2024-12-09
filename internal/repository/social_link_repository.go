package repository

import (
	"personal-api/internal/entity"

	"gorm.io/gorm"
)

type SocialLinkRepository interface {
	Create(link *entity.SocialLink) error
	Update(link *entity.SocialLink) error
	Delete(id uint) error
	GetByID(id uint) (*entity.SocialLink, error)
	GetAll() ([]entity.SocialLink, error)
}

type socialLinkRepository struct {
	db *gorm.DB
}

func NewSocialLinkRepository(db *gorm.DB) SocialLinkRepository {
	return &socialLinkRepository{db: db}
}

func (r *socialLinkRepository) Create(link *entity.SocialLink) error {
	return r.db.Create(link).Error
}

func (r *socialLinkRepository) Update(link *entity.SocialLink) error {
	return r.db.Save(link).Error
}

func (r *socialLinkRepository) Delete(id uint) error {
	return r.db.Delete(&entity.SocialLink{}, id).Error
}

func (r *socialLinkRepository) GetByID(id uint) (*entity.SocialLink, error) {
	var link entity.SocialLink
	err := r.db.First(&link, id).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *socialLinkRepository) GetAll() ([]entity.SocialLink, error) {
	var links []entity.SocialLink
	err := r.db.Find(&links).Error
	return links, err
}
