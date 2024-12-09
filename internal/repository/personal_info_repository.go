package repository

import (
	"personal-api/internal/entity"

	"gorm.io/gorm"
)

type PersonalInfoRepository interface {
	Create(info *entity.PersonalInfo) error
	Update(info *entity.PersonalInfo) error
	GetPersonalInfo() (*entity.PersonalInfo, error)
}

type personalInfoRepository struct {
	db *gorm.DB
}

func NewPersonalInfoRepository(db *gorm.DB) PersonalInfoRepository {
	return &personalInfoRepository{db: db}
}

func (r *personalInfoRepository) Create(info *entity.PersonalInfo) error {
	return r.db.Create(info).Error
}

func (r *personalInfoRepository) Update(info *entity.PersonalInfo) error {
	return r.db.Save(info).Error
}

func (r *personalInfoRepository) GetPersonalInfo() (*entity.PersonalInfo, error) {
	var info entity.PersonalInfo
	err := r.db.First(&info).Error
	if err != nil {
		return nil, err
	}
	return &info, nil
}
