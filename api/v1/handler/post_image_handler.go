package handler

import (
	"net/http"
	"strconv"

	"personal-api/internal/models"
	"personal-api/internal/service"
	"personal-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type PostImageHandler struct {
	postImageService service.PostImageService
}

func NewPostImageHandler(postImageService service.PostImageService) *PostImageHandler {
	return &PostImageHandler{
		postImageService: postImageService,
	}
}

// CreatePostImage godoc
// @Summary Create a new post image
// @Description Create a new image for a post
// @Tags post-images
// @Accept json
// @Produce json
// @Param request body models.CreatePostImageRequest true "Post image creation request"
// @Success 201 {object} models.PostImage
// @Router /post-images [post]
func (h *PostImageHandler) CreatePostImage(c *gin.Context) {
	var req models.CreatePostImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	postImage, err := h.postImageService.CreatePostImage(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, postImage)
}

// GetPostImagesByPostID godoc
// @Summary Get images by post ID
// @Description Get all images for a specific post
// @Tags post-images
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {array} models.PostImage
// @Router /posts/{id}/images [get]
func (h *PostImageHandler) GetPostImagesByPostID(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	images, err := h.postImageService.GetPostImagesByPostID(c.Request.Context(), uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Post images retrieved successfully", images))
}

// UpdatePostImage godoc
// @Summary Update a post image
// @Description Update an existing post image
// @Tags post-images
// @Accept json
// @Produce json
// @Param id path int true "Post Image ID"
// @Param request body models.UpdatePostImageRequest true "Post image update request"
// @Success 200 {object} models.PostImage
// @Router /post-images/{id} [put]
func (h *PostImageHandler) UpdatePostImage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid ID"))
		return
	}

	var req models.UpdatePostImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	postImage, err := h.postImageService.UpdatePostImage(c.Request.Context(), uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Post image updated successfully", postImage))
}

// DeletePostImage godoc
// @Summary Delete a post image
// @Description Delete an existing post image
// @Tags post-images
// @Produce json
// @Param id path int true "Post Image ID"
// @Success 204 "No Content"
// @Router /post-images/{id} [delete]
func (h *PostImageHandler) DeletePostImage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid ID"))
		return
	}

	err = h.postImageService.DeletePostImage(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Post image deleted successfully", nil))
}
