package handler

import (
	"net/http"
	"personal-api/internal/entity"
	"personal-api/internal/service"
	"personal-api/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login credentials"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request format"))
		return
	}

	user, err := h.authService.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error("Invalid email or password"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, response.Error("Invalid email or password"))
		return
	}

	token, err := h.authService.GenerateToken(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to generate token"))
		return
	}

	c.JSON(http.StatusOK, response.Success("Login successful", gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role_id":  user.RoleID,
		},
	}))
}

// Register godoc
// @Summary Register new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param register body RegisterRequest true "Registration details"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request format"))
		return
	}

	// Check if email already exists
	existingUser, _ := h.authService.GetUserByEmail(c.Request.Context(), req.Email)
	if existingUser != nil {
		c.JSON(http.StatusBadRequest, response.Error("Email already registered"))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to process registration"))
		return
	}

	user := &entity.User{
		Username: req.Name,  // Using Name as Username
		Email:    req.Email,
		Password: string(hashedPassword),
		RoleID:   2,  // Default role ID for regular users
	}

	if err := h.authService.CreateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to create user"))
		return
	}

	c.JSON(http.StatusCreated, response.Success("User registered successfully", gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role_id":  user.RoleID,
		},
	}))
}
