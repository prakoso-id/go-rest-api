package service

import (
	"personal-api/internal/entity"
	"personal-api/internal/repository"
)

type SocialLinkService interface {
	CreateSocialLink(link *entity.SocialLink) error
	UpdateSocialLink(link *entity.SocialLink) error
	DeleteSocialLink(id uint) error
	GetSocialLink(id uint) (*entity.SocialLink, error)
	GetAllSocialLinks() ([]entity.SocialLink, error)
}

type socialLinkService struct {
	repo repository.SocialLinkRepository
}

func NewSocialLinkService(repo repository.SocialLinkRepository) SocialLinkService {
	return &socialLinkService{repo: repo}
}

func (s *socialLinkService) CreateSocialLink(link *entity.SocialLink) error {
	return s.repo.Create(link)
}

func (s *socialLinkService) UpdateSocialLink(link *entity.SocialLink) error {
	return s.repo.Update(link)
}

func (s *socialLinkService) DeleteSocialLink(id uint) error {
	return s.repo.Delete(id)
}

func (s *socialLinkService) GetSocialLink(id uint) (*entity.SocialLink, error) {
	return s.repo.GetByID(id)
}

func (s *socialLinkService) GetAllSocialLinks() ([]entity.SocialLink, error) {
	return s.repo.GetAll()
}
