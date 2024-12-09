package service

import (
	"context"
	"errors"
	"personal-api/internal/entity"
	"personal-api/internal/repository"
	"personal-api/pkg/utils"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailTaken         = errors.New("email already registered")
	ErrTokenGeneration    = errors.New("failed to generate token")
)

type AuthService interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
	GenerateToken(ctx context.Context, user *entity.User) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
	jwtKey   []byte
}

func NewAuthService(userRepo repository.UserRepository, jwtKey []byte) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtKey:   jwtKey,
	}
}

func (s *authService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *authService) CreateUser(ctx context.Context, user *entity.User) error {
	existingUser, _ := s.userRepo.FindByEmail(user.Email)
	if existingUser != nil {
		return ErrEmailTaken
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return s.userRepo.Create(user)
}

func (s *authService) GenerateToken(ctx context.Context, user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role.Name,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", ErrTokenGeneration
	}

	return signedToken, nil
}
