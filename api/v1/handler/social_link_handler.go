package handler

import (
	"net/http"
	"personal-api/internal/entity"
	"personal-api/internal/service"
	"personal-api/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SocialLinkHandler struct {
	service service.SocialLinkService
}

func NewSocialLinkHandler(service service.SocialLinkService) *SocialLinkHandler {
	return &SocialLinkHandler{
		service: service,
	}
}

type CreateSocialLinkRequest struct {
	Platform string `json:"platform" binding:"required"`
	URL      string `json:"url" binding:"required"`
	Icon     string `json:"icon"`
}

type UpdateSocialLinkRequest struct {
	Platform string `json:"platform" binding:"required"`
	URL      string `json:"url" binding:"required"`
	Icon     string `json:"icon"`
}

func (h *SocialLinkHandler) CreateSocialLink(c *gin.Context) {
	var req CreateSocialLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	link := &entity.SocialLink{
		Platform: req.Platform,
		URL:      req.URL,
		Icon:     req.Icon,
	}

	if err := h.service.CreateSocialLink(link); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Success("Social link created successfully", link))
}

func (h *SocialLinkHandler) UpdateSocialLink(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid social link ID"))
		return
	}

	var req UpdateSocialLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	link := &entity.SocialLink{
		ID:       uint(id),
		Platform: req.Platform,
		URL:      req.URL,
		Icon:     req.Icon,
	}

	if err := h.service.UpdateSocialLink(link); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Social link updated successfully", link))
}

func (h *SocialLinkHandler) DeleteSocialLink(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid social link ID"))
		return
	}

	if err := h.service.DeleteSocialLink(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Social link deleted successfully", nil))
}

func (h *SocialLinkHandler) GetSocialLink(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid social link ID"))
		return
	}

	link, err := h.service.GetSocialLink(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("Social link not found"))
		return
	}

	c.JSON(http.StatusOK, response.Success("Social link retrieved successfully", link))
}

func (h *SocialLinkHandler) GetAllSocialLinks(c *gin.Context) {
	links, err := h.service.GetAllSocialLinks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Social links retrieved successfully", links))
}
