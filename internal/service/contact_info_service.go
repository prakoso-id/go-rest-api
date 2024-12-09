package service

import (
	"personal-api/internal/entity"
	"personal-api/internal/repository"
)

type ContactInfoService interface {
	UpsertContactInfo(info *entity.ContactInfo) error
	GetContactInfo() (*entity.ContactInfo, error)
}

type contactInfoService struct {
	repo repository.ContactInfoRepository
}

func NewContactInfoService(repo repository.ContactInfoRepository) ContactInfoService {
	return &contactInfoService{repo: repo}
}

func (s *contactInfoService) UpsertContactInfo(info *entity.ContactInfo) error {
	existing, err := s.repo.GetContactInfo()
	if err != nil {
		// If no record exists, create new
		return s.repo.Create(info)
	}
	
	// Update existing record
	info.ID = existing.ID
	return s.repo.Update(info)
}

func (s *contactInfoService) GetContactInfo() (*entity.ContactInfo, error) {
	return s.repo.GetContactInfo()
}
