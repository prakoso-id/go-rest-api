package handler

import (
	"net/http"
	"personal-api/internal/entity"
	"personal-api/internal/service"
	"personal-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type PersonalInfoHandler struct {
	service service.PersonalInfoService
}

func NewPersonalInfoHandler(service service.PersonalInfoService) *PersonalInfoHandler {
	return &PersonalInfoHandler{
		service: service,
	}
}

type UpsertPersonalInfoRequest struct {
	FullName  string `json:"full_name" binding:"required"`
	Title     string `json:"title"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
	ResumeURL string `json:"resume_url"`
}

func (h *PersonalInfoHandler) UpsertPersonalInfo(c *gin.Context) {
	var req UpsertPersonalInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	info := &entity.PersonalInfo{
		FullName:  req.FullName,
		Title:     req.Title,
		Bio:       req.Bio,
		AvatarURL: req.AvatarURL,
		ResumeURL: req.ResumeURL,
	}

	if err := h.service.UpsertPersonalInfo(info); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Personal info updated successfully", info))
}

func (h *PersonalInfoHandler) GetPersonalInfo(c *gin.Context) {
	info, err := h.service.GetPersonalInfo()
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("Personal info not found"))
		return
	}

	c.JSON(http.StatusOK, response.Success("Personal info retrieved successfully", info))
}
