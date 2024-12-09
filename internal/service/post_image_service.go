package service

import (
	"context"
	"errors"

	"personal-api/internal/models"
	"personal-api/internal/repository"
)

type PostImageService interface {
	CreatePostImage(ctx context.Context, req *models.CreatePostImageRequest) (*models.PostImage, error)
	GetPostImageByID(ctx context.Context, id uint) (*models.PostImage, error)
	GetPostImagesByPostID(ctx context.Context, postID uint) ([]models.PostImage, error)
	UpdatePostImage(ctx context.Context, id uint, req *models.UpdatePostImageRequest) (*models.PostImage, error)
	DeletePostImage(ctx context.Context, id uint) error
}

type postImageService struct {
	postImageRepo repository.PostImageRepository
	postRepo     repository.PostRepository
}

func NewPostImageService(postImageRepo repository.PostImageRepository, postRepo repository.PostRepository) PostImageService {
	return &postImageService{
		postImageRepo: postImageRepo,
		postRepo:     postRepo,
	}
}

func (s *postImageService) CreatePostImage(ctx context.Context, req *models.CreatePostImageRequest) (*models.PostImage, error) {
	// Verify post exists
	_, err := s.postRepo.FindByID(req.PostID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	postImage := &models.PostImage{
		PostID:   req.PostID,
		ImageURL: req.ImageURL,
		ImageAlt: req.ImageAlt,
	}

	err = s.postImageRepo.Create(ctx, postImage)
	if err != nil {
		return nil, err
	}

	return postImage, nil
}

func (s *postImageService) GetPostImageByID(ctx context.Context, id uint) (*models.PostImage, error) {
	return s.postImageRepo.GetByID(ctx, id)
}

func (s *postImageService) GetPostImagesByPostID(ctx context.Context, postID uint) ([]models.PostImage, error) {
	return s.postImageRepo.GetByPostID(ctx, postID)
}

func (s *postImageService) UpdatePostImage(ctx context.Context, id uint, req *models.UpdatePostImageRequest) (*models.PostImage, error) {
	postImage, err := s.postImageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.ImageURL != "" {
		postImage.ImageURL = req.ImageURL
	}
	if req.ImageAlt != "" {
		postImage.ImageAlt = req.ImageAlt
	}

	err = s.postImageRepo.Update(ctx, postImage)
	if err != nil {
		return nil, err
	}

	return postImage, nil
}

func (s *postImageService) DeletePostImage(ctx context.Context, id uint) error {
	return s.postImageRepo.Delete(ctx, id)
}
