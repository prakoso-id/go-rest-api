package handler

import (
	"net/http"
	"personal-api/internal/entity"
	"personal-api/internal/service"
	"personal-api/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService service.PostService
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error("User not found in context"))
		return
	}

	post := &entity.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID.(uint),
	}

	if err := h.postService.CreatePost(post); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Success("Post created successfully", post))
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid post ID"))
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error("User not found in context"))
		return
	}

	post := &entity.Post{
		ID:      uint(id),
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID.(uint),
	}

	if err := h.postService.UpdatePost(post, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Post updated successfully", post))
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid post ID"))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error("User not found in context"))
		return
	}

	if err := h.postService.DeletePost(uint(id), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Post deleted successfully", nil))
}

func (h *PostHandler) GetPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid post ID"))
		return
	}

	post, err := h.postService.GetPost(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("Post not found"))
		return
	}

	c.JSON(http.StatusOK, response.Success("Post retrieved successfully", post))
}

func (h *PostHandler) GetAllPosts(c *gin.Context) {
	posts, err := h.postService.GetAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Posts retrieved successfully", posts))
}

func (h *PostHandler) GetUserPosts(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error("User not found in context"))
		return
	}

	posts, err := h.postService.GetUserPosts(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("User posts retrieved successfully", posts))
}
