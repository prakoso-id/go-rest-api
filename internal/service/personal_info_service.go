package service

import (
	"personal-api/internal/entity"
	"personal-api/internal/repository"
)

type PersonalInfoService interface {
	UpsertPersonalInfo(info *entity.PersonalInfo) error
	GetPersonalInfo() (*entity.PersonalInfo, error)
}

type personalInfoService struct {
	repo repository.PersonalInfoRepository
}

func NewPersonalInfoService(repo repository.PersonalInfoRepository) PersonalInfoService {
	return &personalInfoService{repo: repo}
}

func (s *personalInfoService) UpsertPersonalInfo(info *entity.PersonalInfo) error {
	existing, err := s.repo.GetPersonalInfo()
	if err != nil {
		// If no record exists, create new
		return s.repo.Create(info)
	}
	
	// Update existing record
	info.ID = existing.ID
	return s.repo.Update(info)
}

func (s *personalInfoService) GetPersonalInfo() (*entity.PersonalInfo, error) {
	return s.repo.GetPersonalInfo()
}
