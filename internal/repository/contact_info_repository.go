package repository

import (
	"personal-api/internal/entity"

	"gorm.io/gorm"
)

type ContactInfoRepository interface {
	Create(info *entity.ContactInfo) error
	Update(info *entity.ContactInfo) error
	GetContactInfo() (*entity.ContactInfo, error)
}

type contactInfoRepository struct {
	db *gorm.DB
}

func NewContactInfoRepository(db *gorm.DB) ContactInfoRepository {
	return &contactInfoRepository{db: db}
}

func (r *contactInfoRepository) Create(info *entity.ContactInfo) error {
	return r.db.Create(info).Error
}

func (r *contactInfoRepository) Update(info *entity.ContactInfo) error {
	return r.db.Save(info).Error
}

func (r *contactInfoRepository) GetContactInfo() (*entity.ContactInfo, error) {
	var info entity.ContactInfo
	err := r.db.First(&info).Error
	if err != nil {
		return nil, err
	}
	return &info, nil
}
