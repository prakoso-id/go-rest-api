package service

import (
	"errors"
	"personal-api/internal/entity"
	"personal-api/internal/repository"
)

type PostService interface {
	CreatePost(post *entity.Post) error
	UpdatePost(post *entity.Post, userID uint) error
	DeletePost(id uint, userID uint) error
	GetPost(id uint) (*entity.Post, error)
	GetAllPosts() ([]entity.Post, error)
	GetUserPosts(userID uint) ([]entity.Post, error)
}

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) PostService {
	return &postService{
		postRepo: postRepo,
	}
}

func (s *postService) CreatePost(post *entity.Post) error {
	return s.postRepo.Create(post)
}

func (s *postService) UpdatePost(post *entity.Post, userID uint) error {
	existingPost, err := s.postRepo.FindByID(post.ID)
	if err != nil {
		return err
	}

	if existingPost.UserID != userID {
		return errors.New("unauthorized to update this post")
	}

	return s.postRepo.Update(post)
}

func (s *postService) DeletePost(id uint, userID uint) error {
	existingPost, err := s.postRepo.FindByID(id)
	if err != nil {
		return err
	}

	if existingPost.UserID != userID {
		return errors.New("unauthorized to delete this post")
	}

	return s.postRepo.Delete(id)
}

func (s *postService) GetPost(id uint) (*entity.Post, error) {
	return s.postRepo.FindByID(id)
}

func (s *postService) GetAllPosts() ([]entity.Post, error) {
	return s.postRepo.FindAll()
}

func (s *postService) GetUserPosts(userID uint) ([]entity.Post, error) {
	return s.postRepo.FindByUserID(userID)
}
